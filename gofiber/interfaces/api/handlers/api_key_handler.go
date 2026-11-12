package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type APIKeyHandler struct {
	apiKeyService services.APIKeyService
}

func NewAPIKeyHandler(apiKeyService services.APIKeyService) *APIKeyHandler {
	return &APIKeyHandler{
		apiKeyService: apiKeyService,
	}
}

// Generate creates a new API key for the authenticated user
func (h *APIKeyHandler) Generate(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CreateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	apiKey, err := h.apiKeyService.GenerateAPIKey(c.Context(), user.ID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to generate API key", err)
	}

	return utils.SuccessResponse(c, "API key created successfully. Please save this key as it will not be shown again.", apiKey)
}

// List returns all API keys for the authenticated user
func (h *APIKeyHandler) List(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Parse query parameters
	pagination := utils.ParsePaginationParams(c, 20)
	search := c.Query("search")
	showExpired := c.Query("show_expired") == "true"

	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		if val, err := strconv.ParseBool(activeStr); err == nil {
			isActive = &val
		}
	}

	req := dto.ListAPIKeysRequest{
		Offset:      pagination.Offset,
		Limit:       pagination.Limit,
		IsActive:    isActive,
		Search:      search,
		ShowExpired: showExpired,
	}

	result, err := h.apiKeyService.ListAPIKeys(c.Context(), user.ID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list API keys", err)
	}

	return utils.SuccessResponse(c, "API keys retrieved successfully", result)
}

// GetByID returns details of a specific API key
func (h *APIKeyHandler) GetByID(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	keyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid API key ID")
	}

	apiKey, err := h.apiKeyService.GetAPIKey(c.Context(), user.ID, keyID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "API key not found", err)
	}

	return utils.SuccessResponse(c, "API key retrieved successfully", apiKey)
}

// Update updates API key settings
func (h *APIKeyHandler) Update(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	keyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid API key ID")
	}

	var req dto.UpdateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	apiKey, err := h.apiKeyService.UpdateAPIKey(c.Context(), user.ID, keyID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update API key", err)
	}

	return utils.SuccessResponse(c, "API key updated successfully", apiKey)
}

// Revoke revokes (soft deletes) an API key
func (h *APIKeyHandler) Revoke(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	keyID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid API key ID")
	}

	if err := h.apiKeyService.RevokeAPIKey(c.Context(), user.ID, keyID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to revoke API key", err)
	}

	return utils.SuccessResponse(c, "API key revoked successfully", nil)
}

// Validate validates an API key (for testing purposes)
func (h *APIKeyHandler) Validate(c *fiber.Ctx) error {
	var req dto.ValidateAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	result, err := h.apiKeyService.ValidateAPIKey(c.Context(), req.APIKey)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to validate API key", err)
	}

	if !result.Valid {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Invalid API key", nil)
	}

	return utils.SuccessResponse(c, "API key is valid", result)
}