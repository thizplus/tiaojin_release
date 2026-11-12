package routes

import (
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAPIKeyRoutes(api fiber.Router, h *handlers.Handlers) {
	apiKeys := api.Group("/api-keys", middleware.Protected())

	// All API key routes require authentication (agent, sales, admin, super)
	// Only the owner can manage their own API keys

	apiKeys.Post("/", middleware.BookingAccess(), h.APIKeyHandler.Generate)        // POST /api/v1/api-keys
	apiKeys.Get("/", middleware.BookingAccess(), h.APIKeyHandler.List)             // GET /api/v1/api-keys
	apiKeys.Get("/:id", middleware.BookingAccess(), h.APIKeyHandler.GetByID)       // GET /api/v1/api-keys/:id
	apiKeys.Put("/:id", middleware.BookingAccess(), h.APIKeyHandler.Update)        // PUT /api/v1/api-keys/:id
	apiKeys.Delete("/:id", middleware.BookingAccess(), h.APIKeyHandler.Revoke)     // DELETE /api/v1/api-keys/:id

	// Validation endpoint (for testing)
	apiKeys.Post("/validate", h.APIKeyHandler.Validate) // POST /api/v1/api-keys/validate
}