package middleware

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/pkg/utils"
)

// ValidateUUID middleware validates UUID parameters
func ValidateUUID(paramName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params(paramName)
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Missing required parameter: " + paramName,
				"error":   "missing_parameter",
			})
		}

		if !utils.IsValidUUID(id) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid " + paramName + " format",
				"error":   "invalid_uuid",
			})
		}

		return c.Next()
	}
}

// ValidateContentType middleware ensures correct Content-Type for POST/PUT requests
func ValidateContentType() fiber.Handler {
	return func(c *fiber.Ctx) error {
		method := c.Method()
		if method == "POST" || method == "PUT" || method == "PATCH" {
			contentType := c.Get("Content-Type")
			if contentType == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"success": false,
					"message": "Content-Type header is required",
					"error":   "missing_content_type",
				})
			}

			// Allow JSON content type (the most common for APIs)
			if contentType != "application/json" &&
			   !utils.Contains(contentType, "application/json") &&
			   !utils.Contains(contentType, "multipart/form-data") {
				return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
					"success": false,
					"message": "Unsupported Content-Type. Use application/json or multipart/form-data",
					"error":   "unsupported_content_type",
				})
			}
		}

		return c.Next()
	}
}

// ValidateRequestSize middleware limits request body size
func ValidateRequestSize(maxSize int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Body()) > maxSize {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"success": false,
				"message": "Request body too large",
				"error":   "request_too_large",
			})
		}

		return c.Next()
	}
}

// RequireFields middleware ensures specific fields are present in the request
func RequireFields(fields []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body map[string]interface{}
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Invalid request body",
				"error":   "invalid_json",
			})
		}

		missingFields := []string{}
		for _, field := range fields {
			if _, exists := body[field]; !exists {
				missingFields = append(missingFields, field)
			}
		}

		if len(missingFields) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Missing required fields",
				"error":   "missing_fields",
				"missing": missingFields,
			})
		}

		return c.Next()
	}
}