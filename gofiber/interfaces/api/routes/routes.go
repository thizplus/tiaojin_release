package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/domain/services"
	"gofiber-template/interfaces/api/handlers"
)

func SetupRoutes(app *fiber.App, h *handlers.Handlers, apiKeyService services.APIKeyService) {
	// Setup health and root routes
	SetupHealthRoutes(app, h)

	// API version group
	api := app.Group("/api/v1")

	// Setup all route groups
	SetupAuthRoutes(api, h)
	SetupUserRoutes(api, h)
	SetupProfileRoutes(api, h)  // NEW: User profile routes
	SetupCompanyRoutes(api, h)
	SetupCityRoutes(api, h)
	SetupBannerRoutes(api, h)
	SetupSalesRoutes(api, h)   // Sales routes
	SetupTourRoutes(api, h)    // Tour routes
	SetupPeriodRoutes(api, h)  // Period routes
	SetupBookingRoutes(api, h) // Booking routes
	SetupAdminRoutes(api, h)           // NEW: Admin management routes
	SetupAgentRoutes(api, h)           // NEW: Agent operation routes
	SetupSeatAvailabilityRoutes(api, h) // NEW: Seat availability routes
	SetupTaskRoutes(api, h)
	SetupFileRoutes(api, h)
	SetupStorageRoutes(api, h) // Storage file access with presigned URLs
	SetupJobRoutes(api, h)
	SetupConfigRoutes(api, h)      // NEW: System configuration routes (Super admin only)
	SetupBankAccountRoutes(api, h) // NEW: Company bank account routes
	SetupAPIKeyRoutes(api, h)      // NEW: API Key management routes
	SetupPublicRoutes(api, h, apiKeyService) // NEW: Public API routes (API Key auth)
	SetupDebugRoutes(api, h) // Debug routes for development

	// Setup WebSocket routes (needs app, not api group)
	SetupWebSocketRoutes(app)
}