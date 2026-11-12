package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Create creates a new user (admin only)
func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
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

	user, err := h.userService.CreateUser(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "สร้างผู้ใช้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "สร้างผู้ใช้เรียบร้อยแล้ว", user)
}

// GetProfile returns current user profile
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	profile, err := h.userService.GetUserWithProfile(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบข้อมูลผู้ใช้", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลโปรไฟล์สำเร็จ", profile)
}

// GetByID returns user by ID
func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบผู้ใช้", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลผู้ใช้สำเร็จ", user)
}

// GetByEmail returns user by email (for profile lookup)
func (h *UserHandler) GetByEmail(c *fiber.Ctx) error {
	var req dto.UserProfileRequest
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

	user, err := h.userService.GetUserByEmail(c.Context(), req.Email)
	if err != nil {
		return c.JSON(fiber.Map{
			"status": false,
		})
	}

	return c.JSON(fiber.Map{
		"status": true,
		"data":   user,
	})
}

// List returns paginated users list (admin only)
func (h *UserHandler) List(c *fiber.Ctx) error {
	// Parse query parameters
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	search := c.Query("search")
	role := c.Query("role")
	status := c.Query("status")

	// Validate limits
	if offset < 0 {
		offset = 0
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	listReq := dto.UserListRequest{
		Offset: offset,
		Limit:  limit,
	}
	if search != "" {
		listReq.Search = &search
	}
	if role != "" {
		listReq.Role = &role
	}
	if status != "" {
		listReq.Status = &status
	}

	result, err := h.userService.ListUsers(c.Context(), listReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการผู้ใช้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการผู้ใช้สำเร็จ", result)
}

// Update updates user information
func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	var req dto.UpdateUserRequest
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

	user, err := h.userService.UpdateUser(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตข้อมูลผู้ใช้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตข้อมูลผู้ใช้เรียบร้อยแล้ว", user)
}

// UpdateStatus updates user status (admin only)
func (h *UserHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	var req dto.UpdateUserStatusRequest
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

	user, err := h.userService.UpdateUserStatus(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตสถานะผู้ใช้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตสถานะผู้ใช้เรียบร้อยแล้ว", user)
}

// Delete deletes user (admin only)
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	err = h.userService.DeleteUser(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ลบผู้ใช้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ลบผู้ใช้เรียบร้อยแล้ว", nil)
}

// GetStats returns user statistics (admin only)
func (h *UserHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.userService.GetUserStats(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงสถิติผู้ใช้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงสถิติผู้ใช้สำเร็จ", stats)
}

// ApproveUser approves pending user (admin only)
func (h *UserHandler) ApproveUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	user, err := h.userService.ApproveUser(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อนุมัติผู้ใช้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อนุมัติผู้ใช้เรียบร้อยแล้ว", user)
}

// SuspendUser suspends user (admin only)
func (h *UserHandler) SuspendUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	user, err := h.userService.SuspendUser(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ระงับผู้ใช้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ระงับผู้ใช้เรียบร้อยแล้ว", user)
}

// ReactivateUser reactivates suspended user (admin only)
func (h *UserHandler) ReactivateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	user, err := h.userService.ReactivateUser(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "เปิดใช้งานผู้ใช้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "เปิดใช้งานผู้ใช้เรียบร้อยแล้ว", user)
}

// UpdateProfile updates user profile
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.UpdateUserProfileRequest
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

	// For now, return success with basic profile update confirmation
	// This would be implemented based on business requirements
	updatedProfile := fiber.Map{
		"user_id":    user.ID,
		"email":      user.Email,
		"role":       user.Role,
		"updated_at": time.Now(),
	}

	return utils.SuccessResponse(c, "อัปเดตโปรไฟล์เรียบร้อยแล้ว", updatedProfile)
}

// ChangeRole changes user role (admin only)
func (h *UserHandler) ChangeRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid user ID")
	}

	var req dto.ChangeUserRoleRequest
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

	user, err := h.userService.ChangeUserRole(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "เปลี่ยนบทบาทผู้ใช้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "เปลี่ยนบทบาทผู้ใช้เรียบร้อยแล้ว", user)
}