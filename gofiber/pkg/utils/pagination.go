package utils

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gofiber-template/domain/dto"
)

// PaginationParams holds parsed pagination parameters
type PaginationParams struct {
	Page   int
	Limit  int
	Offset int
}

// ParsePaginationParams parses page and limit from query parameters
// Frontend ส่ง page (1-based), backend แปลงเป็น offset ภายใน
func ParsePaginationParams(c *fiber.Ctx, defaultLimit int) PaginationParams {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", strconv.Itoa(defaultLimit)))

	// Validate
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = defaultLimit
	}
	if limit > 100 {
		limit = 100
	}

	// Calculate offset from page
	offset := (page - 1) * limit

	return PaginationParams{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}

// CalculatePagination calculates pagination metadata based on total, offset, and limit
func CalculatePagination(total int64, offset, limit int) dto.PaginationMeta {
	currentPage := (offset / limit) + 1
	totalPages := int((total + int64(limit) - 1) / int64(limit))
	hasMore := offset+limit < int(total)

	return dto.PaginationMeta{
		Total:       total,
		HasMore:     hasMore,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		PageSize:    limit,
		Limit:       limit, // Same as PageSize for backward compatibility
		Offset:      offset,
	}
}

// PaginatedSuccessResponseNew creates a new standardized paginated response
func PaginatedSuccessResponseNew(message string, data interface{}, total int64, offset, limit int) dto.PaginatedResponse {
	pagination := CalculatePagination(total, offset, limit)

	return dto.PaginatedResponse{
		Data:       data,
		Pagination: pagination,
		Success:    true,
		Message:    message,
	}
}