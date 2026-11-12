package serviceimpl

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	websocketInfra "gofiber-template/infrastructure/websocket"
	emailtpl "gofiber-template/infrastructure/email"
	"gofiber-template/domain/ports"
	apperrors "gofiber-template/pkg/errors"
	"gofiber-template/pkg/scheduler"

	"github.com/google/uuid"
)

const (
	MaxImageSize           = 5 * 1024 * 1024 // 5MB
	DefaultServiceFeePerPax = 500.00          // ค่าบริการ/คน (บาท, รวม VAT) — override ด้วย SERVICE_FEE_PER_PAX env
)

// getServiceFeePerPax ดึงค่าบริการ/คนจาก env หรือใช้ค่า default 500 บาท
func getServiceFeePerPax() float64 {
	if val := os.Getenv("SERVICE_FEE_PER_PAX"); val != "" {
		if fee, err := strconv.ParseFloat(val, 64); err == nil && fee > 0 {
			return fee
		}
	}
	return DefaultServiceFeePerPax
}

type BookingServiceImpl struct {
	bookingRepo                  repositories.BookingRepository
	bookingRateTypeRepo          repositories.BookingRateTypeRepository
	bookingRoomTypeRepo          repositories.BookingRoomTypeRepository
	bookingAdditionalServiceRepo repositories.BookingAdditionalServiceRepository
	tourRepo                     repositories.TourRepository
	periodRepo                   repositories.PeriodRepository
	userRepo                     repositories.UserRepository
	paymentSlipRepo              repositories.PaymentSlipRepository
	passengerPassportRepo        repositories.PassengerPassportRepository
	periodService                services.PeriodService
	configService                services.ConfigService
	jobService                   services.JobService    // Optional: nil during initialization
	visionService                services.VisionService // Optional: nil if Vision API disabled
	llmService                   services.LLMService    // Optional: nil if LLM disabled
	seatAvailabilityService      services.SeatAvailabilityService
	salesRepo                    repositories.SalesRepository
	slipVerification             ports.SlipVerificationPort // Optional: nil if SlipOK not configured
	emailPort                    ports.EmailPort            // Optional: nil if SMTP not configured

	// Separate table repositories for price lookup
	periodRateItemRepo          repositories.PeriodRateItemRepository
	periodRoomTypeRepo          repositories.PeriodRoomTypeRepository
	periodAdditionalServiceRepo repositories.PeriodAdditionalServiceRepository
}

func NewBookingService(
	bookingRepo repositories.BookingRepository,
	bookingRateTypeRepo repositories.BookingRateTypeRepository,
	bookingRoomTypeRepo repositories.BookingRoomTypeRepository,
	bookingAdditionalServiceRepo repositories.BookingAdditionalServiceRepository,
	tourRepo repositories.TourRepository,
	periodRepo repositories.PeriodRepository,
	userRepo repositories.UserRepository,
	paymentSlipRepo repositories.PaymentSlipRepository,
	passengerPassportRepo repositories.PassengerPassportRepository,
	periodService services.PeriodService,
	configService services.ConfigService,
	jobService services.JobService,
	visionService services.VisionService,
	llmService services.LLMService,
	periodRateItemRepo repositories.PeriodRateItemRepository,
	periodRoomTypeRepo repositories.PeriodRoomTypeRepository,
	periodAdditionalServiceRepo repositories.PeriodAdditionalServiceRepository,
	seatAvailabilityService services.SeatAvailabilityService,
	salesRepo repositories.SalesRepository,
	slipVerification ports.SlipVerificationPort, // Optional: nil if SlipOK not configured
	emailPort ports.EmailPort, // Optional: nil if SMTP not configured
) services.BookingService {
	return &BookingServiceImpl{
		bookingRepo:                  bookingRepo,
		bookingRateTypeRepo:          bookingRateTypeRepo,
		bookingRoomTypeRepo:          bookingRoomTypeRepo,
		bookingAdditionalServiceRepo: bookingAdditionalServiceRepo,
		tourRepo:                     tourRepo,
		periodRepo:                   periodRepo,
		userRepo:                     userRepo,
		paymentSlipRepo:              paymentSlipRepo,
		passengerPassportRepo:        passengerPassportRepo,
		periodService:                periodService,
		configService:                configService,
		jobService:                   jobService,
		visionService:                visionService,
		llmService:                   llmService,
		periodRateItemRepo:           periodRateItemRepo,
		periodRoomTypeRepo:           periodRoomTypeRepo,
		periodAdditionalServiceRepo:  periodAdditionalServiceRepo,
		seatAvailabilityService:      seatAvailabilityService,
		salesRepo:                    salesRepo,
		slipVerification:             slipVerification,
		emailPort:                    emailPort,
	}
}

func (s *BookingServiceImpl) CreateBooking(ctx context.Context, userID uuid.UUID, req dto.CreateBookingRequest) (*dto.BookingResponse, error) {
	// Business validation - check if user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, fmt.Errorf("failed to validate user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Get period data early for validation and price lookup
	period, err := s.periodRepo.GetByID(ctx, req.PeriodID)
	if err != nil {
		return nil, fmt.Errorf("failed to get period: %w", err)
	}
	if period == nil {
		return nil, fmt.Errorf("period not found")
	}

	// Business validation - check agent booking limits
	totalQuantity := s.calculateTotalQuantity(ctx, req.PeriodID, req.RateTypes)
	err = s.configService.ValidateAgentBookingLimit(ctx, userID, req.PeriodID, totalQuantity)
	if err != nil {
		log.Printf("Agent booking limit validation failed: %v", err)
		return nil, fmt.Errorf("booking limit exceeded: %w", err)
	}

	// Business validation - check if period is available
	available, err := s.CheckPeriodAvailability(ctx, req.PeriodID, totalQuantity)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, apperrors.ErrPeriodNotAvailable
	}

	// Business logic - calculate totals
	sumRoom, sumRate, sumAdditionalServices, sumNet, err := s.CalculateBookingTotal(ctx, req)
	if err != nil {
		return nil, err
	}

	// Calculate commission from period - this becomes the discount amount
	discount, err := s.calculateCommissionForBooking(ctx, req.PeriodID, req.RateTypes)
	if err != nil {
		log.Printf("Error calculating commission: %v", err)
		return nil, fmt.Errorf("failed to calculate commission: %w", err)
	}

	// Generate invoice ID
	invoiceID, err := s.GenerateInvoiceID(ctx)
	if err != nil {
		return nil, err
	}

	// Calculate final net amount after auto-calculated discount
	finalNet := sumNet - discount
	netBeforeTax := finalNet // ยอดก่อนภาษี (หลังหักส่วนลด Agent)

	// ค่าบริการ — รวมอยู่ในราคาที่นั่งแล้ว (เช่น ราคา 19,000 = ค่าตั๋ว 18,500 + ค่าบริการ 500)
	// VAT คิดจากค่าบริการเท่านั้น (500 รวม VAT → ก่อน VAT 467.29, VAT 32.71)
	totalQty := s.calculateTotalQuantity(ctx, req.PeriodID, req.RateTypes)
	// ใช้ค่าบริการจาก Period ถ้ามี, fallback ใช้ค่า default จาก ENV
	serviceFeePerPax := period.ServiceFeePerPax
	if serviceFeePerPax <= 0 {
		serviceFeePerPax = getServiceFeePerPax()
	}
	serviceFeeTotal := serviceFeePerPax * float64(totalQty)

	// Calculate VAT and WHT based on request options
	var vatRate, vatAmount, whtRate, whtAmount float64
	vatRate = 7.00  // Default VAT rate
	whtRate = 3.00  // Default WHT rate

	if req.IncludeVAT {
		// VAT คิดจากค่าบริการเท่านั้น (ค่าบริการเป็นราคารวม VAT → แยก VAT ออก)
		// สูตร: VAT = serviceFeeTotal * vatRate / (100 + vatRate)
		vatAmount = math.Round(serviceFeeTotal*vatRate/(100+vatRate)*100) / 100
	}
	if req.IncludeWHT {
		// WHT คิดจากค่าบริการเท่านั้น (เหมือน VAT)
		whtAmount = math.Round(serviceFeeTotal*(whtRate/100)*100) / 100
	}

	// Calculate grand total
	// VAT เป็น include (รวมในค่าบริการแล้ว) → ไม่ต้องบวกเพิ่ม แค่แสดงเป็นข้อมูล
	// WHT หักจากค่าบริการ → หักจากยอดรวม
	grandTotal := netBeforeTax - whtAmount

	// Calculate deposit based on grand total (ยอดที่ต้องชำระ)
	sumDeposit, err := s.calculateBookingDeposit(ctx, req.PeriodID, grandTotal, s.calculateTotalQuantity(ctx, req.PeriodID, req.RateTypes))
	if err != nil {
		log.Printf("Error calculating deposit: %v", err)
		return nil, fmt.Errorf("failed to calculate deposit: %w", err)
	}

	// Create domain model
	booking := &models.Booking{
		InvoiceID:             invoiceID,
		UserID:                userID,
		TourID:                req.TourID,
		PeriodID:              req.PeriodID,
		SalesID:               req.SalesID,
		CustomerName:          req.CustomerName,
		CustomerPhone:         req.CustomerPhone,
		Note:                  req.Note,
		SumRoom:               sumRoom,
		SumRate:               sumRate,
		SumAdditionalServices: sumAdditionalServices,
		SumQty:                s.calculateTotalQuantity(ctx, req.PeriodID, req.RateTypes),
		SumCommissionQty:      s.calculateCommissionQuantity(ctx, req.PeriodID, req.RateTypes),
		SumNet:                finalNet,
		SumCommission:         discount, // Commission earned becomes the auto-calculated discount
		Discount:              discount, // Auto-calculated discount equals available commission
		SumDeposit:            sumDeposit,
		SumPayment:            0,

		// ค่าบริการ
		ServiceFeePerPax: serviceFeePerPax,
		ServiceFeeTotal:  serviceFeeTotal,

		// VAT fields
		IncludeVAT: req.IncludeVAT,
		VATRate:    vatRate,
		VATAmount:  vatAmount,

		// Withholding Tax fields
		IncludeWHT: req.IncludeWHT,
		WHTRate:    whtRate,
		WHTAmount:  whtAmount,

		// Calculated totals
		NetBeforeTax: netBeforeTax,
		GrandTotal:   grandTotal,

		DateDeposit: time.Now().AddDate(0, 0, 3),          // มัดจำ: 3 วันหลังทำการจอง
		DateFinal:   period.StartDate.AddDate(0, 0, -21), // ชำระเต็ม: 21 วันก่อนเดินทาง
		Status:      models.BookingStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create booking
	if err := s.bookingRepo.Create(ctx, booking); err != nil {
		log.Printf("Error creating booking: %v", err)
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	// Create BookingRoomType records
	if len(req.RoomTypes) > 0 {
		roomTypeRecords := make([]*models.BookingRoomType, len(req.RoomTypes))
		for i, roomTypeReq := range req.RoomTypes {
			// Find room type from separate table
			roomTypeID, err := uuid.Parse(roomTypeReq.ID)
			if err != nil {
				return nil, fmt.Errorf("invalid room type ID %s: %w", roomTypeReq.ID, err)
			}

			roomType, err := s.periodRoomTypeRepo.GetByID(ctx, roomTypeID)
			if err != nil {
				return nil, fmt.Errorf("failed to get room type %s: %w", roomTypeReq.ID, err)
			}
			if roomType == nil {
				return nil, fmt.Errorf("room type with ID %s not found", roomTypeReq.ID)
			}

			// Verify the room type belongs to the correct period
			if roomType.PeriodID != req.PeriodID {
				return nil, fmt.Errorf("room type %s does not belong to period %s", roomTypeReq.ID, req.PeriodID)
			}

			roomTypeRecords[i] = &models.BookingRoomType{
				BookingID: booking.ID,
				Title:     roomType.Title,
				Price:     roomType.Price,
				Quantity:  roomTypeReq.Quantity,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		if err := s.bookingRoomTypeRepo.CreateBatch(ctx, roomTypeRecords); err != nil {
			log.Printf("Error creating booking room types: %v", err)
			return nil, fmt.Errorf("failed to create booking room types: %w", err)
		}
	}

	// Create BookingRateType records
	if len(req.RateTypes) > 0 {
		rateTypeRecords := make([]*models.BookingRateType, len(req.RateTypes))
		for i, rateTypeReq := range req.RateTypes {
			// Find rate item from separate table
			rateItem, err := s.getRateItem(ctx, req.PeriodID, rateTypeReq.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get rate item %s: %w", rateTypeReq.ID, err)
			}

			// Use correct price based on period rate type
			var price float64
			switch period.RateType {
			case "regular":
				price = rateItem.Regular
			case "sale":
				price = rateItem.Sale
			default:
				price = rateItem.Regular
			}

			rateTypeRecords[i] = &models.BookingRateType{
				BookingID: booking.ID,
				Title:     rateItem.Title,
				Price:     price,
				Quantity:  rateTypeReq.Quantity,
				Exclude:   rateItem.Exclude,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		if err := s.bookingRateTypeRepo.CreateBatch(ctx, rateTypeRecords); err != nil {
			log.Printf("Error creating booking rate types: %v", err)
			return nil, fmt.Errorf("failed to create booking rate types: %w", err)
		}
	}

	// Create BookingAdditionalService records
	if len(req.AdditionalServices) > 0 {
		additionalServiceRecords := make([]*models.BookingAdditionalService, len(req.AdditionalServices))
		for i, serviceReq := range req.AdditionalServices {
			// Find additional service from separate table
			service, err := s.getAdditionalService(ctx, req.PeriodID, serviceReq.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get additional service %s: %w", serviceReq.ID, err)
			}
			if service == nil {
				return nil, fmt.Errorf("additional service with ID %s not found in period", serviceReq.ID)
			}

			additionalServiceRecords[i] = &models.BookingAdditionalService{
				BookingID: booking.ID,
				Title:     service.Title,
				Price:     service.Price,
				Quantity:  serviceReq.Quantity,
				RateType:  service.RateType,
				UnitName:  service.UnitName,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		if err := s.bookingAdditionalServiceRepo.CreateBatch(ctx, additionalServiceRecords); err != nil {
			log.Printf("Error creating booking additional services: %v", err)
			return nil, fmt.Errorf("failed to create booking additional services: %w", err)
		}
	}

	// Calculate and set booking expiration based on tiered configuration
	log.Printf("[DEBUG] Calculating expiration for booking %s", booking.ID)
	expirationTime, err := s.CalculateBookingExpiration(ctx, req.PeriodID, booking.CreatedAt)
	if err != nil {
		log.Printf("[WARNING] Failed to calculate booking expiration: %v", err)
		// Don't fail the booking creation, just log the warning
	} else if expirationTime != nil {
		booking.ExpiresAt = expirationTime
		// Update the booking with expiration time
		if err := s.bookingRepo.Update(ctx, booking.ID, booking); err != nil {
			log.Printf("[WARNING] Failed to update booking with expiration: %v", err)
			// Don't fail the booking creation, just log the warning
		} else {
			log.Printf("[SUCCESS] Booking expiration set to: %v", expirationTime)

			// Create dynamic expiry job for this booking
			if err := s.CreateBookingExpiryJob(ctx, booking.ID, *expirationTime); err != nil {
				log.Printf("[WARNING] Failed to create booking expiry job: %v", err)
				// Don't fail the booking creation, just log the warning
			} else {
				log.Printf("[SUCCESS] Booking expiry job created for booking %s", booking.ID)
			}
		}
	}

	response := s.toBookingResponse(ctx, booking)

	// Broadcast seat availability update via WebSocket
	go websocketInfra.BroadcastSeatUpdate(ctx, booking.PeriodID, "booking_created", s.seatAvailabilityService)

	// ส่ง email แจ้ง agent — สร้างการจองสำเร็จ
	s.sendBookingEmail(ctx, booking, "new_booking", "")

	return response, nil
}

func (s *BookingServiceImpl) GetBookingByID(ctx context.Context, bookingID uuid.UUID) (*dto.BookingResponse, error) {
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Printf("Error retrieving booking: %v", err)
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	if booking == nil {
		return nil, apperrors.ErrBookingNotFound
	}

	return s.toBookingResponse(ctx, booking), nil
}

// GetBookingModelByID returns raw booking model with PaymentSlips preloaded
func (s *BookingServiceImpl) GetBookingModelByID(ctx context.Context, bookingID uuid.UUID) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetByIDWithPaymentSlips(ctx, bookingID)
	if err != nil {
		log.Printf("Error retrieving booking with payment slips: %v", err)
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	if booking == nil {
		return nil, apperrors.ErrBookingNotFound
	}

	return booking, nil
}

func (s *BookingServiceImpl) GetBookingByInvoiceID(ctx context.Context, invoiceID string) (*dto.BookingResponse, error) {
	booking, err := s.bookingRepo.GetByInvoiceID(ctx, invoiceID)
	if err != nil {
		log.Printf("Error retrieving booking by invoice ID: %v", err)
		return nil, fmt.Errorf("failed to get booking by invoice ID: %w", err)
	}

	if booking == nil {
		return nil, apperrors.ErrBookingNotFound
	}

	return s.toBookingResponse(ctx, booking), nil
}

func (s *BookingServiceImpl) GetBookingWithDetails(ctx context.Context, bookingID uuid.UUID) (*dto.BookingResponse, error) {
	booking, err := s.bookingRepo.GetWithDetails(ctx, bookingID)
	if err != nil {
		log.Printf("Error retrieving booking with details: %v", err)
		return nil, fmt.Errorf("failed to get booking with details: %w", err)
	}

	if booking == nil {
		return nil, apperrors.ErrBookingNotFound
	}

	return s.toBookingResponse(ctx, booking), nil
}

func (s *BookingServiceImpl) UpdateBooking(ctx context.Context, bookingID uuid.UUID, req dto.UpdateBookingRequest) (*dto.BookingResponse, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Printf("Error retrieving booking: %v", err)
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	if booking == nil {
		return nil, apperrors.ErrBookingNotFound
	}

	// Business rule: Cannot update bookings that are completed or cancelled
	if booking.Status == models.BookingStatusCompleted || booking.Status == models.BookingStatusCancel {
		return nil, fmt.Errorf("cannot update booking with status: %s", booking.Status)
	}

	// Update fields if provided (only safe fields that don't affect pricing or availability)
	if req.CustomerName != nil {
		booking.CustomerName = *req.CustomerName
	}
	if req.CustomerPhone != nil {
		booking.CustomerPhone = *req.CustomerPhone
	}
	if req.Note != nil {
		booking.Note = req.Note
	}

	booking.UpdatedAt = time.Now()

	// Update booking
	if err := s.bookingRepo.Update(ctx, bookingID, booking); err != nil {
		log.Printf("Error updating booking: %v", err)
		return nil, fmt.Errorf("failed to update booking: %w", err)
	}

	// Note: No WebSocket broadcast needed since customer info changes don't affect seat availability
	// Status changes are handled by dedicated APIs (CancelBooking, payment processing, etc.)

	return s.toBookingResponse(ctx, booking), nil
}

func (s *BookingServiceImpl) UpdateBookingStatus(ctx context.Context, req dto.UpdateBookingStatusRequest) (*dto.BookingResponse, error) {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, req.BookingID)
	if err != nil {
		log.Printf("Error retrieving booking: %v", err)
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}

	if booking == nil {
		return nil, apperrors.ErrBookingNotFound
	}

	// Business validation - validate status
	status := models.BookingStatus(req.Status)
	if status != models.BookingStatusPending && status != models.BookingStatusConfirm &&
		status != models.BookingStatusCancel && status != models.BookingStatusCompleted {
		return nil, apperrors.ErrInvalidBookingStatus
	}

	// Store previous status for job management
	previousStatus := booking.Status

	// Update status
	booking.Status = status

	// Add reason to note if provided
	if req.Reason != nil {
		if booking.Note != nil {
			updatedNote := *booking.Note + "\nStatus update reason: " + *req.Reason
			booking.Note = &updatedNote
		} else {
			note := "Status update reason: " + *req.Reason
			booking.Note = &note
		}
	}

	booking.UpdatedAt = time.Now()

	// Update booking
	if err := s.bookingRepo.Update(ctx, req.BookingID, booking); err != nil {
		log.Printf("Error updating booking status: %v", err)
		return nil, fmt.Errorf("failed to update booking status: %w", err)
	}

	// Delete expiry job if status changed to a final state
	if (previousStatus == models.BookingStatusPending || previousStatus == models.BookingStatusProcessing) &&
		(status == models.BookingStatusConfirm || status == models.BookingStatusCompleted ||
			status == models.BookingStatusCancel || status == models.BookingStatusExpired) {
		if err := s.DeleteBookingExpiryJob(ctx, req.BookingID); err != nil {
			log.Printf("[WARNING] Failed to delete booking expiry job: %v", err)
			// Don't fail the status update, just log the warning
		} else {
			log.Printf("[SUCCESS] Booking expiry job deleted for booking %s (status: %s -> %s)", req.BookingID, previousStatus, status)
		}
	}

	response := s.toBookingResponse(ctx, booking)

	// Broadcast seat availability update via WebSocket
	go websocketInfra.BroadcastSeatUpdate(ctx, booking.PeriodID, "booking_updated", s.seatAvailabilityService)

	return response, nil
}

func (s *BookingServiceImpl) CancelBooking(ctx context.Context, bookingID uuid.UUID, reason string) error {
	// Get existing booking
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Printf("Error retrieving booking: %v", err)
		return fmt.Errorf("failed to get booking: %w", err)
	}

	if booking == nil {
		return apperrors.ErrBookingNotFound
	}

	// Business rule - check if booking can be cancelled
	if booking.Status == models.BookingStatusCompleted || booking.Status == models.BookingStatusCancel {
		return apperrors.ErrBookingCannotBeCancelled
	}

	// Update status and add note
	booking.Status = models.BookingStatusCancel
	if booking.Note != nil {
		updatedNote := *booking.Note + "\nCancellation reason: " + reason
		booking.Note = &updatedNote
	} else {
		note := "Cancellation reason: " + reason
		booking.Note = &note
	}
	booking.UpdatedAt = time.Now()

	// Update booking
	if err := s.bookingRepo.Update(ctx, bookingID, booking); err != nil {
		log.Printf("Error cancelling booking: %v", err)
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	// Delete expiry job since booking is cancelled
	if err := s.DeleteBookingExpiryJob(ctx, bookingID); err != nil {
		log.Printf("[WARNING] Failed to delete booking expiry job: %v", err)
		// Don't fail the cancellation, just log the warning
	} else {
		log.Printf("[SUCCESS] Booking expiry job deleted for cancelled booking %s", bookingID)
	}

	// Broadcast seat availability update via WebSocket
	go websocketInfra.BroadcastSeatUpdate(ctx, booking.PeriodID, "booking_cancelled", s.seatAvailabilityService)

	// ส่ง email แจ้ง agent — ยกเลิกการจอง
	s.sendBookingEmail(ctx, booking, "booking_cancelled", reason)

	return nil
}

func (s *BookingServiceImpl) ListBookings(ctx context.Context, req dto.BookingListRequest) (*dto.BookingListResponse, error) {
	// Use repository to get bookings based on filters with search and sorting
	bookings, err := s.bookingRepo.ListWithFilters(ctx, req)
	if err != nil {
		log.Printf("Error listing bookings with filters: %v", err)
		return nil, fmt.Errorf("failed to list bookings: %w", err)
	}

	// Get total count with same filters
	total, err := s.bookingRepo.CountWithFilters(ctx, req)
	if err != nil {
		log.Printf("Error counting bookings with filters: %v", err)
		return nil, fmt.Errorf("failed to count bookings: %w", err)
	}

	// Convert to response DTOs
	responseBookings := make([]*dto.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responseBookings[i] = s.toBookingResponse(ctx, booking)
	}

	pagination := dto.PaginationMeta{
		Total:       total,
		HasMore:     int64(req.Offset+len(bookings)) < total,
		CurrentPage: (req.Offset / req.Limit) + 1,
		TotalPages:  int((total + int64(req.Limit) - 1) / int64(req.Limit)),
		PageSize:    req.Limit,
		Limit:       req.Limit,
		Offset:      req.Offset,
	}

	return &dto.BookingListResponse{
		Bookings:   responseBookings,
		Pagination: pagination,
	}, nil
}

func (s *BookingServiceImpl) GetUserBookings(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.BookingListResponse, error) {
	bookings, err := s.bookingRepo.GetByUserID(ctx, userID, offset, limit)
	if err != nil {
		log.Printf("Error getting user bookings: %v", err)
		return nil, fmt.Errorf("failed to get user bookings: %w", err)
	}

	// Convert to response DTOs
	responseBookings := make([]*dto.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responseBookings[i] = s.toBookingResponse(ctx, booking)
	}

	// Get correct total count from database
	total, err := s.bookingRepo.CountByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error counting user bookings: %v", err)
		return nil, fmt.Errorf("failed to count user bookings: %w", err)
	}

	pagination := dto.PaginationMeta{
		Total:       total,
		HasMore:     len(bookings) == limit && total > int64(offset+limit),
		CurrentPage: (offset / limit) + 1,
		TotalPages:  int((total + int64(limit) - 1) / int64(limit)),
		PageSize:    limit,
		Limit:       limit,
		Offset:      offset,
	}

	return &dto.BookingListResponse{
		Bookings:   responseBookings,
		Pagination: pagination,
	}, nil
}

func (s *BookingServiceImpl) GetSalesBookings(ctx context.Context, salesID uuid.UUID, offset, limit int) (*dto.BookingListResponse, error) {
	bookings, err := s.bookingRepo.GetBySalesID(ctx, salesID, offset, limit)
	if err != nil {
		log.Printf("Error getting sales bookings: %v", err)
		return nil, fmt.Errorf("failed to get sales bookings: %w", err)
	}

	// Convert to response DTOs
	responseBookings := make([]*dto.BookingResponse, len(bookings))
	for i, booking := range bookings {
		responseBookings[i] = s.toBookingResponse(ctx, booking)
	}

	// Simplified total count
	total := int64(len(bookings))

	pagination := dto.PaginationMeta{
		Total:       total,
		HasMore:     len(bookings) == limit,
		CurrentPage: (offset / limit) + 1,
		TotalPages:  int((total + int64(limit) - 1) / int64(limit)),
		PageSize:    limit,
		Limit:       limit,
		Offset:      offset,
	}

	return &dto.BookingListResponse{
		Bookings:   responseBookings,
		Pagination: pagination,
	}, nil
}

// Passport Management
func (s *BookingServiceImpl) UploadPassport(ctx context.Context, req dto.UploadPassportRequest) (*dto.PassengerPassportResponse, error) {
	log.Printf("[INFO] Starting passport upload for booking %s, file: %s", req.BookingID, req.FileName)

	// Business validation - check if booking exists
	booking, err := s.bookingRepo.GetByID(ctx, req.BookingID)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve booking for passport: %v", err)
		return nil, fmt.Errorf("failed to validate booking: %w", err)
	}
	if booking == nil {
		return nil, apperrors.ErrBookingNotFound
	}

	// Business validation - check image size
	if len(req.ImageData) > MaxImageSize {
		log.Printf("[ERROR] Image size too large: %d bytes", len(req.ImageData))
		return nil, apperrors.ErrInvalidImageFormat
	}

	// Use S3 Storage URL if provided, otherwise generate one (for backward compatibility)
	var imageURL string
	if req.ImageURL != "" {
		imageURL = req.ImageURL // Use S3 Storage URL from file service
		log.Printf("[INFO] Using S3 Storage URL: %s", imageURL)
	} else {
		imageURL = s.generateImageURL(req.BookingID, req.FileName) // Fallback for old implementation
		log.Printf("[INFO] Generated fallback URL: %s", imageURL)
	}

	// Create initial passport record
	passport := &models.PassengerPassport{
		BookingID:      req.BookingID,
		PassportNumber: req.PassportNumber,
		Image:          &imageURL,
		Status:         models.PassportStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Create passport record in database first
	if err := s.passengerPassportRepo.Create(ctx, passport); err != nil {
		log.Printf("[ERROR] Failed to create passport record: %v", err)
		return nil, fmt.Errorf("failed to create passport: %w", err)
	}

	log.Printf("[SUCCESS] Passport record created with ID: %s", passport.ID)

	// Process with Vision API and LLM if available
	// Step 1: Extract text with Google Vision API
	if s.visionService != nil {
		log.Printf("[INFO] Processing passport with Google Vision API for OCR")

		ocrResult, err := s.visionService.ExtractTextFromImage(ctx, req.ImageData)
		if err != nil {
			log.Printf("[WARNING] Vision API OCR failed: %v", err)
			passport.Status = models.PassportStatusError
			passport.ErrorMessage = stringPtr(fmt.Sprintf("Vision API OCR error: %v", err))
		} else {
			log.Printf("[SUCCESS] Vision API OCR completed, extracted text length: %d", len(ocrResult.Text))

			// Step 2: Process OCR text with LLM for structured extraction
			if s.llmService != nil && ocrResult.Text != "" {
				log.Printf("[INFO] Processing OCR text with LLM for structured extraction")

				llmResult, err := s.llmService.ExtractPassportFields(ctx, ocrResult.Text)
				if err != nil {
					log.Printf("[WARNING] LLM extraction failed: %v", err)
					passport.Status = models.PassportStatusError
					passport.ErrorMessage = stringPtr(fmt.Sprintf("LLM extraction error: %v", err))
				} else {
					log.Printf("[SUCCESS] LLM extraction completed, confidence: %.2f", llmResult.OverallConfidence)

					// Update passport with LLM extracted data
					s.updatePassportWithExtractionResult(passport, llmResult)

					// Set status based on processing success
					if llmResult.ProcessedSuccessfully && llmResult.OverallConfidence > 0.5 {
						passport.Status = models.PassportStatusProcessed
						passport.ProcessedDate = timePtr(time.Now())
						log.Printf("[SUCCESS] Passport data extracted successfully with LLM")
					} else {
						passport.Status = models.PassportStatusError
						passport.ErrorMessage = llmResult.ErrorMessage
						log.Printf("[WARNING] LLM extraction completed with low confidence or errors")
					}
				}
			} else {
				log.Printf("[INFO] LLM service not available or no OCR text for processing")
				passport.Status = models.PassportStatusError
				passport.ErrorMessage = stringPtr("LLM service not available for text processing")
			}
		}

		// Update passport with processing results
		passport.UpdatedAt = time.Now()
		if err := s.passengerPassportRepo.Update(ctx, passport.ID, passport); err != nil {
			log.Printf("[ERROR] Failed to update passport with processing results: %v", err)
			// Don't fail the upload, but log the error
		}
	} else {
		log.Printf("[INFO] Vision API not available, passport saved for manual processing")
		passport.Status = models.PassportStatusError
		passport.ErrorMessage = stringPtr("Vision API service not available")
	}

	log.Printf("[SUCCESS] Passport upload completed for booking %s", req.BookingID)
	return s.toPassengerPassportResponse(passport), nil
}

func (s *BookingServiceImpl) ProcessPassportOCR(ctx context.Context, passportID uuid.UUID) (*dto.PassengerPassportResponse, error) {
	// Get passport
	passport, err := s.passengerPassportRepo.GetByID(ctx, passportID)
	if err != nil {
		log.Printf("Error retrieving passport: %v", err)
		return nil, fmt.Errorf("failed to get passport: %w", err)
	}
	if passport == nil {
		return nil, apperrors.ErrPassportNotFound
	}

	// TODO: Implement OCR processing
	// For now, just update status to processed
	passport.Status = models.PassportStatusProcessed
	passport.UpdatedAt = time.Now()

	if err := s.passengerPassportRepo.Update(ctx, passportID, passport); err != nil {
		log.Printf("Error updating passport status: %v", err)
		return nil, fmt.Errorf("failed to update passport status: %w", err)
	}

	return s.toPassengerPassportResponse(passport), nil
}

func (s *BookingServiceImpl) UpdatePassportInfo(ctx context.Context, passportID uuid.UUID, passportInfo *dto.PassengerPassportResponse) (*dto.PassengerPassportResponse, error) {
	// Get passport
	passport, err := s.passengerPassportRepo.GetByID(ctx, passportID)
	if err != nil {
		log.Printf("Error retrieving passport: %v", err)
		return nil, fmt.Errorf("failed to get passport: %w", err)
	}
	if passport == nil {
		return nil, apperrors.ErrPassportNotFound
	}

	// Update passport information from DTO
	if err := s.passengerPassportRepo.UpdatePersonalInfo(ctx, passportID, passport); err != nil {
		log.Printf("Error updating passport info: %v", err)
		return nil, fmt.Errorf("failed to update passport info: %w", err)
	}

	return s.toPassengerPassportResponse(passport), nil
}

// Business Logic Methods
func (s *BookingServiceImpl) CalculateBookingTotal(ctx context.Context, req dto.CreateBookingRequest) (float64, float64, float64, float64, error) {
	// Get period data to fetch actual prices
	period, err := s.periodRepo.GetByID(ctx, req.PeriodID)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("failed to get period: %w", err)
	}
	if period == nil {
		return 0, 0, 0, 0, fmt.Errorf("period not found")
	}

	var sumRoom, sumRate, sumAdditionalServices float64

	// Calculate room totals - fetch prices from separate tables
	for _, roomTypeReq := range req.RoomTypes {
		price, err := s.getRoomTypePrice(ctx, req.PeriodID, roomTypeReq.ID)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("room type not found: %s", roomTypeReq.ID)
		}
		sumRoom += price * float64(roomTypeReq.Quantity)
	}

	// Calculate rate totals - fetch prices from separate tables
	// Note: ALL passengers pay rates, including excluded ones (exclude only affects commission)
	for _, rateTypeReq := range req.RateTypes {
		rateItem, err := s.getRateItem(ctx, req.PeriodID, rateTypeReq.ID)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("rate type not found: %s", rateTypeReq.ID)
		}

		// Use regular or sale price based on period rate type
		var price float64
		switch period.RateType {
		case "regular":
			price = rateItem.Regular
		case "sale":
			price = rateItem.Sale
		default:
			price = rateItem.Regular
		}

		// All passengers pay rates (exclude flag only affects commission calculation)
		sumRate += price * float64(rateTypeReq.Quantity)
	}

	// Calculate additional services totals - fetch prices from separate tables
	for _, serviceReq := range req.AdditionalServices {
		service, err := s.getAdditionalService(ctx, req.PeriodID, serviceReq.ID)
		if err != nil {
			return 0, 0, 0, 0, fmt.Errorf("additional service not found: %s", serviceReq.ID)
		}
		sumAdditionalServices += service.Price * float64(serviceReq.Quantity)
	}

	sumNet := sumRoom + sumRate + sumAdditionalServices

	return sumRoom, sumRate, sumAdditionalServices, sumNet, nil
}

func (s *BookingServiceImpl) CheckPeriodAvailability(ctx context.Context, periodID uuid.UUID, quantity int) (bool, error) {
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error checking period availability: %v", err)
		return false, fmt.Errorf("failed to check period availability: %w", err)
	}

	if period == nil {
		return false, apperrors.ErrPeriodNotAvailable
	}

	// Check if period is published and has enough capacity
	if period.Status != models.PeriodStatusPublish {
		return false, nil
	}

	// Calculate real-time active bookings (pending, processing, confirm)
	activeStatuses := []string{"pending", "processing", "confirm"}
	bookings, err := s.bookingRepo.FindByPeriodIDAndStatuses(periodID, activeStatuses)
	if err != nil {
		log.Printf("Warning: Failed to count active bookings for period %s: %v", periodID, err)
		return false, err
	}

	activeBookings := 0
	for _, booking := range bookings {
		activeBookings += booking.SumQty
	}

	availableSeats := period.Seat - activeBookings
	return availableSeats >= quantity, nil
}

func (s *BookingServiceImpl) GenerateInvoiceID(ctx context.Context) (string, error) {
	// Simple invoice ID generation - in real implementation, make it more robust
	now := time.Now()
	count, err := s.bookingRepo.Count(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get booking count for invoice ID: %w", err)
	}

	// Format: INV-YYYYMMDD-NNNNN
	invoiceID := fmt.Sprintf("INV-%s-%05d", now.Format("20060102"), count+1)
	return invoiceID, nil
}

func (s *BookingServiceImpl) GetBookingStats(ctx context.Context, dateFrom, dateTo *time.Time) (*dto.BookingStatsResponse, error) {
	// Get total bookings
	totalBookings, err := s.bookingRepo.Count(ctx)
	if err != nil {
		log.Printf("Error getting booking stats: %v", err)
		return nil, fmt.Errorf("failed to get booking stats: %w", err)
	}

	// Simplified stats - in real implementation, use proper date filtering and status grouping
	return &dto.BookingStatsResponse{
		TotalBookings:     totalBookings,
		PendingBookings:   0, // TODO: Implement status-based counts
		ConfirmedBookings: 0, // TODO: Implement status-based counts
		CancelledBookings: 0, // TODO: Implement status-based counts
		TotalRevenue:      0, // TODO: Implement revenue calculation
		PendingRevenue:    0, // TODO: Implement revenue calculation
	}, nil
}

// Helper methods
func (s *BookingServiceImpl) calculateTotalQuantity(ctx context.Context, periodID uuid.UUID, rateTypes []*dto.CreateBookingRateTypeRequest) int {
	// Count ALL passengers (including excluded ones)
	total := 0
	for _, rateType := range rateTypes {
		total += rateType.Quantity
	}
	return int(math.Max(1, float64(total))) // At least 1
}

func (s *BookingServiceImpl) calculateCommissionQuantity(ctx context.Context, periodID uuid.UUID, rateTypes []*dto.CreateBookingRateTypeRequest) int {
	total := 0

	// Load rate items from separate table to check exclude flag
	rateItems, err := s.periodRateItemRepo.GetByPeriodID(ctx, periodID)
	if err != nil {
		// If we can't get rate items, return 0 for commission quantity
		return 0
	}

	for _, rateType := range rateTypes {
		// Parse rate type ID
		rateTypeID, err := uuid.Parse(rateType.ID)
		if err != nil {
			// If ID is invalid, skip this rate type
			continue
		}

		// Find the rate item in separate table to check exclude flag
		for _, periodRateItem := range rateItems {
			if periodRateItem.ID == rateTypeID {
				// Only count non-excluded passengers for commission
				if !periodRateItem.Exclude {
					total += rateType.Quantity
				}
				break
			}
		}
	}
	return total
}

func (s *BookingServiceImpl) calculateBookingDeposit(ctx context.Context, periodID uuid.UUID, finalNet float64, totalQty int) (float64, error) {
	// Get period to access deposit settings
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		return 0, fmt.Errorf("failed to get period: %w", err)
	}
	if period == nil {
		return 0, fmt.Errorf("period not found")
	}

	var sumDeposit float64

	switch period.DepositType {
	case "percent":
		// Percentage-based deposit: deposit value is percentage (e.g., 30 = 30%)
		sumDeposit = finalNet * (period.Deposit / 100)
		log.Printf("Deposit calculation (percent): %.2f * (%.2f/100) = %.2f", finalNet, period.Deposit, sumDeposit)
	case "per_unit":
		// Per-unit deposit: deposit value is amount per passenger
		sumDeposit = period.Deposit * float64(totalQty)
		log.Printf("Deposit calculation (per_unit): %.2f * %d = %.2f", period.Deposit, totalQty, sumDeposit)
	default:
		// Fallback to percentage for backward compatibility
		sumDeposit = finalNet * (period.Deposit / 100)
		log.Printf("Deposit calculation (fallback percent): %.2f * (%.2f/100) = %.2f", finalNet, period.Deposit, sumDeposit)
	}

	return sumDeposit, nil
}

func (s *BookingServiceImpl) toBookingResponse(ctx context.Context, booking *models.Booking) *dto.BookingResponse {
	response := &dto.BookingResponse{
		ID:                    booking.ID,
		InvoiceID:             booking.InvoiceID,
		UserID:                booking.UserID,
		TourID:                booking.TourID,
		PeriodID:              booking.PeriodID,
		SalesID:               booking.SalesID,
		CustomerName:          booking.CustomerName,
		CustomerPhone:         booking.CustomerPhone,
		Note:                  booking.Note,
		SumRoom:               booking.SumRoom,
		SumRate:               booking.SumRate,
		SumAdditionalServices: booking.SumAdditionalServices,
		SumQty:                booking.SumQty,
		SumCommissionQty:      booking.SumCommissionQty,
		SumNet:                booking.SumNet,
		SumCommission:         booking.SumCommission,
		Discount:              booking.Discount,
		SumDeposit:            booking.SumDeposit,
		SumPayment:            booking.GetActualSumPayment(), // ใช้ computed value แทน stored value

		// ค่าบริการ
		ServiceFeePerPax: booking.ServiceFeePerPax,
		ServiceFeeTotal:  booking.ServiceFeeTotal,

		// VAT fields
		IncludeVAT: booking.IncludeVAT,
		VATRate:    booking.VATRate,
		VATAmount:  booking.VATAmount,

		// Withholding Tax fields
		IncludeWHT: booking.IncludeWHT,
		WHTRate:    booking.WHTRate,
		WHTAmount:  booking.WHTAmount,

		// Calculated totals
		NetBeforeTax: booking.NetBeforeTax,
		GrandTotal:   booking.GrandTotal,

		DateDeposit: booking.DateDeposit.Format("2006-01-02"),
		DateFinal:   booking.DateFinal.Format("2006-01-02"),
		Status:      string(booking.Status),
		CreatedAt:   booking.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   booking.UpdatedAt.Format(time.RFC3339),
	}

	// Add expires_at field if it exists
	if booking.ExpiresAt != nil {
		expiresAtStr := booking.ExpiresAt.Format(time.RFC3339)
		response.ExpiresAt = &expiresAtStr
	}

	// Fetch and include booking room types
	roomTypes, err := s.bookingRoomTypeRepo.GetByBookingID(ctx, booking.ID)
	if err == nil && len(roomTypes) > 0 {
		response.RoomTypes = make([]*dto.BookingRoomTypeResponse, len(roomTypes))
		for i, roomType := range roomTypes {
			response.RoomTypes[i] = &dto.BookingRoomTypeResponse{
				ID:       roomType.ID,
				Title:    roomType.Title,
				Price:    roomType.Price,
				Quantity: roomType.Quantity,
				Total:    roomType.TotalPrice(),
			}
		}
	}

	// Fetch and include booking rate types
	rateTypes, err := s.bookingRateTypeRepo.GetByBookingID(ctx, booking.ID)
	if err == nil && len(rateTypes) > 0 {
		response.RateTypes = make([]*dto.BookingRateTypeResponse, len(rateTypes))
		for i, rateType := range rateTypes {
			response.RateTypes[i] = &dto.BookingRateTypeResponse{
				ID:       rateType.ID,
				Title:    rateType.Title,
				Price:    rateType.Price,
				Quantity: rateType.Quantity,
				Exclude:  rateType.Exclude,
				Total:    rateType.TotalPrice(),
			}
		}
	}

	// Fetch and include booking additional services
	additionalServices, err := s.bookingAdditionalServiceRepo.GetByBookingID(ctx, booking.ID)
	if err == nil && len(additionalServices) > 0 {
		response.AdditionalServices = make([]*dto.BookingAdditionalServiceResponse, len(additionalServices))
		for i, service := range additionalServices {
			response.AdditionalServices[i] = &dto.BookingAdditionalServiceResponse{
				ID:       service.ID,
				Title:    service.Title,
				Price:    service.Price,
				Quantity: service.Quantity,
				RateType: service.RateType,
				UnitName: service.UnitName,
				Total:    service.TotalPrice(),
			}
		}
	}

	// Fetch and include tour details
	if tour, err := s.tourRepo.GetByID(ctx, booking.TourID); err == nil && tour != nil {
		response.TourDetails = &dto.TourResponse{
			ID:        tour.ID,
			Title:     tour.Title,
			Code:      tour.Code,
			CityID:    tour.CityID,
			Detail:    tour.Detail,
			Thumbnail: tour.Thumbnail,
			FilePDF:   tour.FilePDF,
			FileWord:  tour.FileWord,
			Status:    string(tour.Status),
			CreatedAt: tour.CreatedAt.Format(time.RFC3339),
			UpdatedAt: tour.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Fetch and include period details
	if period, err := s.periodRepo.GetByID(ctx, booking.PeriodID); err == nil && period != nil {
		// Calculate real-time booking stats
		realTimeStats, err := s.periodService.CalculateRealTimeBookingStats(ctx, period.ID)
		booking := 0
		heldSeats := 0
		availableSeats := period.Seat
		if err == nil {
			booking = realTimeStats.BookedPassengers
			heldSeats = realTimeStats.HeldSeats
			availableSeats = realTimeStats.AvailableSeats
		}

		// Convert file attachments
		files := make([]*dto.FileAttachmentResponse, len(period.Files))
		for i, file := range period.Files {
			files[i] = &dto.FileAttachmentResponse{
				Name: file.Name,
				Type: file.Type,
				URL:  file.URL,
			}
		}

		response.PeriodDetails = &dto.PeriodSummaryResponse{
			ID:             period.ID,
			StartDate:      period.StartDate.Format("2006-01-02"),
			EndDate:        period.EndDate.Format("2006-01-02"),
			Seat:           period.Seat,
			Booking:        booking,
			HeldSeats:      heldSeats,
			AvailableSeats: availableSeats,
			Deposit:        period.Deposit,
			DepositType:    period.DepositType,
			Files:          files,
		}
	}

	// Include company profile if available
	if booking.CompanyProfile != nil {
		response.CompanyProfile = &dto.CompanyProfileResponse{
			ID:       booking.CompanyProfile.ID,
			UserID:   booking.CompanyProfile.UserID,
			Name:     booking.CompanyProfile.Name,
			TaxID:    booking.CompanyProfile.TaxID,
			Phone:    booking.CompanyProfile.Phone,
			Address:  booking.CompanyProfile.Address,
			Province: booking.CompanyProfile.Province,
			Zipcode:  booking.CompanyProfile.Zipcode,
			ATTA:     booking.CompanyProfile.ATTA,
			DBD:      booking.CompanyProfile.DBD,
		}
	}

	// Include sales details if available
	if booking.Sales != nil {
		response.SalesDetails = &dto.SalesResponse{
			ID:        booking.Sales.ID,
			UserID:    booking.Sales.UserID,
			Code:      booking.Sales.Code,
			FirstName: booking.Sales.FirstName,
			LastName:  booking.Sales.LastName,
			NickName:  booking.Sales.NickName,
			Active:    booking.Sales.Active,
			CreatedAt: booking.Sales.CreatedAt.Format(time.RFC3339),
			UpdatedAt: booking.Sales.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Include user details if available
	if booking.User != nil {
		fmt.Printf("User loaded for booking %s: %s (%s)\n", booking.ID, booking.User.Email, booking.User.Role)
		response.UserDetails = &dto.UserResponse{
			ID:        booking.User.ID,
			Email:     booking.User.Email,
			Role:      string(booking.User.Role),
			Status:    string(booking.User.Status),
			CreatedAt: booking.User.CreatedAt.Format(time.RFC3339),
			UpdatedAt: booking.User.UpdatedAt.Format(time.RFC3339),
		}
	} else {
		fmt.Printf("No user loaded for booking %s (user_id: %s)\n", booking.ID, booking.UserID)
	}

	// Include payment slips if available
	if len(booking.PaymentSlips) > 0 {
		response.PaymentSlips = make([]*dto.PaymentSlipResponse, len(booking.PaymentSlips))
		for i, paymentSlip := range booking.PaymentSlips {
			response.PaymentSlips[i] = s.toPaymentSlipResponse(&paymentSlip)
		}
	}

	// Include passenger passports if available
	if len(booking.PassengerPassports) > 0 {
		response.Passports = make([]*dto.PassengerPassportResponse, len(booking.PassengerPassports))
		for i, passport := range booking.PassengerPassports {
			response.Passports[i] = s.toPassengerPassportResponse(&passport)
		}
	}

	return response
}

func (s *BookingServiceImpl) toPassengerPassportResponse(passport *models.PassengerPassport) *dto.PassengerPassportResponse {
	response := &dto.PassengerPassportResponse{
		ID:             passport.ID,
		BookingID:      passport.BookingID,
		PassportNumber: passport.PassportNumber,
		Image:          passport.Image,
		Status:         string(passport.Status),
		CreatedAt:      passport.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      passport.UpdatedAt.Format(time.RFC3339),
	}

	// Convert optional name fields
	if passport.TitleTH != nil {
		response.TitleTH = passport.TitleTH
	}
	if passport.NameTH != nil {
		response.NameTH = passport.NameTH
	}
	if passport.SurnameTH != nil {
		response.SurnameTH = passport.SurnameTH
	}
	if passport.TitleEN != nil {
		response.TitleEN = passport.TitleEN
	}
	if passport.NameEN != nil {
		response.NameEN = passport.NameEN
	}
	if passport.SurnameEN != nil {
		response.SurnameEN = passport.SurnameEN
	}

	// Convert personal information fields
	if passport.Gender != nil {
		response.Gender = passport.Gender
	}
	if passport.Nationality != nil {
		response.Nationality = passport.Nationality
	}
	if passport.CountryCode != nil {
		response.CountryCode = passport.CountryCode
	}
	if passport.IDNumber != nil {
		response.IDNumber = passport.IDNumber
	}
	if passport.PlaceOfBirth != nil {
		response.PlaceOfBirth = passport.PlaceOfBirth
	}
	if passport.IssuingAuthority != nil {
		response.IssuingAuthority = passport.IssuingAuthority
	}

	// Convert date fields
	if passport.DateOfBirth != nil {
		dateStr := passport.DateOfBirth.Format("2006-01-02")
		response.DateOfBirth = &dateStr
	}
	if passport.DateOfExpiry != nil {
		dateStr := passport.DateOfExpiry.Format("2006-01-02")
		response.DateOfExpiry = &dateStr
	}
	if passport.DateOfIssue != nil {
		dateStr := passport.DateOfIssue.Format("2006-01-02")
		response.DateOfIssue = &dateStr
	}

	// Convert processing information
	if passport.ProcessedDate != nil {
		dateStr := passport.ProcessedDate.Format(time.RFC3339)
		response.ProcessedDate = &dateStr
	}
	if passport.ErrorMessage != nil {
		response.ErrorMessage = passport.ErrorMessage
	}

	// Convert OCR fields
	if passport.OCRText != nil {
		response.OCRText = passport.OCRText
	}
	if passport.OCRConfidence != nil {
		response.OCRConfidence = passport.OCRConfidence
	}

	return response
}

func (s *BookingServiceImpl) generateImageURL(bookingID uuid.UUID, fileName string) string {
	// Generate a unique file path
	ext := filepath.Ext(fileName)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("uploads/bookings/%s/%d%s", bookingID.String(), timestamp, ext)
}

// calculateCommissionForBooking calculates total commission for booking based on commission pool business logic
func (s *BookingServiceImpl) calculateCommissionForBooking(ctx context.Context, periodID uuid.UUID, rateTypes []*dto.CreateBookingRateTypeRequest) (float64, error) {
	// Get period to access commission items and commission type
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		return 0, fmt.Errorf("failed to get period: %w", err)
	}
	if period == nil {
		return 0, fmt.Errorf("period not found")
	}

	// Check if commission items exist
	if len(period.CommissionItems) == 0 {
		return 0, fmt.Errorf("no commission items found for period")
	}

	// Load rate items from separate table to check exclude flags
	rateItems, err := s.periodRateItemRepo.GetByPeriodID(ctx, periodID)
	if err != nil {
		return 0, fmt.Errorf("failed to get rate items from separate table: %w", err)
	}
	if len(rateItems) == 0 {
		return 0, fmt.Errorf("no rate items found for period in separate table")
	}

	// Step 1: Calculate commission pool total from all commission items
	var commissionPoolTotal float64
	for _, commissionItem := range period.CommissionItems {
		// Determine commission value based on commission type from period
		var commissionValue float64
		switch period.CommissionType {
		case "regular":
			commissionValue = commissionItem.Regular
		case "sale":
			commissionValue = commissionItem.Sale
		default:
			// Default to regular if commission type is not specified
			commissionValue = commissionItem.Regular
		}
		commissionPoolTotal += commissionValue
	}

	log.Printf("Commission pool total: %.2f (commission_type: %s)", commissionPoolTotal, period.CommissionType)

	// Step 2: Use existing method to count non-excluded passengers efficiently
	totalNonExcludedPassengers := 0
	for _, bookingRateType := range rateTypes {
		// Parse booking rate type ID
		rateItemID, err := uuid.Parse(bookingRateType.ID)
		if err != nil {
			log.Printf("Warning: Invalid rate item ID: %s", bookingRateType.ID)
			continue
		}

		// Find rate item by ID from separate table to check exclude flag
		var rateItem *models.PeriodRateItem
		for _, ri := range rateItems {
			if ri.ID == rateItemID {
				rateItem = ri
				break
			}
		}

		if rateItem == nil {
			log.Printf("Warning: No matching rate item found for booking rate type ID: %s", bookingRateType.ID)
			continue
		}

		// Only count non-excluded passengers for commission
		if !rateItem.Exclude {
			totalNonExcludedPassengers += bookingRateType.Quantity
			log.Printf("Rate item %s (%s): %d passengers (exclude: %t)",
				rateItem.ID, rateItem.Title, bookingRateType.Quantity, rateItem.Exclude)
		} else {
			log.Printf("Rate item %s (%s): %d passengers EXCLUDED from commission",
				rateItem.ID, rateItem.Title, bookingRateType.Quantity)
		}
	}

	// Step 3: Calculate total commission = commission pool × non-excluded passenger count
	totalCommission := commissionPoolTotal * float64(totalNonExcludedPassengers)

	log.Printf("Commission calculation breakdown:")
	log.Printf("- Commission pool total: %.2f", commissionPoolTotal)
	log.Printf("- Non-excluded passengers: %d", totalNonExcludedPassengers)
	log.Printf("- Total commission: %.2f × %d = %.2f", commissionPoolTotal, totalNonExcludedPassengers, totalCommission)

	return totalCommission, nil
}

// Helper methods to find items by ID from separate tables
func (s *BookingServiceImpl) getRoomTypePrice(ctx context.Context, periodID uuid.UUID, id string) (float64, error) {
	roomTypeID, err := uuid.Parse(id)
	if err != nil {
		return 0, fmt.Errorf("invalid room type ID: %s", id)
	}

	roomType, err := s.periodRoomTypeRepo.GetByID(ctx, roomTypeID)
	if err != nil {
		return 0, fmt.Errorf("failed to get room type: %w", err)
	}
	if roomType == nil {
		return 0, fmt.Errorf("room type with ID %s not found", id)
	}

	// Verify the room type belongs to the correct period
	if roomType.PeriodID != periodID {
		return 0, fmt.Errorf("room type %s does not belong to period %s", id, periodID)
	}

	return roomType.Price, nil
}

func (s *BookingServiceImpl) getRateItem(ctx context.Context, periodID uuid.UUID, id string) (*models.PeriodRateItem, error) {
	rateItemID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid rate item ID: %s", id)
	}

	rateItem, err := s.periodRateItemRepo.GetByID(ctx, rateItemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rate item: %w", err)
	}
	if rateItem == nil {
		return nil, fmt.Errorf("rate item with ID %s not found", id)
	}

	// Verify the rate item belongs to the correct period
	if rateItem.PeriodID != periodID {
		return nil, fmt.Errorf("rate item %s does not belong to period %s", id, periodID)
	}

	return rateItem, nil
}

func (s *BookingServiceImpl) getAdditionalService(ctx context.Context, periodID uuid.UUID, id string) (*models.PeriodAdditionalService, error) {
	serviceID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid additional service ID: %s", id)
	}

	service, err := s.periodAdditionalServiceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get additional service: %w", err)
	}
	if service == nil {
		return nil, fmt.Errorf("additional service with ID %s not found", id)
	}

	// Verify the service belongs to the correct period
	if service.PeriodID != periodID {
		return nil, fmt.Errorf("additional service %s does not belong to period %s", id, periodID)
	}

	return service, nil
}

// ================================================
// TIERED EXPIRATION MANAGEMENT METHODS
// ================================================

// SetBookingExpiration calculates and sets expiration time for a booking based on tiered rules
func (s *BookingServiceImpl) SetBookingExpiration(ctx context.Context, bookingID uuid.UUID, periodID uuid.UUID) error {
	log.Printf("[DEBUG] Setting tiered expiration for booking %s, period %s", bookingID.String(), periodID.String())

	// Get booking to get creation time
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Printf("[ERROR] Failed to get booking %s: %v", bookingID.String(), err)
		return fmt.Errorf("failed to get booking: %w", err)
	}

	if booking == nil {
		log.Printf("[ERROR] Booking not found: %s", bookingID.String())
		return fmt.Errorf("booking not found")
	}

	// Calculate expiration time
	expirationTime, err := s.CalculateBookingExpiration(ctx, periodID, booking.CreatedAt)
	if err != nil {
		log.Printf("[ERROR] Failed to calculate expiration for booking %s: %v", bookingID.String(), err)
		return fmt.Errorf("failed to calculate expiration: %w", err)
	}

	// Update booking with expiration time
	if expirationTime != nil {
		booking.SetExpiration(*expirationTime)
		err = s.bookingRepo.Update(ctx, bookingID, booking)
		if err != nil {
			log.Printf("[ERROR] Failed to update booking expiration %s: %v", bookingID.String(), err)
			return fmt.Errorf("failed to update booking expiration: %w", err)
		}

		log.Printf("[SUCCESS] Booking %s expiration set to %s", bookingID.String(), expirationTime.Format(time.RFC3339))
	} else {
		log.Printf("[INFO] No expiration set for booking %s (default tier has no expiration)", bookingID.String())
	}

	return nil
}

// CalculateBookingExpiration calculates when a booking should expire based on tiered rules
// Priority: 1. Period Override (ทัวร์ไฟไหม้) 2. Global Tiers
func (s *BookingServiceImpl) CalculateBookingExpiration(ctx context.Context, periodID uuid.UUID, createdAt time.Time) (*time.Time, error) {
	log.Printf("[DEBUG] Calculating booking expiration for period %s, created at %s", periodID.String(), createdAt.Format(time.RFC3339))

	// Get period to determine tour start date and check for override
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("[ERROR] Failed to get period %s: %v", periodID.String(), err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		log.Printf("[ERROR] Period not found: %s", periodID.String())
		return nil, fmt.Errorf("period not found")
	}

	// Priority 1: Check for period-level expiration override (ทัวร์ไฟไหม้)
	if period.HasExpirationOverride() {
		expirationTime := period.CalculateOverrideExpiration(createdAt)
		if expirationTime != nil {
			msg := "N/A"
			if period.ExpirationMessage != nil {
				msg = *period.ExpirationMessage
			}
			log.Printf("[SUCCESS] Using period expiration override: %s (message: %s)", expirationTime.Format(time.RFC3339), msg)
			return expirationTime, nil
		}
	}

	// Priority 2: Use global expiration tiers
	daysUntilTour := models.DaysUntilTour(period.StartDate)
	log.Printf("[DEBUG] Days until tour: %d, using global tiers", daysUntilTour)

	// Get applicable tier
	tier, err := s.configService.GetApplicableTier(ctx, daysUntilTour)
	if err != nil {
		log.Printf("[ERROR] Failed to get applicable tier: %v", err)
		return nil, fmt.Errorf("failed to get applicable tier: %w", err)
	}

	// Calculate expiration time using the tier
	expirationTime := tier.CalculateExpirationTime(createdAt)
	if expirationTime != nil {
		log.Printf("[SUCCESS] Calculated expiration using tier %s: %s", tier.Name, expirationTime.Format(time.RFC3339))
	} else {
		log.Printf("[INFO] No expiration calculated (tier %s has no expiration)", tier.Name)
	}

	return expirationTime, nil
}

// GetExpiringBookings returns bookings that will expire within the specified hours
func (s *BookingServiceImpl) GetExpiringBookings(ctx context.Context, withinHours int, offset, limit int) ([]*dto.BookingResponse, error) {
	log.Printf("[DEBUG] Getting bookings expiring within %d hours", withinHours)

	// Calculate the threshold time
	thresholdTime := time.Now().Add(time.Duration(withinHours) * time.Hour)

	// Query bookings with expiration time within the threshold
	// Note: This would require a new repository method - implementing basic logic here
	bookings, err := s.bookingRepo.List(ctx, offset, limit)
	if err != nil {
		log.Printf("[ERROR] Failed to get bookings: %v", err)
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	var expiringBookings []*dto.BookingResponse
	for _, booking := range bookings {
		if booking.ExpiresAt != nil && booking.ExpiresAt.Before(thresholdTime) && booking.CanExpire() {
			bookingResponse := s.toBookingResponse(ctx, booking)
			expiringBookings = append(expiringBookings, bookingResponse)
		}
	}

	log.Printf("[SUCCESS] Found %d bookings expiring within %d hours", len(expiringBookings), withinHours)
	return expiringBookings, nil
}

// ProcessExpiredBookings processes expired bookings (cancel or notify)
func (s *BookingServiceImpl) ProcessExpiredBookings(ctx context.Context, dryRun bool) (*dto.ExpirationProcessResult, error) {
	log.Printf("[DEBUG] Processing expired bookings - dry run: %t", dryRun)

	processedAt := time.Now()
	result := &dto.ExpirationProcessResult{
		ProcessedAt:       processedAt.Format(time.RFC3339),
		DryRun:            dryRun,
		ProcessedBookings: []dto.ProcessedBooking{},
		Errors:            []string{},
	}

	// Get all bookings with expiration times
	allBookings, err := s.bookingRepo.List(ctx, 0, 1000) // Adjust limit as needed
	if err != nil {
		log.Printf("[ERROR] Failed to get bookings for expiration processing: %v", err)
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	result.TotalChecked = len(allBookings)
	now := time.Now()

	for _, booking := range allBookings {
		if booking.ExpiresAt == nil || !booking.CanExpire() {
			continue
		}

		if booking.IsExpired() {
			result.ExpiredFound++

			processedBooking := dto.ProcessedBooking{
				BookingID:    booking.ID.String(),
				InvoiceID:    booking.InvoiceID,
				CustomerName: booking.CustomerName,
				ExpiresAt:    booking.ExpiresAt.Format(time.RFC3339),
				ProcessedAt:  processedAt.Format(time.RFC3339),
				Success:      true,
			}

			if !dryRun {
				// ส่ง email แจ้ง agent ว่าหมดอายุ (ก่อน cancel เพื่อให้ status ยังเป็น pending/processing)
				s.sendBookingEmail(ctx, booking, "booking_expired", "")

				// Actually process the booking (cancel it)
				cancelReason := fmt.Sprintf("Booking expired at %s", booking.ExpiresAt.Format(time.RFC3339))
				err := s.CancelBooking(ctx, booking.ID, cancelReason)
				if err != nil {
					processedBooking.Success = false
					errorMsg := err.Error()
					processedBooking.Error = &errorMsg
					processedBooking.Action = "failed"
					processedBooking.Reason = fmt.Sprintf("Failed to cancel: %s", err.Error())
					result.FailedCount++
					result.Errors = append(result.Errors, fmt.Sprintf("Failed to cancel booking %s: %s", booking.ID.String(), err.Error()))
				} else {
					processedBooking.Action = "canceled"
					processedBooking.Reason = cancelReason
					result.ProcessedCount++
				}
			} else {
				processedBooking.Action = "would_cancel"
				processedBooking.Reason = fmt.Sprintf("Would cancel - expired %s ago", now.Sub(*booking.ExpiresAt).String())
			}

			result.ProcessedBookings = append(result.ProcessedBookings, processedBooking)
		}
	}

	log.Printf("[SUCCESS] Processed expired bookings - checked: %d, expired: %d, processed: %d, failed: %d",
		result.TotalChecked, result.ExpiredFound, result.ProcessedCount, result.FailedCount)

	return result, nil
}

// RefreshBookingExpirations recalculates expiration times for existing bookings
func (s *BookingServiceImpl) RefreshBookingExpirations(ctx context.Context, periodID *uuid.UUID) error {
	log.Printf("[DEBUG] Refreshing booking expirations - period filter: %v", periodID)

	// Get bookings to refresh
	var bookings []*models.Booking
	var err error

	if periodID != nil {
		// Get bookings for specific period
		bookings, err = s.bookingRepo.GetByPeriodID(ctx, *periodID, 0, 1000)
	} else {
		// Get all pending bookings
		bookings, err = s.bookingRepo.List(ctx, 0, 1000)
	}

	if err != nil {
		log.Printf("[ERROR] Failed to get bookings for refresh: %v", err)
		return fmt.Errorf("failed to get bookings: %w", err)
	}

	var refreshed, failed int

	for _, booking := range bookings {
		if !booking.CanExpire() {
			continue // Skip non-pending bookings
		}

		err := s.SetBookingExpiration(ctx, booking.ID, booking.PeriodID)
		if err != nil {
			log.Printf("[ERROR] Failed to refresh expiration for booking %s: %v", booking.ID.String(), err)
			failed++
		} else {
			refreshed++
		}
	}

	log.Printf("[SUCCESS] Refreshed booking expirations - refreshed: %d, failed: %d", refreshed, failed)
	return nil
}

// GetBookingExpirationStats returns statistics about booking expirations
func (s *BookingServiceImpl) GetBookingExpirationStats(ctx context.Context) (*dto.ExpirationStatsResponse, error) {
	log.Printf("[DEBUG] Getting booking expiration statistics")

	// Get all bookings (this is a simplified implementation)
	bookings, err := s.bookingRepo.List(ctx, 0, 10000) // Adjust limit as needed
	if err != nil {
		log.Printf("[ERROR] Failed to get bookings for stats: %v", err)
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	stats := &dto.ExpirationStatsResponse{
		BookingsByTier:   make(map[string]int),
		ExpirationTrends: []dto.ExpirationTrendData{},
	}

	nearExpiringThreshold := 24 * time.Hour

	for _, booking := range bookings {
		if booking.ExpiresAt != nil {
			stats.TotalBookingsWithExpiration++

			if booking.IsExpired() {
				stats.ExpiredBookings++
			} else if booking.IsNearExpiration(nearExpiringThreshold) {
				stats.NearExpiringBookings++
			}
		}
	}

	// Calculate average expiration time (simplified)
	if stats.TotalBookingsWithExpiration > 0 {
		stats.AverageExpirationTime = "24 hours" // Placeholder - would calculate actual average
	}

	log.Printf("[SUCCESS] Generated expiration stats - total: %d, expired: %d, near expiring: %d",
		stats.TotalBookingsWithExpiration, stats.ExpiredBookings, stats.NearExpiringBookings)

	return stats, nil
}

// CreateBookingExpiryJob creates an individual expiry job for a booking
func (s *BookingServiceImpl) CreateBookingExpiryJob(ctx context.Context, bookingID uuid.UUID, expiresAt time.Time) error {
	if s.jobService == nil {
		log.Printf("[WARNING] JobService not available, skipping expiry job creation for booking %s", bookingID)
		return nil
	}

	// Check if expiry time is in the past
	now := time.Now()
	if expiresAt.Before(now) {
		log.Printf("[WARNING] Booking %s has already expired at %s, expiring immediately", bookingID, expiresAt.Format(time.RFC3339))
		// Execute expiry immediately instead of scheduling
		return s.ExpireBooking(ctx, bookingID)
	}

	// Create cron expression for the exact expiry time
	// Standard 5-field cron format: minute hour day month dayofweek
	cronExpr := fmt.Sprintf("%d %d %d %d *",
		expiresAt.Minute(),
		expiresAt.Hour(),
		expiresAt.Day(),
		int(expiresAt.Month()))

	// Create job payload with booking ID and action
	payload := fmt.Sprintf(`{"type":"booking_expiry","booking_id":"%s","action":"expire"}`, bookingID)

	jobName := fmt.Sprintf("expire-booking-%s", bookingID)

	// Validate cron expression before creating job
	if err := scheduler.ValidateCronExpression(cronExpr); err != nil {
		log.Printf("[ERROR] Invalid cron expression %s for booking %s: %v", cronExpr, bookingID, err)
		// Don't return error - let backup cleanup handle it
		log.Printf("[INFO] Booking %s will be handled by backup cleanup job", bookingID)
		return nil
	}

	jobReq := &dto.CreateJobRequest{
		Name:     jobName,
		CronExpr: cronExpr,
		Payload:  payload,
	}

	_, err := s.jobService.CreateJob(ctx, jobReq)
	if err != nil {
		log.Printf("[ERROR] Failed to create expiry job for booking %s: %v", bookingID, err)
		// Don't return error - let backup cleanup handle it
		log.Printf("[INFO] Booking %s will be handled by backup cleanup job", bookingID)
		return nil
	}

	log.Printf("[SUCCESS] Created expiry job '%s' for booking %s, expires at %s (cron: %s)", jobName, bookingID, expiresAt.Format(time.RFC3339), cronExpr)
	return nil
}

// DeleteBookingExpiryJob deletes the expiry job for a booking
func (s *BookingServiceImpl) DeleteBookingExpiryJob(ctx context.Context, bookingID uuid.UUID) error {
	if s.jobService == nil {
		log.Printf("[WARNING] JobService not available, skipping expiry job deletion for booking %s", bookingID)
		return nil
	}

	jobName := fmt.Sprintf("expire-booking-%s", bookingID)

	// Find the job by name first
	jobs, _, err := s.jobService.ListJobs(ctx, 0, 1000)
	if err != nil {
		log.Printf("[ERROR] Failed to list jobs while looking for %s: %v", jobName, err)
		return fmt.Errorf("failed to list jobs: %w", err)
	}

	for _, job := range jobs {
		if job.Name == jobName {
			err := s.jobService.DeleteJob(ctx, job.ID)
			if err != nil {
				log.Printf("[ERROR] Failed to delete expiry job %s for booking %s: %v", job.ID, bookingID, err)
				return fmt.Errorf("failed to delete expiry job: %w", err)
			}
			log.Printf("[SUCCESS] Deleted expiry job for booking %s", bookingID)
			return nil
		}
	}

	log.Printf("[INFO] No expiry job found for booking %s", bookingID)
	return nil
}

// ExpireBooking marks a booking as expired and returns seats to availability
func (s *BookingServiceImpl) ExpireBooking(ctx context.Context, bookingID uuid.UUID) error {
	log.Printf("[INFO] Expiring booking %s", bookingID)

	// Get the booking
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Printf("[ERROR] Failed to get booking %s for expiry: %v", bookingID, err)
		return fmt.Errorf("failed to get booking: %w", err)
	}

	if booking == nil {
		log.Printf("[WARNING] Booking %s not found during expiry", bookingID)
		return apperrors.ErrBookingNotFound
	}

	// Check if booking can be expired
	if !booking.CanExpire() {
		log.Printf("[WARNING] Booking %s cannot be expired (status: %s)", bookingID, booking.Status)
		return fmt.Errorf("booking cannot be expired (current status: %s)", booking.Status)
	}

	// Note: No need to fetch period since we use real-time seat calculation

	// Mark booking as expired
	booking.MarkAsExpired()

	// Update booking in database
	err = s.bookingRepo.Update(ctx, booking.ID, booking)
	if err != nil {
		log.Printf("[ERROR] Failed to update booking %s to expired status: %v", bookingID, err)
		return fmt.Errorf("failed to update booking status: %w", err)
	}

	// Note: No need to manually update period booking count since we now use real-time calculation
	log.Printf("[SUCCESS] Booking %s expired, seats are now available through real-time calculation", bookingID)

	// Broadcast seat update via WebSocket
	go websocketInfra.BroadcastSeatUpdate(ctx, booking.PeriodID, "booking_expired", s.seatAvailabilityService)

	// Delete the expiry job since it's no longer needed
	err = s.DeleteBookingExpiryJob(ctx, bookingID)
	if err != nil {
		log.Printf("[WARNING] Failed to delete expiry job for booking %s: %v", bookingID, err)
		// Don't fail the expiry for this
	}

	log.Printf("[SUCCESS] Booking %s expired successfully", bookingID)
	return nil
}

// Payment Management Methods

// UploadPaymentSlip uploads a payment slip and creates a payment record
func (s *BookingServiceImpl) UploadPaymentSlip(ctx context.Context, req dto.UploadPaymentSlipRequest) (*dto.PaymentSlipResponse, error) {
	log.Printf("[INFO] Uploading payment slip for booking %s with amount %.2f", req.BookingID, req.Amount)

	// Verify booking exists
	booking, err := s.bookingRepo.GetByID(ctx, req.BookingID)
	if err != nil {
		log.Printf("[ERROR] Failed to get booking %s: %v", req.BookingID, err)
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return nil, apperrors.ErrBookingNotFound
	}

	// Validate payment amount doesn't exceed booking total
	if req.Amount > booking.SumNet {
		return nil, fmt.Errorf("payment amount (%.2f) cannot exceed booking total (%.2f)", req.Amount, booking.SumNet)
	}

	// Check current payment total
	currentPayments, err := s.paymentSlipRepo.GetTotalVerifiedAmountByBookingID(ctx, req.BookingID)
	if err != nil {
		log.Printf("[ERROR] Failed to get current payments for booking %s: %v", req.BookingID, err)
		return nil, fmt.Errorf("failed to check current payments: %w", err)
	}

	// Validate total payments won't exceed booking total
	if currentPayments+req.Amount > booking.SumNet {
		return nil, fmt.Errorf("total payments (%.2f + %.2f = %.2f) would exceed booking total (%.2f)",
			currentPayments, req.Amount, currentPayments+req.Amount, booking.SumNet)
	}

	// Create payment slip with retry mechanism for duplicate slip_id handling
	const maxCreationRetries = 3
	var paymentSlip *models.PaymentSlip
	var creationErr error

	for attempt := 1; attempt <= maxCreationRetries; attempt++ {
		// Generate slip ID
		slipID, err := s.paymentSlipRepo.GenerateSlipID(ctx)
		if err != nil {
			log.Printf("[ERROR] Failed to generate slip ID on attempt %d: %v", attempt, err)
			return nil, fmt.Errorf("failed to generate slip ID: %w", err)
		}

		// Create payment slip
		paymentSlip = &models.PaymentSlip{
			BookingID: req.BookingID,
			SlipID:    slipID,
			Image:     &req.Image,
			Amount:    req.Amount,
			Status:    models.PaymentSlipStatusPending,
		}

		// SlipOK verification (ถ้า configured + มี image data)
		if s.slipVerification != nil && len(req.ImageData) > 0 {
			verifyAmount := req.Amount
			verifyResult, verifyErr := s.slipVerification.VerifySlip(ctx, req.ImageData, &verifyAmount)
			if verifyErr != nil {
				log.Printf("[WARNING] SlipOK verification failed: %v (continuing without verification)", verifyErr)
				errMsg := verifyErr.Error()
				paymentSlip.SlipVerifyError = &errMsg
			} else if verifyResult != nil {
				paymentSlip.SlipVerified = verifyResult.Success

				if !verifyResult.TransTimestamp.IsZero() {
					paymentSlip.PaymentDate = &verifyResult.TransTimestamp
				}
				if verifyResult.TransRef != "" {
					paymentSlip.TransRef = &verifyResult.TransRef
				}

				// Map bank codes → ชื่อย่อ
				if verifyResult.SendingBank != "" {
					bankName := verifyResult.SendingBank
					if name, ok := ports.BankCodeToName[verifyResult.SendingBank]; ok {
						bankName = name
					}
					paymentSlip.SendingBank = &bankName
				}
				if verifyResult.ReceivingBank != "" {
					bankName := verifyResult.ReceivingBank
					if name, ok := ports.BankCodeToName[verifyResult.ReceivingBank]; ok {
						bankName = name
					}
					paymentSlip.ReceivingBank = &bankName
				}

				if verifyResult.ReceiverAccount != "" {
					paymentSlip.ReceiverAccount = &verifyResult.ReceiverAccount
				}
				if verifyResult.ReceiverDisplayName != "" {
					paymentSlip.ReceiverName = &verifyResult.ReceiverDisplayName
				}
				if verifyResult.SenderDisplayName != "" {
					paymentSlip.SenderName = &verifyResult.SenderDisplayName
				}

				if !verifyResult.Success {
					errMsg := verifyResult.Message
					paymentSlip.SlipVerifyError = &errMsg
				}

				log.Printf("[SlipOK] Verified=%v Amount=%.2f Bank=%s→%s Ref=%s",
					verifyResult.Success, verifyResult.Amount,
					verifyResult.SendingBank, verifyResult.ReceivingBank, verifyResult.TransRef)
			}
		}

		creationErr = s.paymentSlipRepo.Create(ctx, paymentSlip)
		if creationErr == nil {
			// Success, break the retry loop
			log.Printf("[SUCCESS] Payment slip created successfully with slip_id: %s", slipID)
			break
		}

		// Check if error is due to duplicate slip_id
		if strings.Contains(strings.ToLower(creationErr.Error()), "unique") &&
			strings.Contains(strings.ToLower(creationErr.Error()), "slip_id") {
			log.Printf("[WARNING] Duplicate slip_id detected on attempt %d: %s, retrying...", attempt, slipID)

			if attempt < maxCreationRetries {
				// Small delay before retry
				time.Sleep(time.Millisecond * 100 * time.Duration(attempt))
				continue
			}
		}

		// Non-duplicate error or max retries reached
		log.Printf("[ERROR] Failed to create payment slip on attempt %d: %v", attempt, creationErr)
		return nil, fmt.Errorf("failed to create payment slip after %d attempts: %w", attempt, creationErr)
	}

	if creationErr != nil {
		return nil, fmt.Errorf("failed to create payment slip after %d retries: %w", maxCreationRetries, creationErr)
	}

	// SlipOK verified → auto-confirm + recalculate payment (Agent flow)
	if paymentSlip.SlipVerified {
		paymentSlip.Status = models.PaymentSlipStatusConfirm
		if err := s.paymentSlipRepo.Update(ctx, paymentSlip.ID, paymentSlip); err != nil {
			log.Printf("[WARNING] Failed to auto-confirm payment slip: %v", err)
		} else {
			log.Printf("[SUCCESS] Payment slip %s AUTO-CONFIRMED via SlipOK", paymentSlip.SlipID)
		}

		// Recalculate booking payment total + update status
		totalActivePayments, err := s.paymentSlipRepo.GetTotalActiveAmountByBookingID(ctx, req.BookingID)
		if err != nil {
			log.Printf("[WARNING] Failed to get total active payments: %v", err)
		} else {
			booking.SumPayment = totalActivePayments
			booking.UpdatedAt = time.Now()

			if totalActivePayments >= booking.SumNet {
				booking.Status = models.BookingStatusConfirm
				log.Printf("[SUCCESS] Booking %s FULLY PAID (%.2f/%.2f) via SlipOK auto-confirm", booking.ID, totalActivePayments, booking.SumNet)
			} else if booking.Status == models.BookingStatusPending {
				booking.Status = models.BookingStatusProcessing
			}

			if err := s.bookingRepo.Update(ctx, booking.ID, booking); err != nil {
				log.Printf("[WARNING] Failed to update booking after auto-confirm: %v", err)
			}

			// ถ้า booking confirmed → ลบ expiry job + ส่ง email
			if booking.Status == models.BookingStatusConfirm {
				if err := s.DeleteBookingExpiryJob(ctx, booking.ID); err != nil {
					log.Printf("[WARNING] Failed to delete expiry job for auto-confirmed booking %s: %v", booking.ID, err)
				}
				s.sendBookingEmail(ctx, booking, "payment_confirmed", "")
			}
		}
	} else {
		// SlipOK ไม่ผ่าน หรือไม่ได้ verify → status pending, booking → processing (เหมือนเดิม)
		if booking.Status == models.BookingStatusPending {
			booking.Status = models.BookingStatusProcessing
			booking.UpdatedAt = time.Now()
			if err := s.bookingRepo.Update(ctx, booking.ID, booking); err != nil {
				log.Printf("[WARNING] Failed to update booking status to processing: %v", err)
			} else {
				log.Printf("[SUCCESS] Booking %s status updated to PROCESSING (payment slip uploaded)", booking.ID)
			}
		}
	}

	log.Printf("[SUCCESS] Payment slip created: ID=%s, SlipID=%s, Amount=%.2f, Verified=%v", paymentSlip.ID, paymentSlip.SlipID, paymentSlip.Amount, paymentSlip.SlipVerified)

	return s.toPaymentSlipResponse(paymentSlip), nil
}

// VerifyPaymentSlip verifies or rejects a payment slip and updates booking payment
func (s *BookingServiceImpl) VerifyPaymentSlip(ctx context.Context, req dto.VerifyPaymentSlipRequest) (*dto.PaymentSlipResponse, error) {
	log.Printf("[INFO] Verifying payment slip %s with status %s", req.PaymentSlipID, req.Status)

	// Get payment slip
	paymentSlip, err := s.paymentSlipRepo.GetByID(ctx, req.PaymentSlipID)
	if err != nil {
		log.Printf("[ERROR] Failed to get payment slip %s: %v", req.PaymentSlipID, err)
		return nil, fmt.Errorf("failed to get payment slip: %w", err)
	}
	if paymentSlip == nil {
		return nil, fmt.Errorf("payment slip not found")
	}

	// Validate payment slip is pending
	if paymentSlip.Status != models.PaymentSlipStatusPending {
		return nil, fmt.Errorf("payment slip is not pending (current status: %s)", paymentSlip.Status)
	}

	// Update payment slip status
	var newStatus models.PaymentSlipStatus
	switch req.Status {
	case "confirm":
		newStatus = models.PaymentSlipStatusConfirm
	case "reject":
		newStatus = models.PaymentSlipStatusReject
	default:
		return nil, fmt.Errorf("invalid status: %s", req.Status)
	}

	// Handle amount adjustment if provided
	amountAdjusted := false
	if req.ActualAmount != nil && *req.ActualAmount != paymentSlip.Amount {
		// Get booking for validation
		booking, err := s.bookingRepo.GetByID(ctx, paymentSlip.BookingID)
		if err != nil {
			log.Printf("[ERROR] Failed to get booking for amount validation: %v", err)
			return nil, fmt.Errorf("failed to validate amount adjustment: %w", err)
		}

		// Validate actual amount is reasonable
		if *req.ActualAmount > booking.SumNet {
			return nil, fmt.Errorf("actual amount (%.2f) cannot exceed booking total (%.2f)", *req.ActualAmount, booking.SumNet)
		}

		// Check if adjustment reason is provided for significant changes
		percentDiff := math.Abs((*req.ActualAmount - paymentSlip.Amount) / paymentSlip.Amount * 100)
		if percentDiff > 10 && (req.AdjustmentReason == nil || *req.AdjustmentReason == "") {
			return nil, fmt.Errorf("adjustment reason required for amount changes over 10%% (%.1f%% change)", percentDiff)
		}

		// Store original amount and update
		if paymentSlip.OriginalAmount == nil {
			paymentSlip.OriginalAmount = &paymentSlip.Amount
		}
		paymentSlip.Amount = *req.ActualAmount
		paymentSlip.AdjustmentReason = req.AdjustmentReason
		amountAdjusted = true

		log.Printf("[INFO] Payment slip %s amount adjusted: %.2f → %.2f (%.1f%% change)",
			paymentSlip.SlipID, *paymentSlip.OriginalAmount, paymentSlip.Amount, percentDiff)
	}

	// Update payment slip with new status, note, and amount if adjusted
	paymentSlip.Status = newStatus
	paymentSlip.Note = req.Note
	paymentSlip.UpdatedAt = time.Now()

	err = s.paymentSlipRepo.Update(ctx, paymentSlip.ID, paymentSlip)
	if err != nil {
		log.Printf("[ERROR] Failed to update payment slip: %v", err)
		return nil, fmt.Errorf("failed to update payment slip: %w", err)
	}

	// If confirmed, update booking payment total and status
	if newStatus == models.PaymentSlipStatusConfirm {
		// Get booking
		booking, err := s.bookingRepo.GetByID(ctx, paymentSlip.BookingID)
		if err != nil {
			log.Printf("[ERROR] Failed to get booking %s: %v", paymentSlip.BookingID, err)
			return nil, fmt.Errorf("failed to get booking: %w", err)
		}

		// Store original status for comparison
		originalStatus := booking.Status

		// Validate booking is in a state that allows payment updates
		if booking.Status == models.BookingStatusCancel {
			log.Printf("[WARNING] Cannot update payment for cancelled booking %s", booking.ID)
			return nil, fmt.Errorf("cannot process payment for cancelled booking")
		}
		if booking.Status == models.BookingStatusExpired {
			log.Printf("[WARNING] Cannot update payment for expired booking %s", booking.ID)
			return nil, fmt.Errorf("cannot process payment for expired booking")
		}
		if booking.Status == models.BookingStatusCompleted {
			log.Printf("[INFO] Payment update for already completed booking %s", booking.ID)
			// Allow payment updates for completed bookings (for record keeping)
		}
		if booking.Status == models.BookingStatusProcessing {
			log.Printf("[INFO] Verifying payment for processing booking %s", booking.ID)
			// Processing bookings can be verified - this is the normal workflow
		}

		// Calculate total active payments (confirmed and non-reversed)
		totalActivePayments, err := s.paymentSlipRepo.GetTotalActiveAmountByBookingID(ctx, paymentSlip.BookingID)
		if err != nil {
			log.Printf("[ERROR] Failed to get total active payments: %v", err)
			return nil, fmt.Errorf("failed to calculate total payments: %w", err)
		}

		// Update booking payment total
		booking.SumPayment = totalActivePayments
		booking.UpdatedAt = time.Now()

		// Enhanced business logic for status updates with deposit thresholds
		var statusChanged bool

		if totalActivePayments >= booking.SumNet {
			// Fully paid - confirm booking
			if booking.Status != models.BookingStatusConfirm {
				booking.Status = models.BookingStatusConfirm
				statusChanged = true
				log.Printf("[SUCCESS] Booking %s FULLY PAID (%.2f/%.2f) - status updated from %s to CONFIRMED",
					booking.ID, totalActivePayments, booking.SumNet, originalStatus)
			}
		} else if totalActivePayments >= booking.SumDeposit {
			// Deposit threshold met - move to processing (protected from expiry)
			if booking.Status == models.BookingStatusPending {
				booking.Status = models.BookingStatusProcessing
				statusChanged = true
				log.Printf("[SUCCESS] Booking %s DEPOSIT PAID (%.2f >= %.2f) - status updated from %s to PROCESSING",
					booking.ID, totalActivePayments, booking.SumDeposit, originalStatus)
			} else if booking.Status == models.BookingStatusProcessing {
				log.Printf("[INFO] Booking %s PROCESSING - deposit satisfied (%.2f >= %.2f), remaining: %.2f",
					booking.ID, totalActivePayments, booking.SumDeposit, booking.SumNet-totalActivePayments)
			}
		} else if totalActivePayments > 0 {
			// Some payment received but below deposit threshold
			if booking.Status == models.BookingStatusPending {
				booking.Status = models.BookingStatusProcessing
				statusChanged = true
				log.Printf("[INFO] Booking %s PARTIAL PAYMENT (%.2f < %.2f deposit) - status updated to PROCESSING for protection",
					booking.ID, totalActivePayments, booking.SumDeposit)
			} else if booking.Status == models.BookingStatusProcessing {
				log.Printf("[INFO] Booking %s PROCESSING - partial payment (%.2f < %.2f deposit required)",
					booking.ID, totalActivePayments, booking.SumDeposit)
			}
		}

		// Save booking updates
		err = s.bookingRepo.Update(ctx, booking.ID, booking)
		if err != nil {
			log.Printf("[ERROR] Failed to update booking payment status: %v", err)
			return nil, fmt.Errorf("failed to update booking: %w", err)
		}

		// If booking is now confirmed, cancel the expiry job
		if statusChanged && booking.Status == models.BookingStatusConfirm {
			if err := s.DeleteBookingExpiryJob(ctx, booking.ID); err != nil {
				log.Printf("[WARNING] Failed to delete booking expiry job for confirmed booking %s: %v", booking.ID, err)
				// Don't fail the payment verification if job deletion fails
			} else {
				log.Printf("[SUCCESS] Booking expiry job deleted for confirmed booking %s", booking.ID)
			}
		}

		// Enhanced success logging with amount adjustment information
		if amountAdjusted {
			log.Printf("[SUCCESS] Payment verification completed: Booking %s payment updated to %.2f, status: %s -> %s (amount adjusted from %.2f to %.2f)",
				booking.ID, booking.SumPayment, originalStatus, booking.Status, *paymentSlip.OriginalAmount, paymentSlip.Amount)
		} else {
			log.Printf("[SUCCESS] Payment verification completed: Booking %s payment updated to %.2f, status: %s -> %s",
				booking.ID, booking.SumPayment, originalStatus, booking.Status)
		}

		// ส่ง email แจ้ง agent — ชำระเงินครบ ยืนยันการจอง
		if statusChanged && booking.Status == models.BookingStatusConfirm {
			s.sendBookingEmail(ctx, booking, "payment_confirmed", "")
		}
	} else {
		// Payment slip rejected - need to recalculate booking status based on remaining verified payments
		log.Printf("[INFO] Payment slip %s REJECTED - analyzing payment history for booking status update", req.PaymentSlipID)

		// Get booking for payment history analysis
		booking, err := s.bookingRepo.GetByID(ctx, paymentSlip.BookingID)
		if err != nil {
			log.Printf("[ERROR] Failed to get booking %s for rejection analysis: %v", paymentSlip.BookingID, err)
			// Continue with rejection but without booking status update
		} else {
			// Store original status for comparison
			originalStatus := booking.Status

			// Calculate remaining active payments (excluding this rejected slip and any other rejected slips)
			totalActivePayments, err := s.paymentSlipRepo.GetTotalActiveAmountByBookingID(ctx, paymentSlip.BookingID)
			if err != nil {
				log.Printf("[ERROR] Failed to calculate remaining active payments after rejection: %v", err)
			} else {
				log.Printf("[DEBUG] Rejection analysis for booking %s: remaining active payments = %.2f", booking.ID, totalActivePayments)

				// Update booking payment total based on remaining active payments
				booking.SumPayment = totalActivePayments
				booking.UpdatedAt = time.Now()

				// Apply business logic for status updates based on remaining payments
				var statusChanged bool

				if totalActivePayments >= booking.SumNet {
					// Still fully paid after rejection
					if booking.Status != models.BookingStatusConfirm {
						booking.Status = models.BookingStatusConfirm
						statusChanged = true
						log.Printf("[SUCCESS] Booking %s remains FULLY PAID after rejection (%.2f/%.2f) - status: %s -> CONFIRMED",
							booking.ID, totalActivePayments, booking.SumNet, originalStatus)
					} else {
						log.Printf("[INFO] Booking %s remains FULLY PAID after rejection (%.2f/%.2f) - status unchanged",
							booking.ID, totalActivePayments, booking.SumNet)
					}
				} else if totalActivePayments >= booking.SumDeposit {
					// Deposit threshold still met after rejection - protected from expiry
					if booking.Status == models.BookingStatusPending || booking.Status == models.BookingStatusConfirm {
						booking.Status = models.BookingStatusProcessing
						statusChanged = true
						log.Printf("[SUCCESS] Booking %s DEPOSIT SATISFIED after rejection (%.2f >= %.2f) - status: %s -> PROCESSING",
							booking.ID, totalActivePayments, booking.SumDeposit, originalStatus)
					} else {
						log.Printf("[INFO] Booking %s DEPOSIT SATISFIED after rejection (%.2f >= %.2f) - status unchanged: %s",
							booking.ID, totalActivePayments, booking.SumDeposit, booking.Status)
					}
				} else if totalActivePayments > 0 {
					// Some payments remain but below deposit threshold - still protected but vulnerable
					if booking.Status == models.BookingStatusConfirm {
						booking.Status = models.BookingStatusProcessing
						statusChanged = true
						log.Printf("[WARNING] Booking %s PARTIAL PAYMENT after rejection (%.2f < %.2f deposit) - status: %s -> PROCESSING",
							booking.ID, totalActivePayments, booking.SumDeposit, originalStatus)
					} else {
						log.Printf("[INFO] Booking %s PARTIAL PAYMENT after rejection (%.2f < %.2f deposit) - status unchanged: %s",
							booking.ID, totalActivePayments, booking.SumDeposit, booking.Status)
					}
				} else {
					// No active payments remain - vulnerable to expiry
					if booking.Status != models.BookingStatusPending {
						booking.Status = models.BookingStatusPending
						statusChanged = true
						log.Printf("[WARNING] Booking %s NO ACTIVE PAYMENTS after rejection - status: %s -> PENDING (vulnerable to expiry)",
							booking.ID, originalStatus)
					} else {
						log.Printf("[INFO] Booking %s NO ACTIVE PAYMENTS after rejection - status unchanged: PENDING",
							booking.ID)
					}
				}

				// Save booking updates if status changed
				if statusChanged {
					err = s.bookingRepo.Update(ctx, booking.ID, booking)
					if err != nil {
						log.Printf("[ERROR] Failed to update booking status after rejection: %v", err)
						// Don't fail the rejection process if booking update fails
					} else {
						log.Printf("[SUCCESS] Booking status updated after payment slip rejection")
					}
				}

				// Enhanced rejection logging with payment analysis
				log.Printf("[SUCCESS] Payment rejection analysis completed: Booking %s payment updated to %.2f, status: %s -> %s",
					booking.ID, booking.SumPayment, originalStatus, booking.Status)
			}
		}

		log.Printf("[INFO] Payment slip %s successfully REJECTED", req.PaymentSlipID)
	}

	return s.toPaymentSlipResponse(paymentSlip), nil
}

// GetPendingPayments retrieves all pending payment slips for admin review
func (s *BookingServiceImpl) GetPendingPayments(ctx context.Context, offset, limit int) ([]*dto.PaymentSlipResponse, error) {
	log.Printf("[INFO] Getting pending payments: offset=%d, limit=%d", offset, limit)

	paymentSlips, err := s.paymentSlipRepo.GetPendingSlips(ctx, offset, limit)
	if err != nil {
		log.Printf("[ERROR] Failed to get pending payments: %v", err)
		return nil, fmt.Errorf("failed to get pending payments: %w", err)
	}

	responses := make([]*dto.PaymentSlipResponse, len(paymentSlips))
	for i, slip := range paymentSlips {
		responses[i] = s.toPaymentSlipResponse(slip)

		// Optionally include booking details
		if slip.BookingID != uuid.Nil {
			booking, err := s.GetBookingByID(ctx, slip.BookingID)
			if err == nil {
				responses[i].Booking = booking
			}
		}
	}

	log.Printf("[SUCCESS] Retrieved %d pending payments", len(responses))
	return responses, nil
}

// Helper method to convert PaymentSlip to PaymentSlipResponse
func (s *BookingServiceImpl) toPaymentSlipResponse(paymentSlip *models.PaymentSlip) *dto.PaymentSlipResponse {
	response := &dto.PaymentSlipResponse{
		ID:               paymentSlip.ID,
		BookingID:        paymentSlip.BookingID,
		SlipID:           paymentSlip.SlipID,
		Image:            paymentSlip.Image,
		Amount:           paymentSlip.Amount,
		OriginalAmount:   paymentSlip.OriginalAmount,
		AdjustmentReason: paymentSlip.AdjustmentReason,
		Status:           string(paymentSlip.Status),
		Note:             paymentSlip.Note,

		// Admin payment tracking fields
		PaymentMethod:    string(paymentSlip.PaymentMethod),
		IsAdminCreated:   paymentSlip.IsAdminCreated,
		AdminCreatedNote: paymentSlip.AdminCreatedNote,

		// SlipOK verification fields
		SendingBank:     paymentSlip.SendingBank,
		ReceivingBank:   paymentSlip.ReceivingBank,
		ReceiverAccount: paymentSlip.ReceiverAccount,
		ReceiverName:    paymentSlip.ReceiverName,
		SenderName:      paymentSlip.SenderName,
		TransRef:        paymentSlip.TransRef,
		SlipVerified:    paymentSlip.SlipVerified,
		SlipVerifyError: paymentSlip.SlipVerifyError,

		// Reversal tracking fields
		IsReversed:     paymentSlip.IsReversed,
		ReversalReason: paymentSlip.ReversalReason,
		CreatedAt:      paymentSlip.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      paymentSlip.UpdatedAt.Format(time.RFC3339),
	}

	// Format payment_date if exists
	if paymentSlip.PaymentDate != nil {
		pd := paymentSlip.PaymentDate.Format(time.RFC3339)
		response.PaymentDate = &pd
	}

	// Format reversed_at if exists
	if paymentSlip.ReversedAt != nil {
		reversedAt := paymentSlip.ReversedAt.Format(time.RFC3339)
		response.ReversedAt = &reversedAt
	}

	// Include reversed by user ID if exists
	if paymentSlip.ReversedBy != nil {
		reversedBy := paymentSlip.ReversedBy.String()
		response.ReversedBy = &reversedBy
	}

	// Include created by admin user ID if exists
	if paymentSlip.CreatedByAdmin != nil {
		createdByAdmin := paymentSlip.CreatedByAdmin.String()
		response.CreatedByAdmin = &createdByAdmin
	}

	return response
}

// ReversePaymentSlip reverses a confirmed payment slip and recalculates booking status
func (s *BookingServiceImpl) ReversePaymentSlip(ctx context.Context, req dto.ReversePaymentSlipRequest, adminUserID uuid.UUID) (*dto.PaymentSlipResponse, error) {
	log.Printf("[INFO] Reversing payment slip %s by admin %s", req.PaymentSlipID, adminUserID)

	// Get the payment slip
	paymentSlip, err := s.paymentSlipRepo.GetByID(ctx, req.PaymentSlipID)
	if err != nil {
		log.Printf("[ERROR] Failed to get payment slip %s: %v", req.PaymentSlipID, err)
		return nil, fmt.Errorf("failed to get payment slip: %w", err)
	}

	// Validate that the payment slip can be reversed
	if !paymentSlip.IsReversible() {
		log.Printf("[WARNING] Payment slip %s is not reversible - Status: %s, IsReversed: %t",
			req.PaymentSlipID, paymentSlip.Status, paymentSlip.IsReversed)
		return nil, fmt.Errorf("payment slip cannot be reversed - must be confirmed and not already reversed")
	}

	// Reverse the payment slip in repository
	err = s.paymentSlipRepo.ReversePaymentSlip(ctx, req.PaymentSlipID, adminUserID, req.ReversalReason)
	if err != nil {
		log.Printf("[ERROR] Failed to reverse payment slip %s: %v", req.PaymentSlipID, err)
		return nil, fmt.Errorf("failed to reverse payment slip: %w", err)
	}

	// Get updated payment slip
	paymentSlip, err = s.paymentSlipRepo.GetByID(ctx, req.PaymentSlipID)
	if err != nil {
		log.Printf("[ERROR] Failed to get updated payment slip %s: %v", req.PaymentSlipID, err)
		return nil, fmt.Errorf("failed to get updated payment slip: %w", err)
	}

	// Recalculate booking payment status after reversal
	err = s.recalculateBookingAfterReversal(ctx, paymentSlip.BookingID)
	if err != nil {
		log.Printf("[ERROR] Failed to recalculate booking after reversal: %v", err)
		// Don't fail the reversal if booking recalculation fails, but log the error
	}

	log.Printf("[SUCCESS] Payment slip %s reversed successfully by admin %s", req.PaymentSlipID, adminUserID)
	return s.toPaymentSlipResponse(paymentSlip), nil
}

// GetPaymentSlipHistory returns detailed history and timeline for a payment slip
func (s *BookingServiceImpl) GetPaymentSlipHistory(ctx context.Context, paymentSlipID uuid.UUID) (*dto.PaymentSlipHistoryResponse, error) {
	log.Printf("[INFO] Getting payment slip history for %s", paymentSlipID)

	// Get payment slip with all relations
	paymentSlip, err := s.paymentSlipRepo.GetPaymentSlipHistory(ctx, paymentSlipID)
	if err != nil {
		log.Printf("[ERROR] Failed to get payment slip history %s: %v", paymentSlipID, err)
		return nil, fmt.Errorf("failed to get payment slip history: %w", err)
	}

	// Build timeline
	timeline := s.buildPaymentSlipTimeline(paymentSlip)

	// Build booking info
	bookingInfo := &dto.PaymentSlipBookingInfoResponse{
		BookingID:        paymentSlip.BookingID,
		TotalAmount:      0,
		DepositAmount:    0,
		CurrentPayments:  0,
		RemainingBalance: 0,
		BookingStatus:    "unknown",
	}

	if paymentSlip.Booking != nil {
		// Calculate current active payments
		activePayments, err := s.paymentSlipRepo.GetTotalActiveAmountByBookingID(ctx, paymentSlip.BookingID)
		if err != nil {
			log.Printf("[WARNING] Failed to get active payments for booking %s: %v", paymentSlip.BookingID, err)
		}

		bookingInfo.InvoiceID = paymentSlip.Booking.InvoiceID
		bookingInfo.CustomerName = paymentSlip.Booking.CustomerName
		bookingInfo.TotalAmount = paymentSlip.Booking.SumNet
		bookingInfo.DepositAmount = paymentSlip.Booking.SumDeposit
		bookingInfo.CurrentPayments = activePayments
		bookingInfo.RemainingBalance = paymentSlip.Booking.SumNet - activePayments
		bookingInfo.BookingStatus = string(paymentSlip.Booking.Status)
	}

	return &dto.PaymentSlipHistoryResponse{
		PaymentSlip: s.toPaymentSlipResponse(paymentSlip),
		Timeline:    timeline,
		BookingInfo: bookingInfo,
	}, nil
}

// Helper method to recalculate booking status after payment slip reversal
func (s *BookingServiceImpl) recalculateBookingAfterReversal(ctx context.Context, bookingID uuid.UUID) error {
	log.Printf("[INFO] Recalculating booking %s after payment slip reversal", bookingID)

	// Get booking
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	// Store original values for comparison
	originalStatus := booking.Status
	originalSumPayment := booking.SumPayment

	// Calculate remaining active payments (excluding reversed slips)
	totalActivePayments, err := s.paymentSlipRepo.GetTotalActiveAmountByBookingID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to calculate active payments: %w", err)
	}

	log.Printf("[DEBUG] Reversal recalculation for booking %s: sum_payment %.2f → %.2f (change: %.2f)",
		bookingID, originalSumPayment, totalActivePayments, totalActivePayments-originalSumPayment)

	// Update booking payment total
	booking.SumPayment = totalActivePayments
	booking.UpdatedAt = time.Now()

	// Apply business logic for status updates based on remaining active payments
	var statusChanged bool

	if totalActivePayments >= booking.SumNet {
		// Still fully paid after reversal
		if booking.Status != models.BookingStatusConfirm {
			booking.Status = models.BookingStatusConfirm
			statusChanged = true
			log.Printf("[SUCCESS] Booking %s remains FULLY PAID after reversal (%.2f/%.2f) - status: %s -> CONFIRMED",
				bookingID, totalActivePayments, booking.SumNet, originalStatus)
		}
	} else if totalActivePayments >= booking.SumDeposit {
		// Deposit threshold still met after reversal
		if booking.Status == models.BookingStatusPending || booking.Status == models.BookingStatusConfirm {
			booking.Status = models.BookingStatusProcessing
			statusChanged = true
			log.Printf("[SUCCESS] Booking %s DEPOSIT SATISFIED after reversal (%.2f >= %.2f) - status: %s -> PROCESSING",
				bookingID, totalActivePayments, booking.SumDeposit, originalStatus)
		}
	} else if totalActivePayments > 0 {
		// Some payments remain but below deposit threshold
		if booking.Status == models.BookingStatusConfirm {
			booking.Status = models.BookingStatusProcessing
			statusChanged = true
			log.Printf("[WARNING] Booking %s PARTIAL PAYMENT after reversal (%.2f < %.2f deposit) - status: %s -> PROCESSING",
				bookingID, totalActivePayments, booking.SumDeposit, originalStatus)
		}
	} else {
		// No active payments remain
		if booking.Status != models.BookingStatusPending {
			booking.Status = models.BookingStatusPending
			statusChanged = true
			log.Printf("[WARNING] Booking %s NO ACTIVE PAYMENTS after reversal - status: %s -> PENDING (vulnerable to expiry)",
				bookingID, originalStatus)
		}
	}

	// Save booking updates (always save to update sum_payment, and status if it changed)
	err = s.bookingRepo.Update(ctx, booking.ID, booking)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	if statusChanged {
		log.Printf("[SUCCESS] Booking status updated after reversal: %s -> %s", originalStatus, booking.Status)
	}
	log.Printf("[SUCCESS] Booking sum_payment updated after reversal: %.2f", booking.SumPayment)

	return nil
}

// Helper method to build payment slip timeline
func (s *BookingServiceImpl) buildPaymentSlipTimeline(paymentSlip *models.PaymentSlip) []*dto.PaymentSlipTimelineEntry {
	timeline := []*dto.PaymentSlipTimelineEntry{}

	// Creation entry
	timeline = append(timeline, &dto.PaymentSlipTimelineEntry{
		Action:    "created",
		Timestamp: paymentSlip.CreatedAt.Format(time.RFC3339),
		Details:   fmt.Sprintf("Payment slip created with amount %.2f", paymentSlip.Amount),
		Amount:    &paymentSlip.Amount,
	})

	// Status change entries based on current status
	switch paymentSlip.Status {
	case models.PaymentSlipStatusConfirm:
		timeline = append(timeline, &dto.PaymentSlipTimelineEntry{
			Action:    "confirmed",
			Timestamp: paymentSlip.UpdatedAt.Format(time.RFC3339),
			Details:   "Payment slip verified and confirmed",
			Amount:    &paymentSlip.Amount,
		})
	case models.PaymentSlipStatusReject:
		timeline = append(timeline, &dto.PaymentSlipTimelineEntry{
			Action:    "rejected",
			Timestamp: paymentSlip.UpdatedAt.Format(time.RFC3339),
			Details:   "Payment slip rejected",
		})
	case models.PaymentSlipStatusReversed:
		timeline = append(timeline, &dto.PaymentSlipTimelineEntry{
			Action:    "reversed",
			Timestamp: paymentSlip.UpdatedAt.Format(time.RFC3339),
			Details:   "Payment slip reversed and excluded from calculations",
		})
	}

	// Amount adjustment entry
	if paymentSlip.OriginalAmount != nil {
		timeline = append(timeline, &dto.PaymentSlipTimelineEntry{
			Action:    "amount_adjusted",
			Timestamp: paymentSlip.UpdatedAt.Format(time.RFC3339),
			Details:   fmt.Sprintf("Amount adjusted from %.2f to %.2f", *paymentSlip.OriginalAmount, paymentSlip.Amount),
			Amount:    &paymentSlip.Amount,
		})
	}

	// Reversal entry
	if paymentSlip.IsReversed && paymentSlip.ReversedAt != nil {
		performedBy := "Unknown"
		if paymentSlip.ReversedBy != nil {
			performedBy = paymentSlip.ReversedBy.String()
		}

		details := "Payment reversed"
		if paymentSlip.ReversalReason != nil {
			details = fmt.Sprintf("Payment reversed: %s", *paymentSlip.ReversalReason)
		}

		timeline = append(timeline, &dto.PaymentSlipTimelineEntry{
			Action:      "reversed",
			Timestamp:   paymentSlip.ReversedAt.Format(time.RFC3339),
			PerformedBy: &performedBy,
			Details:     details,
		})
	}

	return timeline
}

// CreateAdminPayment creates a payment slip directly by admin without requiring image upload
func (s *BookingServiceImpl) CreateAdminPayment(ctx context.Context, bookingID uuid.UUID, req dto.CreateAdminPaymentRequest, adminUserID uuid.UUID) (*dto.AdminPaymentResponse, error) {
	log.Printf("[INFO] Admin %s creating direct payment for booking %s: %+v", adminUserID, bookingID, req)

	// Validate payment method
	if !models.IsValidPaymentMethod(req.PaymentMethod) {
		return nil, fmt.Errorf("invalid payment method: %s", req.PaymentMethod)
	}

	// Get booking to verify it exists and calculate status
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Printf("[ERROR] Failed to get booking %s: %v", bookingID, err)
		return nil, fmt.Errorf("booking not found")
	}

	if booking == nil {
		return nil, fmt.Errorf("booking not found")
	}

	// Generate unique slip ID
	slipID, err := s.paymentSlipRepo.GenerateSlipID(ctx)
	if err != nil {
		log.Printf("[ERROR] Failed to generate slip ID: %v", err)
		return nil, fmt.Errorf("failed to generate slip ID: %w", err)
	}

	// Create payment slip with admin-specific fields
	paymentSlip := &models.PaymentSlip{
		BookingID:        bookingID,
		SlipID:           slipID,
		Amount:           req.Amount,
		Status:           models.PaymentSlipStatusConfirm, // Admin payments are auto-confirmed
		PaymentMethod:    models.PaymentMethod(req.PaymentMethod),
		IsAdminCreated:   true,
		CreatedByAdmin:   &adminUserID,
		AdminCreatedNote: req.Note,
		Image:            req.SlipImage, // Optional slip image URL
	}

	// Create payment slip
	err = s.paymentSlipRepo.Create(ctx, paymentSlip)
	if err != nil {
		log.Printf("[ERROR] Failed to create admin payment slip: %v", err)
		return nil, fmt.Errorf("failed to create payment slip: %w", err)
	}

	log.Printf("[SUCCESS] Created admin payment slip %s (ID: %s) for booking %s", slipID, paymentSlip.ID, bookingID)

	// Recalculate booking totals and status
	err = s.recalculateBookingAfterPayment(ctx, bookingID)
	if err != nil {
		log.Printf("[ERROR] Failed to recalculate booking after admin payment: %v", err)
		// Don't fail the payment creation, just log the error
	}

	// Reload fresh payment slip data from database after save
	freshPaymentSlip, err := s.paymentSlipRepo.GetByID(ctx, paymentSlip.ID)
	if err != nil {
		log.Printf("[WARNING] Failed to reload fresh payment slip data: %v", err)
		freshPaymentSlip = paymentSlip // Fallback to original data
	}

	// Get updated booking status
	updatedBooking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		log.Printf("[WARNING] Failed to get updated booking status: %v", err)
	}

	bookingStatus := string(booking.Status)
	if updatedBooking != nil {
		bookingStatus = string(updatedBooking.Status)
	}

	// Generate response message
	message := fmt.Sprintf("Admin payment of %.2f via %s successfully added to booking %s", req.Amount, req.PaymentMethod, booking.InvoiceID)

	response := &dto.AdminPaymentResponse{
		PaymentSlip:   s.toPaymentSlipResponse(freshPaymentSlip),
		BookingStatus: bookingStatus,
		Message:       message,
	}

	log.Printf("[SUCCESS] Admin payment completed. Booking %s status: %s", bookingID, bookingStatus)
	return response, nil
}

// calculateBookingStatusByPayment determines booking status based on payment amounts
// Uses the same logic as Agent payment verification for consistency
func (s *BookingServiceImpl) calculateBookingStatusByPayment(totalVerifiedPayments, sumNet, sumDeposit float64) models.BookingStatus {
	if totalVerifiedPayments >= sumNet {
		// Full payment received - confirm booking (same as Agent logic)
		return models.BookingStatusConfirm
	} else if totalVerifiedPayments >= sumDeposit {
		// Deposit threshold met - move to processing (protected from expiry)
		return models.BookingStatusProcessing
	} else if totalVerifiedPayments > 0 {
		// Some payment received but below deposit threshold
		return models.BookingStatusProcessing
	}
	// No payments - status should remain pending or unchanged
	return models.BookingStatusPending
}

// Helper method to recalculate booking status after admin payment
func (s *BookingServiceImpl) recalculateBookingAfterPayment(ctx context.Context, bookingID uuid.UUID) error {
	log.Printf("[INFO] Recalculating booking %s after admin payment", bookingID)

	// Get booking
	booking, err := s.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}

	// Calculate total verified payments
	totalPayments, err := s.paymentSlipRepo.GetTotalActiveAmountByBookingID(ctx, bookingID)
	if err != nil {
		return fmt.Errorf("failed to calculate total payments: %w", err)
	}

	// Update sum_payment
	oldSumPayment := booking.SumPayment
	booking.SumPayment = totalPayments

	// Determine new status using unified calculation (same as Agent verification)
	oldStatus := booking.Status
	newStatus := s.calculateBookingStatusByPayment(totalPayments, booking.SumNet, booking.SumDeposit)

	// Only update status if it's a valid transition
	if booking.Status == models.BookingStatusPending ||
		booking.Status == models.BookingStatusProcessing ||
		(booking.Status == models.BookingStatusConfirm && newStatus == models.BookingStatusConfirm) {
		booking.Status = newStatus
	}

	statusChanged := oldStatus != booking.Status

	log.Printf("[INFO] Booking %s unified status calculation: payment=%.2f, deposit=%.2f, net=%.2f",
		bookingID, totalPayments, booking.SumDeposit, booking.SumNet)
	log.Printf("[INFO] Booking %s: sum_payment %.2f→%.2f, status %s→%s (changed: %v)",
		bookingID, oldSumPayment, booking.SumPayment, oldStatus, booking.Status, statusChanged)

	// Save booking updates (always save to update sum_payment)
	err = s.bookingRepo.Update(ctx, booking.ID, booking)
	if err != nil {
		return fmt.Errorf("failed to update booking: %w", err)
	}

	log.Printf("[SUCCESS] Booking %s recalculation completed", bookingID)
	return nil
}

// Helper functions for Vision API integration
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// updatePassportWithExtractionResult maps Vision API extraction results to passport model
func (s *BookingServiceImpl) updatePassportWithExtractionResult(passport *models.PassengerPassport, result *dto.PassportExtractionResult) {
	log.Printf("[INFO] Updating passport with Vision API extraction results")

	// Update basic passport information
	if result.PassportNumber != nil && *result.PassportNumber != "" {
		passport.PassportNumber = *result.PassportNumber
		log.Printf("[SUCCESS] Updated passport number: %s", *result.PassportNumber)
	}

	if result.CountryCode != nil {
		passport.CountryCode = result.CountryCode
	}

	if result.Nationality != nil {
		passport.Nationality = result.Nationality
	}

	// Update Thai names
	if result.NameTH != nil {
		passport.NameTH = result.NameTH
	}

	if result.SurnameTH != nil {
		passport.SurnameTH = result.SurnameTH
	}

	if result.TitleTH != nil {
		passport.TitleTH = result.TitleTH
	}

	// Update English names
	if result.NameEN != nil {
		passport.NameEN = result.NameEN
	}

	if result.SurnameEN != nil {
		passport.SurnameEN = result.SurnameEN
	}

	if result.TitleEN != nil {
		passport.TitleEN = result.TitleEN
	}

	// Update personal information
	if result.Gender != nil {
		passport.Gender = result.Gender
	}

	if result.DateOfBirth != nil {
		passport.DateOfBirth = result.DateOfBirth
	}

	if result.DateOfExpiry != nil {
		passport.DateOfExpiry = result.DateOfExpiry
	}

	if result.DateOfIssue != nil {
		passport.DateOfIssue = result.DateOfIssue
	}

	if result.PlaceOfBirth != nil {
		passport.PlaceOfBirth = result.PlaceOfBirth
	}

	if result.IssuingAuthority != nil {
		passport.IssuingAuthority = result.IssuingAuthority
	}

	// Update additional fields
	if result.IDNumber != nil {
		passport.IDNumber = result.IDNumber
	}

	// Store OCR raw text and confidence
	passport.OCRText = &result.RawText
	confidence := result.OverallConfidence
	passport.OCRConfidence = &confidence

	passport.UpdatedAt = time.Now()

	log.Printf("[SUCCESS] Passport data updated with %d extracted fields", len(result.FieldConfidences))
}

// GetBookingDocuments retrieves all documents for a booking (payment slips, passports)
func (s *BookingServiceImpl) GetBookingDocuments(ctx context.Context, bookingID uuid.UUID) (*dto.BookingDocumentsResponse, error) {
	log.Printf("[INFO] Retrieving documents for booking %s", bookingID)

	var documents []dto.BookingDocument

	// Get payment slips
	paymentSlips, err := s.paymentSlipRepo.GetByBookingID(ctx, bookingID)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve payment slips: %v", err)
		return nil, fmt.Errorf("failed to retrieve payment slips: %w", err)
	}

	// Convert payment slips to documents
	for _, slip := range paymentSlips {
		fileSize := int64(0)
		if slip.Image != nil && *slip.Image != "" {
			// Estimate file size or set default
			fileSize = 50000 // Default estimate for image files
		}

		doc := dto.BookingDocument{
			ID:         slip.ID,
			Filename:   fmt.Sprintf("payment_slip_%s.jpg", slip.SlipID),
			Type:       "payment_slip",
			UploadedAt: slip.CreatedAt.Format(time.RFC3339),
			FileURL:    stringValue(slip.Image),
			FileSize:   fileSize,
			Status:     string(slip.Status),
			Amount:     &slip.Amount,
		}
		documents = append(documents, doc)
	}

	// Get passports
	passports, err := s.passengerPassportRepo.GetByBookingID(ctx, bookingID)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve passports: %v", err)
		return nil, fmt.Errorf("failed to retrieve passports: %w", err)
	}

	// Convert passports to documents
	for _, passport := range passports {
		fileSize := int64(0)
		if passport.Image != nil && *passport.Image != "" {
			// Estimate file size or set default
			fileSize = 200000 // Default estimate for passport image files
		}

		filename := fmt.Sprintf("passport_%s.jpg", passport.PassportNumber)
		if passport.NameEN != nil && *passport.NameEN != "" {
			filename = fmt.Sprintf("passport_%s.jpg", *passport.NameEN)
		}

		doc := dto.BookingDocument{
			ID:         passport.ID,
			Filename:   filename,
			Type:       "passport",
			UploadedAt: passport.CreatedAt.Format(time.RFC3339),
			FileURL:    stringValue(passport.Image),
			FileSize:   fileSize,
			Status:     string(passport.Status),
		}
		documents = append(documents, doc)
	}

	// Sort documents by upload date (newest first)
	sort.Slice(documents, func(i, j int) bool {
		timeI, errI := time.Parse(time.RFC3339, documents[i].UploadedAt)
		timeJ, errJ := time.Parse(time.RFC3339, documents[j].UploadedAt)
		if errI != nil || errJ != nil {
			return false
		}
		return timeI.After(timeJ)
	})

	response := &dto.BookingDocumentsResponse{
		BookingID: bookingID,
		Documents: documents,
	}

	log.Printf("[SUCCESS] Retrieved %d documents for booking %s", len(documents), bookingID)
	return response, nil
}

// Helper function to safely get string value from pointer
// GetBookingStatistics calculates summary statistics for filtered bookings
func (s *BookingServiceImpl) GetBookingStatistics(ctx context.Context, req dto.BookingListRequest) (*dto.BookingStatisticsResponse, error) {
	// Get all bookings that match the filter (without pagination)
	statReq := req
	statReq.Offset = 0
	statReq.Limit = 999999 // Set high limit to get all records for statistics

	bookings, err := s.bookingRepo.ListWithFilters(ctx, statReq)
	if err != nil {
		log.Printf("Error getting bookings for statistics: %v", err)
		return nil, fmt.Errorf("failed to get bookings for statistics: %w", err)
	}

	// Calculate statistics
	var totalRevenue, totalDeposit, totalPayments, totalCommission float64
	statusCounts := make(map[string]int64)

	for _, booking := range bookings {
		totalRevenue += booking.SumNet
		totalDeposit += booking.SumDeposit
		totalPayments += booking.GetActualSumPayment()
		totalCommission += booking.SumCommission

		// Count by status
		statusCounts[string(booking.Status)]++
	}

	var avgBookingValue float64
	if len(bookings) > 0 {
		avgBookingValue = totalRevenue / float64(len(bookings))
	}

	return &dto.BookingStatisticsResponse{
		TotalRecords:    int64(len(bookings)),
		TotalRevenue:    totalRevenue,
		TotalDeposit:    totalDeposit,
		TotalPayments:   totalPayments,
		TotalCommission: totalCommission,
		StatusCounts:    statusCounts,
		AvgBookingValue: avgBookingValue,
	}, nil
}

// GetAdminReport returns admin report data for specified month
func (s *BookingServiceImpl) GetAdminReport(ctx context.Context, year int, month int) (*dto.AdminReportResponse, error) {
	log.Printf("📊 Getting admin report for %d-%02d", year, month)

	// Calculate month boundaries
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	// Get all bookings for the month
	monthReq := dto.BookingListRequest{
		StartDate: &startOfMonth,
		EndDate:   &endOfMonth,
		Offset:    0,
		Limit:     999999,
	}

	allBookings, err := s.bookingRepo.ListWithFilters(ctx, monthReq)
	if err != nil {
		log.Printf("Error getting bookings for admin report: %v", err)
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	// Calculate total bookings, seats, sales, and actual revenue
	// ไม่นับ booking ที่มีสถานะ expired และ cancel
	var totalBookings int64
	var totalSeats int64
	var totalSales float64      // ยอดขายทั้งหมด (SumNet)
	var actualRevenue float64   // เงินที่เก็บได้จริง

	for _, booking := range allBookings {
		// ข้าม booking ที่มีสถานะ expired หรือ cancel สำหรับ totals
		if string(booking.Status) != "expired" && string(booking.Status) != "cancel" {
			totalBookings++
			totalSeats += int64(booking.SumQty)
			totalSales += booking.SumNet
		}

		// คำนวณเงินที่เก็บได้จริงจาก PaymentSlips (ทุกสถานะ)
		for _, slip := range booking.PaymentSlips {
			if slip.CanBeIncludedInTotal() { // status = confirm และไม่ถูก reverse
				actualRevenue += slip.Amount
			}
		}
	}

	// Generate daily status report for the month
	dailyReport := make([]*dto.DailyBookingStatusData, 0)

	// Get number of days in month
	daysInMonth := endOfMonth.Day()

	for day := 1; day <= daysInMonth; day++ {
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		dateStr := date.Format("2006-01-02")

		// Initialize counters for this day
		dayData := &dto.DailyBookingStatusData{
			Date:       dateStr,
			Pending:    0,
			Processing: 0,
			Confirm:    0,
			Cancel:     0,
			Completed:  0,
			Expired:    0,
			TotalSeats: 0,
		}

		// Count bookings by status for this day
		for _, booking := range allBookings {
			bookingDate := booking.CreatedAt.Format("2006-01-02")
			if bookingDate == dateStr {
				// Count by status
				switch string(booking.Status) {
				case "pending":
					dayData.Pending++
				case "processing":
					dayData.Processing++
				case "confirm":
					dayData.Confirm++
				case "cancel":
					dayData.Cancel++
				case "completed":
					dayData.Completed++
				case "expired":
					dayData.Expired++
				}

				// Add seats
				dayData.TotalSeats += int64(booking.SumQty)
			}
		}

		dailyReport = append(dailyReport, dayData)
	}

	return &dto.AdminReportResponse{
		TotalBookings:     totalBookings,
		TotalSeats:        totalSeats,
		TotalSales:        totalSales,
		ActualRevenue:     actualRevenue,
		DailyStatusReport: dailyReport,
	}, nil
}

// GetSalesReport returns sales report data for specified month (filtered by sales user)
func (s *BookingServiceImpl) GetSalesReport(ctx context.Context, userID uuid.UUID, year int, month int) (*dto.AdminReportResponse, error) {
	log.Printf("📊 Getting sales report for user %s, %d-%02d", userID.String(), year, month)

	// Find sales record by user_id
	salesRecord, err := s.salesRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error getting sales record for user %s: %v", userID.String(), err)
		return nil, fmt.Errorf("failed to get sales record: %w", err)
	}
	if salesRecord == nil {
		log.Printf("No sales record found for user %s", userID.String())
		return nil, fmt.Errorf("sales record not found for user")
	}

	// Calculate month boundaries
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	// Get bookings for this sales person in the specified month
	monthReq := dto.BookingListRequest{
		StartDate: &startOfMonth,
		EndDate:   &endOfMonth,
		SalesID:   &salesRecord.ID, // Filter by sales_id
		Offset:    0,
		Limit:     999999,
	}

	allBookings, err := s.bookingRepo.ListWithFilters(ctx, monthReq)
	if err != nil {
		log.Printf("Error getting bookings for sales report: %v", err)
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	// Calculate total bookings, seats, sales, and actual revenue
	// ไม่นับ booking ที่มีสถานะ expired และ cancel
	var totalBookings int64
	var totalSeats int64
	var totalSales float64      // ยอดขายทั้งหมด (SumNet)
	var actualRevenue float64   // เงินที่เก็บได้จริง

	for _, booking := range allBookings {
		// ข้าม booking ที่มีสถานะ expired หรือ cancel สำหรับ totals
		if string(booking.Status) != "expired" && string(booking.Status) != "cancel" {
			totalBookings++
			totalSeats += int64(booking.SumQty)
			totalSales += booking.SumNet
		}

		// คำนวณเงินที่เก็บได้จริงจาก PaymentSlips (ทุกสถานะ)
		for _, slip := range booking.PaymentSlips {
			if slip.CanBeIncludedInTotal() { // status = confirm และไม่ถูก reverse
				actualRevenue += slip.Amount
			}
		}
	}

	// Generate daily status report for the month
	dailyReport := make([]*dto.DailyBookingStatusData, 0)

	// Get number of days in month
	daysInMonth := endOfMonth.Day()

	for day := 1; day <= daysInMonth; day++ {
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		dateStr := date.Format("2006-01-02")

		// Initialize counters for this day
		dayData := &dto.DailyBookingStatusData{
			Date:       dateStr,
			Pending:    0,
			Processing: 0,
			Confirm:    0,
			Cancel:     0,
			Completed:  0,
			Expired:    0,
			TotalSeats: 0,
		}

		// Count bookings by status for this day
		for _, booking := range allBookings {
			bookingDate := booking.CreatedAt.Format("2006-01-02")
			if bookingDate == dateStr {
				// Count by status
				switch string(booking.Status) {
				case "pending":
					dayData.Pending++
				case "processing":
					dayData.Processing++
				case "confirm":
					dayData.Confirm++
				case "cancel":
					dayData.Cancel++
				case "completed":
					dayData.Completed++
				case "expired":
					dayData.Expired++
				}

				// Add seats
				dayData.TotalSeats += int64(booking.SumQty)
			}
		}

		dailyReport = append(dailyReport, dayData)
	}

	return &dto.AdminReportResponse{
		TotalBookings:     totalBookings,
		TotalSeats:        totalSeats,
		TotalSales:        totalSales,
		ActualRevenue:     actualRevenue,
		DailyStatusReport: dailyReport,
	}, nil
}

func (s *BookingServiceImpl) GetAgentReport(ctx context.Context, userID uuid.UUID, year int, month int) (*dto.AdminReportResponse, error) {
	log.Printf("📊 Getting agent report for user %s, %d-%02d", userID.String(), year, month)

	// Calculate month boundaries
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	// Get bookings for this agent user in the specified month
	monthReq := dto.BookingListRequest{
		StartDate: &startOfMonth,
		EndDate:   &endOfMonth,
		UserID:    &userID, // Filter by user_id
		Offset:    0,
		Limit:     999999,
	}

	allBookings, err := s.bookingRepo.ListWithFilters(ctx, monthReq)
	if err != nil {
		log.Printf("Error getting bookings for agent report: %v", err)
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	// Calculate total bookings, seats, sales, and actual revenue
	// ไม่นับ booking ที่มีสถานะ expired และ cancel
	var totalBookings int64
	var totalSeats int64
	var totalSales float64      // ยอดขายทั้งหมด (SumNet)
	var actualRevenue float64   // เงินที่เก็บได้จริง

	for _, booking := range allBookings {
		// ข้าม booking ที่มีสถานะ expired หรือ cancel สำหรับ totals
		if string(booking.Status) != "expired" && string(booking.Status) != "cancel" {
			totalBookings++
			totalSeats += int64(booking.SumQty)
			totalSales += booking.SumNet
		}

		// คำนวณเงินที่เก็บได้จริงจาก PaymentSlips (ทุกสถานะ)
		for _, slip := range booking.PaymentSlips {
			if slip.CanBeIncludedInTotal() { // status = confirm และไม่ถูก reverse
				actualRevenue += slip.Amount
			}
		}
	}

	// Generate daily status report for the month
	dailyReport := make([]*dto.DailyBookingStatusData, 0)

	// Get number of days in month
	daysInMonth := endOfMonth.Day()

	for day := 1; day <= daysInMonth; day++ {
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		dateStr := date.Format("2006-01-02")

		// Initialize counters for this day
		dayData := &dto.DailyBookingStatusData{
			Date:       dateStr,
			Pending:    0,
			Processing: 0,
			Confirm:    0,
			Cancel:     0,
			Completed:  0,
			Expired:    0,
			TotalSeats: 0,
		}

		// Count bookings by status for this day
		for _, booking := range allBookings {
			bookingDate := booking.CreatedAt.Format("2006-01-02")
			if bookingDate == dateStr {
				// Count by status
				switch string(booking.Status) {
				case "pending":
					dayData.Pending++
				case "processing":
					dayData.Processing++
				case "confirm":
					dayData.Confirm++
				case "cancel":
					dayData.Cancel++
				case "completed":
					dayData.Completed++
				case "expired":
					dayData.Expired++
				}

				// Add seats
				dayData.TotalSeats += int64(booking.SumQty)
			}
		}

		dailyReport = append(dailyReport, dayData)
	}

	return &dto.AdminReportResponse{
		TotalBookings:     totalBookings,
		TotalSeats:        totalSeats,
		TotalSales:        totalSales,
		ActualRevenue:     actualRevenue,
		DailyStatusReport: dailyReport,
	}, nil
}

func stringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ==================== Email Notification Helpers ====================

// sendBookingEmail ส่ง email แจ้ง agent ตามสถานะ booking (async — ไม่ block main flow)
func (s *BookingServiceImpl) sendBookingEmail(ctx context.Context, booking *models.Booking, emailType string, reason string) {
	if s.emailPort == nil || !s.emailPort.IsConfigured() {
		return
	}

	go func() {
		// ดึง user email
		user, err := s.userRepo.GetByID(ctx, booking.UserID)
		if err != nil || user == nil || user.Email == nil || *user.Email == "" {
			log.Printf("[EMAIL] Skip: no email for user %s", booking.UserID)
			return
		}

		// ดึง tour info
		var tourName, tourCode string
		if tour, err := s.tourRepo.GetByID(ctx, booking.TourID); err == nil && tour != nil {
			tourName = tour.Title
			tourCode = tour.Code
		}

		// ดึง period info
		var periodStart, periodEnd string
		if period, err := s.periodRepo.GetByID(ctx, booking.PeriodID); err == nil && period != nil {
			periodStart = period.StartDate.Format("02/01/2006")
			periodEnd = period.EndDate.Format("02/01/2006")
		}

		data := emailtpl.BookingEmailData{
			InvoiceID:     booking.InvoiceID,
			CustomerName:  booking.CustomerName,
			TourName:      tourName,
			TourCode:      tourCode,
			PeriodStart:   periodStart,
			PeriodEnd:     periodEnd,
			TotalAmount:   fmt.Sprintf("%.2f", booking.SumNet),
			DepositAmount: fmt.Sprintf("%.2f", booking.SumDeposit),
			DateDeposit:   booking.DateDeposit.Format("02/01/2006"),
			DateFinal:     booking.DateFinal.Format("02/01/2006"),
			AgentName:     user.Username,
			Reason:        reason,
		}

		var subject, html string
		switch emailType {
		case "new_booking":
			subject, html = emailtpl.NewBookingEmail(data)
		case "payment_confirmed":
			subject, html = emailtpl.PaymentConfirmedEmail(data)
		case "booking_cancelled":
			subject, html = emailtpl.BookingCancelledEmail(data)
		case "booking_expired":
			subject, html = emailtpl.BookingExpiredEmail(data)
		default:
			return
		}

		err = s.emailPort.SendEmail(ctx, ports.EmailRequest{
			To:      []string{*user.Email},
			Subject: subject,
			HTML:    html,
		})
		if err != nil {
			log.Printf("[EMAIL] Failed to send %s email for booking %s: %v", emailType, booking.InvoiceID, err)
		}
	}()
}
