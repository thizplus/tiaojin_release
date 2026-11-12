package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"gofiber-template/pkg/utils"
)

// RateLimit creates a rate limiting middleware
func RateLimit(max int, duration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: duration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as key for rate limiting
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Rate limit exceeded, please try again later",
				"error":   "too_many_requests",
			})
		},
	})
}

// AuthRateLimit applies rate limiting specifically for authentication endpoints
func AuthRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5, // 5 requests per minute for auth endpoints
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Too many authentication attempts, please try again in a minute",
				"error":   "auth_rate_limit_exceeded",
			})
		},
	})
}

// APIRateLimit applies general rate limiting for API endpoints
func APIRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        100, // 100 requests per minute for general API
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			// For authenticated users, use user ID + IP
			if user := c.Locals("user"); user != nil {
				if userCtx, ok := user.(*utils.UserContext); ok {
					return userCtx.ID.String() + ":" + c.IP()
				}
			}
			// For unauthenticated requests, use IP only
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "API rate limit exceeded, please slow down",
				"error":   "api_rate_limit_exceeded",
			})
		},
	})
}