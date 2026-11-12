package dto

import (
	"gofiber-template/domain/models"
)

// ConfigResponse represents a system configuration in API responses
type ConfigResponse struct {
	ID          string                     `json:"id"`
	ConfigKey   string                     `json:"config_key"`
	ConfigValue models.ConfigValueJSON     `json:"config_value"`
	ConfigType  string                     `json:"config_type"`
	Description *string                    `json:"description,omitempty"`
	IsActive    bool                       `json:"is_active"`
	CreatedAt   string                     `json:"created_at"`
	UpdatedAt   string                     `json:"updated_at"`
}

// CreateConfigRequest represents a request to create a new system configuration
type CreateConfigRequest struct {
	ConfigKey   string                 `json:"config_key" validate:"required,max=100"`
	ConfigValue models.ConfigValueJSON `json:"config_value" validate:"required"`
	ConfigType  string                 `json:"config_type" validate:"required,oneof=general booking payment notification expiration"`
	Description *string                `json:"description,omitempty"`
	IsActive    *bool                  `json:"is_active,omitempty"`
}

// UpdateConfigRequest represents a request to update a system configuration
type UpdateConfigRequest struct {
	ConfigValue models.ConfigValueJSON `json:"config_value,omitempty"`
	ConfigType  *string                `json:"config_type,omitempty" validate:"omitempty,oneof=general booking payment notification expiration"`
	Description *string                `json:"description,omitempty"`
	IsActive    *bool                  `json:"is_active,omitempty"`
}

// ConfigUpdateRequest represents a single config update in bulk operations
type ConfigUpdateRequest struct {
	ConfigKey   string                 `json:"config_key" validate:"required"`
	ConfigValue models.ConfigValueJSON `json:"config_value" validate:"required"`
}

// ConfigListRequest represents parameters for listing configurations
type ConfigListRequest struct {
	ConfigType *string `json:"config_type,omitempty" validate:"omitempty,oneof=general booking payment notification expiration"`
	IsActive   *bool   `json:"is_active,omitempty"`
	Search     *string `json:"search,omitempty"`
	Offset     int     `json:"offset" validate:"min=0"`
	Limit      int     `json:"limit" validate:"min=1,max=100"`
}

// ConfigListResponse represents the response for configuration list
type ConfigListResponse struct {
	Configs    []*ConfigResponse `json:"configs"`
	Pagination PaginationMeta    `json:"pagination"`
}

// ExpirationTierResponse represents an expiration tier in API responses
type ExpirationTierResponse struct {
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	ThresholdDays  *int    `json:"threshold_days"`
	ExpireMinutes  *int    `json:"expire_minutes"`
	ExpireHours    *int    `json:"expire_hours"`
	ExpireDays     *int    `json:"expire_days"`
	SortOrder      int     `json:"sort_order"`
	ExampleMessage string  `json:"example_message"`
}

// ExpirationTiersResponse represents the complete expiration tiers configuration
type ExpirationTiersResponse struct {
	Tiers    []ExpirationTierResponse           `json:"tiers"`
	Metadata models.ExpirationTierMetadata      `json:"metadata"`
}

// UpdateExpirationTiersRequest represents a request to update expiration tiers
type UpdateExpirationTiersRequest struct {
	Tiers []models.ExpirationTier `json:"tiers" validate:"required,min=1"`
}

// TierCalculationRequest represents a request to calculate applicable tier
type TierCalculationRequest struct {
	TourStartDate string `json:"tour_start_date" validate:"required"`
	CreatedAt     string `json:"created_at,omitempty"`
}

// TierCalculationResponse represents the result of tier calculation
type TierCalculationResponse struct {
	ApplicableTier     ExpirationTierResponse `json:"applicable_tier"`
	DaysUntilTour      int                    `json:"days_until_tour"`
	ExpirationTime     *string                `json:"expiration_time,omitempty"`
	TimeUntilExpiry    *string                `json:"time_until_expiry,omitempty"`
	IsDefaultTier      bool                   `json:"is_default_tier"`
	CalculationDetails TierCalculationDetails `json:"calculation_details"`
}

// TierCalculationDetails provides detailed information about the calculation
type TierCalculationDetails struct {
	TourStartDate     string                   `json:"tour_start_date"`
	BookingCreatedAt  string                   `json:"booking_created_at"`
	CurrentTime       string                   `json:"current_time"`
	DaysUntilTour     int                      `json:"days_until_tour"`
	TestedTiers       []TierTestResult         `json:"tested_tiers"`
	SelectedTierName  string                   `json:"selected_tier_name"`
	ExpirationMethod  string                   `json:"expiration_method"` // "hours" or "days"
}

// TierTestResult shows how each tier was evaluated
type TierTestResult struct {
	TierName      string `json:"tier_name"`
	ThresholdDays *int   `json:"threshold_days"`
	Matches       bool   `json:"matches"`
	Reason        string `json:"reason"`
}

// ExpirationStatsResponse provides statistics about booking expirations
type ExpirationStatsResponse struct {
	TotalBookingsWithExpiration int                    `json:"total_bookings_with_expiration"`
	ExpiredBookings             int                    `json:"expired_bookings"`
	NearExpiringBookings        int                    `json:"near_expiring_bookings"`
	BookingsByTier              map[string]int         `json:"bookings_by_tier"`
	AverageExpirationTime       string                 `json:"average_expiration_time"`
	ExpirationTrends            []ExpirationTrendData  `json:"expiration_trends"`
}

// ExpirationTrendData represents trend data for expiration analysis
type ExpirationTrendData struct {
	Date          string `json:"date"`
	Created       int    `json:"created"`
	Expired       int    `json:"expired"`
	NearExpiring  int    `json:"near_expiring"`
}

// ConfigValidationResponse represents the result of configuration validation
type ConfigValidationResponse struct {
	IsValid      bool     `json:"is_valid"`
	Errors       []string `json:"errors,omitempty"`
	Warnings     []string `json:"warnings,omitempty"`
	TotalTiers   int      `json:"total_tiers"`
	DefaultTier  *string  `json:"default_tier,omitempty"`
}

// ExpirationProcessResult represents the result of processing expired bookings
type ExpirationProcessResult struct {
	ProcessedAt      string   `json:"processed_at"`
	DryRun           bool     `json:"dry_run"`
	TotalChecked     int      `json:"total_checked"`
	ExpiredFound     int      `json:"expired_found"`
	ProcessedCount   int      `json:"processed_count"`
	FailedCount      int      `json:"failed_count"`
	ProcessedBookings []ProcessedBooking `json:"processed_bookings"`
	Errors           []string `json:"errors,omitempty"`
}

// ProcessedBooking represents a booking that was processed during expiration
type ProcessedBooking struct {
	BookingID     string    `json:"booking_id"`
	InvoiceID     string    `json:"invoice_id"`
	CustomerName  string    `json:"customer_name"`
	ExpiresAt     string    `json:"expires_at"`
	ProcessedAt   string    `json:"processed_at"`
	Action        string    `json:"action"` // "canceled", "notified", "skipped"
	Reason        string    `json:"reason,omitempty"`
	Success       bool      `json:"success"`
	Error         *string   `json:"error,omitempty"`
}

// AgentBookingLimitResponse represents agent booking limits in API responses
type AgentBookingLimitResponse struct {
	GlobalSettings       AgentBookingGlobalSettingsResponse `json:"global_settings"`
	PeriodSpecificLimits []AgentBookingPeriodLimitResponse   `json:"period_specific_limits"`
	AgentGlobalBypasses  []string                            `json:"agent_global_bypasses"`
	Metadata             AgentBookingLimitMetadataResponse   `json:"metadata"`
}

// AgentBookingGlobalSettingsResponse represents global settings in API responses
type AgentBookingGlobalSettingsResponse struct {
	Enabled               bool     `json:"enabled"`
	DefaultLimitPerPeriod int      `json:"default_limit_per_period"`
	CountAllStatuses      []string `json:"count_all_statuses"`
}

// AgentBookingPeriodLimitResponse represents period-specific limits in API responses
type AgentBookingPeriodLimitResponse struct {
	PeriodID       string   `json:"period_id"`
	Limit          int      `json:"limit"`
	BypassedAgents []string `json:"bypassed_agents"`
}

// AgentBookingLimitMetadataResponse represents metadata in API responses
type AgentBookingLimitMetadataResponse struct {
	Version     string `json:"version"`
	LastUpdated string `json:"last_updated"`
	UpdatedBy   string `json:"updated_by"`
}

// UpdateAgentBookingLimitsRequest represents a request to update agent booking limits
type UpdateAgentBookingLimitsRequest struct {
	GlobalSettings       AgentBookingGlobalSettingsRequest `json:"global_settings" validate:"required"`
	PeriodSpecificLimits []AgentBookingPeriodLimitRequest   `json:"period_specific_limits"`
	AgentGlobalBypasses  []string                           `json:"agent_global_bypasses"`
}

// AgentBookingGlobalSettingsRequest represents global settings in requests
type AgentBookingGlobalSettingsRequest struct {
	Enabled               bool     `json:"enabled"`
	DefaultLimitPerPeriod int      `json:"default_limit_per_period" validate:"min=1,max=1000"`
	CountAllStatuses      []string `json:"count_all_statuses" validate:"required,min=1"`
}

// AgentBookingPeriodLimitRequest represents period-specific limits in requests
type AgentBookingPeriodLimitRequest struct {
	PeriodID       string   `json:"period_id" validate:"required,uuid"`
	Limit          int      `json:"limit" validate:"min=1,max=1000"`
	BypassedAgents []string `json:"bypassed_agents"`
}