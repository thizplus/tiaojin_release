package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupCompanyRoutes(api fiber.Router, h *handlers.Handlers) {
	company := api.Group("/company-profiles")

	// Public endpoints
	company.Get("/validate-tax-id/:tax_id", h.CompanyHandler.ValidateTaxID)
	company.Get("/logos", h.CompanyHandler.GetRandomLogos) // GET /api/v1/company-profiles/logos

	// Protected endpoints (require authentication)
	company.Use(middleware.Protected())

	// Profile management
	company.Post("/", h.CompanyHandler.CreateProfile)           // Create profile
	company.Get("/me", h.CompanyHandler.GetMyProfile)           // Get current user's profile
	company.Put("/me", h.CompanyHandler.UpdateProfile)          // Update profile
	company.Delete("/me", h.CompanyHandler.DeleteProfile)       // Delete profile

	// Logo management (form-data support)
	company.Post("/me/logo", h.CompanyHandler.UploadLogo)       // Upload logo (multipart/form-data)
	company.Delete("/me/logo", h.CompanyHandler.DeleteLogo)     // Delete logo
}