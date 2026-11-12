package serviceimpl

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	apperrors "gofiber-template/pkg/errors"
	"gofiber-template/pkg/utils"
)


type SalesServiceImpl struct {
	salesRepo   repositories.SalesRepository
	bookingRepo repositories.BookingRepository
	userRepo    repositories.UserRepository
}

func NewSalesService(
	salesRepo repositories.SalesRepository,
	bookingRepo repositories.BookingRepository,
	userRepo repositories.UserRepository,
) services.SalesService {
	return &SalesServiceImpl{
		salesRepo:   salesRepo,
		bookingRepo: bookingRepo,
		userRepo:    userRepo,
	}
}

func (s *SalesServiceImpl) CreateSales(ctx context.Context, req dto.CreateSalesRequest) (*dto.SalesResponse, error) {
	// Business validation - validate sales code format
	if !isValidSalesCode(req.Code) {
		return nil, apperrors.ErrInvalidSalesCode
	}

	// Business validation - check if sales code already exists
	if existing, err := s.salesRepo.GetByCode(ctx, req.Code); err != nil {
		log.Printf("Error checking sales code existence: %v", err)
		return nil, fmt.Errorf("failed to validate sales code: %w", err)
	} else if existing != nil {
		return nil, apperrors.ErrSalesCodeExists
	}

	// Create domain model
	sales := &models.Sales{
		UserID:    req.UserID, // เพิ่ม UserID support
		Code:      strings.ToUpper(req.Code),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		NickName:  req.NickName,
		Active:    true, // Default to active
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create sales
	if err := s.salesRepo.Create(ctx, sales); err != nil {
		log.Printf("Error creating sales: %v", err)
		return nil, fmt.Errorf("failed to create sales: %w", err)
	}

	return s.toSalesResponse(sales), nil
}

func (s *SalesServiceImpl) GetSalesByID(ctx context.Context, salesID uuid.UUID) (*dto.SalesResponse, error) {
	sales, err := s.salesRepo.GetByID(ctx, salesID)
	if err != nil {
		log.Printf("Error retrieving sales: %v", err)
		return nil, fmt.Errorf("failed to get sales: %w", err)
	}

	if sales == nil {
		return nil, apperrors.ErrSalesNotFound
	}

	return s.toSalesResponse(sales), nil
}

func (s *SalesServiceImpl) GetSalesByCode(ctx context.Context, code string) (*dto.SalesResponse, error) {
	sales, err := s.salesRepo.GetByCode(ctx, strings.ToUpper(code))
	if err != nil {
		log.Printf("Error retrieving sales by code: %v", err)
		return nil, fmt.Errorf("failed to get sales by code: %w", err)
	}

	if sales == nil {
		return nil, apperrors.ErrSalesNotFound
	}

	return s.toSalesResponse(sales), nil
}

func (s *SalesServiceImpl) UpdateSales(ctx context.Context, salesID uuid.UUID, req dto.UpdateSalesRequest) (*dto.SalesResponse, error) {
	// Get existing sales
	sales, err := s.salesRepo.GetByID(ctx, salesID)
	if err != nil {
		log.Printf("Error retrieving sales: %v", err)
		return nil, fmt.Errorf("failed to get sales: %w", err)
	}

	if sales == nil {
		return nil, apperrors.ErrSalesNotFound
	}

	// Update fields if provided
	if req.UserID != nil {
		log.Printf("Updating UserID from %v to %v", sales.UserID, req.UserID)
		sales.UserID = req.UserID
	}
	if req.FirstName != nil {
		sales.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		sales.LastName = *req.LastName
	}
	if req.NickName != nil {
		sales.NickName = *req.NickName
	}
	if req.Active != nil {
		sales.Active = *req.Active
	}

	sales.UpdatedAt = time.Now()

	// Update sales
	if err := s.salesRepo.Update(ctx, salesID, sales); err != nil {
		log.Printf("Error updating sales: %v", err)
		return nil, fmt.Errorf("failed to update sales: %w", err)
	}

	return s.toSalesResponse(sales), nil
}

func (s *SalesServiceImpl) DeleteSales(ctx context.Context, salesID uuid.UUID) error {
	// Check if sales exists
	sales, err := s.salesRepo.GetByID(ctx, salesID)
	if err != nil {
		log.Printf("Error retrieving sales: %v", err)
		return fmt.Errorf("failed to get sales: %w", err)
	}

	if sales == nil {
		return apperrors.ErrSalesNotFound
	}

	// Business rule - check if sales has bookings (simplified check)
	// In real implementation, you would check for active bookings
	bookings, err := s.bookingRepo.GetBySalesID(ctx, salesID, 0, 1)
	if err != nil {
		log.Printf("Error checking sales bookings: %v", err)
		return fmt.Errorf("failed to check sales bookings: %w", err)
	}

	if len(bookings) > 0 {
		return apperrors.ErrSalesHasBookings
	}

	// Delete sales
	if err := s.salesRepo.Delete(ctx, salesID); err != nil {
		log.Printf("Error deleting sales: %v", err)
		return fmt.Errorf("failed to delete sales: %w", err)
	}

	return nil
}

func (s *SalesServiceImpl) ListSales(ctx context.Context, req dto.SalesListRequest) (*dto.SalesListResponse, error) {
	var salesList []*models.Sales
	var total int64
	var err error

	// Handle different filtering options
	if req.Search != nil && strings.TrimSpace(*req.Search) != "" {
		// Search by code, first_name, last_name, nick_name
		salesList, err = s.salesRepo.SearchByKeyword(ctx, strings.TrimSpace(*req.Search), req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error searching sales: %v", err)
			return nil, fmt.Errorf("failed to search sales: %w", err)
		}
		total = int64(len(salesList)) // Simplified count for search
	} else if req.Active != nil {
		// Filter by active status
		if *req.Active {
			salesList, err = s.salesRepo.GetActive(ctx, req.Offset, req.Limit)
		} else {
			// Get all and filter inactive ones (simplified approach)
			allSales, err := s.salesRepo.List(ctx, req.Offset, req.Limit*2)
			if err != nil {
				return nil, fmt.Errorf("failed to get sales: %w", err)
			}
			salesList = make([]*models.Sales, 0)
			for _, sales := range allSales {
				if !sales.Active {
					salesList = append(salesList, sales)
					if len(salesList) >= req.Limit {
						break
					}
				}
			}
		}
		if err != nil {
			log.Printf("Error getting sales by active status: %v", err)
			return nil, fmt.Errorf("failed to get sales by active status: %w", err)
		}
		total = int64(len(salesList)) // Simplified count
	} else {
		// List all sales
		salesList, err = s.salesRepo.List(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error listing sales: %v", err)
			return nil, fmt.Errorf("failed to list sales: %w", err)
		}

		// Get total count
		total, err = s.salesRepo.Count(ctx)
		if err != nil {
			log.Printf("Error counting sales: %v", err)
			return nil, fmt.Errorf("failed to count sales: %w", err)
		}
	}

	// Convert to response DTOs
	responseSales := make([]*dto.SalesResponse, len(salesList))
	for i, sales := range salesList {
		responseSales[i] = s.toSalesResponse(sales)
	}

	return &dto.SalesListResponse{
		Sales:      responseSales,
		Pagination: utils.CalculatePagination(total, req.Offset, req.Limit),
	}, nil
}

func (s *SalesServiceImpl) GetActiveSales(ctx context.Context, offset, limit int) ([]*dto.SalesResponse, error) {
	salesList, err := s.salesRepo.GetActive(ctx, offset, limit)
	if err != nil {
		log.Printf("Error getting active sales: %v", err)
		return nil, fmt.Errorf("failed to get active sales: %w", err)
	}

	// Convert to response DTOs
	responseSales := make([]*dto.SalesResponse, len(salesList))
	for i, sales := range salesList {
		responseSales[i] = s.toSalesResponse(sales)
	}

	return responseSales, nil
}

func (s *SalesServiceImpl) SearchSales(ctx context.Context, query string, offset, limit int) ([]*dto.SalesResponse, error) {
	// Business validation - clean search query
	query = strings.TrimSpace(query)
	if query == "" {
		return []*dto.SalesResponse{}, nil
	}

	salesList, err := s.salesRepo.SearchByName(ctx, query, offset, limit)
	if err != nil {
		log.Printf("Error searching sales: %v", err)
		return nil, fmt.Errorf("failed to search sales: %w", err)
	}

	// Convert to response DTOs
	responseSales := make([]*dto.SalesResponse, len(salesList))
	for i, sales := range salesList {
		responseSales[i] = s.toSalesResponse(sales)
	}

	return responseSales, nil
}

// Status Management
func (s *SalesServiceImpl) ActivateSales(ctx context.Context, salesID uuid.UUID) (*dto.SalesResponse, error) {
	return s.UpdateActiveStatus(ctx, salesID, true)
}

func (s *SalesServiceImpl) DeactivateSales(ctx context.Context, salesID uuid.UUID) (*dto.SalesResponse, error) {
	return s.UpdateActiveStatus(ctx, salesID, false)
}

func (s *SalesServiceImpl) UpdateActiveStatus(ctx context.Context, salesID uuid.UUID, active bool) (*dto.SalesResponse, error) {
	// Get existing sales
	sales, err := s.salesRepo.GetByID(ctx, salesID)
	if err != nil {
		log.Printf("Error retrieving sales: %v", err)
		return nil, fmt.Errorf("failed to get sales: %w", err)
	}

	if sales == nil {
		return nil, apperrors.ErrSalesNotFound
	}

	// Business validation - check if status change is needed
	if sales.Active == active {
		if active {
			return nil, apperrors.ErrSalesAlreadyActive
		} else {
			return nil, apperrors.ErrSalesAlreadyInactive
		}
	}

	// Update active status
	if err := s.salesRepo.UpdateActiveStatus(ctx, salesID, active); err != nil {
		log.Printf("Error updating sales active status: %v", err)
		return nil, fmt.Errorf("failed to update sales active status: %w", err)
	}

	sales.Active = active
	sales.UpdatedAt = time.Now()

	return s.toSalesResponse(sales), nil
}

// Performance & Statistics
func (s *SalesServiceImpl) GetSalesWithBookingCounts(ctx context.Context, offset, limit int) ([]*dto.SalesWithStatsResponse, error) {
	salesList, err := s.salesRepo.List(ctx, offset, limit)
	if err != nil {
		log.Printf("Error getting sales with booking counts: %v", err)
		return nil, fmt.Errorf("failed to get sales with booking counts: %w", err)
	}

	// Convert to response DTOs with stats
	responseSales := make([]*dto.SalesWithStatsResponse, len(salesList))
	for i, sales := range salesList {
		responseSales[i] = &dto.SalesWithStatsResponse{
			SalesResponse: s.toSalesResponse(sales),
			BookingStats:  s.calculateSalesBookingStats(ctx, sales.ID),
		}
	}

	return responseSales, nil
}

func (s *SalesServiceImpl) GetSalesPerformance(ctx context.Context, req dto.SalesPerformanceRequest) (*dto.SalesPerformanceResponse, error) {
	var salesList []*models.Sales
	var err error

	if req.SalesID != nil {
		// Get specific sales performance
		sales, err := s.salesRepo.GetByID(ctx, *req.SalesID)
		if err != nil {
			log.Printf("Error getting sales performance: %v", err)
			return nil, fmt.Errorf("failed to get sales performance: %w", err)
		}
		if sales != nil {
			salesList = []*models.Sales{sales}
		}
	} else {
		// Get all sales performance
		salesList, err = s.salesRepo.List(ctx, req.Offset, req.Limit)
		if err != nil {
			log.Printf("Error getting sales performance: %v", err)
			return nil, fmt.Errorf("failed to get sales performance: %w", err)
		}
	}

	// Convert to response DTOs with stats
	responseSales := make([]*dto.SalesWithStatsResponse, len(salesList))
	var totalRevenue, totalCommission float64

	for i, sales := range salesList {
		stats := s.calculateSalesBookingStats(ctx, sales.ID)
		responseSales[i] = &dto.SalesWithStatsResponse{
			SalesResponse: s.toSalesResponse(sales),
			BookingStats:  stats,
		}

		if stats != nil {
			totalRevenue += stats.TotalRevenue
			totalCommission += stats.TotalCommission
		}
	}

	response := &dto.SalesPerformanceResponse{
		Sales:           responseSales,
		Pagination:      utils.CalculatePagination(int64(len(responseSales)), req.Offset, req.Limit),
		TotalRevenue:    totalRevenue,
		TotalCommission: totalCommission,
	}

	// Add period information if dates are provided
	if req.StartDate != nil && req.EndDate != nil {
		response.Period = &dto.PerformancePeriod{
			StartDate: req.StartDate.Format("2006-01-02"),
			EndDate:   req.EndDate.Format("2006-01-02"),
		}
	}

	return response, nil
}

func (s *SalesServiceImpl) GetSalesCommissionReport(ctx context.Context, salesID uuid.UUID, startDate, endDate *time.Time) (*dto.SalesCommissionReport, error) {
	// Get sales
	sales, err := s.salesRepo.GetByID(ctx, salesID)
	if err != nil {
		log.Printf("Error getting sales for commission report: %v", err)
		return nil, fmt.Errorf("failed to get sales: %w", err)
	}
	if sales == nil {
		return nil, apperrors.ErrSalesNotFound
	}

	// Calculate commission
	totalCommission, err := s.CalculateCommission(ctx, salesID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get booking count (simplified)
	bookings, err := s.bookingRepo.GetBySalesID(ctx, salesID, 0, 1000)
	if err != nil {
		log.Printf("Error getting bookings for commission report: %v", err)
		return nil, fmt.Errorf("failed to get bookings: %w", err)
	}

	// Create report
	report := &dto.SalesCommissionReport{
		SalesID:         salesID,
		SalesCode:       sales.Code,
		SalesName:       sales.GetFullName(),
		TotalCommission: totalCommission,
		BookingCount:    int64(len(bookings)),
		MonthlyDetails:  []*dto.MonthlyPerformance{}, // TODO: Implement monthly breakdown
	}

	return report, nil
}

func (s *SalesServiceImpl) GetMonthlyPerformance(ctx context.Context, salesID uuid.UUID, year int) ([]*dto.MonthlyPerformance, error) {
	// TODO: Implement monthly performance calculation
	// This would require date-based aggregation of bookings and commissions

	// Placeholder implementation
	monthly := make([]*dto.MonthlyPerformance, 12)
	for i := 0; i < 12; i++ {
		monthly[i] = &dto.MonthlyPerformance{
			Year:       year,
			Month:      i + 1,
			Bookings:   0,
			Revenue:    0,
			Commission: 0,
			Passengers: 0,
		}
	}

	return monthly, nil
}

func (s *SalesServiceImpl) GetTopPerformingSales(ctx context.Context, startDate, endDate *time.Time, limit int) ([]*dto.SalesWithStatsResponse, error) {
	// Get all sales and calculate their performance
	salesList, err := s.salesRepo.List(ctx, 0, limit*2) // Get more to sort
	if err != nil {
		log.Printf("Error getting top performing sales: %v", err)
		return nil, fmt.Errorf("failed to get top performing sales: %w", err)
	}

	// Convert to response DTOs with stats and sort by performance
	responseSales := make([]*dto.SalesWithStatsResponse, 0, len(salesList))
	for _, sales := range salesList {
		stats := s.calculateSalesBookingStats(ctx, sales.ID)
		responseSales = append(responseSales, &dto.SalesWithStatsResponse{
			SalesResponse: s.toSalesResponse(sales),
			BookingStats:  stats,
		})
	}

	// TODO: Sort by performance metrics (revenue, commission, bookings)
	// For now, return first N results
	if len(responseSales) > limit {
		responseSales = responseSales[:limit]
	}

	return responseSales, nil
}

// Business Logic
func (s *SalesServiceImpl) CalculateCommission(ctx context.Context, salesID uuid.UUID, startDate, endDate *time.Time) (float64, error) {
	// Get bookings for sales in date range
	bookings, err := s.bookingRepo.GetBySalesID(ctx, salesID, 0, 1000) // Simplified
	if err != nil {
		log.Printf("Error calculating commission: %v", err)
		return 0, fmt.Errorf("failed to calculate commission: %w", err)
	}

	var totalCommission float64
	for _, booking := range bookings {
		// Filter by date range if provided
		if startDate != nil && booking.CreatedAt.Before(*startDate) {
			continue
		}
		if endDate != nil && booking.CreatedAt.After(*endDate) {
			continue
		}

		totalCommission += booking.SumCommission
	}

	return totalCommission, nil
}

func (s *SalesServiceImpl) ReassignBookings(ctx context.Context, fromSalesID, toSalesID uuid.UUID) (int64, error) {
	// Business validation - check if both sales exist
	fromSales, err := s.salesRepo.GetByID(ctx, fromSalesID)
	if err != nil {
		return 0, fmt.Errorf("failed to validate from sales: %w", err)
	}
	if fromSales == nil {
		return 0, apperrors.ErrSalesNotFound
	}

	toSales, err := s.salesRepo.GetByID(ctx, toSalesID)
	if err != nil {
		return 0, fmt.Errorf("failed to validate to sales: %w", err)
	}
	if toSales == nil {
		return 0, apperrors.ErrSalesNotFound
	}

	// Business validation - check if target sales is active
	if !toSales.Active {
		return 0, fmt.Errorf("target sales is not active")
	}

	// TODO: Implement booking reassignment
	// This would involve updating all bookings from fromSalesID to toSalesID
	log.Printf("Reassigning bookings from %s to %s", fromSalesID, toSalesID)

	return 0, nil // Placeholder return
}

func (s *SalesServiceImpl) GetSalesStats(ctx context.Context) (*dto.SalesStatsResponse, error) {
	// Get total sales count
	totalSales, err := s.salesRepo.Count(ctx)
	if err != nil {
		log.Printf("Error getting sales stats: %v", err)
		return nil, fmt.Errorf("failed to get sales stats: %w", err)
	}

	// Get active sales count
	activeSales, err := s.salesRepo.GetActive(ctx, 0, 1000)
	if err != nil {
		log.Printf("Error getting active sales for stats: %v", err)
		return nil, fmt.Errorf("failed to get sales stats: %w", err)
	}

	activeCount := int64(len(activeSales))
	inactiveCount := totalSales - activeCount

	return &dto.SalesStatsResponse{
		TotalSales:    totalSales,
		ActiveSales:   activeCount,
		InactiveSales: inactiveCount,
	}, nil
}

func (s *SalesServiceImpl) ValidateSalesCode(ctx context.Context, code string, excludeSalesID *uuid.UUID) (bool, error) {
	// Business validation - validate format first
	if !isValidSalesCode(code) {
		return false, nil
	}

	// Check if sales code exists
	existing, err := s.salesRepo.GetByCode(ctx, strings.ToUpper(code))
	if err != nil {
		log.Printf("Error checking sales code existence: %v", err)
		return false, fmt.Errorf("failed to validate sales code: %w", err)
	}

	// If no existing sales found, code is available
	if existing == nil {
		return true, nil
	}

	// If excludeSalesID is provided and matches, code is available for that sales
	if excludeSalesID != nil && existing.ID == *excludeSalesID {
		return true, nil
	}

	// Code already exists for another sales
	return false, nil
}

// Helper methods
func (s *SalesServiceImpl) toSalesResponse(sales *models.Sales) *dto.SalesResponse {
	log.Printf("Converting to SalesResponse: UserID=%v", sales.UserID)
	response := &dto.SalesResponse{
		ID:        sales.ID,
		UserID:    sales.UserID, // เพิ่ม UserID
		Code:      sales.Code,
		FirstName: sales.FirstName,
		LastName:  sales.LastName,
		NickName:  sales.NickName,
		Active:    sales.Active,
		CreatedAt: sales.CreatedAt.Format(time.RFC3339),
		UpdatedAt: sales.UpdatedAt.Format(time.RFC3339),
	}

	// เพิ่มข้อมูล User relationship ถ้ามี
	if sales.User != nil {
		response.UserEmail = sales.User.Email
		role := string(sales.User.Role)
		response.UserRole = &role
	}

	return response
}

func (s *SalesServiceImpl) calculateSalesBookingStats(ctx context.Context, salesID uuid.UUID) *dto.SalesBookingStats {
	// Get bookings for sales (simplified)
	bookings, err := s.bookingRepo.GetBySalesID(ctx, salesID, 0, 1000)
	if err != nil {
		log.Printf("Error calculating booking stats: %v", err)
		return nil
	}

	var totalRevenue, totalCommission float64
	var confirmedBookings int64

	for _, booking := range bookings {
		totalRevenue += booking.SumNet
		totalCommission += booking.SumCommission

		if booking.Status == models.BookingStatusConfirm {
			confirmedBookings++
		}
	}

	var averageBookingValue float64
	if len(bookings) > 0 {
		averageBookingValue = totalRevenue / float64(len(bookings))
	}

	return &dto.SalesBookingStats{
		TotalBookings:       int64(len(bookings)),
		ConfirmedBookings:   confirmedBookings,
		TotalRevenue:        totalRevenue,
		TotalCommission:     totalCommission,
		AverageBookingValue: averageBookingValue,
	}
}

// Additional Analytics & Reporting Methods Implementation
func (s *SalesServiceImpl) GetRevenueStats(ctx context.Context, req dto.RevenueStatsRequest) (*dto.RevenueStatsResponse, error) {
	var totalRevenue, totalCommission float64
	var totalBookings int64
	var monthlyBreakdown []*dto.MonthlyPerformance

	// Get all sales performance - simplified implementation
	salesList, err := s.salesRepo.List(ctx, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales list: %w", err)
	}

	for _, sales := range salesList {
		stats := s.calculateSalesBookingStats(ctx, sales.ID)
		if stats != nil {
			totalRevenue += stats.TotalRevenue
			totalCommission += stats.TotalCommission
			totalBookings += stats.TotalBookings
		}
	}

	var averageBooking float64
	if totalBookings > 0 {
		averageBooking = totalRevenue / float64(totalBookings)
	}

	response := &dto.RevenueStatsResponse{
		TotalRevenue:     totalRevenue,
		TotalCommission:  totalCommission,
		BookingCount:     totalBookings,
		AverageBooking:   averageBooking,
		MonthlyBreakdown: monthlyBreakdown,
	}

	// Add period information if dates provided
	if req.DateFrom != nil && req.DateTo != nil {
		response.Period = &dto.PerformancePeriod{
			StartDate: req.DateFrom.Format("2006-01-02"),
			EndDate:   req.DateTo.Format("2006-01-02"),
		}
	}

	return response, nil
}

func (s *SalesServiceImpl) GetTopTours(ctx context.Context, req dto.TopToursRequest) (*dto.TopToursResponse, error) {
	// TODO: Implement top tours analytics
	// This would require accessing tour/booking data and aggregating by tour

	// Placeholder implementation
	tours := make([]*dto.TopTourStats, 0)

	response := &dto.TopToursResponse{
		Tours: tours,
	}

	// Add period information if dates provided
	if req.DateFrom != nil && req.DateTo != nil {
		response.Period = &dto.PerformancePeriod{
			StartDate: req.DateFrom.Format("2006-01-02"),
			EndDate:   req.DateTo.Format("2006-01-02"),
		}
	}

	return response, nil
}

func (s *SalesServiceImpl) GetCustomerStats(ctx context.Context, req dto.CustomerStatsRequest) (*dto.CustomerStatsResponse, error) {
	// TODO: Implement customer analytics
	// This would require accessing customer/booking data for analytics

	// Placeholder implementation
	response := &dto.CustomerStatsResponse{
		TotalCustomers:        0,
		NewCustomers:          0,
		ReturningCustomers:    0,
		AverageOrderValue:     0,
		CustomerRetentionRate: 0,
	}

	// Add period information if dates provided
	if req.DateFrom != nil && req.DateTo != nil {
		response.Period = &dto.PerformancePeriod{
			StartDate: req.DateFrom.Format("2006-01-02"),
			EndDate:   req.DateTo.Format("2006-01-02"),
		}
	}

	return response, nil
}

func (s *SalesServiceImpl) GetSalesReport(ctx context.Context, req dto.SalesReportRequest) (*dto.SalesReportResponse, error) {
	// Get sales stats
	salesStats, err := s.GetSalesStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales stats: %w", err)
	}

	// Get revenue stats
	revenueReq := dto.RevenueStatsRequest{
		DateFrom: req.DateFrom,
		DateTo:   req.DateTo,
	}
	revenueStats, err := s.GetRevenueStats(ctx, revenueReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue stats: %w", err)
	}

	// Create summary
	summary := &dto.SalesReportSummary{
		TotalRevenue:        revenueStats.TotalRevenue,
		TotalCommission:     revenueStats.TotalCommission,
		TotalBookings:       revenueStats.BookingCount,
		ConfirmedBookings:   revenueStats.BookingCount, // Simplified
		AverageBookingValue: revenueStats.AverageBooking,
		ActiveSalesCount:    salesStats.ActiveSales,
	}

	// Prepare response
	response := &dto.SalesReportResponse{
		ReportType:  "summary",
		Summary:     summary,
		GeneratedAt: time.Now(),
	}

	// Add period if dates provided
	if req.DateFrom != nil && req.DateTo != nil {
		response.Period = &dto.PerformancePeriod{
			StartDate: req.DateFrom.Format("2006-01-02"),
			EndDate:   req.DateTo.Format("2006-01-02"),
		}
	}

	return response, nil
}

func (s *SalesServiceImpl) ExportSalesReport(ctx context.Context, req dto.ExportSalesReportRequest) ([]byte, string, error) {
	// Get report data
	reportReq := dto.SalesReportRequest{
		DateFrom: req.DateFrom,
		DateTo:   req.DateTo,
		TourID:   req.TourID,
		UserID:   req.UserID,
		Status:   req.Status,
		Offset:   0,
		Limit:    1000,
	}

	report, err := s.GetSalesReport(ctx, reportReq)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get sales report: %w", err)
	}

	// Generate filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	var filename string
	var data []byte

	// TODO: Implement actual export functionality based on format
	switch req.Format {
	case "csv":
		filename = fmt.Sprintf("sales_report_%s.csv", timestamp)
		data, err = s.exportToCSV(report)
	case "excel":
		filename = fmt.Sprintf("sales_report_%s.xlsx", timestamp)
		data, err = s.exportToExcel(report)
	default:
		return nil, "", fmt.Errorf("unsupported export format: %s", req.Format)
	}

	if err != nil {
		return nil, "", err
	}

	return data, filename, nil
}

func (s *SalesServiceImpl) GetDashboardData(ctx context.Context) (*dto.DashboardDataResponse, error) {
	// Get sales stats
	salesStats, err := s.GetSalesStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sales stats: %w", err)
	}

	// Get recent revenue (last 30 days)
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	revenueStats, err := s.GetRevenueStats(ctx, dto.RevenueStatsRequest{
		DateFrom: &startDate,
		DateTo:   &endDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get revenue stats: %w", err)
	}

	// Get top performers
	topPerformers, err := s.GetTopPerformingSales(ctx, &startDate, &endDate, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get top performers: %w", err)
	}

	// Get customer stats
	customerStats, err := s.GetCustomerStats(ctx, dto.CustomerStatsRequest{
		DateFrom: &startDate,
		DateTo:   &endDate,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get customer stats: %w", err)
	}

	// TODO: Implement recent activity and monthly trends

	return &dto.DashboardDataResponse{
		SalesStats:     salesStats,
		RecentRevenue:  revenueStats,
		TopPerformers:  topPerformers,
		TopTours:       []*dto.TopTourStats{}, // TODO: Implement
		MonthlyTrends:  []*dto.MonthlyPerformance{}, // TODO: Implement
		CustomerStats:  customerStats,
		RecentActivity: []*dto.RecentActivityItem{}, // TODO: Implement
		LastUpdated:    time.Now(),
	}, nil
}

// User-Sales Relationship Methods
func (s *SalesServiceImpl) CreateSalesFromUser(ctx context.Context, req dto.CreateSalesFromUserRequest) (*dto.SalesResponse, error) {
	// ตรวจสอบว่า User มี role 'sales'
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	if user.Role != models.UserRoleSales {
		return nil, fmt.Errorf("user must have sales role, current role: %s", user.Role)
	}

	// ตรวจสอบว่า User นี้มี Sales record อยู่แล้วหรือไม่
	existingSales, err := s.GetSalesByUserID(ctx, req.UserID)
	if err == nil && existingSales != nil {
		return nil, fmt.Errorf("user already has sales record")
	}

	// สร้าง Sales record พร้อม UserID
	sales := &models.Sales{
		UserID:    &req.UserID,
		Code:      strings.ToUpper(req.Code),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		NickName:  req.NickName,
		Active:    user.Status == models.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// สร้าง Sales record
	if err := s.salesRepo.Create(ctx, sales); err != nil {
		log.Printf("Error creating sales from user: %v", err)
		return nil, fmt.Errorf("failed to create sales from user: %w", err)
	}

	return s.toSalesResponse(sales), nil
}

func (s *SalesServiceImpl) SyncUserToSales(ctx context.Context, req dto.SalesUserSyncRequest) (*dto.SalesUserSyncResponse, error) {
	// ตรวจสอบว่ามี Sales record อยู่แล้วหรือไม่
	existingSales, err := s.GetSalesByUserID(ctx, req.UserID)
	if err == nil && existingSales != nil {
		return &dto.SalesUserSyncResponse{
			Success: true,
			Message: "Sales record already exists",
			Sales:   existingSales,
		}, nil
	}

	// ดึงข้อมูล User
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return &dto.SalesUserSyncResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to get user: %v", err),
		}, err
	}

	if user == nil {
		return &dto.SalesUserSyncResponse{
			Success: false,
			Message: "User not found",
		}, apperrors.ErrUserNotFound
	}

	if user.Role != models.UserRoleSales {
		return &dto.SalesUserSyncResponse{
			Success: false,
			Message: fmt.Sprintf("User must have sales role, current role: %s", user.Role),
		}, nil
	}

	// สร้าง Sales record ใหม่โดยอัตโนมัติ - ใช้ username แทน email
	salesCode := s.generateSalesCodeFromUsername(user.Username)
	firstName, lastName := s.extractNamesFromUsername(user.Username)

	sales := &models.Sales{
		UserID:    &req.UserID,
		Code:      salesCode,
		FirstName: firstName,
		LastName:  lastName,
		NickName:  s.extractNickNameFromUsername(user.Username),
		Active:    user.Status == models.UserStatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// สร้าง Sales record
	if err := s.salesRepo.Create(ctx, sales); err != nil {
		log.Printf("Error syncing user to sales: %v", err)
		return &dto.SalesUserSyncResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create sales record: %v", err),
		}, err
	}

	salesResponse := s.toSalesResponse(sales)
	return &dto.SalesUserSyncResponse{
		Success: true,
		Message: "Sales record created successfully",
		Sales:   salesResponse,
	}, nil
}

func (s *SalesServiceImpl) GetSalesByUserID(ctx context.Context, userID uuid.UUID) (*dto.SalesResponse, error) {
	sales, err := s.salesRepo.GetByUserID(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving sales by user ID: %v", err)
		return nil, fmt.Errorf("failed to get sales by user ID: %w", err)
	}

	if sales == nil {
		return nil, apperrors.ErrSalesNotFound
	}

	return s.toSalesResponse(sales), nil
}

func (s *SalesServiceImpl) GetSalesBookingStats(ctx context.Context, salesID uuid.UUID) (*dto.SalesBookingStatsResponse, error) {
	// Count total bookings for this sales
	totalBookings, err := s.bookingRepo.Count(ctx)
	if err != nil {
		log.Printf("Error counting bookings: %v", err)
		return nil, fmt.Errorf("failed to count bookings: %w", err)
	}

	// For now, return basic stats - this can be expanded with more detailed booking analytics
	stats := &dto.SalesBookingStatsResponse{
		SalesID:           salesID,
		TotalBookings:     totalBookings,
		PendingBookings:   0, // TODO: Implement based on booking status
		ConfirmedBookings: 0, // TODO: Implement based on booking status
		CancelledBookings: 0, // TODO: Implement based on booking status
		TotalRevenue:      0.0, // TODO: Calculate from confirmed bookings
		MonthlyRevenue:    0.0, // TODO: Calculate current month revenue
	}

	return stats, nil
}

func (s *SalesServiceImpl) GetMyProfile(ctx context.Context, userID uuid.UUID) (*dto.SalesMyProfileResponse, error) {
	// ดึงข้อมูล Sales ของ user นี้
	salesResponse, err := s.GetSalesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// ดึงข้อมูล booking stats ของ sales นี้
	salesUUID, err := uuid.Parse(salesResponse.ID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid sales ID: %w", err)
	}

	bookingStatsResponse, err := s.GetSalesBookingStats(ctx, salesUUID)
	if err != nil {
		log.Printf("Warning: Failed to get booking stats: %v", err)
		// ไม่ return error เพราะ stats เป็น optional
	}

	var bookingStats *dto.SalesBookingStats
	if bookingStatsResponse != nil {
		bookingStats = &dto.SalesBookingStats{
			TotalBookings:       bookingStatsResponse.TotalBookings,
			ConfirmedBookings:   bookingStatsResponse.ConfirmedBookings,
			TotalRevenue:        bookingStatsResponse.TotalRevenue,
			TotalCommission:     0.0, // TODO: Calculate commission
			AverageBookingValue: 0.0, // TODO: Calculate average
		}
	}

	return &dto.SalesMyProfileResponse{
		SalesResponse:     salesResponse,
		CanManageBookings: salesResponse.Active,
		BookingStats:      bookingStats,
	}, nil
}

func (s *SalesServiceImpl) UpdateMyProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateSalesProfileRequest) (*dto.SalesMyProfileResponse, error) {
	// ดึงข้อมูล Sales ของ user นี้
	salesResponse, err := s.GetSalesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// สร้าง UpdateSalesRequest จาก UpdateSalesProfileRequest
	updateReq := dto.UpdateSalesRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		NickName:  req.NickName,
	}

	// อัปเดตข้อมูล Sales
	_, err = s.UpdateSales(ctx, salesResponse.ID, updateReq)
	if err != nil {
		return nil, err
	}

	// ดึงข้อมูลที่อัปเดตแล้วกลับมา
	return s.GetMyProfile(ctx, userID)
}

// Helper methods for User-Sales sync
func (s *SalesServiceImpl) generateSalesCodeFromEmail(email string) string {
	// ใช้ prefix จาก email เป็น sales code
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		prefix := strings.ToUpper(parts[0])
		// ตัดให้เหลือแค่ 10 ตัวอักษร และใส่ timestamp suffix
		if len(prefix) > 6 {
			prefix = prefix[:6]
		}
		timestamp := time.Now().Format("0102") // MMDD format
		return fmt.Sprintf("%s-%s", prefix, timestamp)
	}
	// fallback หากไม่สามารถ parse email ได้
	return fmt.Sprintf("SALES-%s", time.Now().Format("010215"))
}

func (s *SalesServiceImpl) extractNamesFromEmail(email string) (firstName, lastName string) {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		namePart := parts[0]
		// ลองแยกชื่อถ้ามี dot หรือ underscore
		if strings.Contains(namePart, ".") {
			nameParts := strings.Split(namePart, ".")
			if len(nameParts) >= 2 {
				return strings.Title(nameParts[0]), strings.Title(nameParts[1])
			}
		} else if strings.Contains(namePart, "_") {
			nameParts := strings.Split(namePart, "_")
			if len(nameParts) >= 2 {
				return strings.Title(nameParts[0]), strings.Title(nameParts[1])
			}
		}
		// ถ้าไม่มี separator ให้ใช้ทั้งหมดเป็น firstName
		return strings.Title(namePart), "Sales"
	}
	return "Sales", "User"
}

func (s *SalesServiceImpl) extractNickNameFromEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) > 0 {
		namePart := parts[0]
		// ใช้ส่วนแรกของ email เป็น nickname
		if len(namePart) > 10 {
			return strings.Title(namePart[:10])
		}
		return strings.Title(namePart)
	}
	return "Sales"
}

// Helper methods for User-Sales sync using username
func (s *SalesServiceImpl) generateSalesCodeFromUsername(username string) string {
	prefix := strings.ToUpper(username)
	// ตัดให้เหลือแค่ 6 ตัวอักษร และใส่ timestamp suffix
	if len(prefix) > 6 {
		prefix = prefix[:6]
	}
	timestamp := time.Now().Format("0102") // MMDD format
	return fmt.Sprintf("%s-%s", prefix, timestamp)
}

func (s *SalesServiceImpl) extractNamesFromUsername(username string) (firstName, lastName string) {
	// ลองแยกชื่อถ้ามี dot หรือ underscore
	if strings.Contains(username, ".") {
		nameParts := strings.Split(username, ".")
		if len(nameParts) >= 2 {
			return strings.Title(nameParts[0]), strings.Title(nameParts[1])
		}
	} else if strings.Contains(username, "_") {
		nameParts := strings.Split(username, "_")
		if len(nameParts) >= 2 {
			return strings.Title(nameParts[0]), strings.Title(nameParts[1])
		}
	}
	// ถ้าไม่มี separator ให้ใช้ทั้งหมดเป็น firstName
	return strings.Title(username), "Sales"
}

func (s *SalesServiceImpl) extractNickNameFromUsername(username string) string {
	// ใช้ username เป็น nickname
	if len(username) > 10 {
		return strings.Title(username[:10])
	}
	return strings.Title(username)
}

// Helper methods for export functionality
func (s *SalesServiceImpl) exportToCSV(report *dto.SalesReportResponse) ([]byte, error) {
	// TODO: Implement CSV export
	return []byte("CSV export not yet implemented"), nil
}

func (s *SalesServiceImpl) exportToExcel(report *dto.SalesReportResponse) ([]byte, error) {
	// TODO: Implement Excel export
	return []byte("Excel export not yet implemented"), nil
}

func isValidSalesCode(code string) bool {
	// Sales code should be 2-20 characters, alphanumeric and hyphens only
	if len(code) < 2 || len(code) > 20 {
		return false
	}

	for _, char := range code {
		if !((char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
			return false
		}
	}

	return true
}