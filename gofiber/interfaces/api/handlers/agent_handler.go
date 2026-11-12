package handlers

import (
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AgentHandler struct {
	tourService    services.TourService
	periodService  services.PeriodService
	bookingService services.BookingService
	userService    services.UserService
	fileService    services.FileService
	companyService services.CompanyService
	salesService   services.SalesService
}

func NewAgentHandler(tourService services.TourService, periodService services.PeriodService, bookingService services.BookingService, userService services.UserService, fileService services.FileService, companyService services.CompanyService, salesService services.SalesService) *AgentHandler {
	return &AgentHandler{
		tourService:    tourService,
		periodService:  periodService,
		bookingService: bookingService,
		userService:    userService,
		fileService:    fileService,
		companyService: companyService,
		salesService:   salesService,
	}
}

// Tour Management
func (h *AgentHandler) CreateTour(c *fiber.Ctx) error {
	// Implementation for agent creating tours
	return utils.SuccessResponse(c, "Agent tour created successfully", nil)
}

func (h *AgentHandler) GetTours(c *fiber.Ctx) error {
	// Implementation for getting agent's tours
	return utils.SuccessResponse(c, "Agent tours retrieved successfully", nil)
}

func (h *AgentHandler) UpdateTour(c *fiber.Ctx) error {
	// Implementation for agent updating their tours
	return utils.SuccessResponse(c, "Agent tour updated successfully", nil)
}

// Period Management
func (h *AgentHandler) CreatePeriod(c *fiber.Ctx) error {
	// Implementation for agent creating periods
	return utils.SuccessResponse(c, "Agent period created successfully", nil)
}

// Booking Management
func (h *AgentHandler) GetBookings(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Parse query parameters (same as admin endpoint)
	pagination := utils.ParsePaginationParams(c, 10)
	status := c.Query("status")
	tourID := c.Query("tour_id")
	search := c.Query("search")
	sort := c.Query("sort")
	order := c.Query("order")

	listReq := dto.BookingListRequest{
		// Filter by UserID for agent - show only their bookings
		UserID: &user.ID,
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
	if status != "" {
		listReq.Status = &status
	}
	if tourID != "" {
		if tourUUID, err := uuid.Parse(tourID); err == nil {
			listReq.TourID = &tourUUID
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

	result, err := h.bookingService.ListBookings(c.Context(), listReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการการจองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "รายการการจองของผู้ใช้ดึงข้อมูลสำเร็จ", result)
}

func (h *AgentHandler) CreateBooking(c *fiber.Ctx) error {
	// Implementation for agent creating customer bookings
	return utils.SuccessResponse(c, "Customer booking created successfully", nil)
}

func (h *AgentHandler) UpdateBookingStatus(c *fiber.Ctx) error {
	// Implementation for agent updating booking status
	return utils.SuccessResponse(c, "Booking status updated successfully", nil)
}

func (h *AgentHandler) CancelBooking(c *fiber.Ctx) error {
	// Implementation for agent cancelling bookings
	return utils.SuccessResponse(c, "Booking cancelled successfully", nil)
}

// Customer Management
func (h *AgentHandler) GetCustomers(c *fiber.Ctx) error {
	// Implementation for getting agent's customers
	return utils.SuccessResponse(c, "Agent customers retrieved successfully", nil)
}

func (h *AgentHandler) GetCustomerDetails(c *fiber.Ctx) error {
	// Implementation for getting customer details
	return utils.SuccessResponse(c, "Customer details retrieved successfully", nil)
}

func (h *AgentHandler) GetCustomerBookings(c *fiber.Ctx) error {
	// Implementation for getting customer booking history
	return utils.SuccessResponse(c, "Customer bookings retrieved successfully", nil)
}

// Commission & Reports
func (h *AgentHandler) GetCommissions(c *fiber.Ctx) error {
	// Implementation for getting agent commissions
	return utils.SuccessResponse(c, "Agent commissions retrieved successfully", nil)
}

func (h *AgentHandler) GetMonthlyReport(c *fiber.Ctx) error {
	// Implementation for getting agent monthly report
	return utils.SuccessResponse(c, "Monthly report retrieved successfully", nil)
}

func (h *AgentHandler) GetPerformanceStats(c *fiber.Ctx) error {
	// Implementation for getting agent performance stats
	return utils.SuccessResponse(c, "Performance stats retrieved successfully", nil)
}

// NEW CUSTOMER-LEVEL METHODS FOR AGENTS

// Profile Management (customer-level)
func (h *AgentHandler) GetProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Get user's company profile (which serves as agent profile)
	profile, err := h.companyService.GetProfileByUserID(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Profile not found", err)
	}

	return utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

func (h *AgentHandler) UpdateProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.UpdateCompanyProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Update company profile
	updatedProfile, err := h.companyService.UpdateProfile(c.Context(), user.ID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update profile", err)
	}

	return utils.SuccessResponse(c, "Profile updated successfully", updatedProfile)
}

func (h *AgentHandler) CreateProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Check if user already has a company profile
	existingProfile, _ := h.companyService.GetProfileByUserID(c.Context(), user.ID)
	if existingProfile != nil {
		return utils.ErrorResponse(c, fiber.StatusConflict, "Profile already exists", fmt.Errorf("agent already has a company profile"))
	}

	var req dto.CreateCompanyProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Set user ID
	req.UserID = user.ID

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Create company profile
	profile, err := h.companyService.CreateProfile(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create profile", err)
	}

	return utils.SuccessResponse(c, "Profile created successfully", profile)
}

func (h *AgentHandler) CheckProfileExists(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Check if user has a company profile
	profile, err := h.companyService.GetProfileByUserID(c.Context(), user.ID)

	response := fiber.Map{
		"has_profile": profile != nil,
	}

	if profile != nil {
		response["profile_id"] = profile.ID
		response["company_name"] = profile.Name
	}

	return utils.SuccessResponse(c, "Profile status checked", response)
}

// Own Booking Management (customer-level)
func (h *AgentHandler) GetMyBookings(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Get only bookings that belong to this user
	bookings, err := h.bookingService.GetUserBookings(c.Context(), user.ID, 0, 100)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve bookings", err)
	}

	return utils.SuccessResponse(c, "My bookings retrieved successfully", bookings)
}

// DEPRECATED: Use BookingHandler.Cancel instead for unified cancellation with WebSocket broadcasting
func (h *AgentHandler) CancelMyBooking_DEPRECATED(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	// Get booking details first
	booking, err := h.bookingService.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", err)
	}

	// Verify booking ownership
	if booking.UserID != user.ID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only cancel your own bookings", nil)
	}

	// Enhanced business logic validation
	switch booking.Status {
	case "completed":
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Cannot cancel completed booking", nil)
	case "cancelled":
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Booking is already cancelled", nil)
	case "expired":
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Cannot cancel expired booking", nil)
	case "pending", "processing":
		// These statuses are allowed to cancel
	default:
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Booking status does not allow cancellation", nil)
	}

	// Get current period availability before cancellation
	periodAvailabilityBefore, err := h.periodService.CheckAvailability(c.Context(), booking.PeriodID, 0)
	if err != nil {
		log.Printf("Warning: Failed to get period availability before cancellation: %v", err)
	}

	// Cancel the booking with enhanced reason
	cancelReason := fmt.Sprintf("Agent cancellation requested by user %s at %s",
		user.Email, time.Now().Format("2006-01-02 15:04:05"))

	err = h.bookingService.CancelBooking(c.Context(), bookingID, cancelReason)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to cancel booking", err)
	}

	// Get updated booking details
	updatedBooking, err := h.bookingService.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		log.Printf("Warning: Failed to get updated booking details: %v", err)
	}

	// Get updated period availability after cancellation (seat recovery)
	periodAvailabilityAfter, err := h.periodService.CheckAvailability(c.Context(), booking.PeriodID, 0)
	if err != nil {
		log.Printf("Warning: Failed to get period availability after cancellation: %v", err)
	}

	// Enhanced response with seat recovery information
	response := fiber.Map{
		"booking": updatedBooking,
		"message": "Booking cancelled successfully",
		"cancellation_details": fiber.Map{
			"cancelled_at":   time.Now().Format(time.RFC3339),
			"cancelled_by":   user.Email,
			"reason":         "Agent cancellation",
			"seats_returned": booking.SumQty,
		},
	}

	// Add period availability information if available
	if periodAvailabilityBefore != nil && periodAvailabilityAfter != nil {
		response["seat_recovery"] = fiber.Map{
			"period_id": booking.PeriodID,
			"before_cancellation": fiber.Map{
				"available_seats": periodAvailabilityBefore.AvailableSeats,
				"booked_seats":    periodAvailabilityBefore.BookedSeats,
				"total_seats":     periodAvailabilityBefore.TotalSeats,
			},
			"after_cancellation": fiber.Map{
				"available_seats": periodAvailabilityAfter.AvailableSeats,
				"booked_seats":    periodAvailabilityAfter.BookedSeats,
				"total_seats":     periodAvailabilityAfter.TotalSeats,
			},
			"seats_recovered": periodAvailabilityAfter.AvailableSeats - periodAvailabilityBefore.AvailableSeats,
		}
	}

	return utils.SuccessResponse(c, "Booking cancelled and seats returned successfully", response)
}

// Document Upload for Own Bookings (customer-level)
func (h *AgentHandler) UploadPaymentSlip(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	// Verify booking belongs to this user
	booking, err := h.bookingService.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", err)
	}

	if booking.UserID != user.ID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only upload documents for your own bookings", nil)
	}

	// Get payment amount from form
	amountStr := c.FormValue("amount")
	if amountStr == "" {
		return utils.ValidationErrorResponse(c, "Payment amount is required")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid payment amount format")
	}

	if amount <= 0 {
		return utils.ValidationErrorResponse(c, "Payment amount must be greater than 0")
	}

	// Handle file upload
	file, err := c.FormFile("payment_slip")
	if err != nil {
		return utils.ValidationErrorResponse(c, "Payment slip file is required")
	}

	// Validate file type and size
	if file.Size > MaxImageSize {
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
		Category: "payment_slips",
		EntityID: bookingID.String(),
		FileType: "slip",
	}

	uploadedFile, err := h.fileService.UploadFile(c.Context(), user.ID, file, uploadOptions)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload file to storage", err)
	}

	// อ่าน image data สำหรับ SlipOK verification
	var imageData []byte
	fileReader, err := file.Open()
	if err == nil {
		imageData, _ = io.ReadAll(fileReader)
		fileReader.Close()
	}

	// Create payment slip request with the S3 Storage URL
	req := dto.UploadPaymentSlipRequest{
		BookingID: bookingID,
		Amount:    amount,
		Image:     uploadedFile.URL, // Store the S3 Storage URL
		ImageData: imageData,        // Raw bytes สำหรับ SlipOK
	}

	// Upload payment slip through service
	paymentSlip, err := h.bookingService.UploadPaymentSlip(c.Context(), req)
	if err != nil {
		// Clean up uploaded file on error by deleting from S3 Storage
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to process payment slip", err)
	}

	// Enhanced response with file information
	response := fiber.Map{
		"payment_slip": paymentSlip,
		"file_info": fiber.Map{
			"file_id":   uploadedFile.ID,
			"url":       uploadedFile.URL,
			"cdn_path":  uploadedFile.CDNPath,
			"file_name": uploadedFile.FileName,
			"file_size": uploadedFile.FileSize,
			"mime_type": uploadedFile.MimeType,
		},
	}

	return utils.SuccessResponse(c, "Payment slip uploaded successfully to S3 Storage", response)
}

func (h *AgentHandler) UploadPassport(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	// Verify booking belongs to this user
	booking, err := h.bookingService.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", err)
	}

	if booking.UserID != user.ID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only upload documents for your own bookings", nil)
	}

	// Handle file upload
	file, err := c.FormFile("passport")
	if err != nil {
		return utils.ValidationErrorResponse(c, "Passport file is required")
	}

	// Validate file type and size
	if file.Size > MaxImageSize {
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

	// Upload file to S3 Storage using FileService first
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
		// Clean up uploaded file on error
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err)
	}
	defer src.Close()

	// Read file data
	fileData, err := io.ReadAll(src)
	if err != nil {
		// Clean up uploaded file on error
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read file", err)
	}

	// Get passport number from form (optional, will be extracted by OCR if not provided)
	passportNumber := c.FormValue("passport_number")
	if passportNumber == "" {
		passportNumber = "AUTO_EXTRACT" // Placeholder, will be replaced by OCR/LLM
	}

	// Create UploadPassportRequest for the booking service with S3 Storage URL
	req := dto.UploadPassportRequest{
		BookingID:      bookingID,
		PassportNumber: passportNumber,
		ImageData:      fileData,
		FileName:       file.Filename,
		ImageURL:       uploadedFile.URL, // Add S3 Storage URL
	}

	// Process passport through booking service (includes Vision API + LLM processing)
	passportResponse, err := h.bookingService.UploadPassport(c.Context(), req)
	if err != nil {
		// Clean up uploaded file on error
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to process passport", err)
	}

	// Enhanced response with OCR and LLM results
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
		"message": "Passport uploaded and processed successfully",
	}

	// Add processing status information
	if passportResponse.Status == "processed" {
		response["processing_status"] = "success"
		response["llm_extraction"] = "enabled"
		response["vision_api"] = "enabled"
	} else if passportResponse.Status == "error" {
		response["processing_status"] = "failed"
		response["llm_extraction"] = "error"
		response["vision_api"] = "error"
	} else {
		response["processing_status"] = "pending"
		response["llm_extraction"] = "disabled"
		response["vision_api"] = "disabled"
	}

	return utils.SuccessResponse(c, "Passport uploaded and processed successfully", response)
}

func (h *AgentHandler) GetMyBookingDocuments(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	// Verify booking belongs to this user
	booking, err := h.bookingService.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", err)
	}

	if booking.UserID != user.ID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only view documents for your own bookings", nil)
	}

	// Get documents from database
	documents, err := h.bookingService.GetBookingDocuments(c.Context(), bookingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve booking documents", err)
	}

	return utils.SuccessResponse(c, "Booking documents retrieved successfully", documents)
}
