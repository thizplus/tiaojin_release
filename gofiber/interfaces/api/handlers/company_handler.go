package handlers

import (
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type CompanyHandler struct {
	companyService services.CompanyService
}

func NewCompanyHandler(companyService services.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

// GetRandomLogos ดึง logo agent แบบ random สำหรับแสดงหน้าแรก (public)
func (h *CompanyHandler) GetRandomLogos(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 20)

	logos, err := h.companyService.GetRandomLogos(c.Context(), limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get logos", err)
	}

	return utils.SuccessResponse(c, "Logos retrieved successfully", fiber.Map{
		"logos": logos,
	})
}

// CreateProfile creates a new company profile
func (h *CompanyHandler) CreateProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.CreateCompanyProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	// Set user ID from context
	req.UserID = user.ID

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	result, err := h.companyService.CreateProfile(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to create company profile", err)
	}

	return utils.SuccessResponse(c, "สร้างโปรไฟล์บริษัทสำเร็จ", result)
}

// GetMyProfile gets current user's company profile
func (h *CompanyHandler) GetMyProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	result, err := h.companyService.GetProfileByUserID(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Company profile not found", err)
	}

	return utils.SuccessResponse(c, "ดึงโปรไฟล์บริษัทสำเร็จ", result)
}

// UpdateProfile updates current user's company profile
func (h *CompanyHandler) UpdateProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	var req dto.UpdateCompanyProfileRequest
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

	result, err := h.companyService.UpdateProfile(c.Context(), user.ID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update company profile", err)
	}

	return utils.SuccessResponse(c, "อัปเดตโปรไฟล์บริษัทสำเร็จ", result)
}

// DeleteProfile deletes current user's company profile
func (h *CompanyHandler) DeleteProfile(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	err = h.companyService.DeleteProfile(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete company profile", err)
	}

	return utils.SuccessResponse(c, "ลบโปรไฟล์บริษัทสำเร็จ", nil)
}

// UploadLogo uploads company logo (supports multipart/form-data)
func (h *CompanyHandler) UploadLogo(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Get file from form
	file, err := c.FormFile("logo")
	if err != nil {
		return utils.ValidationErrorResponse(c, "Logo file is required")
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return utils.ValidationErrorResponse(c, "Only image files are allowed")
	}

	// Validate file size (max 5MB)
	const maxSize = 5 * 1024 * 1024
	if file.Size > maxSize {
		return utils.ValidationErrorResponse(c, "File size exceeds 5MB limit")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err)
	}
	defer src.Close()

	// Read file data
	fileData, err := io.ReadAll(src)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read file", err)
	}

	// Create upload request
	req := dto.UploadLogoRequest{
		UserID:   user.ID,
		LogoData: fileData,
		FileName: file.Filename,
		FileType: contentType,
	}

	result, err := h.companyService.UploadLogo(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to upload logo", err)
	}

	return utils.SuccessResponse(c, "อัปโหลดโลโก้สำเร็จ", result)
}

// DeleteLogo deletes company logo
func (h *CompanyHandler) DeleteLogo(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	result, err := h.companyService.DeleteLogo(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete logo", err)
	}

	return utils.SuccessResponse(c, "ลบโลโก้สำเร็จ", result)
}

// UploadPdfHeaderImage uploads company PDF header image (supports multipart/form-data)
func (h *CompanyHandler) UploadPdfHeaderImage(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Get file from form
	file, err := c.FormFile("pdf_header_image")
	if err != nil {
		return utils.ValidationErrorResponse(c, "PDF header image file is required")
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return utils.ValidationErrorResponse(c, "Only image files are allowed")
	}

	// Validate file size (max 5MB)
	const maxSize = 5 * 1024 * 1024
	if file.Size > maxSize {
		return utils.ValidationErrorResponse(c, "File size exceeds 5MB limit")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err)
	}
	defer src.Close()

	// Read file data
	fileData, err := io.ReadAll(src)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read file", err)
	}

	// Create upload request
	req := dto.UploadPdfHeaderImageRequest{
		UserID:    user.ID,
		ImageData: fileData,
		FileName:  file.Filename,
		FileType:  contentType,
	}

	result, err := h.companyService.UploadPdfHeaderImage(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to upload pdf header image", err)
	}

	return utils.SuccessResponse(c, "อัปโหลดรูปหัวกระดาษ PDF สำเร็จ", result)
}

// DeletePdfHeaderImage deletes company PDF header image
func (h *CompanyHandler) DeletePdfHeaderImage(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	result, err := h.companyService.DeletePdfHeaderImage(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to delete pdf header image", err)
	}

	return utils.SuccessResponse(c, "ลบรูปหัวกระดาษ PDF สำเร็จ", result)
}

// ValidateTaxID validates if tax ID is available (public endpoint)
func (h *CompanyHandler) ValidateTaxID(c *fiber.Ctx) error {
	taxID := c.Params("tax_id")
	if len(taxID) != 13 {
		return utils.ValidationErrorResponse(c, "Tax ID must be 13 digits")
	}

	available, err := h.companyService.ValidateTaxID(c.Context(), taxID, nil)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to validate tax ID", err)
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"available": available,
		"tax_id":    taxID,
	})
}

// Admin endpoints

// ListProfiles lists pending company profiles (admin only)
func (h *CompanyHandler) ListProfiles(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c, 20)
	search := c.Query("search", "")

	req := dto.CompanyProfileListRequest{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}

	if search != "" {
		req.Search = &search
	}

	result, err := h.companyService.ListProfiles(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list company profiles", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการโปรไฟล์บริษัทที่รอการอนุมัติสำเร็จ", result)
}

// SearchProfiles searches company profiles (admin only)
func (h *CompanyHandler) SearchProfiles(c *fiber.Ctx) error {
	query := c.Query("q", "")
	pagination := utils.ParsePaginationParams(c, 20)

	if query == "" {
		return utils.ValidationErrorResponse(c, "Search query is required")
	}

	result, err := h.companyService.SearchProfiles(c.Context(), query, pagination.Offset, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search company profiles", err)
	}

	return utils.SuccessResponse(c, "ค้นหาโปรไฟล์บริษัทสำเร็จ", result)
}

// GetProfileByID gets company profile by ID (admin only)
func (h *CompanyHandler) GetProfileByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	profileID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid profile ID")
	}

	result, err := h.companyService.GetProfileByID(c.Context(), profileID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Company profile not found", err)
	}

	return utils.SuccessResponse(c, "ดึงโปรไฟล์บริษัทสำเร็จ", result)
}

// GetProfileByTaxID gets company profile by tax ID (admin only)
func (h *CompanyHandler) GetProfileByTaxID(c *fiber.Ctx) error {
	taxID := c.Params("tax_id")
	if len(taxID) != 13 {
		return utils.ValidationErrorResponse(c, "Tax ID must be 13 digits")
	}

	result, err := h.companyService.GetProfileByTaxID(c.Context(), taxID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Company profile not found", err)
	}

	return utils.SuccessResponse(c, "ดึงโปรไฟล์บริษัทสำเร็จ", result)
}

// AdminUpdateProfile updates a company profile by ID (admin only — no ownership check)
func (h *CompanyHandler) AdminUpdateProfile(c *fiber.Ctx) error {
	idParam := c.Params("id")
	profileID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid profile ID")
	}

	var req dto.UpdateCompanyProfileRequest
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

	result, err := h.companyService.AdminUpdateProfile(c.Context(), profileID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to update company profile", err)
	}

	return utils.SuccessResponse(c, "อัปเดตโปรไฟล์บริษัทสำเร็จ", result)
}

// ApproveProfile approves a company profile and activates the associated user (admin only)
func (h *CompanyHandler) ApproveProfile(c *fiber.Ctx) error {
	idParam := c.Params("id")
	profileID, err := uuid.Parse(idParam)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid profile ID")
	}

	result, err := h.companyService.ApproveProfile(c.Context(), profileID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to approve company profile", err)
	}

	return utils.SuccessResponse(c, "อนุมัติโปรไฟล์บริษัทสำเร็จ", result)
}

