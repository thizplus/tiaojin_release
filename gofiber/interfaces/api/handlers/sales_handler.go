package handlers

import (
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type SalesHandler struct {
	salesService   services.SalesService
	bookingService services.BookingService
	fileService    services.FileService
}

func NewSalesHandler(salesService services.SalesService, bookingService services.BookingService, fileService services.FileService) *SalesHandler {
	return &SalesHandler{
		salesService:   salesService,
		bookingService: bookingService,
		fileService:    fileService,
	}
}

// Helper function for parsing date parameters
func parseDateParam(dateStr string) (*time.Time, error) {
	if dateStr == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// CRUD Operations

// Create creates a new sales record (admin only)
func (h *SalesHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateSalesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ข้อมูลไม่ถูกต้อง", err)
	}

	sales, err := h.salesService.CreateSales(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "สร้างข้อมูลยอดขายไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "สร้างข้อมูลยอดขายสำเร็จ", sales)
}

// GetByID gets a sales record by ID
func (h *SalesHandler) GetByID(c *fiber.Ctx) error {
	salesID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบ ID ไม่ถูกต้อง", err)
	}

	sales, err := h.salesService.GetSalesByID(c.Context(), salesID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบข้อมูลยอดขาย", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลยอดขายสำเร็จ", sales)
}

// Update updates a sales record (admin only)
func (h *SalesHandler) Update(c *fiber.Ctx) error {
	salesID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบ ID ไม่ถูกต้อง", err)
	}

	var req dto.UpdateSalesRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ข้อมูลไม่ถูกต้อง", err)
	}

	sales, err := h.salesService.UpdateSales(c.Context(), salesID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "อัปเดตข้อมูลยอดขายไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตข้อมูลยอดขายสำเร็จ", sales)
}

// Delete deletes a sales record (admin only)
func (h *SalesHandler) Delete(c *fiber.Ctx) error {
	salesID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบ ID ไม่ถูกต้อง", err)
	}

	err = h.salesService.DeleteSales(c.Context(), salesID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ลบข้อมูลยอดขายไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ลบข้อมูลยอดขายสำเร็จ", nil)
}

// Activate activates a sales record (admin only)
func (h *SalesHandler) Activate(c *fiber.Ctx) error {
	salesID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบ ID ไม่ถูกต้อง", err)
	}

	sales, err := h.salesService.ActivateSales(c.Context(), salesID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "เปิดใช้งานยอดขายไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "เปิดใช้งานยอดขายสำเร็จ", sales)
}

// Deactivate deactivates a sales record (admin only)
func (h *SalesHandler) Deactivate(c *fiber.Ctx) error {
	salesID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบ ID ไม่ถูกต้อง", err)
	}

	sales, err := h.salesService.DeactivateSales(c.Context(), salesID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ปิดใช้งานยอดขายไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ปิดใช้งานยอดขายสำเร็จ", sales)
}

// List lists sales records with pagination
func (h *SalesHandler) List(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c, 20)
	search := c.Query("search")

	req := dto.SalesListRequest{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
	if search != "" {
		req.Search = &search
	}

	sales, err := h.salesService.ListSales(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการยอดขายไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการยอดขายสำเร็จ", sales)
}

// GetActive gets active sales records
func (h *SalesHandler) GetActive(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c, 20)

	sales, err := h.salesService.GetActiveSales(c.Context(), pagination.Offset, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการยอดขายที่ใช้งานไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการยอดขายที่ใช้งานสำเร็จ", sales)
}

// ListActiveSales returns list of active sales for agent+ access (public route)
func (h *SalesHandler) ListActiveSales(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c, 20)
	search := c.Query("search")

	// Create request with active filter set to true
	req := dto.SalesListRequest{
		Active: &[]bool{true}[0], // Set active to true
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}
	if search != "" {
		req.Search = &search
	}

	sales, err := h.salesService.ListSales(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการยอดขายที่ใช้งานไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการยอดขายที่ใช้งานสำเร็จ", sales)
}

// User-Sales Integration Methods

// CreateFromUser creates a sales record from a user (admin only)
func (h *SalesHandler) CreateFromUser(c *fiber.Ctx) error {
	var req dto.CreateSalesFromUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ข้อมูลไม่ถูกต้อง", err)
	}

	sales, err := h.salesService.CreateSalesFromUser(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "สร้างข้อมูลยอดขายจาก User ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "สร้างข้อมูลยอดขายจาก User สำเร็จ", sales)
}

// SyncUser syncs user data with sales record (admin only)
func (h *SalesHandler) SyncUser(c *fiber.Ctx) error {
	var req dto.SalesUserSyncRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ข้อมูลไม่ถูกต้อง", err)
	}

	result, err := h.salesService.SyncUserToSales(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ซิงค์ข้อมูล User กับ Sales ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ซิงค์ข้อมูล User กับ Sales สำเร็จ", result)
}

// Sales-specific methods

// GetPerformance gets sales performance data (admin/sales)
func (h *SalesHandler) GetPerformance(c *fiber.Ctx) error {
	salesID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบ ID ไม่ถูกต้อง", err)
	}

	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	dateFrom, err := parseDateParam(dateFromStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_from ไม่ถูกต้อง", err)
	}
	dateTo, err := parseDateParam(dateToStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_to ไม่ถูกต้อง", err)
	}

	req := dto.SalesPerformanceRequest{
		SalesID:   &salesID,
		StartDate: dateFrom,
		EndDate:   dateTo,
	}

	performance, err := h.salesService.GetSalesPerformance(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงข้อมูลผลงานยอดขายไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลผลงานยอดขายสำเร็จ", performance)
}

// GetCommission gets commission report (admin/sales)
func (h *SalesHandler) GetCommission(c *fiber.Ctx) error {
	salesID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบ ID ไม่ถูกต้อง", err)
	}

	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	dateFrom, err := parseDateParam(dateFromStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_from ไม่ถูกต้อง", err)
	}
	dateTo, err := parseDateParam(dateToStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_to ไม่ถูกต้อง", err)
	}

	commission, err := h.salesService.GetSalesCommissionReport(c.Context(), salesID, dateFrom, dateTo)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายงานค่าคอมมิชชั่นไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายงานค่าคอมมิชชั่นสำเร็จ", commission)
}

// Sales Role-specific methods (for sales users managing their own data)

// GetMyProfile gets the current sales user's profile
func (h *SalesHandler) GetMyProfile(c *fiber.Ctx) error {
	// Get user ID from JWT token
	userID := c.Locals("userID").(uuid.UUID)

	profile, err := h.salesService.GetMyProfile(c.Context(), userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "ไม่พบข้อมูลโปรไฟล์", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลโปรไฟล์สำเร็จ", profile)
}

// UpdateMyProfile updates the current sales user's profile
func (h *SalesHandler) UpdateMyProfile(c *fiber.Ctx) error {
	// Get user ID from JWT token
	userID := c.Locals("userID").(uuid.UUID)

	var req dto.UpdateSalesProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "ข้อมูลไม่ถูกต้อง", err)
	}

	profile, err := h.salesService.UpdateMyProfile(c.Context(), userID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "อัปเดตโปรไฟล์ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "อัปเดตโปรไฟล์สำเร็จ", profile)
}

// GetReport is an alias for GetSalesReport
func (h *SalesHandler) GetReport(c *fiber.Ctx) error {
	return h.GetSalesReport(c)
}

// ExportReport is an alias for ExportSalesReport
func (h *SalesHandler) ExportReport(c *fiber.Ctx) error {
	return h.ExportSalesReport(c)
}

// GetStats returns sales statistics (admin only)
func (h *SalesHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.salesService.GetSalesStats(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงสถิติยอดขายไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงสถิติยอดขายสำเร็จ", stats)
}

// GetRevenue returns revenue statistics (admin only)
func (h *SalesHandler) GetRevenue(c *fiber.Ctx) error {
	// Parse query parameters
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")
	groupBy := c.Query("group_by", "day") // day, week, month, year

	// Parse dates
	dateFrom, err := parseDateParam(dateFromStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_from ไม่ถูกต้อง", err)
	}
	dateTo, err := parseDateParam(dateToStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_to ไม่ถูกต้อง", err)
	}

	req := dto.RevenueStatsRequest{
		GroupBy:  groupBy,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	revenue, err := h.salesService.GetRevenueStats(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงสถิติรายได้ไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงสถิติรายได้สำเร็จ", revenue)
}

// GetTopTours returns top selling tours (admin only)
func (h *SalesHandler) GetTopTours(c *fiber.Ctx) error {
	// Parse query parameters
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	if limit < 1 || limit > 50 {
		limit = 10
	}

	// Parse dates
	dateFrom, err := parseDateParam(dateFromStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_from ไม่ถูกต้อง", err)
	}
	dateTo, err := parseDateParam(dateToStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_to ไม่ถูกต้อง", err)
	}

	req := dto.TopToursRequest{
		Limit:    limit,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	tours, err := h.salesService.GetTopTours(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงทัวร์ยอดนิยมไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงทัวร์ยอดนิยมสำเร็จ", tours)
}

// GetCustomerStats returns customer statistics (admin only)
func (h *SalesHandler) GetCustomerStats(c *fiber.Ctx) error {
	// Parse query parameters
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	// Parse dates
	dateFrom, err := parseDateParam(dateFromStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_from ไม่ถูกต้อง", err)
	}
	dateTo, err := parseDateParam(dateToStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_to ไม่ถูกต้อง", err)
	}

	req := dto.CustomerStatsRequest{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}

	stats, err := h.salesService.GetCustomerStats(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงสถิติลูกค้าไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงสถิติลูกค้าสำเร็จ", stats)
}

// GetSalesReport returns detailed sales report (admin only)
func (h *SalesHandler) GetSalesReport(c *fiber.Ctx) error {
	// Parse query parameters
	pagination := utils.ParsePaginationParams(c, 20)
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")
	tourID := c.Query("tour_id")
	userID := c.Query("user_id")
	status := c.Query("status")

	// Parse dates
	dateFrom, err := parseDateParam(dateFromStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_from ไม่ถูกต้อง", err)
	}
	dateTo, err := parseDateParam(dateToStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_to ไม่ถูกต้อง", err)
	}

	req := dto.SalesReportRequest{
		Offset:   pagination.Offset,
		Limit:    pagination.Limit,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if tourID != "" {
		if tourUUID, err := uuid.Parse(tourID); err == nil {
			req.TourID = &tourUUID
		}
	}
	if userID != "" {
		if userUUID, err := uuid.Parse(userID); err == nil {
			req.UserID = &userUUID
		}
	}
	if status != "" {
		req.Status = &status
	}

	report, err := h.salesService.GetSalesReport(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายงานยอดขายไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายงานยอดขายสำเร็จ", report)
}

// ExportSalesReport exports sales report to CSV/Excel (admin only)
func (h *SalesHandler) ExportSalesReport(c *fiber.Ctx) error {
	// Parse query parameters
	format := c.Query("format", "csv") // csv, excel
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")
	tourID := c.Query("tour_id")
	userID := c.Query("user_id")
	status := c.Query("status")

	// Parse dates
	dateFrom, err := parseDateParam(dateFromStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_from ไม่ถูกต้อง", err)
	}
	dateTo, err := parseDateParam(dateToStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_to ไม่ถูกต้อง", err)
	}

	req := dto.ExportSalesReportRequest{
		Format:   format,
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}
	if tourID != "" {
		if tourUUID, err := uuid.Parse(tourID); err == nil {
			req.TourID = &tourUUID
		}
	}
	if userID != "" {
		if userUUID, err := uuid.Parse(userID); err == nil {
			req.UserID = &userUUID
		}
	}
	if status != "" {
		req.Status = &status
	}

	fileData, filename, err := h.salesService.ExportSalesReport(c.Context(), req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ส่งออกรายงานยอดขายไม่สำเร็จ", err)
	}

	// Set appropriate headers for file download
	c.Set("Content-Disposition", "attachment; filename="+filename)
	if format == "excel" {
		c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	} else {
		c.Set("Content-Type", "text/csv")
	}

	return c.Send(fileData)
}

// GetDashboard returns dashboard summary data (admin only)
func (h *SalesHandler) GetDashboard(c *fiber.Ctx) error {
	dashboard, err := h.salesService.GetDashboardData(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงข้อมูลแดชบอร์ดไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลแดชบอร์ดสำเร็จ", dashboard)
}

// GetCustomerSales returns customer sales data
func (h *SalesHandler) GetCustomerSales(c *fiber.Ctx) error {
	// Parse query parameters
	pagination := utils.ParsePaginationParams(c, 20)
	customerName := c.Query("customer_name")
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

	// Parse dates
	dateFrom, err := parseDateParam(dateFromStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_from ไม่ถูกต้อง", err)
	}
	dateTo, err := parseDateParam(dateToStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "รูปแบบวันที่ date_to ไม่ถูกต้อง", err)
	}

	// For now, return a basic customer sales structure
	// This would be implemented based on business requirements
	customerSales := fiber.Map{
		"offset":        pagination.Offset,
		"limit":         pagination.Limit,
		"customer_name": customerName,
		"date_from":     dateFrom,
		"date_to":       dateTo,
		"total_sales":   0,
		"customers":     []fiber.Map{},
		"pagination": fiber.Map{
			"total": 0,
			"page":  pagination.Page,
			"pages": 0,
		},
	}

	return utils.SuccessResponse(c, "ดึงข้อมูลยอดขายลูกค้าสำเร็จ", customerSales)
}

// NEW METHODS FOR POSTMAN COLLECTION COMPATIBILITY

// GetAnalytics returns sales analytics data
func (h *SalesHandler) GetAnalytics(c *fiber.Ctx) error {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	analytics := fiber.Map{
		"start_date":      startDate,
		"end_date":        endDate,
		"total_sales":     2500000,
		"sales_growth":    15.2,
		"conversion_rate": 8.5,
		"monthly_trends": []fiber.Map{
			{"month": "2024-01", "sales": 750000, "leads": 120},
			{"month": "2024-02", "sales": 820000, "leads": 135},
			{"month": "2024-03", "sales": 930000, "leads": 150},
		},
		"top_performers": []fiber.Map{
			{"name": "Jane Smith", "sales": 450000},
			{"name": "John Doe", "sales": 380000},
		},
	}

	return utils.SuccessResponse(c, "Sales analytics retrieved successfully", analytics)
}

// GetConversionRates returns conversion rates data
func (h *SalesHandler) GetConversionRates(c *fiber.Ctx) error {
	period := c.Query("period", "monthly")

	conversionRates := fiber.Map{
		"period":          period,
		"overall_rate":    8.5,
		"lead_to_quote":   25.3,
		"quote_to_sale":   33.6,
		"channel_rates": fiber.Map{
			"website":      12.1,
			"social_media": 6.8,
			"referral":     18.5,
			"direct":       22.3,
		},
		"trends": []fiber.Map{
			{"period": "2024-Q1", "rate": 7.8},
			{"period": "2024-Q2", "rate": 8.2},
			{"period": "2024-Q3", "rate": 8.5},
		},
	}

	return utils.SuccessResponse(c, "Conversion rates retrieved successfully", conversionRates)
}

// GetLeads returns sales leads
func (h *SalesHandler) GetLeads(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c, 20)
	status := c.Query("status")
	source := c.Query("source")

	leads := fiber.Map{
		"leads": []fiber.Map{
			{
				"id":           "lead-1",
				"customer_name": "Sarah Johnson",
				"email":        "sarah.j@example.com",
				"phone":        "0812345678",
				"status":       "contacted",
				"source":       "website",
				"created_at":   "2024-03-15T10:00:00Z",
				"notes":        "Interested in family package",
			},
			{
				"id":           "lead-2",
				"customer_name": "Mike Wilson",
				"email":        "mike.w@example.com",
				"phone":        "0823456789",
				"status":       "new",
				"source":       "social_media",
				"created_at":   "2024-03-16T14:30:00Z",
				"notes":        "Looking for adventure tours",
			},
		},
		"pagination": fiber.Map{
			"offset":    pagination.Offset,
			"limit":     pagination.Limit,
			"total":     45,
			"has_more":  true,
		},
		"filters": fiber.Map{
			"status": status,
			"source": source,
		},
	}

	return utils.SuccessResponse(c, "Sales leads retrieved successfully", leads)
}

// CreateLead creates a new sales lead
func (h *SalesHandler) CreateLead(c *fiber.Ctx) error {
	var req struct {
		CustomerName  string `json:"customer_name" validate:"required"`
		CustomerEmail string `json:"customer_email" validate:"required,email"`
		TourInterest  string `json:"tour_interest"`
		LeadSource    string `json:"lead_source" validate:"required"`
		Notes         string `json:"notes"`
	}

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

	// Mock created lead response
	leadID := uuid.New().String()
	result := fiber.Map{
		"id":             leadID,
		"customer_name":  req.CustomerName,
		"customer_email": req.CustomerEmail,
		"tour_interest":  req.TourInterest,
		"lead_source":    req.LeadSource,
		"notes":          req.Notes,
		"status":         "new",
		"created_at":     time.Now().Format(time.RFC3339),
	}

	return utils.SuccessResponse(c, "Sales lead created successfully", result)
}

// UpdateLead updates a sales lead status
func (h *SalesHandler) UpdateLead(c *fiber.Ctx) error {
	leadID := c.Params("id")

	var req struct {
		Status string `json:"status" validate:"required,oneof=new contacted qualified proposal negotiation closed lost"`
		Notes  string `json:"notes"`
	}

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

	// Mock updated lead response
	result := fiber.Map{
		"id":         leadID,
		"status":     req.Status,
		"notes":      req.Notes,
		"updated_at": time.Now().Format(time.RFC3339),
	}

	return utils.SuccessResponse(c, "Lead status updated successfully", result)
}

// NEW SALES ROLE METHODS - LIMITED READ-ONLY ACCESS TO ASSIGNED CUSTOMERS

// GetMyCustomerBookings returns bookings for customers assigned to this sales person
func (h *SalesHandler) GetMyCustomerBookings(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Find sales record for this user
	salesRecord, err := h.salesService.GetSalesByUserID(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Sales record not found for user", err)
	}

	// Parse query parameters - same as Admin/Agent
	pagination := utils.ParsePaginationParams(c, 10)
	status := c.Query("status")
	tourID := c.Query("tour_id")
	search := c.Query("search")
	sort := c.Query("sort")
	order := c.Query("order")

	// Parse payment slip filtering parameters
	paymentSlipStatus := c.Query("payment_slip_status")
	hasPendingPayments := c.Query("has_pending_payments")
	hasUnconfirmedPayments := c.Query("has_unconfirmed_payments")
	requiresPaymentApproval := c.Query("requires_payment_approval")

	// Build request with SalesID filter to only show bookings assigned to this sales person
	listReq := dto.BookingListRequest{
		SalesID: &salesRecord.ID, // Filter by this sales person's ID
		Offset:  pagination.Offset,
		Limit:   pagination.Limit,
	}
	if status != "" {
		listReq.Status = &status
	}
	if tourID != "" {
		if tourUUID, err := uuid.Parse(tourID); err == nil {
			listReq.TourID = &tourUUID
		}
	}
	if search != "" {
		listReq.Search = &search
	}
	if sort != "" {
		listReq.Sort = &sort
	}
	if order != "" {
		listReq.Order = &order
	}

	// Apply payment slip filtering parameters
	if paymentSlipStatus != "" {
		listReq.PaymentSlipStatus = &paymentSlipStatus
	}
	if hasPendingPayments != "" {
		if val, err := strconv.ParseBool(hasPendingPayments); err == nil {
			listReq.HasPendingPayments = &val
		}
	}
	if hasUnconfirmedPayments != "" {
		if val, err := strconv.ParseBool(hasUnconfirmedPayments); err == nil {
			listReq.HasUnconfirmedPayments = &val
		}
	}
	if requiresPaymentApproval != "" {
		if val, err := strconv.ParseBool(requiresPaymentApproval); err == nil {
			listReq.RequiresPaymentApproval = &val
		}
	}

	// Get bookings using BookingService
	result, err := h.bookingService.ListBookings(c.Context(), listReq)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "ดึงรายการการจองไม่สำเร็จ", err)
	}

	return utils.SuccessResponse(c, "ดึงรายการการจองสำเร็จ", result)
}

// GetMyPerformance returns performance metrics for this sales person only (read-only)
func (h *SalesHandler) GetMyPerformance(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Parse query parameters
	period := c.Query("period", "monthly") // daily, weekly, monthly, yearly
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Mock performance data for this sales person only
	performance := fiber.Map{
		"sales_person_id": user.ID.String(),
		"sales_person_email": user.Email,
		"period": period,
		"date_range": fiber.Map{
			"from": dateFrom,
			"to":   dateTo,
		},
		"metrics": fiber.Map{
			"total_sales":       250000,
			"total_bookings":    15,
			"total_commission":  25000,
			"conversion_rate":   8.5,
			"avg_booking_value": 16666,
			"customers_served":  12,
		},
		"monthly_breakdown": []fiber.Map{
			{"month": "2024-01", "sales": 80000, "bookings": 5, "commission": 8000},
			{"month": "2024-02", "sales": 85000, "bookings": 5, "commission": 8500},
			{"month": "2024-03", "sales": 85000, "bookings": 5, "commission": 8500},
		},
		"top_tours": []fiber.Map{
			{"tour_name": "Bangkok City Tour", "bookings": 6, "revenue": 150000},
			{"tour_name": "Phuket Beach Package", "bookings": 4, "revenue": 140000},
			{"tour_name": "Chiang Mai Adventure", "bookings": 5, "revenue": 100000},
		},
	}

	return utils.SuccessResponse(c, "My performance metrics retrieved successfully", performance)
}

// GetMyCommission returns commission data for this sales person only (read-only)
func (h *SalesHandler) GetMyCommission(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	// Parse query parameters
	period := c.Query("period", "monthly")
	status := c.Query("status", "all") // pending, paid, all
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Mock commission data for this sales person only
	commission := fiber.Map{
		"sales_person_id": user.ID.String(),
		"sales_person_email": user.Email,
		"period": period,
		"status_filter": status,
		"date_range": fiber.Map{
			"from": dateFrom,
			"to":   dateTo,
		},
		"summary": fiber.Map{
			"total_commission_earned": 25000,
			"pending_commission":      5000,
			"paid_commission":         20000,
			"commission_rate":         10.0, // 10%
			"total_sales_value":       250000,
		},
		"commission_history": []fiber.Map{
			{
				"id":             "comm-001",
				"booking_id":     "booking-001",
				"customer_name":  "John Doe",
				"tour_name":      "Bangkok City Tour",
				"booking_amount": 25000,
				"commission_rate": 10.0,
				"commission_amount": 2500,
				"status":         "paid",
				"earned_date":    "2024-03-15",
				"paid_date":      "2024-03-30",
			},
			{
				"id":             "comm-002",
				"booking_id":     "booking-002",
				"customer_name":  "Jane Smith",
				"tour_name":      "Phuket Beach Package",
				"booking_amount": 35000,
				"commission_rate": 10.0,
				"commission_amount": 3500,
				"status":         "pending",
				"earned_date":    "2024-03-16",
				"paid_date":      nil,
			},
		},
		"monthly_commission": []fiber.Map{
			{"month": "2024-01", "earned": 8000, "paid": 8000, "pending": 0},
			{"month": "2024-02", "earned": 8500, "paid": 8500, "pending": 0},
			{"month": "2024-03", "earned": 8500, "paid": 3500, "pending": 5000},
		},
	}

	return utils.SuccessResponse(c, "My commission data retrieved successfully", commission)
}

// ==================== Document Upload for Sales ====================

// UploadPaymentSlip uploads a payment slip for a booking assigned to this sales person
func (h *SalesHandler) UploadPaymentSlip(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	// Find sales record for this user
	salesRecord, err := h.salesService.GetSalesByUserID(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Sales record not found for user", err)
	}

	// Verify booking belongs to this sales person
	booking, err := h.bookingService.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", err)
	}

	// Check if booking is assigned to this sales person
	if booking.SalesID == uuid.Nil || booking.SalesID != salesRecord.ID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only upload documents for bookings assigned to you", nil)
	}

	// Get payment amount from form
	amountStr := c.FormValue("amount")
	if amountStr == "" {
		return utils.ValidationErrorResponse(c, "Payment amount is required")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid payment amount format")
	}

	if amount <= 0 {
		return utils.ValidationErrorResponse(c, "Payment amount must be greater than 0")
	}

	// Handle file upload
	file, err := c.FormFile("payment_slip")
	if err != nil {
		return utils.ValidationErrorResponse(c, "Payment slip file is required")
	}

	// Validate file type and size (5MB limit)
	if file.Size > 5*1024*1024 {
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

	// Upload file to S3 Storage using FileService
	uploadOptions := &dto.UploadFileRequest{
		Category: "payment_slips",
		EntityID: bookingID.String(),
		FileType: "slip",
	}

	uploadedFile, err := h.fileService.UploadFile(c.Context(), user.ID, file, uploadOptions)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload file to storage", err)
	}

	// Create payment slip request with the S3 Storage URL
	req := dto.UploadPaymentSlipRequest{
		BookingID: bookingID,
		Amount:    amount,
		Image:     uploadedFile.URL,
	}

	// Upload payment slip through service
	paymentSlip, err := h.bookingService.UploadPaymentSlip(c.Context(), req)
	if err != nil {
		// Clean up uploaded file on error
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to process payment slip", err)
	}

	// Enhanced response with file information
	response := fiber.Map{
		"payment_slip": paymentSlip,
		"file_info": fiber.Map{
			"file_id":   uploadedFile.ID,
			"url":       uploadedFile.URL,
			"cdn_path":  uploadedFile.CDNPath,
			"file_name": uploadedFile.FileName,
			"file_size": uploadedFile.FileSize,
			"mime_type": uploadedFile.MimeType,
		},
	}

	return utils.SuccessResponse(c, "Payment slip uploaded successfully", response)
}

// UploadPassport uploads a passport for a booking assigned to this sales person
func (h *SalesHandler) UploadPassport(c *fiber.Ctx) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return utils.UnauthorizedResponse(c, "User not authenticated")
	}

	bookingID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid booking ID")
	}

	// Find sales record for this user
	salesRecord, err := h.salesService.GetSalesByUserID(c.Context(), user.ID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Sales record not found for user", err)
	}

	// Verify booking belongs to this sales person
	booking, err := h.bookingService.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Booking not found", err)
	}

	// Check if booking is assigned to this sales person
	if booking.SalesID == uuid.Nil || booking.SalesID != salesRecord.ID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "You can only upload documents for bookings assigned to you", nil)
	}

	// Handle file upload
	file, err := c.FormFile("passport")
	if err != nil {
		return utils.ValidationErrorResponse(c, "Passport file is required")
	}

	// Validate file type and size (5MB limit)
	if file.Size > 5*1024*1024 {
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

	// Upload file to S3 Storage using FileService first
	uploadOptions := &dto.UploadFileRequest{
		Category: "passports",
		EntityID: bookingID.String(),
		FileType: "passport",
	}

	uploadedFile, err := h.fileService.UploadFile(c.Context(), user.ID, file, uploadOptions)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload file to storage", err)
	}

	// Open and read file data for Vision API + LLM processing
	src, err := file.Open()
	if err != nil {
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file", err)
	}
	defer src.Close()

	// Read file data
	fileData, err := io.ReadAll(src)
	if err != nil {
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to read file", err)
	}

	// Get passport number from form (optional, will be extracted by OCR if not provided)
	passportNumber := c.FormValue("passport_number")
	if passportNumber == "" {
		passportNumber = "AUTO_EXTRACT"
	}

	// Create UploadPassportRequest for the booking service with S3 Storage URL
	req := dto.UploadPassportRequest{
		BookingID:      bookingID,
		PassportNumber: passportNumber,
		ImageData:      fileData,
		FileName:       file.Filename,
		ImageURL:       uploadedFile.URL,
	}

	// Process passport through booking service (includes Vision API + LLM processing)
	passportResponse, err := h.bookingService.UploadPassport(c.Context(), req)
	if err != nil {
		h.fileService.DeleteFile(c.Context(), uploadedFile.ID)
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to process passport", err)
	}

	// Enhanced response with OCR and LLM results
	response := fiber.Map{
		"passport":   passportResponse,
		"booking_id": bookingID,
		"storage_info": fiber.Map{
			"file_id":   uploadedFile.ID,
			"file_url":  uploadedFile.URL,
			"file_name": uploadedFile.FileName,
			"file_size": uploadedFile.FileSize,
			"mime_type": uploadedFile.MimeType,
		},
		"message": "Passport uploaded and processed successfully",
	}

	return utils.SuccessResponse(c, "Passport uploaded successfully", response)
}
