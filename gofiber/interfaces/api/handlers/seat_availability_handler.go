package handlers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type SeatAvailabilityHandler struct {
	seatAvailabilityService services.SeatAvailabilityService
}

func NewSeatAvailabilityHandler(seatAvailabilityService services.SeatAvailabilityService) *SeatAvailabilityHandler {
	return &SeatAvailabilityHandler{
		seatAvailabilityService: seatAvailabilityService,
	}
}

// GetSeatAvailability ดึงข้อมูลที่นั่งว่างสำหรับ period เดียว
func (h *SeatAvailabilityHandler) GetSeatAvailability(c *fiber.Ctx) error {
	log.Printf("🎯 GET /seat-availability/:periodId called")

	// ดึง period ID จาก URL parameter
	periodIDStr := c.Params("periodId")
	if periodIDStr == "" {
		return utils.ValidationErrorResponse(c, "Period ID is required")
	}

	// แปลง string เป็น UUID
	periodID, err := uuid.Parse(periodIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid period ID format")
	}

	// คำนวณจำนวนที่นั่งว่าง
	availability, err := h.seatAvailabilityService.CalculateAvailability(periodID)
	if err != nil {
		log.Printf("❌ Error calculating seat availability: %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to calculate seat availability", err)
	}

	log.Printf("✅ Seat availability calculated successfully for period %s", periodID)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Seat availability retrieved successfully",
		"data":    availability,
	})
}

// GetMultipleSeatAvailability ดึงข้อมูลที่นั่งว่างสำหรับหลาย periods
func (h *SeatAvailabilityHandler) GetMultipleSeatAvailability(c *fiber.Ctx) error {
	log.Printf("🎯 POST /seat-availability/multiple called")

	// รับ request body
	var req struct {
		PeriodIDs []string `json:"period_ids" validate:"required,min=1"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// แปลง string array เป็น UUID array
	periodIDs := make([]uuid.UUID, 0, len(req.PeriodIDs))
	for _, idStr := range req.PeriodIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return utils.ValidationErrorResponse(c, "Invalid period ID format: "+idStr)
		}
		periodIDs = append(periodIDs, id)
	}

	// คำนวณจำนวนที่นั่งว่างสำหรับทุก periods
	availabilities, err := h.seatAvailabilityService.GetAvailabilityForPeriods(periodIDs)
	if err != nil {
		log.Printf("❌ Error calculating multiple seat availabilities: %v", err)
		return utils.InternalServerErrorResponse(c, "Failed to calculate seat availabilities", err)
	}

	log.Printf("✅ Multiple seat availabilities calculated successfully for %d periods", len(periodIDs))

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Seat availabilities retrieved successfully",
		"data":    availabilities,
	})
}