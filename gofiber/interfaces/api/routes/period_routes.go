package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupPeriodRoutes(api fiber.Router, h *handlers.Handlers) {
	periods := api.Group("/periods")

	// Public routes (no authentication required)
	periods.Get("/", h.PeriodHandler.List)                        // GET /api/v1/periods?status=active&tour_id=xxx&date_from=2024-01-01
	periods.Get("/:id", h.PeriodHandler.GetByID)                  // GET /api/v1/periods/:id
	periods.Get("/:id/availability", h.PeriodHandler.GetAvailability) // GET /api/v1/periods/:id/availability
	periods.Get("/test-deployment", h.PeriodHandler.TestDeployment) // GET /api/v1/periods/test-deployment

	// Admin only routes (require admin authentication)
	// Note: These are legacy routes. Use /api/v1/admin/periods for new admin functionality
	adminPeriods := periods.Group("/", middleware.Protected(), middleware.AdminOnly())
	{
		adminPeriods.Post("/", h.PeriodHandler.Create)             // POST /api/v1/periods (legacy)
		adminPeriods.Put("/:id", h.PeriodHandler.Update)           // PUT /api/v1/periods/:id (legacy)
		adminPeriods.Delete("/:id", h.PeriodHandler.Delete)        // DELETE /api/v1/periods/:id (legacy)
		adminPeriods.Put("/:id/status", h.PeriodHandler.UpdateStatus) // PUT /api/v1/periods/:id/status (legacy)
		adminPeriods.Get("/stats", h.PeriodHandler.GetStats)       // GET /api/v1/periods/stats (legacy)
	}
}

