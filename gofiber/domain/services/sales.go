package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)


type SalesService interface {
	CreateSales(ctx context.Context, req dto.CreateSalesRequest) (*dto.SalesResponse, error)
	GetSalesByID(ctx context.Context, salesID uuid.UUID) (*dto.SalesResponse, error)
	GetSalesByCode(ctx context.Context, code string) (*dto.SalesResponse, error)
	UpdateSales(ctx context.Context, salesID uuid.UUID, req dto.UpdateSalesRequest) (*dto.SalesResponse, error)
	DeleteSales(ctx context.Context, salesID uuid.UUID) error
	ListSales(ctx context.Context, req dto.SalesListRequest) (*dto.SalesListResponse, error)
	GetActiveSales(ctx context.Context, offset, limit int) ([]*dto.SalesResponse, error)
	SearchSales(ctx context.Context, query string, offset, limit int) ([]*dto.SalesResponse, error)

	// Status Management
	ActivateSales(ctx context.Context, salesID uuid.UUID) (*dto.SalesResponse, error)
	DeactivateSales(ctx context.Context, salesID uuid.UUID) (*dto.SalesResponse, error)
	UpdateActiveStatus(ctx context.Context, salesID uuid.UUID, active bool) (*dto.SalesResponse, error)

	// Performance & Statistics
	GetSalesWithBookingCounts(ctx context.Context, offset, limit int) ([]*dto.SalesWithStatsResponse, error)
	GetSalesPerformance(ctx context.Context, req dto.SalesPerformanceRequest) (*dto.SalesPerformanceResponse, error)
	GetSalesCommissionReport(ctx context.Context, salesID uuid.UUID, startDate, endDate *time.Time) (*dto.SalesCommissionReport, error)
	GetMonthlyPerformance(ctx context.Context, salesID uuid.UUID, year int) ([]*dto.MonthlyPerformance, error)
	GetTopPerformingSales(ctx context.Context, startDate, endDate *time.Time, limit int) ([]*dto.SalesWithStatsResponse, error)

	// Business Logic
	CalculateCommission(ctx context.Context, salesID uuid.UUID, startDate, endDate *time.Time) (float64, error)
	ReassignBookings(ctx context.Context, fromSalesID, toSalesID uuid.UUID) (int64, error)
	GetSalesStats(ctx context.Context) (*dto.SalesStatsResponse, error)
	ValidateSalesCode(ctx context.Context, code string, excludeSalesID *uuid.UUID) (bool, error)

	// Additional Analytics & Reporting Methods
	GetRevenueStats(ctx context.Context, req dto.RevenueStatsRequest) (*dto.RevenueStatsResponse, error)
	GetTopTours(ctx context.Context, req dto.TopToursRequest) (*dto.TopToursResponse, error)
	GetCustomerStats(ctx context.Context, req dto.CustomerStatsRequest) (*dto.CustomerStatsResponse, error)
	GetSalesReport(ctx context.Context, req dto.SalesReportRequest) (*dto.SalesReportResponse, error)
	ExportSalesReport(ctx context.Context, req dto.ExportSalesReportRequest) ([]byte, string, error)
	GetDashboardData(ctx context.Context) (*dto.DashboardDataResponse, error)

	// User-Sales Integration Methods
	CreateSalesFromUser(ctx context.Context, req dto.CreateSalesFromUserRequest) (*dto.SalesResponse, error)
	SyncUserToSales(ctx context.Context, req dto.SalesUserSyncRequest) (*dto.SalesUserSyncResponse, error)
	GetSalesByUserID(ctx context.Context, userID uuid.UUID) (*dto.SalesResponse, error)
	GetSalesBookingStats(ctx context.Context, salesID uuid.UUID) (*dto.SalesBookingStatsResponse, error)

	// Sales Profile Management
	GetMyProfile(ctx context.Context, userID uuid.UUID) (*dto.SalesMyProfileResponse, error)
	UpdateMyProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateSalesProfileRequest) (*dto.SalesMyProfileResponse, error)
}