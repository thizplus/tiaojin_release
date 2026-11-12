package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type CompanyBankAccountHandler struct {
	service services.CompanyBankAccountService
}

func NewCompanyBankAccountHandler(service services.CompanyBankAccountService) *CompanyBankAccountHandler {
	return &CompanyBankAccountHandler{service: service}
}

// Create สร้างบัญชีธนาคารใหม่ (admin only)
func (h *CompanyBankAccountHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCompanyBankAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}
	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "errors": errors})
	}

	result, err := h.service.Create(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create bank account", err)
	}

	return utils.SuccessResponse(c, "Bank account created successfully", result)
}

// GetByID ดึงข้อมูลบัญชีธนาคาร
func (h *CompanyBankAccountHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid bank account ID")
	}

	result, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return utils.NotFoundResponse(c, "Bank account not found")
	}

	return utils.SuccessResponse(c, "Bank account retrieved successfully", result)
}

// Update อัปเดตบัญชีธนาคาร (admin only)
func (h *CompanyBankAccountHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid bank account ID")
	}

	var req dto.UpdateCompanyBankAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	result, err := h.service.Update(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update bank account", err)
	}

	return utils.SuccessResponse(c, "Bank account updated successfully", result)
}

// Delete ลบบัญชีธนาคาร (admin only)
func (h *CompanyBankAccountHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid bank account ID")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete bank account", err)
	}

	return utils.SuccessResponse(c, "Bank account deleted successfully", nil)
}

// List ดึงบัญชีทั้งหมด (admin only)
func (h *CompanyBankAccountHandler) List(c *fiber.Ctx) error {
	result, err := h.service.List(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list bank accounts", err)
	}

	return utils.SuccessResponse(c, "Bank accounts retrieved successfully", result)
}

// ListActive ดึงเฉพาะบัญชีที่ใช้งาน (public — agent ใช้เลือกบัญชีโอนเงิน)
func (h *CompanyBankAccountHandler) ListActive(c *fiber.Ctx) error {
	result, err := h.service.ListActive(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list bank accounts", err)
	}

	return utils.SuccessResponse(c, "Active bank accounts retrieved successfully", result)
}
