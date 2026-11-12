package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupProfileRoutes(api fiber.Router, h *handlers.Handlers) {
	user := api.Group("/user")

	// All profile routes require authentication
	user.Use(middleware.Protected())

	// Profile management routes (all authenticated users can access their own profile)
	user.Get("/profile", h.UserHandler.GetProfile)       // GET /api/v1/user/profile
	user.Put("/profile", h.UserHandler.UpdateProfile)    // PUT /api/v1/user/profile
}