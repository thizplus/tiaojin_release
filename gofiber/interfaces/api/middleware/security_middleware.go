package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

// Security applies various security headers
func Security() fiber.Handler {
	return helmet.New(helmet.Config{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		ReferrerPolicy:        "strict-origin-when-cross-origin",
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none';",
	})
}

// BlockSuspiciousRequests middleware blocks requests with suspicious patterns
func BlockSuspiciousRequests() fiber.Handler {
	suspiciousPatterns := []string{
		"../",
		"..\\",
		"<script",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"eval(",
		"document.cookie",
		"union select",
		"drop table",
		"insert into",
		"delete from",
		"update set",
	}

	return func(c *fiber.Ctx) error {
		// Check URL path
		path := strings.ToLower(c.Path())
		for _, pattern := range suspiciousPatterns {
			if strings.Contains(path, pattern) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"message": "Suspicious request detected",
					"error":   "blocked_request",
				})
			}
		}

		// Check query parameters
		c.Context().QueryArgs().VisitAll(func(key, value []byte) {
			queryValue := strings.ToLower(string(value))
			for _, pattern := range suspiciousPatterns {
				if strings.Contains(queryValue, pattern) {
					c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"success": false,
						"message": "Suspicious request detected",
						"error":   "blocked_request",
					})
					return
				}
			}
		})

		// Check request body for POST/PUT requests
		if c.Method() == "POST" || c.Method() == "PUT" || c.Method() == "PATCH" {
			bodyStr := strings.ToLower(string(c.Body()))
			for _, pattern := range suspiciousPatterns {
				if strings.Contains(bodyStr, pattern) {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"success": false,
						"message": "Suspicious request detected",
						"error":   "blocked_request",
					})
				}
			}
		}

		return c.Next()
	}
}

// RequireHTTPS middleware redirects HTTP to HTTPS in production
func RequireHTTPS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Protocol() == "http" && c.Get("X-Forwarded-Proto") != "https" {
			// In production, redirect to HTTPS
			httpsURL := "https://" + c.Hostname() + c.OriginalURL()
			return c.Redirect(httpsURL, fiber.StatusMovedPermanently)
		}
		return c.Next()
	}
}

// ValidateUserAgent middleware blocks requests with suspicious or missing User-Agent
func ValidateUserAgent() fiber.Handler {
	blockedUserAgents := []string{
		"sqlmap",
		"nmap",
		"nikto",
		"masscan",
		"zap",
		"burp",
		"gobuster",
		"dirb",
		"dirbuster",
	}

	return func(c *fiber.Ctx) error {
		userAgent := strings.ToLower(c.Get("User-Agent"))

		// Block completely missing user agent
		if userAgent == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "User-Agent header is required",
				"error":   "missing_user_agent",
			})
		}

		// Block known malicious user agents
		for _, blocked := range blockedUserAgents {
			if strings.Contains(userAgent, blocked) {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"message": "Access denied",
					"error":   "blocked_user_agent",
				})
			}
		}

		return c.Next()
	}
}