package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type AuthHandler struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthHandler(authService services.AuthService, userService services.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
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

	result, err := h.authService.Register(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Registration failed", err)
	}

	return utils.SuccessResponse(c, "ลงทะเบียนสำเร็จ กรุณาเข้าสู่ระบบ", result)
}

// Login handles user authentication
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
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

	result, err := h.authService.Login(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "เข้าสู่ระบบไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "เข้าสู่ระบบสำเร็จ", result)
}

// ValidateToken validates JWT token
func (h *AuthHandler) ValidateToken(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  false,
			"message": "Authorization header is missing",
		})
	}

	// Remove "Bearer " prefix
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	req := dto.ValidateTokenRequest{Token: token}
	result, err := h.authService.ValidateToken(c.Context(), req.Token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
		})
	}

	return c.JSON(fiber.Map{
		"status":  true,
		"decoded": result,
	})
}

// RefreshToken refreshes JWT token
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	result, err := h.authService.RefreshToken(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Token refresh failed", err)
	}

	return utils.SuccessResponse(c, "Token refreshed successfully", result)
}

// ResetPassword handles password reset
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordRequest
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

	err := h.authService.ResetPassword(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รีเซ็ตรหัสผ่านไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "รีเซ็ตรหัสผ่านสำเร็จ", nil)
}

// UpdatePassword handles password change
func (h *AuthHandler) UpdatePassword(c *fiber.Ctx) error {
	// Get authenticated user from context
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Determine target user_id and check if it's admin changing other user's password
	targetUserID := user.ID // Default: change own password
	isChangingOther := false

	// If user_id is provided in request
	if req.UserID != nil {
		// Only admin/super can change other user's password
		if user.Role != "admin" && user.Role != "super" {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "Only admin and super users can change other users' passwords", nil)
		}
		targetUserID = *req.UserID
		isChangingOther = (*req.UserID != user.ID)
	}

	// Validate: Must provide old_password when changing own password (regardless of role)
	// Admin/Super don't need old_password only when changing OTHER user's password
	isAdmin := user.Role == "admin" || user.Role == "super"
	if !isChangingOther && req.OldPassword == "" {
		// Changing own password - old_password required
		return utils.ValidationErrorResponse(c, "old_password is required when changing your own password")
	} else if isChangingOther && !isAdmin {
		// Non-admin trying to change other's password - should never reach here due to check above
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Only admin and super users can change other users' passwords", nil)
	}

	// Set the target user ID
	req.UserID = &targetUserID

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	// Pass current user's role to service for authorization logic
	err = h.authService.UpdatePassword(c.Context(), req, user.Role)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "เปลี่ยนรหัสผ่านไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "เปลี่ยนรหัสผ่านเรียบร้อยแล้ว", nil)
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	err = h.authService.Logout(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Logout failed", err)
	}

	return utils.SuccessResponse(c, "ออกจากระบบเรียบร้อยแล้ว", nil)
}

// RevokeAllTokens revokes all user tokens
func (h *AuthHandler) RevokeAllTokens(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	err = h.authService.RevokeAllTokens(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to revoke tokens", err)
	}

	return utils.SuccessResponse(c, "ยกเลิกการเข้าสู่ระบบทั้งหมดแล้ว", nil)
}