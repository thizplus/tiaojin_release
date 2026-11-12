package routes

import (
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAgentRoutes(api fiber.Router, h *handlers.Handlers) {
	agent := api.Group("/agent")

	// All agent routes require authentication and agent role only
	agent.Use(middleware.Protected())
	agent.Use(middleware.RequireRole("agent"))

	// Profile management (customer-level)
	agent.Get("/profile", h.AgentHandler.GetProfile)               // GET /api/v1/agent/profile - view own profile
	agent.Put("/profile", h.AgentHandler.UpdateProfile)            // PUT /api/v1/agent/profile - edit own profile
	agent.Post("/profile", h.AgentHandler.CreateProfile)           // POST /api/v1/agent/profile - create profile
	agent.Get("/profile/check", h.AgentHandler.CheckProfileExists) // GET /api/v1/agent/profile/check - check if profile exists

	// Tour browsing (read-only, customer-level)
	agent.Get("/tours", h.TourHandler.GetPublicTours) // GET /api/v1/agent/tours - view available tours only

	// Own booking management (customer-level)
	agent.Get("/bookings", h.AgentHandler.GetBookings)     // GET /api/v1/agent/bookings - view own bookings only
	agent.Post("/bookings", h.AgentHandler.CreateBooking)  // POST /api/v1/agent/bookings - create booking
	agent.Delete("/bookings/:id", h.BookingHandler.Cancel) // DELETE /api/v1/agent/bookings/:id - cancel booking (unified)

	// Document upload for own bookings (customer-level)
	agent.Post("/bookings/:id/payment-slip", h.AgentHandler.UploadPaymentSlip) // POST /api/v1/agent/bookings/:id/payment-slip
	agent.Post("/bookings/:id/passport", h.AgentHandler.UploadPassport)        // POST /api/v1/agent/bookings/:id/passport
	agent.Get("/bookings/:id/documents", h.AgentHandler.GetMyBookingDocuments) // GET /api/v1/agent/bookings/:id/documents - view own documents
}
