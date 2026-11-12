package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/pkg/utils"
)

func SetupHealthRoutes(app *fiber.App, h *handlers.Handlers) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is running",
			"service": "GoFiber Template API",
		})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Welcome to GoFiber Template API",
			"version": "1.0.0",
			"docs":    "/api/v1",
			"health":  "/health",
		})
	})

	// Test endpoint - simplified booking call for testing
	app.Get("/debug/bookings", func(c *fiber.Ctx) error {
		// Set basic query parameters for testing
		c.Query("offset", "0")
		c.Query("limit", "5")

		// Create a mock user context for testing
		userCtx := &utils.UserContext{
			ID:       uuid.New(),
			Email:    "debug@test.com",
			Role:     "admin",
			Status:   "active",
		}
		c.Locals("user", userCtx)

		// Call the handler List method
		return h.BookingHandler.List(c)
	})

	// Check sales records in database
	app.Get("/debug/sales-list", func(c *fiber.Ctx) error {
		// Create a mock admin user context
		userCtx := &utils.UserContext{
			ID:       uuid.New(),
			Email:    "admin@test.com",
			Role:     "admin",
			Status:   "active",
		}
		c.Locals("user", userCtx)

		// Call the sales handler List method to see all sales records
		return h.SalesHandler.List(c)
	})

	// Test agent endpoint - requires a user that has a sales record
	app.Get("/debug/agent-bookings/:user_id", func(c *fiber.Ctx) error {
		// Get user_id from URL parameter
		userIDStr := c.Params("user_id")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(fiber.Map{
				"success": false,
				"message": "Invalid user_id format",
				"error":   err.Error(),
			})
		}

		// Create test user context with the provided ID
		userCtx := &utils.UserContext{
			ID:       userID,
			Email:    "test-agent@example.com",
			Role:     "agent",
			Status:   "active",
		}
		c.Locals("user", userCtx)

		// Call the agent handler GetBookings method
		return h.AgentHandler.GetBookings(c)
	})




}
