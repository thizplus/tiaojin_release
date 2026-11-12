package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupUserRoutes(api fiber.Router, h *handlers.Handlers) {
	users := api.Group("/users")

	// Public routes (for profile lookup)
	users.Post("/by-email", h.UserHandler.GetByEmail)

	// Protected routes (all authenticated users can access own profile)
	protected := users.Group("", middleware.Protected())
	protected.Get("/profile", h.UserHandler.GetProfile)
	protected.Put("/profile", h.UserHandler.UpdateProfile) // Users can update own profile

	// Super admin only routes (full user management)
	super := users.Group("", middleware.Protected(), middleware.SuperAdminOnly())
	super.Delete("/:id", h.UserHandler.Delete)               // Only super can delete users
	super.Post("/:id/change-role", h.UserHandler.ChangeRole) // Only super can change roles

	// Management routes (admin+ can manage users but not delete or change roles)
	admin := users.Group("", middleware.Protected(), middleware.ManagementOnly())
	admin.Get("/", h.UserHandler.List)
	admin.Get("/stats", h.UserHandler.GetStats)
	admin.Get("/:id", h.UserHandler.GetByID)
	admin.Put("/:id", h.UserHandler.Update)
	admin.Put("/:id/status", h.UserHandler.UpdateStatus)
	admin.Post("/:id/approve", h.UserHandler.ApproveUser)
	admin.Post("/:id/suspend", h.UserHandler.SuspendUser)
	admin.Post("/:id/reactivate", h.UserHandler.ReactivateUser)
}