package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SystemConfig represents system-wide configuration settings
type SystemConfig struct {
	ID           uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ConfigKey    string          `gorm:"type:varchar(100);uniqueIndex;not null" json:"config_key"`
	ConfigValue  ConfigValueJSON `gorm:"type:jsonb;not null" json:"config_value"`
	ConfigType   ConfigType      `gorm:"type:varchar(50);not null;default:'general'" json:"config_type"`
	Description  *string         `gorm:"type:text" json:"description"`
	IsActive     bool            `gorm:"not null;default:true" json:"is_active"`
	CreatedAt    time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// ConfigType defines the type of configuration
type ConfigType string

const (
	ConfigTypeGeneral      ConfigType = "general"
	ConfigTypeBooking      ConfigType = "booking"
	ConfigTypePayment      ConfigType = "payment"
	ConfigTypeNotification ConfigType = "notification"
	ConfigTypeExpiration   ConfigType = "expiration"
)

// ConfigValueJSON handles JSONB data type
type ConfigValueJSON map[string]interface{}

// Scan implements the Scanner interface for JSONB
func (c *ConfigValueJSON) Scan(value interface{}) error {
	if value == nil {
		*c = make(ConfigValueJSON)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, c)
}

// Value implements the driver Valuer interface for JSONB
func (c ConfigValueJSON) Value() (driver.Value, error) {
	if c == nil {
		return nil, nil
	}
	return json.Marshal(c)
}

// TableName sets the table name for GORM
func (SystemConfig) TableName() string {
	return "system_configs"
}

// BeforeUpdate hook to update UpdatedAt timestamp
func (sc *SystemConfig) BeforeUpdate(tx *gorm.DB) error {
	sc.UpdatedAt = time.Now()
	return nil
}

// ExpirationTier represents a tier in the expiration configuration
type ExpirationTier struct {
	Name           string  `json:"name"`
	Description    string  `json:"description"`
	ThresholdDays  *int    `json:"threshold_days"`  // Days until tour start
	ExpireMinutes  *int    `json:"expire_minutes"`  // Expiration in minutes
	ExpireHours    *int    `json:"expire_hours"`    // Expiration in hours
	ExpireDays     *int    `json:"expire_days"`     // Expiration in days
	SortOrder      int     `json:"sort_order"`      // Order for processing
}

// ExpirationTierConfig represents the complete expiration tier configuration
type ExpirationTierConfig struct {
	Tiers    []ExpirationTier        `json:"tiers"`
	Metadata ExpirationTierMetadata  `json:"metadata"`
}

// ExpirationTierMetadata contains metadata about the configuration
type ExpirationTierMetadata struct {
	Version     string    `json:"version"`
	LastUpdated time.Time `json:"last_updated"`
	UpdatedBy   string    `json:"updated_by"`
}

// GetTiers extracts expiration tiers from config value
func (sc *SystemConfig) GetTiers() ([]ExpirationTier, error) {
	if sc.ConfigKey != "booking.expiration_tiers" {
		return nil, errors.New("not an expiration tiers configuration")
	}

	tiersData, ok := sc.ConfigValue["tiers"]
	if !ok {
		return nil, errors.New("tiers data not found in configuration")
	}

	// Convert to JSON and back to get proper type
	jsonData, err := json.Marshal(tiersData)
	if err != nil {
		return nil, err
	}

	var tiers []ExpirationTier
	err = json.Unmarshal(jsonData, &tiers)
	if err != nil {
		return nil, err
	}

	return tiers, nil
}

// SetTiers updates the expiration tiers in config value
func (sc *SystemConfig) SetTiers(tiers []ExpirationTier, updatedBy string) error {
	if sc.ConfigKey != "booking.expiration_tiers" {
		return errors.New("not an expiration tiers configuration")
	}

	metadata := ExpirationTierMetadata{
		Version:     "1.0",
		LastUpdated: time.Now(),
		UpdatedBy:   updatedBy,
	}

	sc.ConfigValue = ConfigValueJSON{
		"tiers":    tiers,
		"metadata": metadata,
	}

	return nil
}

// IsExpired checks if this tier indicates the booking should be expired
func (et *ExpirationTier) IsExpired(createdAt time.Time) bool {
	now := time.Now()

	var expirationTime time.Time

	if et.ExpireMinutes != nil {
		expirationTime = createdAt.Add(time.Duration(*et.ExpireMinutes) * time.Minute)
	} else if et.ExpireHours != nil {
		expirationTime = createdAt.Add(time.Duration(*et.ExpireHours) * time.Hour)
	} else if et.ExpireDays != nil {
		expirationTime = createdAt.Add(time.Duration(*et.ExpireDays) * 24 * time.Hour)
	} else {
		// No expiration defined for this tier
		return false
	}

	return now.After(expirationTime)
}

// CalculateExpirationTime calculates when a booking should expire based on this tier
func (et *ExpirationTier) CalculateExpirationTime(createdAt time.Time) *time.Time {
	if et.ExpireMinutes != nil {
		expirationTime := createdAt.Add(time.Duration(*et.ExpireMinutes) * time.Minute)
		return &expirationTime
	} else if et.ExpireHours != nil {
		expirationTime := createdAt.Add(time.Duration(*et.ExpireHours) * time.Hour)
		return &expirationTime
	} else if et.ExpireDays != nil {
		expirationTime := createdAt.Add(time.Duration(*et.ExpireDays) * 24 * time.Hour)
		return &expirationTime
	}

	// No expiration defined
	return nil
}

// DaysUntilTour calculates days from current time to tour start date
func DaysUntilTour(tourStartDate time.Time) int {
	now := time.Now()
	duration := tourStartDate.Sub(now)
	days := int(duration.Hours() / 24)

	// Ensure we don't return negative days
	if days < 0 {
		return 0
	}

	return days
}

// MatchesTier checks if the given days until tour matches this tier's threshold
func (et *ExpirationTier) MatchesTier(daysUntilTour int) bool {
	// Default tier (no threshold) matches everything
	if et.ThresholdDays == nil {
		return true
	}

	// Check if days until tour is less than threshold
	return daysUntilTour < *et.ThresholdDays
}

// GetExpirationDuration returns the expiration duration for this tier
func (et *ExpirationTier) GetExpirationDuration() time.Duration {
	if et.ExpireMinutes != nil {
		return time.Duration(*et.ExpireMinutes) * time.Minute
	} else if et.ExpireHours != nil {
		return time.Duration(*et.ExpireHours) * time.Hour
	} else if et.ExpireDays != nil {
		return time.Duration(*et.ExpireDays) * 24 * time.Hour
	}

	// Return zero duration if no expiration defined
	return 0
}

// Validate validates the tier configuration
func (et *ExpirationTier) Validate() error {
	if et.Name == "" {
		return errors.New("tier name is required")
	}

	// Count how many expiration methods are specified
	expirationMethods := 0
	if et.ExpireMinutes != nil {
		expirationMethods++
	}
	if et.ExpireHours != nil {
		expirationMethods++
	}
	if et.ExpireDays != nil {
		expirationMethods++
	}

	if expirationMethods == 0 {
		return errors.New("one of expire_minutes, expire_hours, or expire_days must be specified")
	}

	if expirationMethods > 1 {
		return errors.New("only one of expire_minutes, expire_hours, or expire_days should be specified")
	}

	// Validate positive values
	if et.ExpireMinutes != nil && *et.ExpireMinutes <= 0 {
		return errors.New("expire_minutes must be positive")
	}

	if et.ExpireHours != nil && *et.ExpireHours <= 0 {
		return errors.New("expire_hours must be positive")
	}

	if et.ExpireDays != nil && *et.ExpireDays <= 0 {
		return errors.New("expire_days must be positive")
	}

	if et.ThresholdDays != nil && *et.ThresholdDays <= 0 {
		return errors.New("threshold_days must be positive")
	}

	return nil
}

// String returns a string representation of the tier
func (et *ExpirationTier) String() string {
	expiration := ""
	if et.ExpireMinutes != nil {
		expiration = fmt.Sprintf("%d minutes", *et.ExpireMinutes)
	} else if et.ExpireHours != nil {
		expiration = fmt.Sprintf("%d hours", *et.ExpireHours)
	} else if et.ExpireDays != nil {
		expiration = fmt.Sprintf("%d days", *et.ExpireDays)
	}

	threshold := "default"
	if et.ThresholdDays != nil {
		threshold = fmt.Sprintf("< %d days", *et.ThresholdDays)
	}

	return fmt.Sprintf("Tier %s: %s → %s", et.Name, threshold, expiration)
}

// AgentBookingLimitConfig represents the complete agent booking limit configuration
type AgentBookingLimitConfig struct {
	GlobalSettings       AgentBookingGlobalSettings    `json:"global_settings"`
	PeriodSpecificLimits []AgentBookingPeriodLimit     `json:"period_specific_limits"`
	AgentGlobalBypasses  []string                      `json:"agent_global_bypasses"`
	Metadata             AgentBookingLimitMetadata     `json:"metadata"`
}

// AgentBookingGlobalSettings contains global settings for agent booking limits
type AgentBookingGlobalSettings struct {
	Enabled                bool     `json:"enabled"`
	DefaultLimitPerPeriod  int      `json:"default_limit_per_period"`
	CountAllStatuses       []string `json:"count_all_statuses"`
}

// AgentBookingPeriodLimit represents period-specific booking limits
type AgentBookingPeriodLimit struct {
	PeriodID        string   `json:"period_id"`
	Limit           int      `json:"limit"`
	BypassedAgents  []string `json:"bypassed_agents"`
}

// AgentBookingLimitMetadata contains metadata about the configuration
type AgentBookingLimitMetadata struct {
	Version     string    `json:"version"`
	LastUpdated time.Time `json:"last_updated"`
	UpdatedBy   string    `json:"updated_by"`
}

// GetAgentBookingLimits extracts agent booking limits from config value
func (sc *SystemConfig) GetAgentBookingLimits() (*AgentBookingLimitConfig, error) {
	if sc.ConfigKey != "booking.agent_period_limits" {
		return nil, errors.New("not an agent booking limits configuration")
	}

	// Convert to JSON and back to get proper type
	jsonData, err := json.Marshal(sc.ConfigValue)
	if err != nil {
		return nil, err
	}

	var limits AgentBookingLimitConfig
	err = json.Unmarshal(jsonData, &limits)
	if err != nil {
		return nil, err
	}

	return &limits, nil
}

// SetAgentBookingLimits updates the agent booking limits in config value
func (sc *SystemConfig) SetAgentBookingLimits(limits *AgentBookingLimitConfig, updatedBy string) error {
	if sc.ConfigKey != "booking.agent_period_limits" {
		return errors.New("not an agent booking limits configuration")
	}

	limits.Metadata = AgentBookingLimitMetadata{
		Version:     "1.0",
		LastUpdated: time.Now(),
		UpdatedBy:   updatedBy,
	}

	// Convert to ConfigValueJSON
	jsonData, err := json.Marshal(limits)
	if err != nil {
		return err
	}

	var configValue ConfigValueJSON
	err = json.Unmarshal(jsonData, &configValue)
	if err != nil {
		return err
	}

	sc.ConfigValue = configValue
	return nil
}

// IsAgentBypassed checks if an agent is bypassed for booking limits
func (limits *AgentBookingLimitConfig) IsAgentBypassed(agentID string, periodID string) bool {
	// Check global bypasses
	for _, bypassedAgent := range limits.AgentGlobalBypasses {
		if bypassedAgent == agentID {
			return true
		}
	}

	// Check period-specific bypasses
	for _, periodLimit := range limits.PeriodSpecificLimits {
		if periodLimit.PeriodID == periodID {
			for _, bypassedAgent := range periodLimit.BypassedAgents {
				if bypassedAgent == agentID {
					return true
				}
			}
		}
	}

	return false
}

// GetPeriodLimit gets the booking limit for a specific period
func (limits *AgentBookingLimitConfig) GetPeriodLimit(periodID string) int {
	// Check period-specific limits first
	for _, periodLimit := range limits.PeriodSpecificLimits {
		if periodLimit.PeriodID == periodID {
			return periodLimit.Limit
		}
	}

	// Return default limit
	return limits.GlobalSettings.DefaultLimitPerPeriod
}