package handlers

import (
	"log"
	"time"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PeriodHandler struct {
	periodService services.PeriodService
}

func NewPeriodHandler(periodService services.PeriodService) *PeriodHandler {
	return &PeriodHandler{
		periodService: periodService,
	}
}

// Create creates a new period (admin only)
func (h *PeriodHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePeriodRequest
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

	period, err := h.periodService.CreatePeriod(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "สร้างช่วงเวลาไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "สร้างช่วงเวลาเรียบร้อยแล้ว", period)
}

// GetByID returns period by ID
func (h *PeriodHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid period ID")
	}

	period, err := h.periodService.GetPeriodByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบช่วงเวลา", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลช่วงเวลาสำเร็จ", period)
}

// List returns paginated periods list
func (h *PeriodHandler) List(c *fiber.Ctx) error {
	// Parse query parameters - support both direct params and params[...] format
	pagination := utils.ParsePaginationParams(c, 10)
	tourID := getQueryParam(c, "tour_id", "")
	status := getQueryParam(c, "status", "")
	dateFrom := getQueryParam(c, "date_from", "")
	dateTo := getQueryParam(c, "date_to", "")

	listReq := dto.PeriodListRequest{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
	if tourID != "" {
		if tourUUID, err := uuid.Parse(tourID); err == nil {
			listReq.TourID = &tourUUID
		}
	}
	if status != "" {
		listReq.Status = &status
	}
	// Parse dates
	if dateFrom != "" {
		parsedDateFrom, err := parseDateParam(dateFrom)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_from ไม่ถูกต้อง", err)
		}
		listReq.DateFrom = parsedDateFrom
	}
	if dateTo != "" {
		parsedDateTo, err := parseDateParam(dateTo)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_to ไม่ถูกต้อง", err)
		}
		listReq.DateTo = parsedDateTo
	}

	result, err := h.periodService.ListPeriods(c.Context(), listReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการช่วงเวลาไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการช่วงเวลาสำเร็จ", result)
}

// Update updates period information (admin only)
func (h *PeriodHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid period ID")
	}

	var req dto.UpdatePeriodRequest

	// Log raw request body for debugging
	log.Printf("UpdatePeriod raw body: %s", string(c.Body()))

	if err := c.BodyParser(&req); err != nil {
		log.Printf("UpdatePeriod body parsing error: %v", err)
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Log parsed request for debugging
	log.Printf("UpdatePeriod parsed request: %+v", req)

	if err := utils.ValidateStruct(&req); err != nil {
		log.Printf("UpdatePeriod validation error: %v", err)
		errors := utils.GetValidationErrors(err)
		log.Printf("UpdatePeriod validation errors detail: %+v", errors)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	period, err := h.periodService.UpdatePeriod(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตช่วงเวลาไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตช่วงเวลาเรียบร้อยแล้ว", period)
}

// UpdateStatus updates period status (admin only)
func (h *PeriodHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid period ID")
	}

	var req dto.UpdatePeriodStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Set the ID from URL parameter
	req.ID = id

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	period, err := h.periodService.UpdatePeriodStatus(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตสถานะช่วงเวลาไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตสถานะช่วงเวลาเรียบร้อยแล้ว", period)
}

// Delete deletes period (admin only)
func (h *PeriodHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid period ID")
	}

	err = h.periodService.DeletePeriod(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ลบช่วงเวลาไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ลบช่วงเวลาเรียบร้อยแล้ว", nil)
}

// GetAvailability returns period availability
func (h *PeriodHandler) GetAvailability(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid period ID")
	}

	availability, err := h.periodService.GetPeriodAvailability(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงข้อมูลว่างช่วงเวลาไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลว่างช่วงเวลาสำเร็จ", availability)
}

// GetStats returns period statistics
func (h *PeriodHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.periodService.GetPeriodStats(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงสถิติช่วงเวลาไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงสถิติช่วงเวลาสำเร็จ", stats)
}

// DebugBookingStats returns detailed debug information for period booking calculations
func (h *PeriodHandler) DebugBookingStats(c *fiber.Ctx) error {
	log.Printf("[DEBUG] DebugBookingStats endpoint called")

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		log.Printf("[ERROR] Invalid period ID in debug endpoint: %s", c.Params("id"))
		return utils.ValidationErrorResponse(c, "Invalid period ID")
	}

	log.Printf("[DEBUG] Testing real-time booking stats for period: %s", id.String())

	// Test the real-time booking calculation
	stats, err := h.periodService.CalculateRealTimeBookingStats(c.Context(), id)
	if err != nil {
		log.Printf("[ERROR] Debug: Failed to calculate real-time stats: %v", err)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": false,
			"message": "Debug: Failed to calculate real-time booking stats",
			"error":   err.Error(),
			"period_id": id.String(),
		})
	}

	log.Printf("[SUCCESS] Debug: Real-time stats calculated successfully")

	// Also get regular period info for comparison
	period, err := h.periodService.GetPeriodByID(c.Context(), id)
	if err != nil {
		log.Printf("[WARNING] Debug: Failed to get period details: %v", err)
	}

	debugInfo := fiber.Map{
		"success":   true,
		"message":   "Debug information for period booking calculation",
		"period_id": id.String(),
		"real_time_stats": stats,
	}

	if period != nil {
		debugInfo["period_static_data"] = fiber.Map{
			"seat":            period.Seat,
			"static_booking":  period.Booking,
			"available_seats": period.AvailableSeats,
			"status":          period.Status,
		}

		debugInfo["comparison"] = fiber.Map{
			"static_vs_realtime_booking": fiber.Map{
				"static":    period.Booking,
				"realtime":  stats.BookedPassengers,
				"difference": stats.BookedPassengers - period.Booking,
			},
			"static_vs_realtime_available": fiber.Map{
				"static":    period.AvailableSeats,
				"realtime":  stats.AvailableSeats,
				"difference": stats.AvailableSeats - period.AvailableSeats,
			},
		}
	}

	return c.Status(fiber.StatusOK).JSON(debugInfo)
}

// TestDeployment simple endpoint to test if period service deployment is working
func (h *PeriodHandler) TestDeployment(c *fiber.Ctx) error {
	log.Printf("[TEST] TestDeployment endpoint called")

	// Test basic service functionality
	ctx := c.Context()

	// Test period service availability
	if h.periodService == nil {
		log.Printf("[ERROR] PeriodService is nil in handler")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "PeriodService not available",
			"test":    "deployment",
		})
	}

	log.Printf("[SUCCESS] PeriodService is available")

	// Try to get period stats (simple test)
	_, err := h.periodService.GetPeriodStats(ctx)
	if err != nil {
		log.Printf("[WARNING] GetPeriodStats failed: %v", err)
	} else {
		log.Printf("[SUCCESS] GetPeriodStats working")
	}

	// Test if we can create a fake period calculation (with test data)
	testPeriodID := "00000000-0000-0000-0000-000000000000"
	id, parseErr := uuid.Parse(testPeriodID)
	if parseErr == nil {
		_, calcErr := h.periodService.CalculateRealTimeBookingStats(ctx, id)
		if calcErr != nil {
			log.Printf("[EXPECTED] CalculateRealTimeBookingStats failed for test ID (expected): %v", calcErr)
		} else {
			log.Printf("[UNEXPECTED] CalculateRealTimeBookingStats succeeded for test ID")
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success":     true,
		"message":     "Period service deployment test completed",
		"test":        "deployment",
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
		"endpoints": fiber.Map{
			"debug_booking_stats": "/api/v1/periods/{id}/debug-booking-stats",
			"test_deployment":     "/api/v1/periods/test-deployment",
			"availability":        "/api/v1/periods/{id}/availability",
			"stats":               "/api/v1/periods/stats",
		},
		"real_time_features": fiber.Map{
			"booking_calculation": "enabled",
			"fallback_mechanism":  "enabled",
			"debug_logging":       "enabled",
		},
	})
}

// Helper function to get query parameter supporting both direct and params[...] format
func getQueryParam(c *fiber.Ctx, key, defaultValue string) string {
	// Try direct parameter first
	if value := c.Query(key); value != "" {
		return value
	}

	// Try params[key] format
	if value := c.Query("params[" + key + "]"); value != "" {
		return value
	}

	return defaultValue
}
