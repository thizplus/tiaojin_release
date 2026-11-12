package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

// APIKeyAuth creates middleware that authenticates requests using API keys
// This is an alternative to JWT authentication for external integrations
func APIKeyAuth(apiKeyService services.APIKeyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Check for API key in X-API-Key header
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			// Also check Authorization header with "Bearer" prefix
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "ApiKey ") {
				apiKey = strings.TrimPrefix(authHeader, "ApiKey ")
			}
		}

		if apiKey == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "API key is required. Provide it in X-API-Key header or Authorization header as 'ApiKey <key>'",
			})
		}

		// Validate API key
		result, err := apiKeyService.ValidateAPIKey(c.Context(), apiKey)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Failed to validate API key",
				"error":   err.Error(),
			})
		}

		if !result.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired API key",
			})
		}

		// Store user info in context for later use
		c.Locals("user_id", result.UserID)
		c.Locals("api_key_id", result.KeyID)
		c.Locals("api_key_permissions", result.Permissions)
		c.Locals("auth_method", "api_key")

		return c.Next()
	}
}

// APIKeyOrJWT creates middleware that accepts either API key or JWT authentication
func APIKeyOrJWT(apiKeyService services.APIKeyService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try API key first
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "ApiKey ") {
				apiKey = strings.TrimPrefix(authHeader, "ApiKey ")
			}
		}

		// If API key is provided, validate it
		if apiKey != "" {
			result, err := apiKeyService.ValidateAPIKey(c.Context(), apiKey)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"success": false,
					"message": "Failed to validate API key",
					"error":   err.Error(),
				})
			}

			if !result.Valid {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "Invalid or expired API key",
				})
			}

			// Store user info in context
			c.Locals("user_id", result.UserID)
			c.Locals("api_key_id", result.KeyID)
			c.Locals("api_key_permissions", result.Permissions)
			c.Locals("auth_method", "api_key")

			return c.Next()
		}

		// If no API key, try JWT
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Authentication required. Provide either X-API-Key header or JWT token",
			})
		}

		// JWT is valid, store auth method
		c.Locals("auth_method", "jwt")
		c.Locals("user", user)

		return c.Next()
	}
}

// RequirePermission checks if the authenticated API key has the required permission
func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authMethod := c.Locals("auth_method")

		// If authenticated via JWT, allow access (JWT users have full permissions)
		if authMethod == "jwt" {
			return c.Next()
		}

		// If authenticated via API key, check permissions
		if authMethod == "api_key" {
			permissions, ok := c.Locals("api_key_permissions").([]string)
			if !ok {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"message": "No permissions found",
				})
			}

			// Check if permission exists
			hasPermission := false
			for _, p := range permissions {
				if p == permission || p == "*" {
					hasPermission = true
					break
				}
			}

			if !hasPermission {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"message": "Insufficient permissions. Required: " + permission,
				})
			}

			return c.Next()
		}

		// No valid authentication method
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Authentication required",
		})
	}
}