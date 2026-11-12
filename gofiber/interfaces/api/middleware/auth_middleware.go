package middleware

import (
	"gofiber-template/pkg/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Protected middleware validates JWT tokens and sets user context
func Protected() fiber.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.UnauthorizedResponse(c, "Missing authorization header")
		}

		// Extract token from header
		token := utils.ExtractTokenFromHeader(authHeader)
		if token == "" {
			return utils.UnauthorizedResponse(c, "Invalid authorization header format")
		}

		// Validate token and get user context
		userCtx, err := utils.ValidateTokenStringToUUID(token, jwtSecret)
		if err != nil {
			log.Printf("❌ Token validation failed: %v", err)
			switch err {
			case utils.ErrExpiredToken:
				return utils.UnauthorizedResponse(c, "Token has expired")
			case utils.ErrInvalidToken:
				return utils.UnauthorizedResponse(c, "Invalid token")
			case utils.ErrMissingToken:
				return utils.UnauthorizedResponse(c, "Missing token")
			default:
				return utils.UnauthorizedResponse(c, "Token validation failed")
			}
		}

		log.Printf("✅ Token validated for user: %s (%s)", userCtx.Email, userCtx.ID)

		// Set user context in fiber locals
		c.Locals("user", userCtx)

		return c.Next()
	}
}

// RequireRole middleware checks if user has specific role
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			log.Printf("❌ RequireRole: User not authenticated: %v", err)
			return utils.UnauthorizedResponse(c, "User not authenticated")
		}

		log.Printf("🔍 RequireRole: User role='%s', Required role='%s'", user.Role, role)

		if user.Role != role {
			log.Printf("❌ RequireRole: Access denied - user role '%s' != required role '%s'", user.Role, role)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Insufficient permissions",
				"error":   "Access denied",
			})
		}

		log.Printf("✅ RequireRole: Access granted for role '%s'", user.Role)
		return c.Next()
	}
}

// SuperAdminOnly middleware ensures only super users can access (full system access)
func SuperAdminOnly() fiber.Handler {
	return RequireRole("super")
}

// AdminOnly middleware ensures only admin+ users can access (admin and super)
func AdminOnly() fiber.Handler {
	return RequireRoles([]string{"super", "admin"})
}

// AgentOnly middleware ensures only agent+ users can access (agent, admin, super)
func AgentOnly() fiber.Handler {
	return RequireRoles([]string{"super", "admin", "agent"})
}

// SalesOnly middleware ensures only sales+ users can access (sales, agent, admin, super)
func SalesOnly() fiber.Handler {
	return RequireRoles([]string{"super", "admin", "agent", "sales"})
}

// ManagementOnly middleware for super and admin roles only
func ManagementOnly() fiber.Handler {
	log.Printf("🏢 ManagementOnly middleware created - allowing roles: [super, admin]")
	return RequireRoles([]string{"super", "admin"})
}

// BookingAccess middleware for roles that can manage bookings (agent+)
func BookingAccess() fiber.Handler {
	return RequireRoles([]string{"super", "admin", "agent", "sales"})
}

// SalesManagement middleware for roles that can view sales data
func SalesManagement() fiber.Handler {
	return RequireRoles([]string{"super", "admin", "agent", "sales"})
}

// RequireRoles middleware checks if user has any of the specified roles
func RequireRoles(roles []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			log.Printf("❌ RequireRoles: User not authenticated: %v", err)
			return utils.UnauthorizedResponse(c, "User not authenticated")
		}

		log.Printf("🔍 RequireRoles: User role='%s', Required roles=%v", user.Role, roles)

		hasValidRole := false
		for _, role := range roles {
			log.Printf("🔍 RequireRoles: Comparing user.Role='%s' with required='%s'", user.Role, role)
			if user.Role == role {
				hasValidRole = true
				log.Printf("✅ RequireRoles: Role match found: '%s'", role)
				break
			}
		}

		if !hasValidRole {
			log.Printf("❌ RequireRoles: Access denied - user role '%s' not in allowed roles %v", user.Role, roles)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Insufficient permissions",
				"error":   "Access denied",
			})
		}

		log.Printf("✅ RequireRoles: Access granted for role '%s'", user.Role)
		return c.Next()
	}
}

// RequireActiveStatus middleware ensures user status is active
func RequireActiveStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			return utils.UnauthorizedResponse(c, "User not authenticated")
		}

		if user.Status != "active" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "บัญชีของคุณไม่ได้รับการอนุมัติหรือถูกระงับ",
				"error":   "Account not active",
			})
		}

		return c.Next()
	}
}

// OwnerOnly middleware checks if user is the owner of the resource
func OwnerOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.GetUserFromContext(c)
		if err != nil {
			return utils.UnauthorizedResponse(c, "User not authenticated")
		}

		c.Locals("requireOwnership", true)
		c.Locals("ownerUserID", user.ID)

		return c.Next()
	}
}

// Optional middleware that doesn't require authentication but sets user context if token is present
func Optional() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Next()
		}

		token := utils.ExtractTokenFromHeader(authHeader)
		if token == "" {
			return c.Next()
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		userCtx, err := utils.ValidateTokenStringToUUID(token, jwtSecret)
		if err != nil {
			return c.Next()
		}

		c.Locals("user", userCtx)
		return c.Next()
	}
}
