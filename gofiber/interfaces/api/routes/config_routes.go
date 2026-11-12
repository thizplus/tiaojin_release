package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupConfigRoutes(api fiber.Router, h *handlers.Handlers) {
	config := api.Group("/config")
	config.Use(middleware.Protected())
	config.Use(middleware.SuperAdminOnly()) // Super admin only

	// Root config endpoint - summary/status
	config.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message":          "System Configuration API",
			"version":          "1.0",
			"expiration_tiers": "/api/v1/config/expiration/tiers",
			"configs":          "/api/v1/config/configs",
			"stats":            "/api/v1/config/expiration/stats",
		})
	})

	// Expiration Tiers Management
	expiration := config.Group("/expiration")
	expiration.Get("/tiers", h.ConfigHandler.GetExpirationTiers)
	expiration.Put("/tiers", h.ConfigHandler.UpdateExpirationTiers)
	expiration.Post("/tiers/validate", h.ConfigHandler.ValidateExpirationConfig)
	expiration.Post("/tiers/calculate", h.ConfigHandler.CalculateTier)
	expiration.Get("/stats", h.ConfigHandler.GetExpirationStats)
	expiration.Post("/process-expired", h.ConfigHandler.ProcessExpiredBookings)

	// General Config Management
	configs := config.Group("/configs")
	configs.Get("/", h.ConfigHandler.ListConfigs)
	configs.Post("/", h.ConfigHandler.CreateConfig)
	configs.Get("/:id", h.ConfigHandler.GetConfig)
	configs.Put("/:id", h.ConfigHandler.UpdateConfig)
	configs.Delete("/:id", h.ConfigHandler.DeleteConfig)

	// Agent Booking Limits Management
	agent := config.Group("/agent-booking-limits")
	agent.Get("/", h.ConfigHandler.GetAgentBookingLimits)
	agent.Put("/", h.ConfigHandler.UpdateAgentBookingLimits)
}