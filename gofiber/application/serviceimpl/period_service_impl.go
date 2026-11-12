package serviceimpl

import (
	"context"
	"fmt"
	"log"
	"time"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/cdn"
	apperrors "gofiber-template/pkg/errors"

	"github.com/google/uuid"
)

const (
	MaxPeriodFileSize = 10 * 1024 * 1024 // 10MB
)

type PeriodServiceImpl struct {
	periodRepo                  repositories.PeriodRepository
	tourRepo                    repositories.TourRepository
	periodRateItemRepo          repositories.PeriodRateItemRepository
	periodRoomTypeRepo          repositories.PeriodRoomTypeRepository
	periodAdditionalServiceRepo repositories.PeriodAdditionalServiceRepository
	bookingRepo                 repositories.BookingRepository
	seatHoldRepo                repositories.SeatHoldRepository
}

func NewPeriodService(
	periodRepo repositories.PeriodRepository,
	tourRepo repositories.TourRepository,
	periodRateItemRepo repositories.PeriodRateItemRepository,
	periodRoomTypeRepo repositories.PeriodRoomTypeRepository,
	periodAdditionalServiceRepo repositories.PeriodAdditionalServiceRepository,
	bookingRepo repositories.BookingRepository,
	seatHoldRepo repositories.SeatHoldRepository,
) services.PeriodService {
	return &PeriodServiceImpl{
		periodRepo:                  periodRepo,
		tourRepo:                    tourRepo,
		periodRateItemRepo:          periodRateItemRepo,
		periodRoomTypeRepo:          periodRoomTypeRepo,
		periodAdditionalServiceRepo: periodAdditionalServiceRepo,
		bookingRepo:                 bookingRepo,
		seatHoldRepo:                seatHoldRepo,
	}
}

func (s *PeriodServiceImpl) CreatePeriod(ctx context.Context, req dto.CreatePeriodRequest) (*dto.PeriodResponse, error) {
	// Business validation - check if tour exists
	tour, err := s.tourRepo.GetByID(ctx, req.TourID)
	if err != nil {
		log.Printf("Error retrieving tour: %v", err)
		return nil, fmt.Errorf("failed to validate tour: %w", err)
	}
	if tour == nil {
		return nil, fmt.Errorf("tour not found")
	}

	// Business validation - validate dates
	if err := s.ValidatePeriodDates(ctx, req.StartDate, req.EndDate); err != nil {
		return nil, err
	}

	// Business validation - validate rate type
	if req.RateType != "regular" && req.RateType != "sale" {
		return nil, apperrors.ErrInvalidRateType
	}

	// Convert DTOs to domain models
	sector1 := s.convertFlightSectors(req.Sector1)
	sector2 := s.convertFlightSectors(req.Sector2)
	commissionItems := s.convertCommissionItems(req.CommissionItems)

	// Convert to new separate table models
	rateItems := s.convertToPeriodRateItems(req.RateItems)
	roomTypes := s.convertToPeriodRoomTypes(req.RoomTypes)
	additionalServices := s.convertToPeriodAdditionalServices(req.AdditionalServices)

	// Determine status - use request status or default to draft
	status := models.PeriodStatusDraft // Default
	if req.Status != nil {
		switch *req.Status {
		case "publish":
			status = models.PeriodStatusPublish
		case "draft":
			status = models.PeriodStatusDraft
		default:
			return nil, apperrors.ErrInvalidPeriodStatus
		}
	}

	// Handle expiration override - ทัวร์ไฟไหม้
	expirationOverride := false
	if req.ExpirationOverride != nil {
		expirationOverride = *req.ExpirationOverride
	}

	// Create domain model (without JSONB arrays)
	period := &models.Period{
		TourID:      req.TourID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Seat:        req.Seat,
		Deposit:     req.Deposit,
		DepositType: req.DepositType,
		//SingleStay:      req.SingleStay,
		Airline:         req.Airline,
		Airport:         req.Airport,
		Note:            req.Note,
		Sector1:         sector1,
		Sector2:         sector2,
		RateType:        models.RateType(req.RateType),
		CommissionType:  models.RateType(req.CommissionType),
		CommissionItems: commissionItems,
		ServiceFeePerPax: req.ServiceFeePerPax,
		Status:           status, // Use requested status or default
		CreatedAt:        time.Now(),
		UpdatedAt:       time.Now(),
		// JSONB arrays will be empty - data goes to separate tables
		RateItems:          models.RateItems{},
		RoomTypes:          models.RoomTypes{},
		AdditionalServices: models.AdditionalServices{},
		// Expiration Override - สำหรับทัวร์ไฟไหม้
		ExpirationOverride: expirationOverride,
		ExpirationMinutes:  req.ExpirationMinutes,
		ExpirationHours:    req.ExpirationHours,
		ExpirationDays:     req.ExpirationDays,
		ExpirationMessage:  req.ExpirationMessage,
	}

	// Create period and related items in transaction
	err = s.periodRepo.Create(ctx, period)
	if err != nil {
		log.Printf("Error creating period: %v", err)
		return nil, fmt.Errorf("failed to create period: %w", err)
	}

	// Set period IDs for related items
	for _, item := range rateItems {
		item.PeriodID = period.ID
	}
	for _, item := range roomTypes {
		item.PeriodID = period.ID
	}
	for _, item := range additionalServices {
		item.PeriodID = period.ID
	}

	// Create related items in separate tables
	if len(rateItems) > 0 {
		if err := s.periodRateItemRepo.CreateBatch(ctx, rateItems); err != nil {
			log.Printf("Error creating rate items: %v", err)
			return nil, fmt.Errorf("failed to create rate items: %w", err)
		}
	}

	if len(roomTypes) > 0 {
		if err := s.periodRoomTypeRepo.CreateBatch(ctx, roomTypes); err != nil {
			log.Printf("Error creating room types: %v", err)
			return nil, fmt.Errorf("failed to create room types: %w", err)
		}
	}

	if len(additionalServices) > 0 {
		if err := s.periodAdditionalServiceRepo.CreateBatch(ctx, additionalServices); err != nil {
			log.Printf("Error creating additional services: %v", err)
			return nil, fmt.Errorf("failed to create additional services: %w", err)
		}
	}

	return s.toPeriodResponseWithSeparateTables(ctx, period), nil
}

func (s *PeriodServiceImpl) GetPeriodByID(ctx context.Context, periodID uuid.UUID) (*dto.PeriodResponse, error) {
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period: %v", err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return nil, apperrors.ErrPeriodNotFound
	}

	return s.toPeriodResponseWithSeparateTables(ctx, period), nil
}

func (s *PeriodServiceImpl) GetPeriodWithTour(ctx context.Context, periodID uuid.UUID) (*dto.PeriodResponse, error) {
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period: %v", err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return nil, apperrors.ErrPeriodNotFound
	}

	response := s.toPeriodResponseWithSeparateTables(ctx, period)

	// Get tour details
	if tour, err := s.tourRepo.GetByID(ctx, period.TourID); err == nil && tour != nil {
		response.TourDetails = &dto.TourResponse{
			ID:        tour.ID,
			Title:     tour.Title,
			Code:      tour.Code,
			Status:    string(tour.Status),
			CreatedAt: tour.CreatedAt.Format(time.RFC3339),
			UpdatedAt: tour.UpdatedAt.Format(time.RFC3339),
		}
	}

	return response, nil
}

func (s *PeriodServiceImpl) UpdatePeriod(ctx context.Context, periodID uuid.UUID, req dto.UpdatePeriodRequest) (*dto.PeriodResponse, error) {
	// Get existing period
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period: %v", err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return nil, apperrors.ErrPeriodNotFound
	}

	// Business validation - validate dates if provided
	if req.StartDate != nil && req.EndDate != nil {
		if err := s.ValidatePeriodDates(ctx, *req.StartDate, *req.EndDate); err != nil {
			return nil, err
		}
	}

	// Business validation - validate status if provided
	if req.Status != nil {
		status := models.PeriodStatus(*req.Status)
		if status != models.PeriodStatusPublish && status != models.PeriodStatusDraft {
			return nil, apperrors.ErrInvalidPeriodStatus
		}
	}

	// Update fields if provided
	if req.StartDate != nil {
		period.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		period.EndDate = *req.EndDate
	}
	if req.Seat != nil {
		period.Seat = *req.Seat
	}
	if req.Deposit != nil {
		period.Deposit = *req.Deposit
	}
	if req.DepositType != nil {
		period.DepositType = *req.DepositType
	}
	/*
		if req.SingleStay != nil {
			period.SingleStay = *req.SingleStay
		}
	*/
	if req.Airline != nil {
		period.Airline = *req.Airline
	}
	if req.Airport != nil {
		period.Airport = *req.Airport
	}
	if req.Note != nil {
		period.Note = req.Note
	}
	if req.Sector1 != nil {
		period.Sector1 = s.convertFlightSectors(req.Sector1)
	}
	if req.Sector2 != nil {
		period.Sector2 = s.convertFlightSectors(req.Sector2)
	}
	if req.RateType != nil {
		period.RateType = models.RateType(*req.RateType)
	}
	// Commission items are still stored in JSONB (for now)
	if req.CommissionType != nil {
		period.CommissionType = models.RateType(*req.CommissionType)
	}
	if req.CommissionItems != nil {
		period.CommissionItems = s.convertCommissionItems(req.CommissionItems)
	}
	if req.Status != nil {
		period.Status = models.PeriodStatus(*req.Status)
	}
	if req.Files != nil {
		// Convert FileAttachmentRequest to FileAttachment model
		files := make(models.FileAttachments, len(req.Files))
		for i, file := range req.Files {
			files[i] = models.FileAttachment{
				Name: file.Name,
				Type: file.Type,
				URL:  file.URL,
			}
		}
		period.Files = files
	}

	// Handle expiration override fields - ทัวร์ไฟไหม้
	if req.ExpirationOverride != nil {
		period.ExpirationOverride = *req.ExpirationOverride
	}
	if req.ExpirationMinutes != nil {
		period.ExpirationMinutes = req.ExpirationMinutes
	}
	if req.ExpirationHours != nil {
		period.ExpirationHours = req.ExpirationHours
	}
	if req.ExpirationDays != nil {
		period.ExpirationDays = req.ExpirationDays
	}
	if req.ExpirationMessage != nil {
		period.ExpirationMessage = req.ExpirationMessage
	}
	if req.ServiceFeePerPax != nil {
		period.ServiceFeePerPax = *req.ServiceFeePerPax
	}

	period.UpdatedAt = time.Now()

	// Update period
	if err := s.periodRepo.Update(ctx, periodID, period); err != nil {
		log.Printf("Error updating period: %v", err)
		return nil, fmt.Errorf("failed to update period: %w", err)
	}

	// Update separate tables if provided
	if req.RateItems != nil {
		log.Printf("UpdatePeriod: Converting %d rate items", len(req.RateItems))
		for i, item := range req.RateItems {
			log.Printf("UpdatePeriod: RateItem[%d] - ID: %s, Title: %s, Regular: %.2f, Sale: %.2f",
				i, item.ID, item.Title, item.Regular, item.Sale)
		}
		rateItems := s.convertToPeriodRateItems(req.RateItems)
		log.Printf("UpdatePeriod: Converted to %d models", len(rateItems))
		for i, item := range rateItems {
			log.Printf("UpdatePeriod: Model[%d] - ID: %s, Title: %s, Regular: %.2f, Sale: %.2f",
				i, item.ID.String(), item.Title, item.Regular, item.Sale)
		}
		if err := s.periodRateItemRepo.UpdateBatch(ctx, periodID, rateItems); err != nil {
			log.Printf("Error updating rate items: %v", err)
			return nil, fmt.Errorf("failed to update rate items: %w", err)
		}
		log.Printf("UpdatePeriod: Rate items updated successfully")
	}

	if req.RoomTypes != nil {
		log.Printf("UpdatePeriod: Converting %d room types", len(req.RoomTypes))
		roomTypes := s.convertToPeriodRoomTypes(req.RoomTypes)
		if err := s.periodRoomTypeRepo.UpdateBatch(ctx, periodID, roomTypes); err != nil {
			log.Printf("Error updating room types: %v", err)
			return nil, fmt.Errorf("failed to update room types: %w", err)
		}
		log.Printf("UpdatePeriod: Room types updated successfully")
	}

	if req.AdditionalServices != nil {
		log.Printf("UpdatePeriod: Converting %d additional services", len(req.AdditionalServices))
		additionalServices := s.convertToPeriodAdditionalServices(req.AdditionalServices)
		if err := s.periodAdditionalServiceRepo.UpdateBatch(ctx, periodID, additionalServices); err != nil {
			log.Printf("Error updating additional services: %v", err)
			return nil, fmt.Errorf("failed to update additional services: %w", err)
		}
		log.Printf("UpdatePeriod: Additional services updated successfully")
	}

	return s.toPeriodResponseWithSeparateTables(ctx, period), nil
}

func (s *PeriodServiceImpl) DeletePeriod(ctx context.Context, periodID uuid.UUID) error {
	// Check if period exists
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period: %v", err)
		return fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return apperrors.ErrPeriodNotFound
	}

	// Business rule - check if period has active bookings (real-time calculation)
	activeStatuses := []string{"pending", "processing", "confirm"}
	bookings, err := s.bookingRepo.FindByPeriodIDAndStatuses(periodID, activeStatuses)
	if err != nil {
		log.Printf("Error checking active bookings for period %s: %v", periodID, err)
		return fmt.Errorf("failed to check active bookings: %w", err)
	}

	activeBookingCount := 0
	for _, booking := range bookings {
		activeBookingCount += booking.SumQty
	}

	if activeBookingCount > 0 {
		return apperrors.ErrPeriodHasBookings
	}

	// Delete period
	if err := s.periodRepo.Delete(ctx, periodID); err != nil {
		log.Printf("Error deleting period: %v", err)
		return fmt.Errorf("failed to delete period: %w", err)
	}

	return nil
}

func (s *PeriodServiceImpl) ListPeriods(ctx context.Context, req dto.PeriodListRequest) (*dto.PeriodListResponse, error) {
	// Use new ListWithFilters method that supports multiple filters at once
	periods, total, err := s.periodRepo.ListWithFilters(ctx, req)
	if err != nil {
		log.Printf("Error listing periods with filters: %v", err)
		return nil, fmt.Errorf("failed to list periods: %w", err)
	}

	// Convert to response DTOs
	responsePeriods := make([]*dto.PeriodResponse, len(periods))
	for i, period := range periods {
		responsePeriods[i] = s.toPeriodResponseWithSeparateTables(ctx, period)
	}

	pagination := dto.PaginationMeta{
		Total:       total,
		HasMore:     int64(req.Offset+len(periods)) < total,
		CurrentPage: (req.Offset / req.Limit) + 1,
		TotalPages:  int((total + int64(req.Limit) - 1) / int64(req.Limit)),
		PageSize:    req.Limit,
		Offset:      req.Offset,
	}

	return &dto.PeriodListResponse{
		Periods:    responsePeriods,
		Pagination: pagination,
	}, nil
}

func (s *PeriodServiceImpl) GetPeriodsByTour(ctx context.Context, tourID uuid.UUID, offset, limit int) ([]*dto.PeriodResponse, error) {
	periods, err := s.periodRepo.GetByTour(ctx, tourID, offset, limit)
	if err != nil {
		log.Printf("Error getting periods by tour: %v", err)
		return nil, fmt.Errorf("failed to get periods by tour: %w", err)
	}

	// Convert to response DTOs
	responsePeriods := make([]*dto.PeriodResponse, len(periods))
	for i, period := range periods {
		responsePeriods[i] = s.toPeriodResponseWithSeparateTables(ctx, period)
	}

	return responsePeriods, nil
}

func (s *PeriodServiceImpl) GetAvailablePeriods(ctx context.Context, offset, limit int) ([]*dto.PeriodResponse, error) {
	periods, err := s.periodRepo.GetAvailable(ctx, offset, limit)
	if err != nil {
		log.Printf("Error getting available periods: %v", err)
		return nil, fmt.Errorf("failed to get available periods: %w", err)
	}

	// Convert to response DTOs
	responsePeriods := make([]*dto.PeriodResponse, len(periods))
	for i, period := range periods {
		responsePeriods[i] = s.toPeriodResponseWithSeparateTables(ctx, period)
	}

	return responsePeriods, nil
}

func (s *PeriodServiceImpl) GetUpcomingPeriods(ctx context.Context, days int, offset, limit int) ([]*dto.PeriodResponse, error) {
	periods, err := s.periodRepo.GetUpcoming(ctx, offset, limit)
	if err != nil {
		log.Printf("Error getting upcoming periods: %v", err)
		return nil, fmt.Errorf("failed to get upcoming periods: %w", err)
	}

	// Filter by days if needed (simplified approach)
	filteredPeriods := make([]*models.Period, 0)
	cutoffDate := time.Now().AddDate(0, 0, days)

	for _, period := range periods {
		if period.StartDate.Before(cutoffDate) {
			filteredPeriods = append(filteredPeriods, period)
		}
	}

	// Convert to response DTOs
	responsePeriods := make([]*dto.PeriodResponse, len(filteredPeriods))
	for i, period := range filteredPeriods {
		responsePeriods[i] = s.toPeriodResponseWithSeparateTables(ctx, period)
	}

	return responsePeriods, nil
}

// Availability Management
func (s *PeriodServiceImpl) CheckAvailability(ctx context.Context, periodID uuid.UUID, requestedSeats int) (*dto.AvailabilityCheckResponse, error) {
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error checking period availability: %v", err)
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}

	if period == nil {
		return nil, apperrors.ErrPeriodNotFound
	}

	// Calculate real-time available seats
	realTimeStats, err := s.CalculateRealTimeBookingStats(ctx, periodID)
	var availableSeats int
	if err != nil {
		log.Printf("Warning: Failed to calculate real-time stats for period %s: %v", periodID, err)
		availableSeats = period.Seat // Fallback: assume all seats available
	} else {
		availableSeats = realTimeStats.AvailableSeats
	}

	available := availableSeats >= requestedSeats && period.Status == models.PeriodStatusPublish

	return &dto.AvailabilityCheckResponse{
		Available:      available,
		RequestedSeats: requestedSeats,
		AvailableSeats: availableSeats,
		TotalSeats:     period.Seat,
		BookedSeats:    period.Seat - availableSeats, // Calculate from available seats
	}, nil
}

// Note: UpdateBookingCount function removed - no longer needed with real-time calculation

func (s *PeriodServiceImpl) ReserveSeats(ctx context.Context, periodID uuid.UUID, seats int) error {
	// Check availability first
	available, err := s.CheckAvailability(ctx, periodID, seats)
	if err != nil {
		return err
	}

	if !available.Available {
		return apperrors.ErrInsufficientSeats
	}

	// Seat reservation is now handled by booking creation in database
	// No need to update period booking counter - using real-time calculation
	return nil
}

func (s *PeriodServiceImpl) ReleaseSeats(ctx context.Context, periodID uuid.UUID, seats int) error {
	// Seat release is now handled by booking status updates in database
	// No need to update period booking counter - using real-time calculation
	return nil
}

// File Management
func (s *PeriodServiceImpl) UploadPeriodFile(ctx context.Context, periodID uuid.UUID, fileData []byte, fileName, fileType string) (*dto.PeriodResponse, error) {
	// Business validation - check file size
	if len(fileData) > MaxPeriodFileSize {
		return nil, apperrors.ErrInvalidPeriodFileType
	}

	// Get existing period
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period: %v", err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return nil, apperrors.ErrPeriodNotFound
	}

	// Generate file URL (in real implementation, upload to storage service)
	fileURL := s.generateFileURL(periodID, fileName, fileType)

	// Add file to period's file attachments
	fileAttachment := models.FileAttachment{
		Name: fileName,
		Type: fileType,
		URL:  fileURL,
	}

	period.Files = append(period.Files, fileAttachment)
	period.UpdatedAt = time.Now()

	// Update period
	if err := s.periodRepo.Update(ctx, periodID, period); err != nil {
		log.Printf("Error updating period with file: %v", err)
		return nil, fmt.Errorf("failed to update period with file: %w", err)
	}

	return s.toPeriodResponseWithSeparateTables(ctx, period), nil
}

func (s *PeriodServiceImpl) DeletePeriodFile(ctx context.Context, periodID uuid.UUID, fileURL string) (*dto.PeriodResponse, error) {
	// Get existing period
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period: %v", err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return nil, apperrors.ErrPeriodNotFound
	}

	// Remove file from attachments
	updatedFiles := make(models.FileAttachments, 0)
	for _, file := range period.Files {
		if file.URL != fileURL {
			updatedFiles = append(updatedFiles, file)
		}
	}

	period.Files = updatedFiles
	period.UpdatedAt = time.Now()

	// Update period
	if err := s.periodRepo.Update(ctx, periodID, period); err != nil {
		log.Printf("Error removing period file: %v", err)
		return nil, fmt.Errorf("failed to remove period file: %w", err)
	}

	return s.toPeriodResponseWithSeparateTables(ctx, period), nil
}

// Business Logic
func (s *PeriodServiceImpl) PublishPeriod(ctx context.Context, periodID uuid.UUID) (*dto.PeriodResponse, error) {
	// Get existing period
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period: %v", err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return nil, apperrors.ErrPeriodNotFound
	}

	// Business validation - check if period is in future
	if period.StartDate.Before(time.Now()) {
		return nil, apperrors.ErrPeriodInPast
	}

	// Update status to published
	if err := s.periodRepo.UpdateStatus(ctx, periodID, models.PeriodStatusPublish); err != nil {
		log.Printf("Error publishing period: %v", err)
		return nil, fmt.Errorf("failed to publish period: %w", err)
	}

	period.Status = models.PeriodStatusPublish
	period.UpdatedAt = time.Now()

	return s.toPeriodResponseWithSeparateTables(ctx, period), nil
}

func (s *PeriodServiceImpl) UnpublishPeriod(ctx context.Context, periodID uuid.UUID) (*dto.PeriodResponse, error) {
	// Get existing period
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period: %v", err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return nil, apperrors.ErrPeriodNotFound
	}

	// Update status to draft
	if err := s.periodRepo.UpdateStatus(ctx, periodID, models.PeriodStatusDraft); err != nil {
		log.Printf("Error unpublishing period: %v", err)
		return nil, fmt.Errorf("failed to unpublish period: %w", err)
	}

	period.Status = models.PeriodStatusDraft
	period.UpdatedAt = time.Now()

	return s.toPeriodResponseWithSeparateTables(ctx, period), nil
}

func (s *PeriodServiceImpl) GetPeriodStats(ctx context.Context) (*dto.PeriodStatsResponse, error) {
	// Get total periods count
	totalPeriods, err := s.periodRepo.Count(ctx)
	if err != nil {
		log.Printf("Error getting period stats: %v", err)
		return nil, fmt.Errorf("failed to get period stats: %w", err)
	}

	// Simplified stats - in real implementation, use aggregate queries
	availablePeriods, err := s.periodRepo.GetAvailable(ctx, 0, 1000)
	if err != nil {
		log.Printf("Error getting available periods for stats: %v", err)
		return nil, fmt.Errorf("failed to get period stats: %w", err)
	}

	availableCount := int64(len(availablePeriods))
	publishedCount := availableCount // Simplified
	fullyBookedCount := int64(0)     // TODO: Implement
	totalSeats := int64(0)           // TODO: Implement
	bookedSeats := int64(0)          // TODO: Implement
	averageOccupancy := 0.0          // TODO: Implement

	return &dto.PeriodStatsResponse{
		TotalPeriods:       totalPeriods,
		PublishedPeriods:   publishedCount,
		AvailablePeriods:   availableCount,
		FullyBookedPeriods: fullyBookedCount,
		TotalSeats:         totalSeats,
		BookedSeats:        bookedSeats,
		AverageOccupancy:   averageOccupancy,
	}, nil
}

func (s *PeriodServiceImpl) CalculateRealTimeBookingStats(ctx context.Context, periodID uuid.UUID) (*dto.PeriodBookingStatsResponse, error) {
	log.Printf("[DEBUG] CalculateRealTimeBookingStats started for period ID: %s", periodID.String())

	// Validate input
	if periodID == uuid.Nil {
		log.Printf("[ERROR] Invalid period ID provided: %s", periodID.String())
		return nil, fmt.Errorf("invalid period ID: cannot be nil")
	}

	// Get period details
	log.Printf("[DEBUG] Fetching period details for ID: %s", periodID.String())
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("[ERROR] Failed to get period %s: %v", periodID.String(), err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}
	if period == nil {
		log.Printf("[ERROR] Period not found for ID: %s", periodID.String())
		return nil, fmt.Errorf("period not found")
	}

	log.Printf("[DEBUG] Period found - ID: %s, Seat: %d, Status: %s",
		period.ID.String(), period.Seat, string(period.Status))

	// Validate booking repository
	if s.bookingRepo == nil {
		log.Printf("[ERROR] BookingRepository is nil - dependency injection issue")
		return nil, fmt.Errorf("booking repository not initialized")
	}

	// Get all active bookings for this period (exclude cancelled)
	activeStatuses := []string{"pending", "processing", "confirm"}
	log.Printf("[DEBUG] Checking active booking statuses: %v", activeStatuses)
	log.Printf("[DEBUG] FIXED: Now including 'confirm' status in real-time calculations")

	// Count total bookings by status
	bookingsByStatus := make(map[string]int)
	totalBookings := 0
	bookedPassengers := 0
	queryErrors := []error{}

	// Query bookings by status
	for _, status := range activeStatuses {
		log.Printf("[DEBUG] Querying bookings for period %s with status: %s", periodID.String(), status)

		count, passengers, err := s.countBookingsByPeriodAndStatus(ctx, periodID, status)
		if err != nil {
			log.Printf("[ERROR] Failed to count bookings for status %s: %v", status, err)
			queryErrors = append(queryErrors, fmt.Errorf("status %s: %w", status, err))
			continue
		}

		log.Printf("[DEBUG] Status %s results - Count: %d, Passengers: %d", status, count, passengers)

		bookingsByStatus[status] = count
		totalBookings += count
		bookedPassengers += passengers
	}

	// Check if all queries failed
	if len(queryErrors) == len(activeStatuses) {
		log.Printf("[ERROR] All booking status queries failed for period %s", periodID.String())
		return nil, fmt.Errorf("failed to query any booking statuses: %v", queryErrors)
	}

	// Log partial errors but continue
	if len(queryErrors) > 0 {
		log.Printf("[WARNING] Some booking queries failed for period %s: %v", periodID.String(), queryErrors)
	}

	// Get held seats if seat hold repository is available
	heldSeats := 0
	if s.seatHoldRepo != nil {
		activeHeldSeats, err := s.seatHoldRepo.GetActiveHeldSeats(ctx, periodID)
		if err != nil {
			log.Printf("[WARNING] Failed to get held seats for period %s: %v", periodID.String(), err)
			// Continue with calculation but log the warning
		} else {
			heldSeats = activeHeldSeats
			log.Printf("[DEBUG] Active held seats for period %s: %d", periodID.String(), heldSeats)
		}
	} else {
		log.Printf("[DEBUG] SeatHoldRepository not available - held seats calculation skipped")
	}

	// Calculate available seats (accounting for both bookings and holds)
	availableSeats := period.Seat - bookedPassengers - heldSeats
	isAvailable := availableSeats > 0 && period.Status == models.PeriodStatusPublish

	log.Printf("[DEBUG] Enhanced calculation results for period %s:", periodID.String())
	log.Printf("[DEBUG] - Total Seats: %d", period.Seat)
	log.Printf("[DEBUG] - Total Bookings: %d", totalBookings)
	log.Printf("[DEBUG] - Booked Passengers: %d", bookedPassengers)
	log.Printf("[DEBUG] - Held Seats: %d", heldSeats)
	log.Printf("[DEBUG] - Available Seats: %d (Total - Booked - Held)", availableSeats)
	log.Printf("[DEBUG] - Is Available: %t", isAvailable)
	log.Printf("[DEBUG] - Period Status: %s", string(period.Status))
	log.Printf("[DEBUG] - Bookings by Status: %+v", bookingsByStatus)

	response := &dto.PeriodBookingStatsResponse{
		PeriodID:         periodID,
		TotalSeats:       period.Seat,
		TotalBookings:    totalBookings,
		BookedPassengers: bookedPassengers,
		HeldSeats:        heldSeats,
		AvailableSeats:   availableSeats,
		BookingsByStatus: bookingsByStatus,
		IsAvailable:      isAvailable,
	}

	log.Printf("[DEBUG] CalculateRealTimeBookingStats completed successfully for period %s", periodID.String())
	return response, nil
}

// Helper method to count bookings and passengers by period and status
func (s *PeriodServiceImpl) countBookingsByPeriodAndStatus(ctx context.Context, periodID uuid.UUID, status string) (int, int, error) {
	log.Printf("[DEBUG] Helper countBookingsByPeriodAndStatus called - Period: %s, Status: %s", periodID.String(), status)

	// Validate inputs at service level
	if periodID == uuid.Nil {
		log.Printf("[ERROR] Service level validation failed - nil period ID")
		return 0, 0, fmt.Errorf("invalid period ID at service level")
	}

	if status == "" {
		log.Printf("[ERROR] Service level validation failed - empty status")
		return 0, 0, fmt.Errorf("invalid status at service level")
	}

	// Validate repository dependency
	if s.bookingRepo == nil {
		log.Printf("[ERROR] BookingRepository dependency is nil")
		return 0, 0, fmt.Errorf("booking repository not available")
	}

	// Call repository method with error wrapping
	count, passengers, err := s.bookingRepo.CountBookingsByPeriodAndStatus(ctx, periodID, status)
	if err != nil {
		log.Printf("[ERROR] Repository call failed for period %s, status %s: %v", periodID.String(), status, err)
		return 0, 0, fmt.Errorf("repository error for period %s status %s: %w", periodID.String(), status, err)
	}

	// Validate results
	if count < 0 {
		log.Printf("[ERROR] Invalid booking count returned: %d", count)
		return 0, 0, fmt.Errorf("invalid booking count: %d", count)
	}

	if passengers < 0 {
		log.Printf("[ERROR] Invalid passenger count returned: %d", passengers)
		return 0, 0, fmt.Errorf("invalid passenger count: %d", passengers)
	}

	log.Printf("[SUCCESS] Helper method completed - Period: %s, Status: %s, Count: %d, Passengers: %d",
		periodID.String(), status, count, passengers)

	return count, passengers, nil
}

func (s *PeriodServiceImpl) ValidatePeriodDates(ctx context.Context, startDate, endDate time.Time) error {
	// Business rule - start date must be before end date
	if !startDate.Before(endDate) {
		return apperrors.ErrInvalidPeriodDates
	}

	// Business rule - start date must be in the future
	if startDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return apperrors.ErrPeriodInPast
	}

	return nil
}

func (s *PeriodServiceImpl) CalculateCommission(ctx context.Context, periodID uuid.UUID, rateType string) (float64, error) {
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period for commission calculation: %v", err)
		return 0, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return 0, apperrors.ErrPeriodNotFound
	}

	// Calculate commission based on rate type
	var totalCommission float64
	for _, commission := range period.CommissionItems {
		switch models.RateType(rateType) {
		case models.RateTypeRegular:
			totalCommission += commission.Regular
		case models.RateTypeSale:
			totalCommission += commission.Sale
		default:
			return 0, apperrors.ErrInvalidRateType
		}
	}

	return totalCommission, nil
}

// CalculateCommissionForBooking calculates commission for booking with exclude field support
func (s *PeriodServiceImpl) CalculateCommissionForBooking(ctx context.Context, periodID uuid.UUID, rateType string, passengers []models.RateItem) (float64, error) {
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period for commission calculation: %v", err)
		return 0, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return 0, apperrors.ErrPeriodNotFound
	}

	// Check if any passengers have exclude = true
	hasExcludedPassengers := false
	for _, passenger := range passengers {
		if passenger.Exclude {
			hasExcludedPassengers = true
			break
		}
	}

	// If exclude passengers exist, commission might be 0 (business rule flexibility)
	// This can be configured based on business requirements
	if hasExcludedPassengers {
		// Option 1: No commission for bookings with excluded passengers
		// return 0, nil

		// Option 2: Partial commission (calculate only for non-excluded)
		// Could be implemented based on business rules
	}

	// Calculate total commission from all commission items
	var totalCommission float64
	for _, commission := range period.CommissionItems {
		switch models.RateType(rateType) {
		case models.RateTypeRegular:
			totalCommission += commission.Regular
		case models.RateTypeSale:
			totalCommission += commission.Sale
		default:
			return 0, apperrors.ErrInvalidRateType
		}
	}

	return totalCommission, nil
}

// validateCommissionStructure validates commission items structure
func (s *PeriodServiceImpl) validateCommissionStructure(commissionItems models.CommissionItems) error {
	if len(commissionItems) == 0 {
		return fmt.Errorf("commission items cannot be empty")
	}

	// Support multiple commission items like old system (primary + bonus)
	// Example: [{"regular": 800, "sale": 800}, {"regular": 200, "sale": 500}]
	for i, item := range commissionItems {
		if item.Regular < 0 {
			return fmt.Errorf("commission item %d regular amount cannot be negative", i+1)
		}
		if item.Sale < 0 {
			return fmt.Errorf("commission item %d sale amount cannot be negative", i+1)
		}
	}

	// Business rule: At least one commission item should have positive value
	hasPositiveCommission := false
	for _, item := range commissionItems {
		if item.Regular > 0 || item.Sale > 0 {
			hasPositiveCommission = true
			break
		}
	}

	if !hasPositiveCommission {
		return fmt.Errorf("at least one commission item must have positive value")
	}

	return nil
}

// UpdatePeriodStatus updates period status
func (s *PeriodServiceImpl) UpdatePeriodStatus(ctx context.Context, req dto.UpdatePeriodStatusRequest) (*dto.PeriodResponse, error) {
	// Get existing period
	period, err := s.periodRepo.GetByID(ctx, req.ID)
	if err != nil {
		log.Printf("Error retrieving period: %v", err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return nil, apperrors.ErrPeriodNotFound
	}

	// Business validation - check if status change is needed
	currentStatus := string(period.Status)
	if currentStatus == req.Status {
		return nil, fmt.Errorf("period is already in %s status", req.Status)
	}

	// Validate status value
	var newStatus models.PeriodStatus
	switch req.Status {
	case "publish":
		newStatus = models.PeriodStatusPublish
	case "draft":
		newStatus = models.PeriodStatusDraft
	default:
		return nil, fmt.Errorf("invalid period status: %s", req.Status)
	}

	// Additional business validation for publishing
	if newStatus == models.PeriodStatusPublish {
		// Check if period has required data
		if len(period.RateItems) == 0 {
			return nil, fmt.Errorf("cannot publish period without rate items")
		}
		if len(period.CommissionItems) == 0 {
			return nil, fmt.Errorf("cannot publish period without commission items")
		}

		// Validate commission structure (support 1-N items like old system)
		if err := s.validateCommissionStructure(period.CommissionItems); err != nil {
			return nil, fmt.Errorf("invalid commission structure: %w", err)
		}

		// Check if dates are valid
		if err := s.ValidatePeriodDates(ctx, period.StartDate, period.EndDate); err != nil {
			return nil, fmt.Errorf("cannot publish period with invalid dates: %w", err)
		}
	}

	// Update status
	if err := s.periodRepo.UpdateStatus(ctx, req.ID, newStatus); err != nil {
		log.Printf("Error updating period status: %v", err)
		return nil, fmt.Errorf("failed to update period status: %w", err)
	}

	period.Status = newStatus
	period.UpdatedAt = time.Now()

	return s.toPeriodResponseWithSeparateTables(ctx, period), nil
}

// GetPeriodAvailability returns detailed availability information for a period
func (s *PeriodServiceImpl) GetPeriodAvailability(ctx context.Context, periodID uuid.UUID) (*dto.PeriodAvailabilityResponse, error) {
	// Get period
	period, err := s.periodRepo.GetByID(ctx, periodID)
	if err != nil {
		log.Printf("Error retrieving period: %v", err)
		return nil, fmt.Errorf("failed to get period: %w", err)
	}

	if period == nil {
		return nil, apperrors.ErrPeriodNotFound
	}

	// Calculate real-time availability
	realTimeStats, err := s.CalculateRealTimeBookingStats(ctx, periodID)
	if err != nil {
		log.Printf("Error calculating real-time booking stats: %v", err)
		// Fallback to static values (no bookings assumed)
		return &dto.PeriodAvailabilityResponse{
			PeriodID:       periodID,
			TotalSeats:     period.Seat,
			BookedSeats:    0, // Fallback: assume no bookings
			AvailableSeats: period.Seat, // Fallback: all seats available
			IsAvailable:    period.Status == models.PeriodStatusPublish,
		}, nil
	}

	return &dto.PeriodAvailabilityResponse{
		PeriodID:       periodID,
		TotalSeats:     realTimeStats.TotalSeats,
		BookedSeats:    realTimeStats.BookedPassengers,
		AvailableSeats: realTimeStats.AvailableSeats,
		IsAvailable:    realTimeStats.IsAvailable,
	}, nil
}

// Helper methods
func (s *PeriodServiceImpl) toPeriodResponse(ctx context.Context, period *models.Period) *dto.PeriodResponse {
	log.Printf("[DEBUG] toPeriodResponse called for period ID: %s", period.ID.String())

	// Calculate real-time booking stats
	realTimeStats, err := s.CalculateRealTimeBookingStats(ctx, period.ID)
	if err != nil {
		log.Printf("[ERROR] Failed to calculate real-time booking stats for period %s: %v", period.ID, err)
		log.Printf("[FALLBACK] Using static booking values for period %s", period.ID.String())
		log.Printf("[FALLBACK] Static values - Using fallback with no bookings, TotalSeats: %d",
			period.Seat)

		// Fallback to static values with proper structure (no bookings assumed)
		realTimeStats = &dto.PeriodBookingStatsResponse{
			PeriodID:         period.ID,
			TotalSeats:       period.Seat,
			TotalBookings:    0,              // Fallback: assume no bookings
			BookedPassengers: 0,              // Fallback: assume no passengers
			HeldSeats:        0,              // No held seats in fallback mode
			AvailableSeats:   period.Seat,    // Fallback: all seats available
			BookingsByStatus: map[string]int{}, // Empty status map for fallback
			IsAvailable:      period.Status == models.PeriodStatusPublish,
		}

		log.Printf("[FALLBACK] Fallback stats created - TotalSeats: %d, BookedPassengers: %d, AvailableSeats: %d, IsAvailable: %t",
			realTimeStats.TotalSeats, realTimeStats.BookedPassengers, realTimeStats.AvailableSeats, realTimeStats.IsAvailable)
	} else {
		log.Printf("[SUCCESS] Real-time booking stats calculated successfully for period %s", period.ID.String())
	}

	response := &dto.PeriodResponse{
		ID:             period.ID,
		TourID:         period.TourID,
		StartDate:      period.StartDate.Format("2006-01-02"),
		EndDate:        period.EndDate.Format("2006-01-02"),
		Seat:           period.Seat,
		Booking:        realTimeStats.BookedPassengers, // Use real-time booked passengers count
		HeldSeats:      realTimeStats.HeldSeats,        // Use real-time held seats count
		AvailableSeats: realTimeStats.AvailableSeats,   // Use real-time available seats
		Deposit:        period.Deposit,
		DepositType:    period.DepositType,
		//SingleStay:     period.SingleStay,
		Airline:          period.Airline,
		Airport:          period.Airport,
		Note:             period.Note,
		RateType:         string(period.RateType),
		CommissionType:   string(period.CommissionType),
		Status:           string(period.Status),
		ServiceFeePerPax: period.ServiceFeePerPax,
		CreatedAt:        period.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        period.UpdatedAt.Format(time.RFC3339),
	}

	// Convert flight sectors
	response.Sector1 = s.toFlightSectorResponses(period.Sector1)
	response.Sector2 = s.toFlightSectorResponses(period.Sector2)

	// Convert rate items
	response.RateItems = s.toRateItemResponses(period.RateItems)

	// Convert commission items
	response.CommissionItems = s.toCommissionItemResponses(period.CommissionItems)

	// Convert room types
	response.RoomTypes = s.toRoomTypeResponses(period.RoomTypes)

	// Convert additional services
	response.AdditionalServices = s.toAdditionalServiceResponses(period.AdditionalServices)

	// Convert file attachments
	response.Files = s.toFileAttachmentResponses(period.Files)

	return response
}

func (s *PeriodServiceImpl) convertFlightSectors(sectors []*dto.CreateFlightSectorRequest) models.FlightSectors {
	result := make(models.FlightSectors, len(sectors))
	for i, sector := range sectors {
		result[i] = models.FlightSector{
			Date:   sector.Date,
			Flight: sector.Flight,
			Route:  sector.Route,
		}
	}
	return result
}

func (s *PeriodServiceImpl) convertRateItems(items []*dto.CreateRateItemRequest) models.RateItems {
	result := make(models.RateItems, len(items))
	for i, item := range items {
		// Generate UUID if ID is empty
		id := item.ID
		if id == "" {
			id = uuid.New().String()
		}
		result[i] = models.RateItem{
			ID:      id,
			Title:   item.Title,
			Regular: item.Regular,
			Sale:    item.Sale,
			Exclude: item.Exclude,
		}
	}
	return result
}

func (s *PeriodServiceImpl) convertCommissionItems(items []*dto.CreateCommissionItemRequest) models.CommissionItems {
	result := make(models.CommissionItems, len(items))
	for i, item := range items {
		result[i] = models.CommissionItem{
			Regular: item.Regular,
			Sale:    item.Sale,
		}
	}
	return result
}

func (s *PeriodServiceImpl) convertRoomTypes(items []*dto.CreateRoomTypeRequest) models.RoomTypes {
	if items == nil {
		return models.RoomTypes{}
	}
	result := make(models.RoomTypes, len(items))
	for i, item := range items {
		// Generate UUID if ID is empty
		id := item.ID
		if id == "" {
			id = uuid.New().String()
		}
		result[i] = models.RoomType{
			ID:    id,
			Title: item.Title,
			Price: item.Price,
		}
	}
	return result
}

func (s *PeriodServiceImpl) toFlightSectorResponses(sectors models.FlightSectors) []*dto.FlightSectorResponse {
	responses := make([]*dto.FlightSectorResponse, len(sectors))
	for i, sector := range sectors {
		responses[i] = &dto.FlightSectorResponse{
			Date:   sector.Date,
			Flight: sector.Flight,
			Route:  sector.Route,
		}
	}
	return responses
}

func (s *PeriodServiceImpl) toRateItemResponses(items models.RateItems) []*dto.RateItemResponse {
	responses := make([]*dto.RateItemResponse, len(items))
	for i, item := range items {
		responses[i] = &dto.RateItemResponse{
			ID:      item.ID,
			Title:   item.Title,
			Regular: item.Regular,
			Sale:    item.Sale,
			Exclude: item.Exclude,
		}
	}
	return responses
}

func (s *PeriodServiceImpl) toCommissionItemResponses(items models.CommissionItems) []*dto.CommissionItemResponse {
	responses := make([]*dto.CommissionItemResponse, len(items))
	for i, item := range items {
		responses[i] = &dto.CommissionItemResponse{
			Regular: item.Regular,
			Sale:    item.Sale,
		}
	}
	return responses
}

func (s *PeriodServiceImpl) toRoomTypeResponses(items models.RoomTypes) []*dto.RoomTypeResponse {
	responses := make([]*dto.RoomTypeResponse, len(items))
	for i, item := range items {
		responses[i] = &dto.RoomTypeResponse{
			ID:    item.ID,
			Title: item.Title,
			Price: item.Price,
		}
	}
	return responses
}

func (s *PeriodServiceImpl) convertAdditionalServices(items []*dto.CreateAdditionalServiceRequest) models.AdditionalServices {
	if items == nil {
		return models.AdditionalServices{}
	}
	result := make(models.AdditionalServices, len(items))
	for i, item := range items {
		// Generate UUID if ID is empty
		id := item.ID
		if id == "" {
			id = uuid.New().String()
		}
		result[i] = models.AdditionalService{
			ID:       id,
			Title:    item.Title,
			Price:    item.Price,
			RateType: item.RateType,
			UnitName: item.UnitName,
		}
	}
	return result
}

func (s *PeriodServiceImpl) toAdditionalServiceResponses(items models.AdditionalServices) []*dto.AdditionalServiceResponse {
	responses := make([]*dto.AdditionalServiceResponse, len(items))
	for i, item := range items {
		responses[i] = &dto.AdditionalServiceResponse{
			ID:       item.ID,
			Title:    item.Title,
			Price:    item.Price,
			RateType: item.RateType,
			UnitName: item.UnitName,
		}
	}
	return responses
}

func (s *PeriodServiceImpl) toFileAttachmentResponses(files models.FileAttachments) []*dto.FileAttachmentResponse {
	responses := make([]*dto.FileAttachmentResponse, len(files))
	for i, file := range files {
		responses[i] = &dto.FileAttachmentResponse{
			Name: file.Name,
			Type: file.Type,
			URL:  cdn.GetFullFileURL(file.URL), // Construct full URL for frontend
		}
	}
	return responses
}

func (s *PeriodServiceImpl) generateFileURL(periodID uuid.UUID, fileName, fileType string) string {
	// Generate a unique file path
	timestamp := time.Now().Unix()
	return fmt.Sprintf("periods/%s/%s/%d_%s", periodID.String(), fileType, timestamp, fileName)
}

// New converter methods for separate tables
func (s *PeriodServiceImpl) convertToPeriodRateItems(items []*dto.CreateRateItemRequest) []*models.PeriodRateItem {
	result := make([]*models.PeriodRateItem, len(items))
	for i, item := range items {
		// Use provided sort order or default to index
		sortOrder := i
		if item.SortOrder != nil {
			sortOrder = *item.SortOrder
		}

		rateItem := &models.PeriodRateItem{
			Title:     item.Title,
			Regular:   item.Regular,
			Sale:      item.Sale,
			Exclude:   item.Exclude,
			SortOrder: sortOrder,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Handle existing ID if provided
		if item.ID != "" {
			if id, err := uuid.Parse(item.ID); err == nil {
				rateItem.ID = id
			}
		}

		result[i] = rateItem
	}
	return result
}

func (s *PeriodServiceImpl) convertToPeriodRoomTypes(items []*dto.CreateRoomTypeRequest) []*models.PeriodRoomType {
	if items == nil {
		return []*models.PeriodRoomType{}
	}
	result := make([]*models.PeriodRoomType, len(items))
	for i, item := range items {
		// Use provided sort order or default to index
		sortOrder := i
		if item.SortOrder != nil {
			sortOrder = *item.SortOrder
		}

		roomType := &models.PeriodRoomType{
			Title:     item.Title,
			Price:     item.Price,
			SortOrder: sortOrder,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Handle existing ID if provided
		if item.ID != "" {
			if id, err := uuid.Parse(item.ID); err == nil {
				roomType.ID = id
			}
		}

		result[i] = roomType
	}
	return result
}

func (s *PeriodServiceImpl) convertToPeriodAdditionalServices(items []*dto.CreateAdditionalServiceRequest) []*models.PeriodAdditionalService {
	if items == nil {
		return []*models.PeriodAdditionalService{}
	}
	result := make([]*models.PeriodAdditionalService, len(items))
	for i, item := range items {
		// Use provided sort order or default to index
		sortOrder := i
		if item.SortOrder != nil {
			sortOrder = *item.SortOrder
		}

		additionalService := &models.PeriodAdditionalService{
			Title:     item.Title,
			Price:     item.Price,
			RateType:  item.RateType,
			UnitName:  item.UnitName,
			SortOrder: sortOrder,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Handle existing ID if provided
		if item.ID != "" {
			if id, err := uuid.Parse(item.ID); err == nil {
				additionalService.ID = id
			}
		}

		result[i] = additionalService
	}
	return result
}

// New method to create period response using separate tables
func (s *PeriodServiceImpl) toPeriodResponseWithSeparateTables(ctx context.Context, period *models.Period) *dto.PeriodResponse {
	// Calculate real-time booking stats (same as toPeriodResponse)
	realTimeStats, err := s.CalculateRealTimeBookingStats(ctx, period.ID)
	if err != nil {
		log.Printf("[ERROR] Failed to calculate real-time booking stats for period %s: %v", period.ID, err)
		log.Printf("[FALLBACK] Using static booking values for period %s", period.ID.String())

		// Fallback to static values with proper structure (no bookings assumed)
		realTimeStats = &dto.PeriodBookingStatsResponse{
			PeriodID:         period.ID,
			TotalSeats:       period.Seat,
			TotalBookings:    0,              // Fallback: assume no bookings
			BookedPassengers: 0,              // Fallback: assume no passengers
			HeldSeats:        0,              // No held seats in fallback mode
			AvailableSeats:   period.Seat,    // Fallback: all seats available
			BookingsByStatus: map[string]int{},
			IsAvailable:      period.Status == models.PeriodStatusPublish,
		}
	}

	response := &dto.PeriodResponse{
		ID:             period.ID,
		TourID:         period.TourID,
		StartDate:      period.StartDate.Format("2006-01-02"),
		EndDate:        period.EndDate.Format("2006-01-02"),
		Seat:           period.Seat,
		Booking:        realTimeStats.BookedPassengers, // Use real-time booked passengers count
		HeldSeats:      realTimeStats.HeldSeats,        // Use real-time held seats count
		AvailableSeats: realTimeStats.AvailableSeats,   // Use real-time available seats
		Deposit:        period.Deposit,
		DepositType:    period.DepositType,
		//	SingleStay:     period.SingleStay,
		Airline:          period.Airline,
		Airport:          period.Airport,
		Note:             period.Note,
		RateType:         string(period.RateType),
		CommissionType:   string(period.CommissionType),
		Status:           string(period.Status),
		ServiceFeePerPax: period.ServiceFeePerPax,
		CreatedAt:        period.CreatedAt.Format(time.RFC3339),
		UpdatedAt:        period.UpdatedAt.Format(time.RFC3339),

		// Expiration Override fields
		ExpirationOverride: period.ExpirationOverride,
		ExpirationMinutes:  period.ExpirationMinutes,
		ExpirationHours:    period.ExpirationHours,
		ExpirationDays:     period.ExpirationDays,
		ExpirationMessage:  period.ExpirationMessage,
		ExpirationDisplay:  period.GetExpirationDisplay(),
	}

	// Convert flight sectors (still from Period model)
	response.Sector1 = s.toFlightSectorResponses(period.Sector1)
	response.Sector2 = s.toFlightSectorResponses(period.Sector2)

	// Convert commission items (still from Period model)
	response.CommissionItems = s.toCommissionItemResponses(period.CommissionItems)

	// Convert file attachments (still from Period model)
	response.Files = s.toFileAttachmentResponses(period.Files)

	// Load and convert data from separate tables
	rateItems, err := s.periodRateItemRepo.GetByPeriodID(ctx, period.ID)
	if err != nil {
		log.Printf("Error loading rate items for period %s: %v", period.ID, err)
		response.RateItems = []*dto.RateItemResponse{}
	} else {
		response.RateItems = s.toPeriodRateItemResponses(rateItems)
	}

	roomTypes, err := s.periodRoomTypeRepo.GetByPeriodID(ctx, period.ID)
	if err != nil {
		log.Printf("Error loading room types for period %s: %v", period.ID, err)
		response.RoomTypes = []*dto.RoomTypeResponse{}
	} else {
		response.RoomTypes = s.toPeriodRoomTypeResponses(roomTypes)
	}

	additionalServices, err := s.periodAdditionalServiceRepo.GetByPeriodID(ctx, period.ID)
	if err != nil {
		log.Printf("Error loading additional services for period %s: %v", period.ID, err)
		response.AdditionalServices = []*dto.AdditionalServiceResponse{}
	} else {
		response.AdditionalServices = s.toPeriodAdditionalServiceResponses(additionalServices)
	}

	// Check WL: มี booking active ที่จ่ายเงินไม่ครบหรือไม่
	response.HasWaitingList = s.checkHasWaitingList(ctx, period.ID)

	return response
}

// checkHasWaitingList เช็คว่า period นี้มี booking active ที่จ่ายเงินไม่ครบหรือไม่
// ใช้ SumPayment กับ GrandTotal/SumNet ที่เก็บใน DB ไม่ต้อง preload relations
func (s *PeriodServiceImpl) checkHasWaitingList(ctx context.Context, periodID uuid.UUID) bool {
	if s.bookingRepo == nil {
		return false
	}

	bookings, err := s.bookingRepo.GetByPeriodID(ctx, periodID, 0, 1000)
	if err != nil {
		return false
	}

	for _, booking := range bookings {
		// เฉพาะ booking ที่ active
		if booking.Status != "pending" && booking.Status != "processing" && booking.Status != "confirm" {
			continue
		}

		// เช็คจ่ายครบหรือไม่ (ใช้ field ที่เก็บใน DB ตรงๆ)
		total := booking.GrandTotal
		if total == 0 {
			total = booking.SumNet
		}
		if booking.SumPayment < total {
			return true
		}
	}

	return false
}

// New converter methods for separate table responses
func (s *PeriodServiceImpl) toPeriodRateItemResponses(items []*models.PeriodRateItem) []*dto.RateItemResponse {
	responses := make([]*dto.RateItemResponse, len(items))
	for i, item := range items {
		responses[i] = &dto.RateItemResponse{
			ID:        item.ID.String(),
			Title:     item.Title,
			Regular:   item.Regular,
			Sale:      item.Sale,
			Exclude:   item.Exclude,
			SortOrder: item.SortOrder,
		}
	}
	return responses
}

func (s *PeriodServiceImpl) toPeriodRoomTypeResponses(items []*models.PeriodRoomType) []*dto.RoomTypeResponse {
	responses := make([]*dto.RoomTypeResponse, len(items))
	for i, item := range items {
		responses[i] = &dto.RoomTypeResponse{
			ID:        item.ID.String(),
			Title:     item.Title,
			Price:     item.Price,
			SortOrder: item.SortOrder,
		}
	}
	return responses
}

func (s *PeriodServiceImpl) toPeriodAdditionalServiceResponses(items []*models.PeriodAdditionalService) []*dto.AdditionalServiceResponse {
	responses := make([]*dto.AdditionalServiceResponse, len(items))
	for i, item := range items {
		responses[i] = &dto.AdditionalServiceResponse{
			ID:        item.ID.String(),
			Title:     item.Title,
			Price:     item.Price,
			RateType:  item.RateType,
			UnitName:  item.UnitName,
			SortOrder: item.SortOrder,
		}
	}
	return responses
}
