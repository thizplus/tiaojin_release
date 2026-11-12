package routes

import (
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupSalesRoutes(api fiber.Router, h *handlers.Handlers) {
	sales := api.Group("/sales")

	// All sales routes require authentication
	sales.Use(middleware.Protected())

	// Public route for agents to get active sales list (agent+ access)
	sales.Get("/active", middleware.AgentOnly(), h.SalesHandler.ListActiveSales) // GET /api/v1/sales/active

	// Management routes for sales administration (admin, super only)
	adminSales := sales.Group("/admin")
	adminSales.Use(middleware.ManagementOnly())
	adminSales.Get("/", h.SalesHandler.List)                       // List all sales personnel
	adminSales.Post("/", h.SalesHandler.Create)                    // Create new sales person
	adminSales.Get("/:id", h.SalesHandler.GetByID)                 // Get sales person details
	adminSales.Put("/:id", h.SalesHandler.Update)                  // Update sales person
	adminSales.Delete("/:id", h.SalesHandler.Delete)               // Delete sales person
	adminSales.Patch("/:id/activate", h.SalesHandler.Activate)     // Activate sales person
	adminSales.Patch("/:id/deactivate", h.SalesHandler.Deactivate) // Deactivate sales person
	adminSales.Get("/stats", h.SalesHandler.GetStats)              // System-wide sales statistics
	adminSales.Get("/export", h.SalesHandler.ExportReport)         // Export system-wide sales report

	// Sales user routes (sales role only) - read-only access to own data
	salesUser := sales.Group("/")
	salesUser.Use(middleware.RequireRole("sales"))
	salesUser.Get("/bookings", h.SalesHandler.GetMyCustomerBookings)              // View bookings assigned to this sales person (same format as admin/agent)
	salesUser.Get("/my-customers/bookings", h.SalesHandler.GetMyCustomerBookings) // Legacy endpoint - redirect to /bookings
	salesUser.Get("/my-performance", h.SalesHandler.GetMyPerformance)             // View own performance metrics
	salesUser.Get("/my-commission", h.SalesHandler.GetMyCommission)               // View own commission data

	// Document upload for bookings assigned to this sales person
	salesUser.Post("/bookings/:id/payment-slip", h.SalesHandler.UploadPaymentSlip) // POST /api/v1/sales/bookings/:id/payment-slip
	salesUser.Post("/bookings/:id/passport", h.SalesHandler.UploadPassport)        // POST /api/v1/sales/bookings/:id/passport
}
