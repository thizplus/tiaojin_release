package routes

import (
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupBookingRoutes(api fiber.Router, h *handlers.Handlers) {
	bookings := api.Group("/bookings", middleware.Protected())

	// All booking routes in order of specificity (specific paths first, dynamic params last)

	// Admin-only routes
	bookings.Get("/stats", middleware.ManagementOnly(), h.BookingHandler.GetStats)                    // GET /api/v1/bookings/stats
	bookings.Get("/report", middleware.ManagementOnly(), h.BookingHandler.GetAdminReport)             // GET /api/v1/bookings/report?year=2024&month=4

	// Sales+ routes
	bookings.Get("/commission", middleware.SalesManagement(), h.BookingHandler.GetCommissionData)     // GET /api/v1/bookings/commission
	bookings.Get("/sales-report", middleware.SalesManagement(), h.BookingHandler.GetSalesReport)      // GET /api/v1/bookings/sales-report?year=2024&month=4

	// Agent+ routes (reports)
	bookings.Get("/agent-report", middleware.BookingAccess(), h.BookingHandler.GetAgentReport)        // GET /api/v1/bookings/agent-report?year=2024&month=4

	// Agent+ routes (basic CRUD)
	bookings.Post("/", middleware.BookingAccess(), h.BookingHandler.Create)                           // POST /api/v1/bookings
	bookings.Get("/", middleware.BookingAccess(), h.BookingHandler.List)                              // GET /api/v1/bookings

	// Dynamic routes (must be last to avoid conflicts)
	bookings.Get("/:id", middleware.BookingAccess(), h.BookingHandler.GetByID)                        // GET /api/v1/bookings/:id
	bookings.Put("/:id", middleware.BookingAccess(), h.BookingHandler.Update)                         // PUT /api/v1/bookings/:id
	bookings.Put("/:id/status", middleware.ManagementOnly(), h.BookingHandler.UpdateStatus)          // PUT /api/v1/bookings/:id/status
	bookings.Delete("/:id", middleware.BookingAccess(), h.BookingHandler.Cancel)                      // DELETE /api/v1/bookings/:id
}