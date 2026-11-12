package routes

import (
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupTourRoutes(api fiber.Router, h *handlers.Handlers) {
	tours := api.Group("/tours")

	// Public routes (no authentication required)
	tours.Get("/", h.TourHandler.List)                  // GET /api/v1/tours?status=publish&page=1&limit=20
	tours.Get("/regular", h.TourHandler.ListWithPeriods) // GET /api/v1/tours/regular - tours with periods included
	tours.Get("/special", h.TourHandler.ListSpecialTours) // GET /api/v1/tours/special - tours with sale rate/commission periods
	tours.Get("/:id", h.TourHandler.GetByID)            // GET /api/v1/tours/:id
	tours.Get("/:id/periods", h.TourHandler.GetPeriods) // GET /api/v1/tours/:id/periods

	// Admin only routes (require admin authentication)
	adminTours := tours.Group("/", middleware.Protected(), middleware.AdminOnly())
	{
		adminTours.Post("/", h.TourHandler.Create)                // POST /api/v1/tours
		adminTours.Get("/stats", h.TourHandler.GetStats)          // GET /api/v1/tours/stats
		adminTours.Put("/:id", h.TourHandler.Update)              // PUT /api/v1/tours/:id
		adminTours.Delete("/:id", h.TourHandler.Delete)           // DELETE /api/v1/tours/:id
		adminTours.Put("/:id/status", h.TourHandler.UpdateStatus) // PUT /api/v1/tours/:id/status

	}
}
