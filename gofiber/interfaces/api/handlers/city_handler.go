package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type CityHandler struct {
	cityService services.CityService
}

func NewCityHandler(cityService services.CityService) *CityHandler {
	return &CityHandler{
		cityService: cityService,
	}
}

// Create creates a new city (admin only)
func (h *CityHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateCityRequest
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

	city, err := h.cityService.CreateCity(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "สร้างเมืองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "สร้างเมืองเรียบร้อยแล้ว", city)
}

// GetByID returns city by ID
func (h *CityHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid city ID")
	}

	city, err := h.cityService.GetCityByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบเมือง", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลเมืองสำเร็จ", city)
}

// List returns paginated cities list
func (h *CityHandler) List(c *fiber.Ctx) error {
	// Parse pagination (frontend ส่ง page, backend แปลงเป็น offset)
	pagination := utils.ParsePaginationParams(c, 10)
	search := c.Query("search")
	status := c.Query("status")
	stat := c.Query("stat") == "true"

	listReq := dto.CityListRequest{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
	if search != "" {
		listReq.Search = &search
	}
	if status != "" {
		listReq.Status = &status
	}
	if stat {
		listReq.Stat = &stat
	}

	result, err := h.cityService.ListCities(c.Context(), listReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการเมืองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการเมืองสำเร็จ", result)
}

// Update updates city information (admin only)
func (h *CityHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid city ID")
	}

	var req dto.UpdateCityRequest
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

	city, err := h.cityService.UpdateCity(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตเมืองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตเมืองเรียบร้อยแล้ว", city)
}

// UpdateStatus updates city status (admin only)
func (h *CityHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid city ID")
	}

	var req dto.UpdateCityStatusRequest
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

	city, err := h.cityService.UpdateCityStatus(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตสถานะเมืองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตสถานะเมืองเรียบร้อยแล้ว", city)
}

// Delete deletes city (admin only)
func (h *CityHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid city ID")
	}

	err = h.cityService.DeleteCity(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ลบเมืองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ลบเมืองเรียบร้อยแล้ว", nil)
}

// GetTours returns tours in a city
func (h *CityHandler) GetTours(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid city ID")
	}

	// Parse pagination (frontend ส่ง page, backend แปลงเป็น offset)
	pagination := utils.ParsePaginationParams(c, 10)

	tours, err := h.cityService.GetCityTours(c.Context(), id, pagination.Offset, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงทัวร์ในเมืองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงทัวร์ในเมืองสำเร็จ", tours)
}