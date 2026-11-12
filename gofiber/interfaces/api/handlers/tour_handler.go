package handlers

import (
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TourHandler struct {
	tourService services.TourService
}

func NewTourHandler(tourService services.TourService) *TourHandler {
	return &TourHandler{
		tourService: tourService,
	}
}

// Create creates a new tour (admin only)
func (h *TourHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateTourRequest
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

	tour, err := h.tourService.CreateTour(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "สร้างทัวร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "สร้างทัวร์เรียบร้อยแล้ว", tour)
}

// GetByID returns tour by ID
func (h *TourHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tour ID")
	}

	tour, err := h.tourService.GetTourByID(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบทัวร์", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลทัวร์สำเร็จ", tour)
}

// List returns paginated tours list
func (h *TourHandler) List(c *fiber.Ctx) error {
	// Parse pagination (frontend ส่ง page, backend แปลงเป็น offset)
	pagination := utils.ParsePaginationParams(c, 10)
	search := c.Query("search")
	status := c.Query("status")
	cityID := c.Query("city_id")

	listReq := dto.TourListRequest{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
	if search != "" {
		listReq.Search = &search
	}
	if status != "" {
		listReq.Status = &status
	}
	if cityID != "" {
		if cityUUID, err := uuid.Parse(cityID); err == nil {
			listReq.CityID = &cityUUID
		}
	}

	result, err := h.tourService.ListTours(c.Context(), listReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการทัวร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการทัวร์สำเร็จ", result)
}

// ListWithPeriods returns tours with their periods included
func (h *TourHandler) ListWithPeriods(c *fiber.Ctx) error {
	// Parse pagination (frontend ส่ง page, backend แปลงเป็น offset)
	pagination := utils.ParsePaginationParams(c, 10)
	search := c.Query("search")
	status := c.Query("status", "publish") // Default to published tours
	cityID := c.Query("city_id")
	availableOnly := c.Query("available_only") == "true"

	listReq := dto.TourListRequest{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
	if search != "" {
		listReq.Search = &search
	}
	if status != "" {
		listReq.Status = &status
	}
	if cityID != "" {
		if cityUUID, err := uuid.Parse(cityID); err == nil {
			listReq.CityID = &cityUUID
		}
	}

	// Call service method to get tours with periods
	result, err := h.tourService.ListToursWithPeriods(c.Context(), listReq, availableOnly)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการทัวร์พร้อม periods ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการทัวร์พร้อม periods สำเร็จ", result)
}

// ListSpecialTours returns tours with periods that have sale rate type or commission type
func (h *TourHandler) ListSpecialTours(c *fiber.Ctx) error {
	// Parse pagination (frontend ส่ง page, backend แปลงเป็น offset)
	pagination := utils.ParsePaginationParams(c, 10)
	search := c.Query("search")
	status := c.Query("status", "publish") // Default to published tours
	cityID := c.Query("city_id")
	availableOnly := c.Query("available_only") == "true"

	listReq := dto.TourListRequest{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
	if search != "" {
		listReq.Search = &search
	}
	if status != "" {
		listReq.Status = &status
	}
	if cityID != "" {
		if cityUUID, err := uuid.Parse(cityID); err == nil {
			listReq.CityID = &cityUUID
		}
	}

	// Call service method to get special tours (with sale rate/commission)
	result, err := h.tourService.ListSpecialTours(c.Context(), listReq, availableOnly)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการทัวร์พิเศษไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการทัวร์พิเศษสำเร็จ", result)
}

// Update updates tour information (admin only)
func (h *TourHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tour ID")
	}

	var req dto.UpdateTourRequest
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

	tour, err := h.tourService.UpdateTour(c.Context(), id, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตทัวร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตทัวร์เรียบร้อยแล้ว", tour)
}

// UpdateStatus updates tour status (admin only)
func (h *TourHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tour ID")
	}

	var req dto.UpdateTourStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Set the ID from URL parameter
	req.ID = id

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	tour, err := h.tourService.UpdateTourStatus(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "อัปเดตสถานะทัวร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตสถานะทัวร์เรียบร้อยแล้ว", tour)
}

// Delete deletes tour (admin only)
func (h *TourHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tour ID")
	}

	err = h.tourService.DeleteTour(c.Context(), id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ลบทัวร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ลบทัวร์เรียบร้อยแล้ว", nil)
}

// GetStats returns tour statistics (admin only)
func (h *TourHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.tourService.GetTourStats(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงสถิติทัวร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงสถิติทัวร์สำเร็จ", stats)
}

// GetPeriods returns available periods for a tour with pagination
func (h *TourHandler) GetPeriods(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid tour ID")
	}

	// Parse pagination (frontend ส่ง page, backend แปลงเป็น offset)
	pagination := utils.ParsePaginationParams(c, 20)

	periods, err := h.tourService.GetTourPeriodsWithPagination(c.Context(), id, pagination.Offset, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงช่วงเวลาทัวร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงช่วงเวลาทัวร์สำเร็จ", periods)
}

// GetPublicTours returns public tours for agents/customers to browse (read-only)
func (h *TourHandler) GetPublicTours(c *fiber.Ctx) error {
	// Parse pagination (frontend ส่ง page, backend แปลงเป็น offset)
	pagination := utils.ParsePaginationParams(c, 20)
	search := c.Query("search", "")
	cityID := c.Query("city_id", "")

	// Create request for public tours (only active/published tours)
	req := dto.TourListRequest{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
		Search: &search,
		Status: stringPtr("active"), // Only show active tours to customers
	}

	if cityID != "" {
		if cityUUID, err := uuid.Parse(cityID); err == nil {
			req.CityID = &cityUUID
		}
	}

	tours, err := h.tourService.ListTours(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการทัวร์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการทัวร์สาธารณะสำเร็จ", tours)
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
