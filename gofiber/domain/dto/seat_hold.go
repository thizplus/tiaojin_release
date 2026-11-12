package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateSeatHoldRequest represents a request to create a seat hold
type CreateSeatHoldRequest struct {
	PeriodID  uuid.UUID `json:"period_id" validate:"required"`
	HeldSeats int       `json:"held_seats" validate:"required,min=1,max=50"`
	Reason    *string   `json:"reason,omitempty"`
	ExpiresAt *string   `json:"expires_at,omitempty"` // ISO 8601 format
}

// UpdateSeatHoldRequest represents a request to update a seat hold
type UpdateSeatHoldRequest struct {
	HeldSeats *int    `json:"held_seats,omitempty" validate:"omitempty,min=1,max=50"`
	Reason    *string `json:"reason,omitempty"`
	ExpiresAt *string `json:"expires_at,omitempty"` // ISO 8601 format, or "never" to remove expiry
}

// SeatHoldResponse represents a seat hold in responses
type SeatHoldResponse struct {
	ID        uuid.UUID `json:"id"`
	PeriodID  uuid.UUID `json:"period_id"`
	AdminID   uuid.UUID `json:"admin_id"`
	HeldSeats int       `json:"held_seats"`
	Reason    *string   `json:"reason"`
	ExpiresAt *string   `json:"expires_at,omitempty"`
	Status    string    `json:"status"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`

	// Optional relations
	PeriodDetails *PeriodSummaryResponse `json:"period_details,omitempty"`
	AdminDetails  *UserResponse          `json:"admin_details,omitempty"`

	// Calculated fields
	TimeRemaining *string `json:"time_remaining,omitempty"` // Human readable like "2h 30m"
	IsExpired     bool    `json:"is_expired"`
	CanExtend     bool    `json:"can_extend"`
	CanRelease    bool    `json:"can_release"`
}

// SeatHoldListResponse represents a paginated list of seat holds
type SeatHoldListResponse struct {
	SeatHolds  []*SeatHoldResponse `json:"seat_holds"`
	Pagination PaginationMeta       `json:"pagination"`
}

// SeatAvailabilityResponse represents seat availability information
type SeatAvailabilityResponse struct {
	PeriodID         uuid.UUID `json:"period_id"`
	TotalSeats       int       `json:"total_seats"`
	BookedSeats      int       `json:"booked_seats"`
	HeldSeats        int       `json:"held_seats"`
	AvailableSeats   int       `json:"available_seats"`
	RequestedSeats   int       `json:"requested_seats"`
	CanHold          bool      `json:"can_hold"`
	MaxHoldable      int       `json:"max_holdable"`
	RecommendedHold  int       `json:"recommended_hold"`
	PercentAvailable float64   `json:"percent_available"`
}

// ExpireHoldsResponse represents the result of expiring old holds
type ExpireHoldsResponse struct {
	ExpiredCount     int                  `json:"expired_count"`
	ExpiredHolds     []*SeatHoldResponse  `json:"expired_holds,omitempty"`
	TotalSeatsFreed  int                  `json:"total_seats_freed"`
	AffectedPeriods  []uuid.UUID          `json:"affected_periods"`
	ProcessedAt      string               `json:"processed_at"`
	Message          string               `json:"message"`
}

// ReleaseHoldsResponse represents the result of releasing holds
type ReleaseHoldsResponse struct {
	ReleasedCount    int                  `json:"released_count"`
	ReleasedHolds    []*SeatHoldResponse  `json:"released_holds,omitempty"`
	TotalSeatsFreed  int                  `json:"total_seats_freed"`
	AffectedPeriods  []uuid.UUID          `json:"affected_periods"`
	Reason           string               `json:"reason"`
	ProcessedAt      string               `json:"processed_at"`
	Message          string               `json:"message"`
}

// SeatHoldSummaryResponse represents aggregated seat hold information for a period
type SeatHoldSummaryResponse struct {
	PeriodID             uuid.UUID `json:"period_id"`
	TotalActiveHolds     int       `json:"total_active_holds"`
	TotalHeldSeats       int       `json:"total_held_seats"`
	TotalExpiredHolds    int       `json:"total_expired_holds"`
	TotalReleasedHolds   int       `json:"total_released_holds"`
	TotalConvertedHolds  int       `json:"total_converted_holds"`
	EarliestExpiry       *string   `json:"earliest_expiry,omitempty"`
	LatestExpiry         *string   `json:"latest_expiry,omitempty"`
	ExpiringWithin24H    int       `json:"expiring_within_24h"`
	ExpiringWithin1H     int       `json:"expiring_within_1h"`
	PermanentHolds       int       `json:"permanent_holds"` // Holds without expiry
	AverageHoldDuration  *string   `json:"average_hold_duration,omitempty"`
}

// AdminSeatHoldStatsResponse represents statistics for an admin's seat holds
type AdminSeatHoldStatsResponse struct {
	AdminID              uuid.UUID                  `json:"admin_id"`
	TotalHolds           int64                      `json:"total_holds"`
	ActiveHolds          int64                      `json:"active_holds"`
	ReleasedHolds        int64                      `json:"released_holds"`
	ExpiredHolds         int64                      `json:"expired_holds"`
	ConvertedHolds       int64                      `json:"converted_holds"`
	TotalSeatsHeld       int64                      `json:"total_seats_held"`
	AverageSeatsPerHold  float64                    `json:"average_seats_per_hold"`
	ConversionRate       float64                    `json:"conversion_rate"` // converted/total ratio
	ExpirationRate       float64                    `json:"expiration_rate"` // expired/total ratio
	HoldsByPeriod        map[string]int             `json:"holds_by_period,omitempty"`
	RecentActivity       []*SeatHoldResponse        `json:"recent_activity,omitempty"`
	PerformanceMetrics   map[string]interface{}     `json:"performance_metrics"`
}

// SeatHoldAnalyticsRequest represents a request for seat hold analytics
type SeatHoldAnalyticsRequest struct {
	StartDate    *string     `json:"start_date,omitempty"`    // ISO 8601 format
	EndDate      *string     `json:"end_date,omitempty"`      // ISO 8601 format
	PeriodIDs    []uuid.UUID `json:"period_ids,omitempty"`
	AdminIDs     []uuid.UUID `json:"admin_ids,omitempty"`
	Status       *string     `json:"status,omitempty"`
	IncludeStats bool        `json:"include_stats" default:"true"`
	GroupBy      *string     `json:"group_by,omitempty"` // period, admin, date, status
}

// SeatHoldAnalyticsResponse represents analytics data for seat holds
type SeatHoldAnalyticsResponse struct {
	Summary struct {
		TotalHolds      int64   `json:"total_holds"`
		TotalSeatsHeld  int64   `json:"total_seats_held"`
		AverageHoldTime string  `json:"average_hold_time"`
		ConversionRate  float64 `json:"conversion_rate"`
		ExpirationRate  float64 `json:"expiration_rate"`
	} `json:"summary"`

	ByPeriod map[string]struct {
		PeriodID   uuid.UUID `json:"period_id"`
		HoldCount  int       `json:"hold_count"`
		SeatCount  int       `json:"seat_count"`
		PeriodName string    `json:"period_name,omitempty"`
	} `json:"by_period,omitempty"`

	ByAdmin map[string]struct {
		AdminID     uuid.UUID `json:"admin_id"`
		HoldCount   int       `json:"hold_count"`
		SeatCount   int       `json:"seat_count"`
		AdminName   string    `json:"admin_name,omitempty"`
		Performance string    `json:"performance"` // excellent, good, average, poor
	} `json:"by_admin,omitempty"`

	ByStatus map[string]int `json:"by_status"`

	Timeline []struct {
		Date      string `json:"date"`
		HoldCount int    `json:"hold_count"`
		SeatCount int    `json:"seat_count"`
	} `json:"timeline,omitempty"`

	Trends struct {
		HoldGrowthRate     float64 `json:"hold_growth_rate"`     // % change
		SeatGrowthRate     float64 `json:"seat_growth_rate"`     // % change
		EfficiencyTrend    string  `json:"efficiency_trend"`     // improving, stable, declining
		PeakHoldHour       int     `json:"peak_hold_hour"`       // 0-23
		PeakHoldDay        string  `json:"peak_hold_day"`        // monday, tuesday, etc.
		PopularHoldSizes   []int   `json:"popular_hold_sizes"`   // most common seat counts
	} `json:"trends"`

	Recommendations []string `json:"recommendations,omitempty"`
	GeneratedAt     string   `json:"generated_at"`
	DataRange       struct {
		From string `json:"from"`
		To   string `json:"to"`
	} `json:"data_range"`
}

// Utility DTOs for internal operations
type SeatHoldFilter struct {
	PeriodIDs     []uuid.UUID `json:"period_ids,omitempty"`
	AdminIDs      []uuid.UUID `json:"admin_ids,omitempty"`
	Status        *string     `json:"status,omitempty"`
	CreatedAfter  *time.Time  `json:"created_after,omitempty"`
	CreatedBefore *time.Time  `json:"created_before,omitempty"`
	ExpiringBefore *time.Time `json:"expiring_before,omitempty"`
	MinSeats      *int        `json:"min_seats,omitempty"`
	MaxSeats      *int        `json:"max_seats,omitempty"`
}

