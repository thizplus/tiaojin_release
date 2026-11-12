package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"gofiber-template/domain/ports"
	"gofiber-template/pkg/utils"
)

type StorageHandler struct {
	storage ports.StoragePort
}

func NewStorageHandler(storage ports.StoragePort) *StorageHandler {
	return &StorageHandler{
		storage: storage,
	}
}

// ServeFile serves a file from private S3 bucket by proxying through our server
// Route: GET /api/v1/storage/*path
// This avoids CORS issues by fetching the file server-side and streaming to client
func (h *StorageHandler) ServeFile(c *fiber.Ctx) error {
	// Check if storage is configured
	if h.storage == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "Storage service not available", nil)
	}

	// Get file path from wildcard parameter
	path := c.Params("*")
	if path == "" {
		return utils.ValidationErrorResponse(c, "File path is required")
	}

	// Generate presigned URL (valid for 5 minutes - just for fetching)
	presignedURL, err := h.storage.GetPresignedURL(path, 5*time.Minute)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate file URL", err)
	}

	// Fetch file from S3 using presigned URL
	resp, err := http.Get(presignedURL)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to fetch file", err)
	}
	defer resp.Body.Close()

	// Check if file exists
	if resp.StatusCode == http.StatusNotFound {
		return utils.NotFoundResponse(c, "File not found")
	}
	if resp.StatusCode != http.StatusOK {
		return utils.ErrorResponse(c, fiber.StatusBadGateway, "Failed to fetch file from storage", nil)
	}

	// Set content type from response or detect from extension
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = getContentType(path)
	}
	c.Set("Content-Type", contentType)

	// Set Content-Disposition if filename is provided via query param
	// ใช้สำหรับตั้งชื่อไฟล์ตอนดาวน์โหลด แทน UUID
	if filename := c.Query("filename"); filename != "" {
		// RFC 5987: ใช้ UTF-8 encoding สำหรับชื่อไฟล์ภาษาไทย/อื่นๆ
		encodedFilename := url.PathEscape(filename)
		c.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"; filename*=UTF-8''%s`, filename, encodedFilename))
	}

	// Set cache headers (cache for 1 hour)
	c.Set("Cache-Control", "public, max-age=3600")

	// Stream the file content
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read file", err)
	}

	return c.Send(body)
}

// getContentType returns MIME type based on file extension
func getContentType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	default:
		return "application/octet-stream"
	}
}

// GetPresignedURL returns a presigned URL as JSON instead of redirecting
// Route: GET /api/v1/storage/url?path=xxx
// Useful when frontend needs the URL for other purposes (e.g., download link)
func (h *StorageHandler) GetPresignedURL(c *fiber.Ctx) error {
	// Check if storage is configured
	if h.storage == nil {
		return utils.ErrorResponse(c, fiber.StatusServiceUnavailable, "Storage service not available", nil)
	}

	// Get file path from query parameter
	path := c.Query("path")
	if path == "" {
		return utils.ValidationErrorResponse(c, "File path is required")
	}

	// Get expiration from query (default: 15 minutes, max: 7 days)
	expirationMinutes := c.QueryInt("expiration", 15)
	if expirationMinutes < 1 {
		expirationMinutes = 1
	}
	if expirationMinutes > 10080 { // 7 days max
		expirationMinutes = 10080
	}

	// Generate presigned URL
	presignedURL, err := h.storage.GetPresignedURL(path, time.Duration(expirationMinutes)*time.Minute)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to generate presigned URL", err)
	}

	return utils.SuccessResponse(c, "Presigned URL generated", fiber.Map{
		"url":        presignedURL,
		"expires_in": expirationMinutes * 60, // in seconds
		"path":       path,
	})
}
