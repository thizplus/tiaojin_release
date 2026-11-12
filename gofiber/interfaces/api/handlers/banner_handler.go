package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type BannerHandler struct {
	bannerService services.BannerService
}

func NewBannerHandler(bannerService services.BannerService) *BannerHandler {
	return &BannerHandler{
		bannerService: bannerService,
	}
}

// Create creates a new banner (admin only)
func (h *BannerHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateBannerRequest
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

	banner, err := h.bannerService.CreateBanner(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "สร้างแบนเนอร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "สร้างแบนเนอร์เรียบร้อยแล้ว", banner)
}

// GetByID returns banner by ID
func (h *BannerHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid banner ID")
	}

	banner, err := h.bannerService.GetBannerByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบแบนเนอร์", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลแบนเนอร์สำเร็จ", banner)
}

// List returns paginated banners list
func (h *BannerHandler) List(c *fiber.Ctx) error {
	// Parse pagination (frontend ส่ง page, backend แปลงเป็น offset)
	pagination := utils.ParsePaginationParams(c, 10)
	search := c.Query("search")
	status := c.Query("status")
	position := c.Query("position")

	listReq := dto.BannerListRequest{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
	if search != "" {
		listReq.Search = &search
	}
	if status != "" {
		listReq.Status = &status
	}
	if position != "" {
		listReq.Position = &position
	}

	result, err := h.bannerService.ListBanners(c.Context(), listReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการแบนเนอร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการแบนเนอร์สำเร็จ", result)
}

// GetActive returns active banners for public display
func (h *BannerHandler) GetActive(c *fiber.Ctx) error {
	position := c.Query("position")

	banners, err := h.bannerService.GetActiveBanners(c.Context(), position)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงแบนเนอร์ที่ใช้งานไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงแบนเนอร์ที่ใช้งานสำเร็จ", banners)
}

// Update updates banner information (admin only)
func (h *BannerHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid banner ID")
	}

	var req dto.UpdateBannerRequest
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

	banner, err := h.bannerService.UpdateBanner(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตแบนเนอร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตแบนเนอร์เรียบร้อยแล้ว", banner)
}

// UpdateStatus updates banner status (admin only)
func (h *BannerHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid banner ID")
	}

	var req dto.UpdateBannerStatusRequest
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

	banner, err := h.bannerService.UpdateBannerStatus(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตสถานะแบนเนอร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตสถานะแบนเนอร์เรียบร้อยแล้ว", banner)
}

// UpdateOrder updates banner display order (admin only)
func (h *BannerHandler) UpdateOrder(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid banner ID")
	}

	var req dto.UpdateBannerOrderRequest
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

	banner, err := h.bannerService.UpdateBannerOrder(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตลำดับแบนเนอร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตลำดับแบนเนอร์เรียบร้อยแล้ว", banner)
}

// UpdateBulkOrder updates multiple banner display orders (admin only)
func (h *BannerHandler) UpdateBulkOrder(c *fiber.Ctx) error {
	var req dto.UpdateBannerBulkOrderRequest
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

	// Convert bulk request to the format expected by service
	orderReq := dto.UpdateBannerOrderRequest{
		BannerOrders: req.BannerOrders,
	}
	err := h.bannerService.UpdateOrder(c.Context(), orderReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตลำดับแบนเนอร์จำนวนมากไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตลำดับแบนเนอร์จำนวนมากเรียบร้อยแล้ว", nil)
}

// Delete deletes banner (admin only)
func (h *BannerHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid banner ID")
	}

	err = h.bannerService.DeleteBanner(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ลบแบนเนอร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ลบแบนเนอร์เรียบร้อยแล้ว", nil)
}