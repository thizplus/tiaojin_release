package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupSeatAvailabilityRoutes(api fiber.Router, h *handlers.Handlers) {
	// สร้าง handler สำหรับ seat availability
	seatHandler := h.SeatAvailabilityHandler

	// Seat availability routes (ใช้ optional authentication)
	seatRoutes := api.Group("/seat-availability")

	// GET /api/v1/seat-availability/:periodId - ดูจำนวนที่นั่งว่างสำหรับ period เดียว
	seatRoutes.Get("/:periodId", middleware.Optional(), seatHandler.GetSeatAvailability)

	// POST /api/v1/seat-availability/multiple - ดูจำนวนที่นั่งว่างสำหรับหลาย periods
	seatRoutes.Post("/multiple", middleware.Optional(), seatHandler.GetMultipleSeatAvailability)
}