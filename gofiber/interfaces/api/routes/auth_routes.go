package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupAuthRoutes(api fiber.Router, h *handlers.Handlers) {
	auth := api.Group("/auth")

	// Public authentication endpoints
	auth.Post("/register", h.AuthHandler.Register)
	auth.Post("/login", h.AuthHandler.Login)
	auth.Post("/refresh", h.AuthHandler.RefreshToken)
	auth.Post("/validate", h.AuthHandler.ValidateToken)
	auth.Post("/reset-password", h.AuthHandler.ResetPassword)

	// Protected authentication endpoints (require valid token)
	auth.Post("/logout", middleware.Protected(), h.AuthHandler.Logout)
	auth.Post("/update-password", middleware.Protected(), h.AuthHandler.UpdatePassword)
	auth.Post("/revoke-all", middleware.Protected(), h.AuthHandler.RevokeAllTokens)
}