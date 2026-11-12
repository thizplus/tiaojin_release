package handlers

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)


type BookingHandler struct {
	bookingService services.BookingService
	fileService    services.FileService
}

func NewBookingHandler(bookingService services.BookingService, fileService services.FileService) *BookingHandler {
	return &BookingHandler{
		bookingService: bookingService,
		fileService:    fileService,
	}
}

// Create creates a new booking
func (h *BookingHandler) Create(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CreateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Determine which user_id to use for the booking
	bookingUserID := user.ID // Default: use authenticated user's ID

	// If user_id is provided in payload and user is sales/admin/super, use that instead
	if req.UserID != nil {
		if user.Role == "sales" || user.Role == "admin" || user.Role == "super" {
			bookingUserID = *req.UserID
		} else {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Only sales, admin and super users can create bookings for other users", nil)
		}
	}

	booking, err := h.bookingService.CreateBooking(c.Context(), bookingUserID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "สร้างการจองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "สร้างการจองเรียบร้อยแล้ว", booking)
}

// GetByID returns booking by ID
func (h *BookingHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	booking, err := h.bookingService.GetBookingByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบการจอง", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลการจองสำเร็จ", booking)
}

// List returns paginated bookings list
func (h *BookingHandler) List(c *fiber.Ctx) error {
	_, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Parse pagination (frontend ส่ง page, backend แปลงเป็น offset)
	pagination := utils.ParsePaginationParams(c, 10)

	// Parse other query parameters
	status := c.Query("status")
	tourID := c.Query("tour_id")
	periodID := c.Query("period_id")
	salesID := c.Query("sales_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	search := c.Query("search")
	sort := c.Query("sort")
	order := c.Query("order")
	includeStat := c.Query("stat") == "true"

	// Parse payment slip filtering parameters
	paymentSlipStatus := c.Query("payment_slip_status")
	hasPendingPayments := c.Query("has_pending_payments")
	hasUnconfirmedPayments := c.Query("has_unconfirmed_payments")
	requiresPaymentApproval := c.Query("requires_payment_approval")

	listReq := dto.BookingListRequest{
		// Don't filter by UserID for admin routes - show all bookings
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
	if status != "" {
		// Support comma-separated status values for multiple status filtering
		if strings.Contains(status, ",") {
			statuses := strings.Split(status, ",")
			// Trim whitespace from each status
			for i := range statuses {
				statuses[i] = strings.TrimSpace(statuses[i])
			}
			listReq.Statuses = statuses
		} else {
			// Single status (backward compatibility)
			listReq.Status = &status
		}
	}
	if tourID != "" {
		if tourUUID, err := uuid.Parse(tourID); err == nil {
			listReq.TourID = &tourUUID
		}
	}
	if periodID != "" {
		if periodUUID, err := uuid.Parse(periodID); err == nil {
			listReq.PeriodID = &periodUUID
		}
	}
	if salesID != "" {
		if salesUUID, err := uuid.Parse(salesID); err == nil {
			listReq.SalesID = &salesUUID
		}
	}
	if startDate != "" {
		if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
			listReq.StartDate = &parsed
		}
	}
	if endDate != "" {
		if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
			listReq.EndDate = &parsed
		}
	}
	if search != "" {
		listReq.Search = &search
	}
	if sort != "" {
		listReq.Sort = &sort
	}
	if order != "" {
		listReq.Order = &order
	}

	// Apply payment slip filtering parameters
	if paymentSlipStatus != "" {
		listReq.PaymentSlipStatus = &paymentSlipStatus
	}
	if hasPendingPayments != "" {
		if val, err := strconv.ParseBool(hasPendingPayments); err == nil {
			listReq.HasPendingPayments = &val
		}
	}
	if hasUnconfirmedPayments != "" {
		if val, err := strconv.ParseBool(hasUnconfirmedPayments); err == nil {
			listReq.HasUnconfirmedPayments = &val
		}
	}
	if requiresPaymentApproval != "" {
		if val, err := strconv.ParseBool(requiresPaymentApproval); err == nil {
			listReq.RequiresPaymentApproval = &val
		}
	}

	result, err := h.bookingService.ListBookings(c.Context(), listReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการการจองไม่สำเร็จ", err)
	}

	// If stat=true, calculate statistics for all filtered records
	if includeStat {
		stats, err := h.bookingService.GetBookingStatistics(c.Context(), listReq)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "คำนวณสถิติการจองไม่สำเร็จ", err)
		}

		// Add statistics to response
		response := fiber.Map{
			"bookings":   result.Bookings,
			"pagination": result.Pagination,
			"statistics": stats,
		}
		return utils.SuccessResponse(c, "ดึงรายการการจองพร้อมสถิติสำเร็จ", response)
	}

	return utils.SuccessResponse(c, "ดึงรายการการจองสำเร็จ", result)
}

// Update updates booking information
func (h *BookingHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	var req dto.UpdateBookingRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	booking, err := h.bookingService.UpdateBooking(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตการจองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตการจองเรียบร้อยแล้ว", booking)
}

// UpdateStatus updates booking status (admin only)
func (h *BookingHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	var req dto.UpdateBookingStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	req.BookingID = id
	booking, err := h.bookingService.UpdateBookingStatus(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตสถานะการจองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตสถานะการจองเรียบร้อยแล้ว", booking)
}

// Cancel cancels a booking
func (h *BookingHandler) Cancel(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	var req struct {
		Reason string `json:"reason"` // Optional for agents, required for admin
	}

	// Try to parse request body (may be empty for agent routes)
	c.BodyParser(&req)

	// If reason is empty, set a default
	if req.Reason == "" {
		req.Reason = "Cancelled by user"
	}

	// Get current user to check permissions
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// For agent and sale routes, verify booking ownership unless user is admin/super
	if user.Role == "agent" || user.Role == "sale" {
		booking, err := h.bookingService.GetBookingByID(c.Context(), id)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", err)
		}

		// Check if agent/sale owns this booking
		if booking.UserID != user.ID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only cancel your own bookings", nil)
		}
	}

	err = h.bookingService.CancelBooking(c.Context(), id, req.Reason)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ยกเลิกการจองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ยกเลิกการจองเรียบร้อยแล้ว", nil)
}

// GetStats returns booking statistics (admin only)
func (h *BookingHandler) GetStats(c *fiber.Ctx) error {
	// Parse query parameters for date range
	var dateFrom, dateTo *time.Time

	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = &parsed
		}
	}

	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = &parsed
		}
	}

	stats, err := h.bookingService.GetBookingStats(c.Context(), dateFrom, dateTo)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงสถิติการจองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงสถิติการจองสำเร็จ", stats)
}

// GetCommissionData returns commission data for bookings
func (h *BookingHandler) GetCommissionData(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Parse query parameters
	var dateFrom, dateTo *time.Time
	if dateFromStr := c.Query("date_from"); dateFromStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			dateFrom = &parsed
		}
	}

	if dateToStr := c.Query("date_to"); dateToStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateToStr); err == nil {
			dateTo = &parsed
		}
	}

	// For now, return a basic commission data structure
	// This would be implemented based on business requirements
	commissionData := fiber.Map{
		"user_id":    user.ID,
		"date_from":  dateFrom,
		"date_to":    dateTo,
		"total_commission": 0,
		"bookings_count":   0,
		"commission_rate":  0.05, // 5% default
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลค่าคอมมิชชั่นสำเร็จ", commissionData)
}

// Admin Payment Management Endpoints

// GetPendingPayments retrieves all pending payment slips for admin review
func (h *BookingHandler) GetPendingPayments(c *fiber.Ctx) error {
	// Parse pagination (frontend ส่ง page, backend แปลงเป็น offset)
	pagination := utils.ParsePaginationParams(c, 20)

	// Get pending payments
	payments, err := h.bookingService.GetPendingPayments(c.Context(), pagination.Offset, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get pending payments", err)
	}

	return utils.SuccessResponse(c, "Pending payments retrieved successfully", fiber.Map{
		"payments": payments,
		"page":     pagination.Page,
		"limit":    pagination.Limit,
		"count":    len(payments),
	})
}

// VerifyPaymentSlip allows admin to approve or reject payment slips
func (h *BookingHandler) VerifyPaymentSlip(c *fiber.Ctx) error {
	// Parse payment slip ID from URL
	paymentSlipID, err := uuid.Parse(c.Params("paymentSlipId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid payment slip ID")
	}

	// Parse request body
	var req dto.VerifyPaymentSlipRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request format")
	}

	// Set the payment slip ID from URL
	req.PaymentSlipID = paymentSlipID

	// Validate request
	if req.Status != "confirm" && req.Status != "reject" {
		return utils.ValidationErrorResponse(c, "Status must be either 'confirm' or 'reject'")
	}

	// Validate amount adjustment if provided
	if req.ActualAmount != nil {
		if req.Status == "reject" {
			return utils.ValidationErrorResponse(c, "Cannot adjust amount when rejecting payment slip")
		}
		if *req.ActualAmount <= 0 {
			return utils.ValidationErrorResponse(c, "Actual amount must be greater than 0")
		}
	}

	// Process verification
	paymentSlip, err := h.bookingService.VerifyPaymentSlip(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to verify payment slip", err)
	}

	// Get updated booking information for enhanced response
	booking, err := h.bookingService.GetBookingModelByID(c.Context(), paymentSlip.BookingID)
	if err != nil {
		// If booking fetch fails, still return payment slip response
		return utils.SuccessResponse(c, "Payment slip verified successfully", paymentSlip)
	}

	// Enhanced response with booking status information
	actualPayment := booking.GetActualSumPayment()
	response := fiber.Map{
		"payment_slip": paymentSlip,
		"booking_update": fiber.Map{
			"booking_id":        booking.ID,
			"status":            booking.Status,
			"sum_payment":       actualPayment, // ใช้ computed value
			"sum_deposit":       booking.SumDeposit,
			"sum_net":           booking.SumNet,
			"payment_ratio":     actualPayment / booking.SumNet,
			"deposit_ratio":     actualPayment / booking.SumDeposit,
			"fully_paid":        actualPayment >= booking.SumNet,
			"deposit_satisfied": actualPayment >= booking.SumDeposit,
			"remaining_amount":  booking.SumNet - actualPayment,
			"updated_at":        booking.UpdatedAt,
		},
	}

	// Add adjustment information if amount was adjusted
	if req.ActualAmount != nil {
		response["amount_adjustment"] = fiber.Map{
			"original_amount":   paymentSlip.OriginalAmount,
			"adjusted_amount":   paymentSlip.Amount,
			"adjustment_reason": paymentSlip.AdjustmentReason,
			"was_adjusted":      paymentSlip.OriginalAmount != nil,
		}
	}

	var message string
	if req.Status == "confirm" {
		if actualPayment >= booking.SumNet {
			message = "Payment verified and booking confirmed (fully paid)"
		} else if actualPayment >= booking.SumDeposit {
			message = "Payment verified - deposit satisfied, booking protected from expiry"
		} else {
			message = "Payment verified - partial payment received, booking protected from expiry"
		}

		// Add adjustment information to message if applicable
		if req.ActualAmount != nil && paymentSlip.OriginalAmount != nil {
			message += fmt.Sprintf(" (amount adjusted: %.2f → %.2f)", *paymentSlip.OriginalAmount, paymentSlip.Amount)
		}
	} else {
		message = "Payment slip rejected"
	}

	return utils.SuccessResponse(c, message, response)
}

// GetPaymentSlipDetails retrieves details of a specific payment slip
func (h *BookingHandler) GetPaymentSlipDetails(c *fiber.Ctx) error {
	paymentSlipID, err := uuid.Parse(c.Params("paymentSlipId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid payment slip ID")
	}

	// Get payment slip history which includes all details
	history, err := h.bookingService.GetPaymentSlipHistory(c.Context(), paymentSlipID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Failed to get payment slip details", err)
	}

	return utils.SuccessResponse(c, "Payment slip details retrieved successfully", history.PaymentSlip)
}

// ReversePaymentSlip allows admin to reverse a confirmed payment slip
func (h *BookingHandler) ReversePaymentSlip(c *fiber.Ctx) error {
	// Parse payment slip ID from URL
	paymentSlipID, err := uuid.Parse(c.Params("paymentSlipId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid payment slip ID")
	}

	// Parse request body
	var req dto.ReversePaymentSlipRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request format")
	}

	// Set the payment slip ID from URL
	req.PaymentSlipID = paymentSlipID

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Get admin user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Process reversal
	paymentSlip, err := h.bookingService.ReversePaymentSlip(c.Context(), req, user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to reverse payment slip", err)
	}

	// Get updated booking information for enhanced response
	booking, err := h.bookingService.GetBookingModelByID(c.Context(), paymentSlip.BookingID)
	if err != nil {
		// If booking fetch fails, still return payment slip response
		return utils.SuccessResponse(c, "Payment slip reversed successfully", paymentSlip)
	}

	// Enhanced response with booking status information
	actualPayment := booking.GetActualSumPayment()
	response := fiber.Map{
		"payment_slip": paymentSlip,
		"booking_update": fiber.Map{
			"booking_id":        booking.ID,
			"status":            booking.Status,
			"sum_payment":       actualPayment, // ใช้ computed value
			"sum_deposit":       booking.SumDeposit,
			"sum_net":           booking.SumNet,
			"remaining_amount":  booking.SumNet - actualPayment,
			"updated_at":        booking.UpdatedAt,
		},
	}

	return utils.SuccessResponse(c, "Payment slip reversed successfully", response)
}

// GetPaymentSlipHistory retrieves complete history and timeline for a payment slip
func (h *BookingHandler) GetPaymentSlipHistory(c *fiber.Ctx) error {
	paymentSlipID, err := uuid.Parse(c.Params("paymentSlipId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid payment slip ID")
	}

	// Get payment slip history
	history, err := h.bookingService.GetPaymentSlipHistory(c.Context(), paymentSlipID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Failed to get payment slip history", err)
	}

	return utils.SuccessResponse(c, "Payment slip history retrieved successfully", history)
}

// CreateAdminPayment allows admins to create direct payments with optional slip upload
func (h *BookingHandler) CreateAdminPayment(c *fiber.Ctx) error {
	// Get booking ID from path
	bookingID, err := uuid.Parse(c.Params("bookingId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	// Get admin user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Determine content type and parse accordingly
	contentType := c.Get("Content-Type")
	isFormData := strings.Contains(strings.ToLower(contentType), "multipart/form-data")

	var req dto.CreateAdminPaymentRequest
	var slipImageURL *string

	if isFormData {
		// Handle form-data (with potential file upload)
		req, slipImageURL, err = h.parseFormDataRequest(c, bookingID, user.ID)
		if err != nil {
			return err // Error already formatted by parseFormDataRequest
		}
	} else {
		// Handle JSON request (no file upload)
		if err := c.BodyParser(&req); err != nil {
			return utils.ValidationErrorResponse(c, "Invalid JSON request body")
		}
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ValidationErrorResponse(c, err.Error())
	}

	// Set slip image URL if file was uploaded
	if slipImageURL != nil {
		req.SlipImage = slipImageURL
	}

	// Create admin payment
	response, err := h.bookingService.CreateAdminPayment(c.Context(), bookingID, req, user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create admin payment", err)
	}

	return utils.SuccessResponse(c, "Admin payment created successfully", response)
}



// parseFormDataRequest parses form-data request and handles optional file upload
func (h *BookingHandler) parseFormDataRequest(c *fiber.Ctx, bookingID uuid.UUID, userID uuid.UUID) (dto.CreateAdminPaymentRequest, *string, error) {
	var req dto.CreateAdminPaymentRequest
	var slipImageURL *string

	// Parse required fields from form
	amountStr := c.FormValue("amount")
	if amountStr == "" {
		return req, nil, utils.ValidationErrorResponse(c, "Amount is required")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return req, nil, utils.ValidationErrorResponse(c, "Invalid amount format")
	}

	paymentMethod := c.FormValue("payment_method")
	if paymentMethod == "" {
		return req, nil, utils.ValidationErrorResponse(c, "Payment method is required")
	}

	// Optional note
	note := c.FormValue("note")
	var notePtr *string
	if note != "" {
		notePtr = &note
	}

	// Handle optional file upload
	file, err := c.FormFile("slip_image")
	if err == nil && file != nil {
		// File was provided, validate and upload
		if file.Size > 5*1024*1024 { // 5MB limit
			return req, nil, utils.ValidationErrorResponse(c, "File size exceeds 5MB limit")
		}

		fileExt := strings.ToLower(filepath.Ext(file.Filename))
		allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
		isValidExt := false
		for _, ext := range allowedExtensions {
			if fileExt == ext {
				isValidExt = true
				break
			}
		}

		if !isValidExt {
			return req, nil, utils.ValidationErrorResponse(c, "Only JPG, JPEG, PNG, and PDF files are allowed")
		}

		// Upload file to S3 Storage
		uploadOptions := &dto.UploadFileRequest{
			Category: "payment_slips",
			EntityID: bookingID.String(),
			FileType: "admin_slip",
		}

		uploadedFile, err := h.fileService.UploadFile(c.Context(), userID, file, uploadOptions)
		if err != nil {
			return req, nil, utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload slip image", err)
		}

		slipImageURL = &uploadedFile.URL
	} else if err != nil && strings.Contains(err.Error(), "no such file") {
		// No file provided, which is fine for admin payments
	} else if err != nil {
		// Other file upload error
		return req, nil, utils.ValidationErrorResponse(c, "Error processing slip image: "+err.Error())
	}

	// Build request object
	req = dto.CreateAdminPaymentRequest{
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Note:          notePtr,
	}

	return req, slipImageURL, nil
}

// GetAdminReport returns admin booking report for specified month
func (h *BookingHandler) GetAdminReport(c *fiber.Ctx) error {
	// Get user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Check authorization: only admin and super users can access admin report
	if user.Role != "admin" && user.Role != "super" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied. Admin or super role required", nil)
	}

	// Parse query parameters
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil || year < 2020 || year > 2030 {
		return utils.ValidationErrorResponse(c, "Invalid year parameter")
	}

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 1 || month > 12 {
		return utils.ValidationErrorResponse(c, "Invalid month parameter")
	}

	// Get admin report data from service
	report, err := h.bookingService.GetAdminReport(c.Context(), year, month)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve admin report", err)
	}

	return utils.SuccessResponse(c, "Admin report retrieved successfully", report)
}

// GetSalesReport returns sales booking report for specified month (for current sales user)
func (h *BookingHandler) GetSalesReport(c *fiber.Ctx) error {
	// Get user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Check authorization: only sales, admin and super users can access sales report
	if user.Role != "sales" && user.Role != "admin" && user.Role != "super" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied. Sales, admin or super role required", nil)
	}

	// Parse query parameters
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil || year < 2020 || year > 2030 {
		return utils.ValidationErrorResponse(c, "Invalid year parameter")
	}

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 1 || month > 12 {
		return utils.ValidationErrorResponse(c, "Invalid month parameter")
	}

	// Get sales report data from service using current user's ID
	report, err := h.bookingService.GetSalesReport(c.Context(), user.ID, year, month)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve sales report", err)
	}

	return utils.SuccessResponse(c, "Sales report retrieved successfully", report)
}

// GetAgentReport returns agent booking report for specified month (for current agent user)
func (h *BookingHandler) GetAgentReport(c *fiber.Ctx) error {
	// Get user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Check authorization: only agent, admin and super users can access agent report
	if user.Role != "agent" && user.Role != "admin" && user.Role != "super" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Access denied. Agent, admin or super role required", nil)
	}

	// Parse query parameters
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil || year < 2020 || year > 2030 {
		return utils.ValidationErrorResponse(c, "Invalid year parameter")
	}

	month, err := strconv.Atoi(c.Query("month"))
	if err != nil || month < 1 || month > 12 {
		return utils.ValidationErrorResponse(c, "Invalid month parameter")
	}

	// Get agent report data from service using current user's ID
	report, err := h.bookingService.GetAgentReport(c.Context(), user.ID, year, month)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve agent report", err)
	}

	return utils.SuccessResponse(c, "Agent report retrieved successfully", report)
}

// AdminUploadPassport allows admin to upload passport for any booking
func (h *BookingHandler) AdminUploadPassport(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	bookingID, err := uuid.Parse(c.Params("bookingId"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	// Verify booking exists
	_, err = h.bookingService.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", err)
	}

	// Handle file upload
	file, err := c.FormFile("passport")
	if err != nil {
		return utils.ValidationErrorResponse(c, "Passport file is required")
	}

	// Validate file type and size (5MB limit)
	if file.Size > 5*1024*1024 {
		return utils.ValidationErrorResponse(c, "File size exceeds 5MB limit")
	}

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
	isValidExt := false
	for _, ext := range allowedExtensions {
		if fileExt == ext {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return utils.ValidationErrorResponse(c, "Invalid file type. Only JPG, PNG, and PDF files are allowed")
	}

	// Upload file to S3 Storage using FileService
	uploadOptions := &dto.UploadFileRequest{
		Category: "passports",
		EntityID: bookingID.String(),
		FileType: "passport",
	}

	uploadedFile, err := h.fileService.UploadFile(c.Context(), user.ID, file, uploadOptions)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload file to storage", err)
	}

	// Open and read file data for Vision API + LLM processing
	src, err := file.Open()
	if err != nil {
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err)
	}
	defer src.Close()

	// Read file data
	fileData, err := io.ReadAll(src)
	if err != nil {
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read file", err)
	}

	// Get passport number from form (optional, will be extracted by OCR if not provided)
	passportNumber := c.FormValue("passport_number")
	if passportNumber == "" {
		passportNumber = "AUTO_EXTRACT"
	}

	// Create UploadPassportRequest for the booking service with S3 Storage URL
	req := dto.UploadPassportRequest{
		BookingID:      bookingID,
		PassportNumber: passportNumber,
		ImageData:      fileData,
		FileName:       file.Filename,
		ImageURL:       uploadedFile.URL,
	}

	// Process passport through booking service (includes Vision API + LLM processing)
	passportResponse, err := h.bookingService.UploadPassport(c.Context(), req)
	if err != nil {
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to process passport", err)
	}

	// Enhanced response
	response := fiber.Map{
		"passport":   passportResponse,
		"booking_id": bookingID,
		"storage_info": fiber.Map{
			"file_id":   uploadedFile.ID,
			"file_url":  uploadedFile.URL,
			"file_name": uploadedFile.FileName,
			"file_size": uploadedFile.FileSize,
			"mime_type": uploadedFile.MimeType,
		},
		"message": "Passport uploaded and processed successfully by admin",
	}

	return utils.SuccessResponse(c, "Passport uploaded successfully", response)
}

