package routes

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/interfaces/api/middleware"
)

func SetupDebugRoutes(api fiber.Router, h *handlers.Handlers) {
	debug := api.Group("/debug")

	// Apply authentication middleware for security
	debug.Use(middleware.Protected())

	// Vision API debug endpoints - only available in development
	vision := debug.Group("/vision")
	vision.Post("/ocr", h.DebugHandler.DebugVisionOCR)             // Test OCR functionality
	vision.Get("/config", h.DebugHandler.DebugVisionConfig)        // Check Vision API config
	vision.Post("/parse", h.DebugHandler.DebugVisionParse)         // Test passport parsing logic
	vision.Post("/full", h.DebugHandler.DebugVisionFullProcess)    // Test complete pipeline
}