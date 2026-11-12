package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
)

func SetupStorageRoutes(api fiber.Router, h *handlers.Handlers) {
	storage := api.Group("/storage")

	// Public endpoint - no auth required for serving files
	// Files are accessed via presigned URLs which provide their own security
	storage.Get("/url", h.StorageHandler.GetPresignedURL)  // Get presigned URL as JSON
	storage.Get("/*", h.StorageHandler.ServeFile)          // Redirect to presigned URL
}
