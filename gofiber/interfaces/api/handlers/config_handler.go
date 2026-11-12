package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type ConfigHandler struct {
	configService  services.ConfigService
	bookingService services.BookingService
}

func NewConfigHandler(configService services.ConfigService, bookingService services.BookingService) *ConfigHandler {
	return &ConfigHandler{
		configService:  configService,
		bookingService: bookingService,
	}
}

// GetExpirationTiers gets the current expiration tiers configuration
// GET /api/admin/config/expiration-tiers
func (h *ConfigHandler) GetExpirationTiers(c *fiber.Ctx) error {
	tiers, err := h.configService.GetExpirationTiers(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get expiration tiers", err)
	}

	// Convert to response format
	response := dto.ExpirationTiersResponse{
		Tiers: make([]dto.ExpirationTierResponse, len(tiers)),
		Metadata: models.ExpirationTierMetadata{
			LastUpdated: time.Now(),
			Version:     "1.0",
			UpdatedBy:   "system",
		},
	}

	for i, tier := range tiers {
		response.Tiers[i] = dto.ExpirationTierResponse{
			Name:           tier.Name,
			Description:    tier.Description,
			ThresholdDays:  tier.ThresholdDays,
			ExpireMinutes:  tier.ExpireMinutes,
			ExpireHours:    tier.ExpireHours,
			ExpireDays:     tier.ExpireDays,
			SortOrder:      tier.SortOrder,
			ExampleMessage: "Example: Book " + tier.Name + " tier",
		}
	}

	return utils.SuccessResponse(c, "Expiration tiers retrieved successfully", response)
}

// UpdateExpirationTiers updates the expiration tiers configuration
// PUT /api/admin/config/expiration-tiers
func (h *ConfigHandler) UpdateExpirationTiers(c *fiber.Ctx) error {
	var req dto.UpdateExpirationTiersRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Get user ID from context for audit trail
	userID := "admin"
	if uid := c.Locals("user_id"); uid != nil {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	// Update tiers
	err := h.configService.UpdateExpirationTiers(c.Context(), req.Tiers, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update expiration tiers", err)
	}

	return utils.SuccessResponse(c, "Expiration tiers updated successfully", nil)
}

// ValidateExpirationConfig validates expiration configuration
// POST /api/admin/config/expiration-tiers/validate
func (h *ConfigHandler) ValidateExpirationConfig(c *fiber.Ctx) error {
	var req dto.UpdateExpirationTiersRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate the tiers
	err := h.configService.ValidateExpirationConfig(c.Context(), req.Tiers)

	response := dto.ConfigValidationResponse{
		IsValid:     err == nil,
		TotalTiers:  len(req.Tiers),
	}

	if err != nil {
		response.Errors = []string{err.Error()}
	} else {
		// Find default tier
		for _, tier := range req.Tiers {
			if tier.ThresholdDays == nil {
				response.DefaultTier = &tier.Name
				break
			}
		}
	}

	return utils.SuccessResponse(c, "Configuration validation completed", response)
}

// CalculateTier calculates which tier applies for a given scenario
// POST /api/admin/config/expiration-tiers/calculate
func (h *ConfigHandler) CalculateTier(c *fiber.Ctx) error {
	var req dto.TierCalculationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Use current time if created_at not provided
	createdAt := req.CreatedAt
	if createdAt == "" {
		createdAt = time.Now().Format("2006-01-02T15:04:05Z")
	}

	// Calculate applicable tier
	tier, err := h.configService.CalculateTieredExpiration(c.Context(), req.TourStartDate, createdAt)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to calculate tier", err)
	}

	// Parse dates for calculation details
	tourDate, _ := time.Parse("2006-01-02", req.TourStartDate)
	bookingCreated, _ := time.Parse("2006-01-02T15:04:05Z", createdAt)
	currentTime := time.Now()

	daysUntilTour := models.DaysUntilTour(tourDate)

	// Calculate expiration time
	var expirationTime *string
	var timeUntilExpiry *string

	if tier.ExpireHours != nil {
		expTime := bookingCreated.Add(time.Duration(*tier.ExpireHours) * time.Hour)
		expTimeStr := expTime.Format(time.RFC3339)
		expirationTime = &expTimeStr

		timeLeft := time.Until(expTime)
		if timeLeft > 0 {
			timeLeftStr := timeLeft.String()
			timeUntilExpiry = &timeLeftStr
		}
	} else if tier.ExpireDays != nil {
		expTime := bookingCreated.AddDate(0, 0, *tier.ExpireDays)
		expTimeStr := expTime.Format(time.RFC3339)
		expirationTime = &expTimeStr

		timeLeft := time.Until(expTime)
		if timeLeft > 0 {
			timeLeftStr := timeLeft.String()
			timeUntilExpiry = &timeLeftStr
		}
	}

	// Get all tiers for testing details
	allTiers, _ := h.configService.GetExpirationTiers(c.Context())
	testedTiers := make([]dto.TierTestResult, len(allTiers))

	for i, t := range allTiers {
		matches := t.MatchesTier(daysUntilTour)
		reason := "Does not match criteria"
		if matches {
			if t.ThresholdDays == nil {
				reason = "Default tier (no threshold)"
			} else {
				reason = "Days until tour matches threshold"
			}
		}

		testedTiers[i] = dto.TierTestResult{
			TierName:      t.Name,
			ThresholdDays: t.ThresholdDays,
			Matches:       matches,
			Reason:        reason,
		}
	}

	// Determine expiration method
	expirationMethod := "days"
	if tier.ExpireHours != nil {
		expirationMethod = "hours"
	}

	response := dto.TierCalculationResponse{
		ApplicableTier: dto.ExpirationTierResponse{
			Name:           tier.Name,
			Description:    tier.Description,
			ThresholdDays:  tier.ThresholdDays,
			ExpireHours:    tier.ExpireHours,
			ExpireDays:     tier.ExpireDays,
			SortOrder:      tier.SortOrder,
			ExampleMessage: "Example: " + tier.Description,
		},
		DaysUntilTour:      daysUntilTour,
		ExpirationTime:     expirationTime,
		TimeUntilExpiry:    timeUntilExpiry,
		IsDefaultTier:      tier.ThresholdDays == nil,
		CalculationDetails: dto.TierCalculationDetails{
			TourStartDate:    req.TourStartDate,
			BookingCreatedAt: createdAt,
			CurrentTime:      currentTime.Format(time.RFC3339),
			DaysUntilTour:    daysUntilTour,
			TestedTiers:      testedTiers,
			SelectedTierName: tier.Name,
			ExpirationMethod: expirationMethod,
		},
	}

	return utils.SuccessResponse(c, "Tier calculation completed", response)
}

// GetExpirationStats gets expiration statistics
// GET /api/admin/config/expiration-stats
func (h *ConfigHandler) GetExpirationStats(c *fiber.Ctx) error {
	response, err := h.bookingService.GetBookingExpirationStats(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get expiration statistics", err)
	}

	return utils.SuccessResponse(c, "Expiration statistics retrieved", response)
}

// ProcessExpiredBookings processes expired bookings
// POST /api/admin/config/process-expired
func (h *ConfigHandler) ProcessExpiredBookings(c *fiber.Ctx) error {
	// Parse dry_run parameter
	dryRun := c.Query("dry_run", "true") == "true"

	response, err := h.bookingService.ProcessExpiredBookings(c.Context(), dryRun)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to process expired bookings", err)
	}

	return utils.SuccessResponse(c, "Expired bookings processing completed", response)
}

// ListConfigs lists system configurations
// GET /api/admin/config
func (h *ConfigHandler) ListConfigs(c *fiber.Ctx) error {
	// Parse query parameters
	configTypeStr := c.Query("config_type")
	isActiveStr := c.Query("is_active")
	search := c.Query("search")
	pagination := utils.ParsePaginationParams(c, 10)

	// Create request
	req := dto.ConfigListRequest{
		Offset: pagination.Offset,
		Limit:  pagination.Limit,
	}

	if configTypeStr != "" {
		req.ConfigType = &configTypeStr
	}

	if isActiveStr != "" {
		isActive := isActiveStr == "true"
		req.IsActive = &isActive
	}

	if search != "" {
		req.Search = &search
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Convert configType if provided
	var configType *models.ConfigType
	if req.ConfigType != nil {
		ct := models.ConfigType(*req.ConfigType)
		configType = &ct
	}

	// Get configurations
	configs, err := h.configService.ListConfigs(c.Context(), configType, req.Offset, req.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to list configurations", err)
	}

	// Convert to response format
	configResponses := make([]*dto.ConfigResponse, len(configs))
	for i, config := range configs {
		configResponses[i] = &dto.ConfigResponse{
			ID:          config.ID.String(),
			ConfigKey:   config.ConfigKey,
			ConfigValue: config.ConfigValue,
			ConfigType:  string(config.ConfigType),
			Description: config.Description,
			IsActive:    config.IsActive,
			CreatedAt:   config.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   config.UpdatedAt.Format(time.RFC3339),
		}
	}

	response := dto.ConfigListResponse{
		Configs:    configResponses,
		Pagination: utils.CalculatePagination(int64(len(configResponses)), req.Offset, req.Limit),
	}

	return utils.SuccessResponse(c, "Configurations retrieved successfully", response)
}

// CreateConfig creates a new configuration
// POST /api/admin/config
func (h *ConfigHandler) CreateConfig(c *fiber.Ctx) error {
	var req dto.CreateConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Set default values
	description := ""
	if req.Description != nil {
		description = *req.Description
	}

	// Create configuration
	err := h.configService.SetConfig(
		c.Context(),
		req.ConfigKey,
		req.ConfigValue,
		models.ConfigType(req.ConfigType),
		description,
	)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create configuration", err)
	}

	return utils.SuccessResponse(c, "Configuration created successfully", nil)
}

// UpdateConfig updates an existing configuration
// PUT /api/admin/config/:key
func (h *ConfigHandler) UpdateConfig(c *fiber.Ctx) error {
	configID := c.Params("id")
	if configID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Config ID is required", nil)
	}

	var req dto.UpdateConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Update configuration by ID
	err := h.configService.UpdateConfigByID(c.Context(), configID, req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update configuration", err)
	}

	return utils.SuccessResponse(c, "Configuration updated successfully", nil)
}

// GetConfig gets a single configuration by ID
// GET /api/v1/config/configs/:id
func (h *ConfigHandler) GetConfig(c *fiber.Ctx) error {
	configID := c.Params("id")
	if configID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Config ID is required", nil)
	}

	config, err := h.configService.GetConfigByID(c.Context(), configID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Configuration not found", err)
	}

	response := &dto.ConfigResponse{
		ID:          config.ID.String(),
		ConfigKey:   config.ConfigKey,
		ConfigValue: config.ConfigValue,
		ConfigType:  string(config.ConfigType),
		Description: config.Description,
		IsActive:    config.IsActive,
		CreatedAt:   config.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   config.UpdatedAt.Format(time.RFC3339),
	}

	return utils.SuccessResponse(c, "Configuration retrieved successfully", response)
}

// DeleteConfig deletes a configuration
// DELETE /api/v1/config/configs/:id
func (h *ConfigHandler) DeleteConfig(c *fiber.Ctx) error {
	configID := c.Params("id")
	if configID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Config ID is required", nil)
	}

	err := h.configService.DeleteConfigByID(c.Context(), configID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete configuration", err)
	}

	return utils.SuccessResponse(c, "Configuration deleted successfully", nil)
}

// GetAgentBookingLimits gets the current agent booking limits configuration
// GET /api/admin/config/agent-booking-limits
func (h *ConfigHandler) GetAgentBookingLimits(c *fiber.Ctx) error {
	limits, err := h.configService.GetAgentBookingLimits(c.Context())
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get agent booking limits", err)
	}

	// Convert to response format
	response := dto.AgentBookingLimitResponse{
		GlobalSettings: dto.AgentBookingGlobalSettingsResponse{
			Enabled:               limits.GlobalSettings.Enabled,
			DefaultLimitPerPeriod: limits.GlobalSettings.DefaultLimitPerPeriod,
			CountAllStatuses:      limits.GlobalSettings.CountAllStatuses,
		},
		PeriodSpecificLimits: make([]dto.AgentBookingPeriodLimitResponse, len(limits.PeriodSpecificLimits)),
		AgentGlobalBypasses:  limits.AgentGlobalBypasses,
		Metadata: dto.AgentBookingLimitMetadataResponse{
			Version:     limits.Metadata.Version,
			LastUpdated: limits.Metadata.LastUpdated.Format(time.RFC3339),
			UpdatedBy:   limits.Metadata.UpdatedBy,
		},
	}

	for i, periodLimit := range limits.PeriodSpecificLimits {
		response.PeriodSpecificLimits[i] = dto.AgentBookingPeriodLimitResponse{
			PeriodID:       periodLimit.PeriodID,
			Limit:          periodLimit.Limit,
			BypassedAgents: periodLimit.BypassedAgents,
		}
	}

	return utils.SuccessResponse(c, "Agent booking limits retrieved successfully", response)
}

// UpdateAgentBookingLimits updates the agent booking limits configuration
// PUT /api/admin/config/agent-booking-limits
func (h *ConfigHandler) UpdateAgentBookingLimits(c *fiber.Ctx) error {
	var req dto.UpdateAgentBookingLimitsRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body", err)
	}

	// Validate request
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Validation failed", err)
	}

	// Get user ID from context for audit trail
	userID := "admin"
	if uid := c.Locals("user_id"); uid != nil {
		if uidStr, ok := uid.(string); ok {
			userID = uidStr
		}
	}

	// Convert request to domain model
	limits := &models.AgentBookingLimitConfig{
		GlobalSettings: models.AgentBookingGlobalSettings{
			Enabled:               req.GlobalSettings.Enabled,
			DefaultLimitPerPeriod: req.GlobalSettings.DefaultLimitPerPeriod,
			CountAllStatuses:      req.GlobalSettings.CountAllStatuses,
		},
		PeriodSpecificLimits: make([]models.AgentBookingPeriodLimit, len(req.PeriodSpecificLimits)),
		AgentGlobalBypasses:  req.AgentGlobalBypasses,
	}

	for i, periodLimit := range req.PeriodSpecificLimits {
		limits.PeriodSpecificLimits[i] = models.AgentBookingPeriodLimit{
			PeriodID:       periodLimit.PeriodID,
			Limit:          periodLimit.Limit,
			BypassedAgents: periodLimit.BypassedAgents,
		}
	}

	// Update limits
	err := h.configService.UpdateAgentBookingLimits(c.Context(), limits, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update agent booking limits", err)
	}

	return utils.SuccessResponse(c, "Agent booking limits updated successfully", nil)
}