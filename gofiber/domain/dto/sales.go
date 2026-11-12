package dto

import (
	"time"

	"github.com/google/uuid"
)

type SalesResponse struct {
	ID        uuid.UUID  `json:"id"`
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	Code      string     `json:"code"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	NickName  string     `json:"nick_name"`
	Active    bool       `json:"active"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`

	// User relationship data
	UserEmail *string `json:"user_email,omitempty"`
	UserRole  *string `json:"user_role,omitempty"`
}

type CreateSalesRequest struct {
	UserID    *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,uuid"`
	Code      string     `json:"code" validate:"required,max=20,uppercase"`
	FirstName string     `json:"first_name" validate:"required,max=100"`
	LastName  string     `json:"last_name" validate:"required,max=100"`
	NickName  string     `json:"nick_name" validate:"required,max=50"`
}

type UpdateSalesRequest struct {
	UserID    *uuid.UUID `json:"user_id,omitempty" validate:"omitempty,uuid"`
	FirstName *string    `json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName  *string    `json:"last_name,omitempty" validate:"omitempty,max=100"`
	NickName  *string    `json:"nick_name,omitempty" validate:"omitempty,max=50"`
	Active    *bool      `json:"active,omitempty"`
}

type SalesListRequest struct {
	Active *bool   `json:"active,omitempty"`
	Search *string `json:"search,omitempty"`
	Offset int     `json:"offset" validate:"min=0"`
	Limit  int     `json:"limit" validate:"min=1,max=100"`
}

type SalesListResponse struct {
	Sales      []*SalesResponse `json:"sales"`
	Pagination PaginationMeta   `json:"pagination"`
}

type SalesWithStatsResponse struct {
	*SalesResponse
	BookingStats *SalesBookingStats `json:"booking_stats,omitempty"`
}

type SalesBookingStats struct {
	TotalBookings       int64   `json:"total_bookings"`
	ConfirmedBookings   int64   `json:"confirmed_bookings"`
	TotalRevenue        float64 `json:"total_revenue"`
	TotalCommission     float64 `json:"total_commission"`
	AverageBookingValue float64 `json:"average_booking_value"`
}

type SalesPerformanceRequest struct {
	SalesID   *uuid.UUID `json:"sales_id,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Offset    int        `json:"offset" validate:"min=0"`
	Limit     int        `json:"limit" validate:"min=1,max=100"`
}

type SalesPerformanceResponse struct {
	Sales           []*SalesWithStatsResponse `json:"sales"`
	Pagination      PaginationMeta            `json:"pagination"`
	TotalRevenue    float64                   `json:"total_revenue"`
	TotalCommission float64                   `json:"total_commission"`
	Period          *PerformancePeriod        `json:"period,omitempty"`
}

type PerformancePeriod struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type MonthlyPerformance struct {
	Year       int     `json:"year"`
	Month      int     `json:"month"`
	Bookings   int64   `json:"bookings"`
	Revenue    float64 `json:"revenue"`
	Commission float64 `json:"commission"`
	Passengers int64   `json:"passengers"`
}

type SalesCommissionReport struct {
	SalesID         uuid.UUID             `json:"sales_id"`
	SalesCode       string                `json:"sales_code"`
	SalesName       string                `json:"sales_name"`
	TotalCommission float64               `json:"total_commission"`
	BookingCount    int64                 `json:"booking_count"`
	MonthlyDetails  []*MonthlyPerformance `json:"monthly_details"`
}

type SalesStatsResponse struct {
	TotalSales    int64 `json:"total_sales"`
	ActiveSales   int64 `json:"active_sales"`
	InactiveSales int64 `json:"inactive_sales"`
}

// Revenue Analytics DTOs
type RevenueStatsRequest struct {
	GroupBy  string     `json:"groupBy,omitempty" validate:"omitempty,oneof=day week month year"`
	DateFrom *time.Time `json:"dateFrom,omitempty"`
	DateTo   *time.Time `json:"dateTo,omitempty"`
}

type RevenueStatsResponse struct {
	TotalRevenue     float64                `json:"total_revenue"`
	TotalCommission  float64                `json:"total_commission"`
	BookingCount     int64                  `json:"booking_count"`
	AverageBooking   float64                `json:"average_booking"`
	Period           *PerformancePeriod     `json:"period,omitempty"`
	MonthlyBreakdown []*MonthlyPerformance  `json:"monthly_breakdown,omitempty"`
}

// Top Tours DTOs
type TopToursRequest struct {
	DateFrom *time.Time `json:"dateFrom,omitempty"`
	DateTo   *time.Time `json:"dateTo,omitempty"`
	Limit    int        `json:"limit" validate:"min=1,max=50"`
}

type TopToursResponse struct {
	Tours  []*TopTourStats    `json:"tours"`
	Period *PerformancePeriod `json:"period,omitempty"`
}

type TopTourStats struct {
	TourID      uuid.UUID `json:"tour_id"`
	TourName    string    `json:"tour_name"`
	BookingCount int64    `json:"booking_count"`
	TotalRevenue float64  `json:"total_revenue"`
	TotalCommission float64 `json:"total_commission"`
	AveragePrice float64  `json:"average_price"`
}

// Customer Analytics DTOs
type CustomerStatsRequest struct {
	DateFrom *time.Time `json:"dateFrom,omitempty"`
	DateTo   *time.Time `json:"dateTo,omitempty"`
}

type CustomerStatsResponse struct {
	TotalCustomers    int64   `json:"total_customers"`
	NewCustomers      int64   `json:"new_customers"`
	ReturningCustomers int64  `json:"returning_customers"`
	AverageOrderValue float64 `json:"average_order_value"`
	CustomerRetentionRate float64 `json:"customer_retention_rate"`
	Period            *PerformancePeriod `json:"period,omitempty"`
}

// Sales Report DTOs
type SalesReportRequest struct {
	DateFrom *time.Time `json:"dateFrom,omitempty"`
	DateTo   *time.Time `json:"dateTo,omitempty"`
	TourID   *uuid.UUID `json:"tourId,omitempty"`
	UserID   *uuid.UUID `json:"userId,omitempty"`
	Status   *string    `json:"status,omitempty" validate:"omitempty,oneof=pending confirmed cancelled"`
	Offset   int        `json:"offset" validate:"min=0"`
	Limit    int        `json:"limit" validate:"min=1,max=100"`
}

type SalesReportResponse struct {
	ReportType       string                 `json:"report_type"`
	Period           *PerformancePeriod     `json:"period"`
	Summary          *SalesReportSummary    `json:"summary"`
	SalesPerformance []*SalesWithStatsResponse `json:"sales_performance,omitempty"`
	MonthlyData      []*MonthlyPerformance  `json:"monthly_data,omitempty"`
	TopTours         []*TopTourStats        `json:"top_tours,omitempty"`
	GeneratedAt      time.Time              `json:"generated_at"`
}

type SalesReportSummary struct {
	TotalRevenue      float64 `json:"total_revenue"`
	TotalCommission   float64 `json:"total_commission"`
	TotalBookings     int64   `json:"total_bookings"`
	ConfirmedBookings int64   `json:"confirmed_bookings"`
	AverageBookingValue float64 `json:"average_booking_value"`
	ActiveSalesCount  int64   `json:"active_sales_count"`
	TopPerformerId    *uuid.UUID `json:"top_performer_id,omitempty"`
	TopPerformerName  string  `json:"top_performer_name,omitempty"`
}

// Export Report DTOs
type ExportSalesReportRequest struct {
	DateFrom *time.Time `json:"dateFrom,omitempty"`
	DateTo   *time.Time `json:"dateTo,omitempty"`
	TourID   *uuid.UUID `json:"tourId,omitempty"`
	UserID   *uuid.UUID `json:"userId,omitempty"`
	Status   *string    `json:"status,omitempty" validate:"omitempty,oneof=pending confirmed cancelled"`
	Format   string     `json:"format" validate:"required,oneof=excel csv"`
}

// Dashboard DTOs
type DashboardDataResponse struct {
	SalesStats       *SalesStatsResponse    `json:"sales_stats"`
	RecentRevenue    *RevenueStatsResponse  `json:"recent_revenue"`
	TopPerformers    []*SalesWithStatsResponse `json:"top_performers"`
	TopTours         []*TopTourStats        `json:"top_tours"`
	MonthlyTrends    []*MonthlyPerformance  `json:"monthly_trends"`
	CustomerStats    *CustomerStatsResponse `json:"customer_stats"`
	RecentActivity   []*RecentActivityItem  `json:"recent_activity"`
	LastUpdated      time.Time              `json:"last_updated"`
}

type RecentActivityItem struct {
	ID          uuid.UUID `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	SalesID     *uuid.UUID `json:"sales_id,omitempty"`
	SalesName   string     `json:"sales_name,omitempty"`
	Amount      *float64   `json:"amount,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// User-Sales relationship DTOs
type CreateSalesFromUserRequest struct {
	UserID    uuid.UUID `json:"user_id" validate:"required,uuid"`
	Code      string    `json:"code" validate:"required,max=20,uppercase"`
	FirstName string    `json:"first_name" validate:"required,max=100"`
	LastName  string    `json:"last_name" validate:"required,max=100"`
	NickName  string    `json:"nick_name" validate:"required,max=50"`
}

type SalesUserSyncRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
}

type SalesUserSyncResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Sales   *SalesResponse `json:"sales,omitempty"`
}

type SalesMyProfileResponse struct {
	*SalesResponse
	CanManageBookings bool `json:"can_manage_bookings"`
	BookingStats      *SalesBookingStats `json:"booking_stats,omitempty"`
}

type SalesBookingStatsResponse struct {
	SalesID           uuid.UUID `json:"sales_id"`
	TotalBookings     int64     `json:"total_bookings"`
	PendingBookings   int64     `json:"pending_bookings"`
	ConfirmedBookings int64     `json:"confirmed_bookings"`
	CancelledBookings int64     `json:"cancelled_bookings"`
	TotalRevenue      float64   `json:"total_revenue"`
	MonthlyRevenue    float64   `json:"monthly_revenue"`
}

type UpdateSalesProfileRequest struct {
	FirstName *string `json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName  *string `json:"last_name,omitempty" validate:"omitempty,max=100"`
	NickName  *string `json:"nick_name,omitempty" validate:"omitempty,max=50"`
}