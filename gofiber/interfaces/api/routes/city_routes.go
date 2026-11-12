package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupCityRoutes(api fiber.Router, h *handlers.Handlers) {
	cities := api.Group("/cities")

	// Public endpoints
	cities.Get("/", h.CityHandler.List)           // Get all cities with optional status filter
	cities.Get("/active", h.CityHandler.List) // Get only active/published cities (will filter by status=publish)
	cities.Get("/:id", h.CityHandler.GetByID)       // Get city by ID

	// Management only endpoints (super, admin can manage cities)
	adminCities := cities.Group("/")
	adminCities.Use(middleware.Protected())
	adminCities.Use(middleware.ManagementOnly()) // super, admin only

	adminCities.Post("/", h.CityHandler.Create)           // Create new city
	adminCities.Put("/:id", h.CityHandler.Update)        // Update city
	adminCities.Delete("/:id", h.CityHandler.Delete)     // Delete city
	adminCities.Put("/:id/status", h.CityHandler.UpdateStatus) // Update city status only
}