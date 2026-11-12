package routes

import (
	"gofiber-template/domain/services"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetupPublicRoutes sets up public API routes that require API Key authentication
func SetupPublicRoutes(api fiber.Router, h *handlers.Handlers, apiKeyService services.APIKeyService) {
	// Public API group - requires API Key authentication
	public := api.Group("/public")

	// Apply API Key authentication middleware
	public.Use(middleware.APIKeyAuth(apiKeyService))

	// Tours endpoints
	public.Get("/tours", h.PublicTourHandler.ListAvailableTours)       // GET /api/v1/public/tours
	public.Get("/tours/:id", h.PublicTourHandler.GetByID)              // GET /api/v1/public/tours/:id
}