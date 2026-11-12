package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupBannerRoutes(api fiber.Router, h *handlers.Handlers) {
	banners := api.Group("/banners")

	// Public endpoints (no authentication required)
	banners.Get("/", h.BannerHandler.GetActive)         // Get active banners for public
	banners.Get("/active", h.BannerHandler.GetActive)   // Alternative endpoint
	banners.Get("/:id", h.BannerHandler.GetByID)        // Get banner by ID (public)

	// Admin only endpoints (require authentication and admin role)
	adminBanners := banners.Group("/")
	adminBanners.Use(middleware.Protected())
	adminBanners.Use(middleware.AdminOnly())

	// These routes are handled by admin routes now
	// Keep minimal routes here for compatibility
}