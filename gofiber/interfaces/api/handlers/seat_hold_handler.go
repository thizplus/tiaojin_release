package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type SeatHoldHandler struct {
	seatHoldService services.SeatHoldService
}

func NewSeatHoldHandler(seatHoldService services.SeatHoldService) *SeatHoldHandler {
	return &SeatHoldHandler{
		seatHoldService: seatHoldService,
	}
}

// CreateSeatHold creates a new seat hold for a period
// POST /admin/periods/{id}/hold-seats
func (h *SeatHoldHandler) CreateSeatHold(c *fiber.Ctx) error {
	// Get admin user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		log.Printf("Failed to get user from context: %v", err)
		return utils.UnauthorizedResponse(c, "Admin authentication required")
	}

	log.Printf("Admin user authenticated: ID=%s, Email=%s, Role=%s", user.ID, user.Email, user.Role)

	// Get period ID from URL
	periodID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid period ID")
	}

	// Parse request body
	var req dto.CreateSeatHoldRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Set period ID from URL
	req.PeriodID = periodID

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Add admin ID to context for service layer
	ctx := context.WithValue(c.Context(), "admin_id", user.ID)
	log.Printf("Added admin ID to context: %s", user.ID)

	// Create seat hold
	seatHold, err := h.seatHoldService.CreateSeatHold(ctx, req)
	if err != nil {
		log.Printf("Error creating seat hold: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create seat hold", err)
	}

	return utils.SuccessResponse(c, "Seat hold created successfully", seatHold)
}

// GetSeatHold retrieves a specific seat hold
// GET /admin/seat-holds/{hold_id}
func (h *SeatHoldHandler) GetSeatHold(c *fiber.Ctx) error {
	// Get hold ID from URL
	holdID, err := uuid.Parse(c.Params("hold_id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid hold ID")
	}

	// Get seat hold
	seatHold, err := h.seatHoldService.GetSeatHoldByID(c.Context(), holdID)
	if err != nil {
		log.Printf("Error getting seat hold: %v", err)
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Seat hold not found", err)
	}

	return utils.SuccessResponse(c, "Seat hold retrieved successfully", seatHold)
}

// GetPeriodSeatHolds retrieves all seat holds for a period
// GET /admin/periods/{id}/seat-holds
func (h *SeatHoldHandler) GetPeriodSeatHolds(c *fiber.Ctx) error {
	// Get period ID from URL
	periodID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid period ID")
	}

	// Get query parameter for including inactive holds
	includeInactive := c.Query("include_inactive", "false") == "true"

	// Get seat holds
	seatHolds, err := h.seatHoldService.GetSeatHoldsByPeriod(c.Context(), periodID, includeInactive)
	if err != nil {
		log.Printf("Error getting period seat holds: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve seat holds", err)
	}

	return utils.SuccessResponse(c, "Period seat holds retrieved successfully", fiber.Map{
		"period_id":   periodID,
		"seat_holds":  seatHolds,
		"total":       len(seatHolds),
		"include_inactive": includeInactive,
	})
}

// UpdateSeatHold updates an existing seat hold
// PATCH /admin/seat-holds/{hold_id}
func (h *SeatHoldHandler) UpdateSeatHold(c *fiber.Ctx) error {
	// Get admin user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		log.Printf("Failed to get user from context: %v", err)
		return utils.UnauthorizedResponse(c, "Admin authentication required")
	}

	// Get hold ID from URL
	holdID, err := uuid.Parse(c.Params("hold_id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid hold ID")
	}

	// Parse request body
	var req dto.UpdateSeatHoldRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Add admin ID to context for service layer
	ctx := context.WithValue(c.Context(), "admin_id", user.ID)

	// Update seat hold
	seatHold, err := h.seatHoldService.UpdateSeatHold(ctx, holdID, req)
	if err != nil {
		log.Printf("Error updating seat hold: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update seat hold", err)
	}

	return utils.SuccessResponse(c, "Seat hold updated successfully", seatHold)
}

// ReleaseSeatHold releases a seat hold
// DELETE /admin/seat-holds/{hold_id}/release
func (h *SeatHoldHandler) ReleaseSeatHold(c *fiber.Ctx) error {
	// Get admin user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		log.Printf("Failed to get user from context: %v", err)
		return utils.UnauthorizedResponse(c, "Admin authentication required")
	}

	// Get hold ID from URL
	holdID, err := uuid.Parse(c.Params("hold_id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid hold ID")
	}

	// Get optional reason from request body or query
	type ReleaseRequest struct {
		Reason string `json:"reason,omitempty"`
	}

	var req ReleaseRequest
	if err := c.BodyParser(&req); err != nil {
		// If body parsing fails, try query parameter
		req.Reason = c.Query("reason", "Manual release")
	}

	// Add admin ID to context for service layer
	ctx := context.WithValue(c.Context(), "admin_id", user.ID)

	// Release seat hold
	seatHold, err := h.seatHoldService.ReleaseSeatHold(ctx, holdID, req.Reason)
	if err != nil {
		log.Printf("Error releasing seat hold: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to release seat hold", err)
	}

	return utils.SuccessResponse(c, "Seat hold released successfully", seatHold)
}

// ExtendSeatHold extends the expiry time of a seat hold
// POST /admin/seat-holds/{hold_id}/extend
func (h *SeatHoldHandler) ExtendSeatHold(c *fiber.Ctx) error {
	// Get hold ID from URL
	holdID, err := uuid.Parse(c.Params("hold_id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid hold ID")
	}

	// Parse request body
	type ExtendRequest struct {
		Hours   *int `json:"hours,omitempty"`
		Minutes *int `json:"minutes,omitempty"`
		Days    *int `json:"days,omitempty"`
	}

	var req ExtendRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Calculate total duration
	totalMinutes := 0
	if req.Days != nil {
		totalMinutes += *req.Days * 24 * 60
	}
	if req.Hours != nil {
		totalMinutes += *req.Hours * 60
	}
	if req.Minutes != nil {
		totalMinutes += *req.Minutes
	}

	if totalMinutes <= 0 {
		return utils.ValidationErrorResponse(c, "Extension duration must be positive")
	}

	duration := time.Duration(totalMinutes) * time.Minute

	// Get admin user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		log.Printf("Failed to get user from context: %v", err)
		return utils.UnauthorizedResponse(c, "Admin authentication required")
	}

	// Add admin ID to context for service layer
	ctx := context.WithValue(c.Context(), "admin_id", user.ID)

	// Extend seat hold
	seatHold, err := h.seatHoldService.ExtendSeatHold(ctx, holdID, duration)
	if err != nil {
		log.Printf("Error extending seat hold: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to extend seat hold", err)
	}

	return utils.SuccessResponse(c, fmt.Sprintf("Seat hold extended by %s", duration.String()), seatHold)
}

// CheckSeatAvailability checks seat availability for a potential hold
// GET /admin/periods/{id}/seat-availability
func (h *SeatHoldHandler) CheckSeatAvailability(c *fiber.Ctx) error {
	// Get period ID from URL
	periodID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid period ID")
	}

	// Get requested seats from query parameter
	requestedSeatsStr := c.Query("seats", "1")
	requestedSeats, err := strconv.Atoi(requestedSeatsStr)
	if err != nil || requestedSeats <= 0 {
		return utils.ValidationErrorResponse(c, "Invalid seats parameter")
	}

	// Check availability
	availability, err := h.seatHoldService.CheckSeatAvailability(c.Context(), periodID, requestedSeats)
	if err != nil {
		log.Printf("Error checking seat availability: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check seat availability", err)
	}

	return utils.SuccessResponse(c, "Seat availability checked successfully", availability)
}

// GetActiveSeatHolds retrieves all active seat holds with pagination
// GET /admin/seat-holds
func (h *SeatHoldHandler) GetActiveSeatHolds(c *fiber.Ctx) error {
	// Parse pagination parameters
	pagination := utils.ParsePaginationParams(c, 20)

	// Get active seat holds
	result, err := h.seatHoldService.GetActiveSeatHolds(c.Context(), pagination.Offset, pagination.Limit)
	if err != nil {
		log.Printf("Error getting active seat holds: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve active seat holds", err)
	}

	return utils.SuccessResponse(c, "Active seat holds retrieved successfully", result)
}

// GetAdminSeatHolds retrieves seat holds for a specific admin with pagination
// GET /admin/admins/{admin_id}/seat-holds
func (h *SeatHoldHandler) GetAdminSeatHolds(c *fiber.Ctx) error {
	// Get admin ID from URL
	adminID, err := uuid.Parse(c.Params("admin_id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid admin ID")
	}

	// Parse pagination parameters
	pagination := utils.ParsePaginationParams(c, 20)

	// Get admin seat holds
	result, err := h.seatHoldService.GetSeatHoldsByAdmin(c.Context(), adminID, pagination.Offset, pagination.Limit)
	if err != nil {
		log.Printf("Error getting admin seat holds: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve admin seat holds", err)
	}

	return utils.SuccessResponse(c, "Admin seat holds retrieved successfully", result)
}

// ExpireOldSeatHolds manually triggers expiry of old seat holds
// POST /admin/seat-holds/expire-old
func (h *SeatHoldHandler) ExpireOldSeatHolds(c *fiber.Ctx) error {
	// This should typically be called by a background job, but exposing it for manual triggers
	result, err := h.seatHoldService.ExpireOldSeatHolds(c.Context())
	if err != nil {
		log.Printf("Error expiring old seat holds: %v", err)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to expire old seat holds", err)
	}

	return utils.SuccessResponse(c, "Old seat holds expired successfully", result)
}

// DeleteSeatHold permanently deletes a seat hold (super admin only)
// DELETE /admin/seat-holds/{hold_id}
func (h *SeatHoldHandler) DeleteSeatHold(c *fiber.Ctx) error {
	// Get admin user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		log.Printf("Failed to get user from context: %v", err)
		return utils.UnauthorizedResponse(c, "Admin authentication required")
	}

	// Get hold ID from URL
	holdID, err := uuid.Parse(c.Params("hold_id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid hold ID")
	}

	// Add admin ID to context for service layer
	ctx := context.WithValue(c.Context(), "admin_id", user.ID)

	// Delete seat hold
	err = h.seatHoldService.DeleteSeatHold(ctx, holdID)
	if err != nil {
		log.Printf("Error deleting seat hold: %v", err)
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Failed to delete seat hold", err)
	}

	return utils.SuccessResponse(c, "Seat hold deleted successfully", fiber.Map{
		"hold_id": holdID,
		"deleted": true,
	})
}

// ConvertSeatHoldToBooking converts a seat hold to a booking
// POST /admin/seat-holds/{hold_id}/convert-to-booking
func (h *SeatHoldHandler) ConvertSeatHoldToBooking(c *fiber.Ctx) error {
	// Get hold ID from URL
	holdID, err := uuid.Parse(c.Params("hold_id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid hold ID")
	}

	// Parse booking request
	var bookingReq dto.CreateBookingRequest
	if err := c.BodyParser(&bookingReq); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking request body")
	}

	// Validate booking request
	if err := utils.ValidateStruct(&bookingReq); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Booking validation failed",
			"errors":  errors,
		})
	}

	// Convert seat hold to booking
	seatHold, err := h.seatHoldService.ConvertSeatHoldToBooking(c.Context(), holdID, bookingReq)
	if err != nil {
		log.Printf("Error converting seat hold to booking: %v", err)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to convert seat hold to booking", err)
	}

	return utils.SuccessResponse(c, "Seat hold converted to booking successfully", seatHold)
}