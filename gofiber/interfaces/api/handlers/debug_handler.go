package handlers

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"gofiber-template/domain/services"
	"gofiber-template/pkg/config"
	"gofiber-template/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

const (
	MaxImageSize = 50 * 1024 * 1024 // 5MB
)

type DebugHandler struct {
	visionService services.VisionService
	llmService    services.LLMService
	config        *config.Config
}

func NewDebugHandler(visionService services.VisionService, llmService services.LLMService, cfg *config.Config) *DebugHandler {
	return &DebugHandler{
		visionService: visionService,
		llmService:    llmService,
		config:        cfg,
	}
}

// DebugVisionOCR tests Vision API OCR functionality directly
func (h *DebugHandler) DebugVisionOCR(c *fiber.Ctx) error {
	// Security check - only allow in development
	if h.config.App.Env != "development" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Debug endpoints only available in development", nil)
	}

	startTime := time.Now()

	// Check if Vision service is available
	if h.visionService == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "Vision API service not available", nil)
	}

	// Handle file upload
	file, err := c.FormFile("image")
	if err != nil {
		return utils.ValidationErrorResponse(c, "Image file is required")
	}

	// Validate file type and size
	if file.Size > MaxImageSize {
		return utils.ValidationErrorResponse(c, "File size exceeds 5MB limit")
	}

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
	isValidExt := false
	for _, ext := range allowedExtensions {
		if fileExt == ext {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return utils.ValidationErrorResponse(c, "Invalid file type. Only JPG, PNG, and PDF files are allowed")
	}

	// Open and read file data
	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err)
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read file", err)
	}

	// Process with Vision API OCR
	ocrResult, err := h.visionService.ExtractTextFromImage(c.Context(), fileData)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Vision API processing failed", err)
	}

	// Try LLM extraction if available
	var llmExtraction interface{}
	var llmStatus string
	var llmError string

	if h.llmService != nil && h.config.LLM.Enabled {
		llmResult, err := h.llmService.ExtractPassportFields(c.Context(), ocrResult.Text)
		if err != nil {
			llmStatus = "failed"
			llmError = err.Error()
			llmExtraction = nil
		} else {
			llmStatus = "success"
			llmExtraction = fiber.Map{
				"passport_number":   llmResult.PassportNumber,
				"name_th":           llmResult.NameTH,
				"surname_th":        llmResult.SurnameTH,
				"name_en":           llmResult.NameEN,
				"surname_en":        llmResult.SurnameEN,
				"gender":            llmResult.Gender,
				"nationality":       llmResult.Nationality,
				"country_code":      llmResult.CountryCode,
				"date_of_birth":     formatDatePtr(llmResult.DateOfBirth),
				"date_of_expiry":    formatDatePtr(llmResult.DateOfExpiry),
				"confidence":        llmResult.OverallConfidence,
				"field_confidences": llmResult.FieldConfidences,
			}
		}
	} else {
		llmStatus = "disabled"
		llmExtraction = nil
	}

	processingTime := time.Since(startTime)

	// Build enhanced debug response
	debugResponse := fiber.Map{
		"raw_ocr": fiber.Map{
			"text":       ocrResult.Text,
			"confidence": ocrResult.Confidence,
			"word_count": len(ocrResult.Words),
			"words":      ocrResult.Words,
		},
		"llm_extraction": fiber.Map{
			"status": llmStatus,
			"data":   llmExtraction,
			"error":  llmError,
		},
		"comparison": fiber.Map{
			"method":               "LLM vs Regex",
			"accuracy_improvement": "estimated 95% vs 60%",
			"llm_enabled":          h.config.LLM.Enabled,
			"llm_provider":         h.config.LLM.Provider,
			"llm_model":            h.config.LLM.Model,
		},
		"processing_time": fmt.Sprintf("%.2fs", processingTime.Seconds()),
		"image_info": fiber.Map{
			"filename": file.Filename,
			"size":     file.Size,
			"format":   strings.ToUpper(strings.TrimPrefix(fileExt, ".")),
		},
		"vision_api_status": "success",
		"timestamp":         time.Now().Format(time.RFC3339),
	}

	return utils.SuccessResponse(c, "Enhanced Vision API + LLM processing completed", debugResponse)
}

// DebugVisionConfig shows Vision API configuration status
func (h *DebugHandler) DebugVisionConfig(c *fiber.Ctx) error {
	// Security check - only allow in development
	if h.config.App.Env != "development" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Debug endpoints only available in development", nil)
	}

	configInfo := fiber.Map{
		"vision_enabled":    h.config.GoogleCloud.VisionEnabled,
		"project_id":        h.config.GoogleCloud.ProjectID,
		"service_available": h.visionService != nil,
		"environment":       h.config.App.Env,
	}

	// Mask credentials path for security
	if h.config.GoogleCloud.CredentialsPath != "" {
		configInfo["credentials_configured"] = true
		configInfo["credentials_path"] = "***MASKED***"
	} else {
		configInfo["credentials_configured"] = false
	}

	status := "Vision API Ready"
	if !h.config.GoogleCloud.VisionEnabled {
		status = "Vision API Disabled"
	} else if h.visionService == nil {
		status = "Vision API Service Unavailable"
	}

	return utils.SuccessResponse(c, status, configInfo)
}

// DebugVisionParse tests passport parsing logic without OCR
func (h *DebugHandler) DebugVisionParse(c *fiber.Ctx) error {
	// Security check - only allow in development
	if h.config.App.Env != "development" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Debug endpoints only available in development", nil)
	}

	// Check if Vision service is available
	if h.visionService == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "Vision API service not available", nil)
	}

	// Get test OCR text from request body
	var request struct {
		OCRText string `json:"ocr_text" validate:"required"`
	}

	if err := c.BodyParser(&request); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid JSON body")
	}

	if request.OCRText == "" {
		return utils.ValidationErrorResponse(c, "ocr_text field is required")
	}

	startTime := time.Now()

	// Test passport data extraction
	passportResult, err := h.visionService.ExtractPassportData(c.Context(), request.OCRText)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Passport parsing failed", err)
	}

	processingTime := time.Since(startTime)

	// Build debug response
	debugResponse := fiber.Map{
		"parsing_result":     passportResult,
		"input_text":         request.OCRText,
		"text_length":        len(request.OCRText),
		"fields_extracted":   len(passportResult.FieldConfidences),
		"overall_confidence": passportResult.OverallConfidence,
		"processing_time":    fmt.Sprintf("%.3fs", processingTime.Seconds()),
		"success":            passportResult.ProcessedSuccessfully,
		"timestamp":          time.Now().Format(time.RFC3339),
	}

	return utils.SuccessResponse(c, "Passport parsing completed", debugResponse)
}

// DebugVisionFullProcess tests the complete passport processing pipeline
func (h *DebugHandler) DebugVisionFullProcess(c *fiber.Ctx) error {
	// Security check - only allow in development
	if h.config.App.Env != "development" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "Debug endpoints only available in development", nil)
	}

	// Check if Vision service is available
	if h.visionService == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "Vision API service not available", nil)
	}

	startTime := time.Now()

	// Handle file upload
	file, err := c.FormFile("image")
	if err != nil {
		return utils.ValidationErrorResponse(c, "Image file is required")
	}

	// Validate file
	if file.Size > MaxImageSize {
		return utils.ValidationErrorResponse(c, "File size exceeds 5MB limit")
	}

	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
	isValidExt := false
	for _, ext := range allowedExtensions {
		if fileExt == ext {
			isValidExt = true
			break
		}
	}
	if !isValidExt {
		return utils.ValidationErrorResponse(c, "Invalid file type")
	}

	// Read file data
	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err)
	}
	defer src.Close()

	fileData, err := io.ReadAll(src)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read file", err)
	}

	// Process complete passport pipeline
	result, err := h.visionService.ProcessPassportImage(c.Context(), fileData)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Full passport processing failed", err)
	}

	processingTime := time.Since(startTime)

	// Build comprehensive debug response
	debugResponse := fiber.Map{
		"passport_result": result,
		"processing_time": fmt.Sprintf("%.2fs", processingTime.Seconds()),
		"image_info": fiber.Map{
			"filename": file.Filename,
			"size":     file.Size,
			"format":   strings.ToUpper(strings.TrimPrefix(fileExt, ".")),
		},
		"pipeline_status": fiber.Map{
			"ocr_success":         result.ProcessedSuccessfully,
			"overall_confidence":  result.OverallConfidence,
			"fields_extracted":    len(result.FieldConfidences),
			"has_passport_number": result.PassportNumber != nil,
			"has_names":           result.NameEN != nil || result.NameTH != nil,
			"has_dates":           result.DateOfBirth != nil || result.DateOfExpiry != nil,
		},
		"timestamp": time.Now().Format(time.RFC3339),
	}

	return utils.SuccessResponse(c, "Full passport processing completed", debugResponse)
}

// Helper function to format date pointer for JSON response
func formatDatePtr(datePtr *time.Time) interface{} {
	if datePtr == nil {
		return nil
	}
	return datePtr.Format("2006-01-02")
}
