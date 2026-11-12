package serviceimpl

import (
	"context"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type ConfigServiceImpl struct {
	configRepo repositories.SystemConfigRepository
	bookingRepo repositories.BookingRepository
}

func NewConfigService(configRepo repositories.SystemConfigRepository, bookingRepo repositories.BookingRepository) services.ConfigService {
	return &ConfigServiceImpl{
		configRepo: configRepo,
		bookingRepo: bookingRepo,
	}
}

// Basic configuration operations
func (s *ConfigServiceImpl) GetConfig(ctx context.Context, configKey string) (*models.SystemConfig, error) {
	log.Printf("[DEBUG] Getting config for key: %s", configKey)

	config, err := s.configRepo.GetByKey(ctx, configKey)
	if err != nil {
		log.Printf("[ERROR] Failed to get config %s: %v", configKey, err)
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	if config == nil {
		log.Printf("[WARNING] Config not found for key: %s", configKey)
		return nil, fmt.Errorf("config not found: %s", configKey)
	}

	log.Printf("[SUCCESS] Config retrieved for key: %s", configKey)
	return config, nil
}

func (s *ConfigServiceImpl) SetConfig(ctx context.Context, configKey string, configValue models.ConfigValueJSON, configType models.ConfigType, description string) error {
	log.Printf("[DEBUG] Setting config for key: %s, type: %s", configKey, configType)

	config := &models.SystemConfig{
		ConfigKey:   configKey,
		ConfigValue: configValue,
		ConfigType:  configType,
		Description: &description,
		IsActive:    true,
	}

	err := s.configRepo.Create(ctx, config)
	if err != nil {
		log.Printf("[ERROR] Failed to set config %s: %v", configKey, err)
		return fmt.Errorf("failed to set config: %w", err)
	}

	log.Printf("[SUCCESS] Config set for key: %s", configKey)
	return nil
}

func (s *ConfigServiceImpl) UpdateConfig(ctx context.Context, configKey string, configValue models.ConfigValueJSON) error {
	log.Printf("[DEBUG] Updating config for key: %s", configKey)

	err := s.configRepo.UpdateByKey(ctx, configKey, configValue)
	if err != nil {
		log.Printf("[ERROR] Failed to update config %s: %v", configKey, err)
		return fmt.Errorf("failed to update config: %w", err)
	}

	log.Printf("[SUCCESS] Config updated for key: %s", configKey)
	return nil
}

func (s *ConfigServiceImpl) DeleteConfig(ctx context.Context, configKey string) error {
	log.Printf("[DEBUG] Deleting config for key: %s", configKey)

	// First get the config to get its ID
	config, err := s.configRepo.GetByKey(ctx, configKey)
	if err != nil {
		return fmt.Errorf("failed to find config: %w", err)
	}

	if config == nil {
		return fmt.Errorf("config not found: %s", configKey)
	}

	err = s.configRepo.Delete(ctx, config.ID.String())
	if err != nil {
		log.Printf("[ERROR] Failed to delete config %s: %v", configKey, err)
		return fmt.Errorf("failed to delete config: %w", err)
	}

	log.Printf("[SUCCESS] Config deleted for key: %s", configKey)
	return nil
}

func (s *ConfigServiceImpl) ListConfigs(ctx context.Context, configType *models.ConfigType, offset, limit int) ([]*models.SystemConfig, error) {
	log.Printf("[DEBUG] Listing configs - type: %v, offset: %d, limit: %d", configType, offset, limit)

	var configs []*models.SystemConfig
	var err error

	if configType != nil {
		configs, err = s.configRepo.GetByType(ctx, *configType, offset, limit)
	} else {
		configs, err = s.configRepo.List(ctx, offset, limit)
	}

	if err != nil {
		log.Printf("[ERROR] Failed to list configs: %v", err)
		return nil, fmt.Errorf("failed to list configs: %w", err)
	}

	log.Printf("[SUCCESS] Listed %d configs", len(configs))
	return configs, nil
}

// Expiration tier specific operations
func (s *ConfigServiceImpl) GetExpirationTiers(ctx context.Context) ([]models.ExpirationTier, error) {
	log.Printf("[DEBUG] Getting expiration tiers configuration")

	config, err := s.configRepo.GetByKey(ctx, "booking.expiration_tiers")
	if err != nil {
		log.Printf("[ERROR] Failed to get expiration tiers config: %v", err)
		return nil, fmt.Errorf("failed to get expiration tiers: %w", err)
	}

	if config == nil {
		log.Printf("[ERROR] Expiration tiers configuration not found")
		return nil, fmt.Errorf("expiration tiers configuration not found")
	}

	tiers, err := config.GetTiers()
	if err != nil {
		log.Printf("[ERROR] Failed to parse expiration tiers: %v", err)
		return nil, fmt.Errorf("failed to parse expiration tiers: %w", err)
	}

	// Sort tiers by sort_order
	sort.Slice(tiers, func(i, j int) bool {
		return tiers[i].SortOrder < tiers[j].SortOrder
	})

	log.Printf("[SUCCESS] Retrieved %d expiration tiers", len(tiers))
	return tiers, nil
}

func (s *ConfigServiceImpl) UpdateExpirationTiers(ctx context.Context, tiers []models.ExpirationTier, updatedBy string) error {
	log.Printf("[DEBUG] Updating expiration tiers - count: %d, updated by: %s", len(tiers), updatedBy)

	// Validate the tiers first
	if err := s.ValidateExpirationConfig(ctx, tiers); err != nil {
		log.Printf("[ERROR] Expiration tiers validation failed: %v", err)
		return fmt.Errorf("validation failed: %w", err)
	}

	// Get existing config
	config, err := s.configRepo.GetByKey(ctx, "booking.expiration_tiers")
	if err != nil {
		log.Printf("[ERROR] Failed to get existing expiration tiers config: %v", err)
		return fmt.Errorf("failed to get existing config: %w", err)
	}

	if config == nil {
		log.Printf("[ERROR] Expiration tiers configuration not found")
		return fmt.Errorf("expiration tiers configuration not found")
	}

	// Update the tiers
	err = config.SetTiers(tiers, updatedBy)
	if err != nil {
		log.Printf("[ERROR] Failed to set tiers in config: %v", err)
		return fmt.Errorf("failed to set tiers: %w", err)
	}

	// Save the updated config
	err = s.configRepo.UpdateByKey(ctx, "booking.expiration_tiers", config.ConfigValue)
	if err != nil {
		log.Printf("[ERROR] Failed to save updated expiration tiers: %v", err)
		return fmt.Errorf("failed to save updated tiers: %w", err)
	}

	log.Printf("[SUCCESS] Expiration tiers updated successfully")
	return nil
}

func (s *ConfigServiceImpl) ValidateExpirationConfig(ctx context.Context, tiers []models.ExpirationTier) error {
	log.Printf("[DEBUG] Validating expiration configuration - tier count: %d", len(tiers))

	if len(tiers) == 0 {
		return fmt.Errorf("at least one tier must be defined")
	}

	var hasDefaultTier bool
	tierNames := make(map[string]bool)

	for i, tier := range tiers {
		// Validate individual tier
		if err := tier.Validate(); err != nil {
			return fmt.Errorf("tier %d (%s) validation failed: %w", i+1, tier.Name, err)
		}

		// Check for duplicate names
		if tierNames[tier.Name] {
			return fmt.Errorf("duplicate tier name: %s", tier.Name)
		}
		tierNames[tier.Name] = true

		// Check for default tier
		if tier.ThresholdDays == nil {
			if hasDefaultTier {
				return fmt.Errorf("only one default tier (no threshold) is allowed")
			}
			hasDefaultTier = true
		}
	}

	if !hasDefaultTier {
		return fmt.Errorf("a default tier (no threshold_days) must be defined")
	}

	log.Printf("[SUCCESS] Expiration configuration validation passed")
	return nil
}

func (s *ConfigServiceImpl) GetDefaultExpirationTier(ctx context.Context) (*models.ExpirationTier, error) {
	log.Printf("[DEBUG] Getting default expiration tier")

	tiers, err := s.GetExpirationTiers(ctx)
	if err != nil {
		return nil, err
	}

	for _, tier := range tiers {
		if tier.ThresholdDays == nil {
			log.Printf("[SUCCESS] Found default tier: %s", tier.Name)
			return &tier, nil
		}
	}

	log.Printf("[ERROR] No default tier found")
	return nil, fmt.Errorf("no default tier found")
}

// Tier calculation operations
func (s *ConfigServiceImpl) CalculateTieredExpiration(ctx context.Context, tourStartDate, createdAt string) (*models.ExpirationTier, error) {
	log.Printf("[DEBUG] Calculating tiered expiration - tour: %s, created: %s", tourStartDate, createdAt)

	// Parse tour start date
	tourDate, err := time.Parse("2006-01-02", tourStartDate)
	if err != nil {
		log.Printf("[ERROR] Invalid tour start date format: %s", tourStartDate)
		return nil, fmt.Errorf("invalid tour start date format: %w", err)
	}

	// Calculate days until tour
	daysUntilTour := models.DaysUntilTour(tourDate)
	log.Printf("[DEBUG] Days until tour: %d", daysUntilTour)

	return s.GetApplicableTier(ctx, daysUntilTour)
}

func (s *ConfigServiceImpl) GetApplicableTier(ctx context.Context, daysUntilTour int) (*models.ExpirationTier, error) {
	log.Printf("[DEBUG] Getting applicable tier for %d days until tour", daysUntilTour)

	tiers, err := s.GetExpirationTiers(ctx)
	if err != nil {
		return nil, err
	}

	// Loop through tiers from lowest sort_order to highest
	for _, tier := range tiers {
		if tier.MatchesTier(daysUntilTour) {
			log.Printf("[SUCCESS] Found applicable tier: %s for %d days", tier.Name, daysUntilTour)
			return &tier, nil
		}
	}

	// This should not happen if we have a default tier, but fallback
	defaultTier, err := s.GetDefaultExpirationTier(ctx)
	if err != nil {
		log.Printf("[ERROR] No applicable tier found and no default tier available")
		return nil, fmt.Errorf("no applicable tier found")
	}

	log.Printf("[WARNING] Using fallback default tier: %s", defaultTier.Name)
	return defaultTier, nil
}

// Bulk operations
func (s *ConfigServiceImpl) GetMultipleConfigs(ctx context.Context, keys []string) (map[string]*models.SystemConfig, error) {
	log.Printf("[DEBUG] Getting multiple configs - keys: %v", keys)

	configs, err := s.configRepo.GetMultipleByKeys(ctx, keys)
	if err != nil {
		log.Printf("[ERROR] Failed to get multiple configs: %v", err)
		return nil, fmt.Errorf("failed to get multiple configs: %w", err)
	}

	// Convert to map
	configMap := make(map[string]*models.SystemConfig)
	for _, config := range configs {
		configMap[config.ConfigKey] = config
	}

	log.Printf("[SUCCESS] Retrieved %d configs out of %d requested", len(configMap), len(keys))
	return configMap, nil
}

// GetConfigByID gets a configuration by ID
func (s *ConfigServiceImpl) GetConfigByID(ctx context.Context, configID string) (*models.SystemConfig, error) {
	log.Printf("[DEBUG] Getting config by ID: %s", configID)

	config, err := s.configRepo.GetByID(ctx, configID)
	if err != nil {
		log.Printf("[ERROR] Failed to get config by ID %s: %v", configID, err)
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	if config == nil {
		log.Printf("[ERROR] Config not found: %s", configID)
		return nil, fmt.Errorf("config not found")
	}

	log.Printf("[SUCCESS] Retrieved config: %s", config.ConfigKey)
	return config, nil
}

// DeleteConfigByID deletes a configuration by ID
func (s *ConfigServiceImpl) DeleteConfigByID(ctx context.Context, configID string) error {
	log.Printf("[DEBUG] Deleting config by ID: %s", configID)

	// Check if config exists first
	config, err := s.configRepo.GetByID(ctx, configID)
	if err != nil {
		log.Printf("[ERROR] Failed to get config for deletion %s: %v", configID, err)
		return fmt.Errorf("failed to get config: %w", err)
	}

	if config == nil {
		log.Printf("[ERROR] Config not found for deletion: %s", configID)
		return fmt.Errorf("config not found")
	}

	// Don't allow deletion of critical configs
	criticalConfigs := []string{"booking.expiration_tiers"}
	for _, critical := range criticalConfigs {
		if config.ConfigKey == critical {
			log.Printf("[ERROR] Attempted to delete critical config: %s", config.ConfigKey)
			return fmt.Errorf("cannot delete critical configuration: %s", config.ConfigKey)
		}
	}

	err = s.configRepo.Delete(ctx, configID)
	if err != nil {
		log.Printf("[ERROR] Failed to delete config %s: %v", configID, err)
		return fmt.Errorf("failed to delete config: %w", err)
	}

	log.Printf("[SUCCESS] Deleted config: %s (%s)", config.ConfigKey, configID)
	return nil
}

// UpdateConfigByID updates a configuration by ID
func (s *ConfigServiceImpl) UpdateConfigByID(ctx context.Context, configID string, req dto.UpdateConfigRequest) error {
	log.Printf("[DEBUG] Updating config by ID: %s", configID)

	// Get existing config
	config, err := s.configRepo.GetByID(ctx, configID)
	if err != nil {
		log.Printf("[ERROR] Failed to get config for update %s: %v", configID, err)
		return fmt.Errorf("failed to get config: %w", err)
	}

	if config == nil {
		log.Printf("[ERROR] Config not found for update: %s", configID)
		return fmt.Errorf("config not found")
	}

	// Update fields if provided
	if req.ConfigValue != nil {
		config.ConfigValue = req.ConfigValue
	}
	if req.ConfigType != nil {
		config.ConfigType = models.ConfigType(*req.ConfigType)
	}
	if req.Description != nil {
		config.Description = req.Description
	}
	if req.IsActive != nil {
		config.IsActive = *req.IsActive
	}

	// Update in database
	err = s.configRepo.Update(ctx, configID, config)
	if err != nil {
		log.Printf("[ERROR] Failed to update config %s: %v", configID, err)
		return fmt.Errorf("failed to update config: %w", err)
	}

	log.Printf("[SUCCESS] Updated config: %s (%s)", config.ConfigKey, configID)
	return nil
}

func (s *ConfigServiceImpl) BulkUpdateConfigs(ctx context.Context, updates []dto.ConfigUpdateRequest) error {
	log.Printf("[DEBUG] Bulk updating %d configs", len(updates))

	for i, update := range updates {
		err := s.UpdateConfig(ctx, update.ConfigKey, update.ConfigValue)
		if err != nil {
			log.Printf("[ERROR] Failed to update config %s (item %d): %v", update.ConfigKey, i+1, err)
			return fmt.Errorf("failed to update config %s: %w", update.ConfigKey, err)
		}
	}

	log.Printf("[SUCCESS] Bulk updated %d configs", len(updates))
	return nil
}

// Cache and performance (placeholder implementations)
func (s *ConfigServiceImpl) RefreshConfigCache(ctx context.Context) error {
	log.Printf("[DEBUG] Config cache refresh requested (not implemented)")
	// TODO: Implement caching mechanism
	return nil
}

func (s *ConfigServiceImpl) GetCachedConfig(ctx context.Context, configKey string) (*models.SystemConfig, error) {
	log.Printf("[DEBUG] Getting cached config for key: %s (falling back to database)", configKey)
	// TODO: Implement caching mechanism
	return s.GetConfig(ctx, configKey)
}

// Agent booking limit specific operations
func (s *ConfigServiceImpl) GetAgentBookingLimits(ctx context.Context) (*models.AgentBookingLimitConfig, error) {
	log.Printf("[DEBUG] Getting agent booking limits configuration")

	config, err := s.configRepo.GetByKey(ctx, "booking.agent_period_limits")
	if err != nil {
		log.Printf("[ERROR] Failed to get agent booking limits config: %v", err)
		return nil, fmt.Errorf("failed to get agent booking limits: %w", err)
	}

	if config == nil {
		log.Printf("[WARNING] Agent booking limits configuration not found, returning default")
		// Return default configuration
		return &models.AgentBookingLimitConfig{
			GlobalSettings: models.AgentBookingGlobalSettings{
				Enabled:               false,
				DefaultLimitPerPeriod: 10,
				CountAllStatuses:      []string{"pending", "processing", "confirm"},
			},
			PeriodSpecificLimits: []models.AgentBookingPeriodLimit{},
			AgentGlobalBypasses:  []string{},
			Metadata: models.AgentBookingLimitMetadata{
				Version:     "1.0",
				LastUpdated: time.Now(),
				UpdatedBy:   "system",
			},
		}, nil
	}

	limits, err := config.GetAgentBookingLimits()
	if err != nil {
		log.Printf("[ERROR] Failed to parse agent booking limits: %v", err)
		return nil, fmt.Errorf("failed to parse agent booking limits: %w", err)
	}

	log.Printf("[SUCCESS] Retrieved agent booking limits configuration")
	return limits, nil
}

func (s *ConfigServiceImpl) UpdateAgentBookingLimits(ctx context.Context, limits *models.AgentBookingLimitConfig, updatedBy string) error {
	log.Printf("[DEBUG] Updating agent booking limits - updated by: %s", updatedBy)

	// Get existing config or create new one
	config, err := s.configRepo.GetByKey(ctx, "booking.agent_period_limits")
	if err != nil {
		log.Printf("[ERROR] Failed to get existing agent booking limits config: %v", err)
		return fmt.Errorf("failed to get existing config: %w", err)
	}

	if config == nil {
		// Create new config
		config = &models.SystemConfig{
			ConfigKey:   "booking.agent_period_limits",
			ConfigType:  models.ConfigTypeBooking,
			Description: new(string),
			IsActive:    true,
		}
		*config.Description = "Agent booking period limits configuration"
	}

	// Set the limits
	err = config.SetAgentBookingLimits(limits, updatedBy)
	if err != nil {
		log.Printf("[ERROR] Failed to set limits in config: %v", err)
		return fmt.Errorf("failed to set limits: %w", err)
	}

	// Save or update the config
	if config.ID == uuid.Nil {
		err = s.configRepo.Create(ctx, config)
	} else {
		err = s.configRepo.UpdateByKey(ctx, "booking.agent_period_limits", config.ConfigValue)
	}

	if err != nil {
		log.Printf("[ERROR] Failed to save agent booking limits: %v", err)
		return fmt.Errorf("failed to save limits: %w", err)
	}

	log.Printf("[SUCCESS] Agent booking limits updated successfully")
	return nil
}

func (s *ConfigServiceImpl) ValidateAgentBookingLimit(ctx context.Context, agentID uuid.UUID, periodID uuid.UUID, requestedQuantity int) error {
	log.Printf("[DEBUG] Validating agent booking limit - agent: %s, period: %s, quantity: %d", agentID, periodID, requestedQuantity)

	// Get agent booking limits configuration
	limits, err := s.GetAgentBookingLimits(ctx)
	if err != nil {
		log.Printf("[ERROR] Failed to get agent booking limits: %v", err)
		// If we can't get the config, allow the booking (fail open)
		log.Printf("[WARNING] Allowing booking due to config error")
		return nil
	}

	// Check if limits are enabled
	if !limits.GlobalSettings.Enabled {
		log.Printf("[DEBUG] Agent booking limits are disabled, allowing booking")
		return nil
	}

	// Check if agent is bypassed
	if limits.IsAgentBypassed(agentID.String(), periodID.String()) {
		log.Printf("[DEBUG] Agent %s is bypassed for period %s, allowing booking", agentID, periodID)
		return nil
	}

	// Get the limit for this period
	limit := limits.GetPeriodLimit(periodID.String())
	log.Printf("[DEBUG] Booking limit for period %s: %d", periodID, limit)

	// Count current bookings by this agent for this period
	statuses := limits.GlobalSettings.CountAllStatuses
	currentCount, err := s.countAgentBookingsForPeriod(ctx, agentID, periodID, statuses)
	if err != nil {
		log.Printf("[ERROR] Failed to count agent bookings: %v", err)
		// If we can't count bookings, allow the booking (fail open)
		log.Printf("[WARNING] Allowing booking due to counting error")
		return nil
	}

	log.Printf("[DEBUG] Agent %s current bookings for period %s: %d", agentID, periodID, currentCount)

	// Check if adding the requested quantity would exceed the limit
	if currentCount+requestedQuantity > limit {
		log.Printf("[WARNING] Agent booking limit exceeded - current: %d, requested: %d, limit: %d", currentCount, requestedQuantity, limit)
		return fmt.Errorf("agent booking limit exceeded: current %d + requested %d > limit %d", currentCount, requestedQuantity, limit)
	}

	log.Printf("[SUCCESS] Agent booking limit validation passed")
	return nil
}

// Helper function to count agent bookings for a period
func (s *ConfigServiceImpl) countAgentBookingsForPeriod(ctx context.Context, agentID uuid.UUID, periodID uuid.UUID, statuses []string) (int, error) {
	// Use the booking repository to count bookings
	bookings, err := s.bookingRepo.FindByPeriodIDAndStatuses(periodID, statuses)
	if err != nil {
		return 0, fmt.Errorf("failed to get bookings: %w", err)
	}

	// Count bookings for this specific agent
	agentBookingCount := 0
	for _, booking := range bookings {
		if booking.UserID == agentID {
			agentBookingCount += booking.SumQty
		}
	}

	return agentBookingCount, nil
}