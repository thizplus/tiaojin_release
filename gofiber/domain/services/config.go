package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
)

type ConfigService interface {
	// Basic configuration operations
	GetConfig(ctx context.Context, configKey string) (*models.SystemConfig, error)
	GetConfigByID(ctx context.Context, configID string) (*models.SystemConfig, error)
	SetConfig(ctx context.Context, configKey string, configValue models.ConfigValueJSON, configType models.ConfigType, description string) error
	UpdateConfig(ctx context.Context, configKey string, configValue models.ConfigValueJSON) error
	UpdateConfigByID(ctx context.Context, configID string, req dto.UpdateConfigRequest) error
	DeleteConfig(ctx context.Context, configKey string) error
	DeleteConfigByID(ctx context.Context, configID string) error
	ListConfigs(ctx context.Context, configType *models.ConfigType, offset, limit int) ([]*models.SystemConfig, error)

	// Expiration tier specific operations
	GetExpirationTiers(ctx context.Context) ([]models.ExpirationTier, error)
	UpdateExpirationTiers(ctx context.Context, tiers []models.ExpirationTier, updatedBy string) error
	ValidateExpirationConfig(ctx context.Context, tiers []models.ExpirationTier) error
	GetDefaultExpirationTier(ctx context.Context) (*models.ExpirationTier, error)

	// Tier calculation operations
	CalculateTieredExpiration(ctx context.Context, tourStartDate, createdAt string) (*models.ExpirationTier, error)
	GetApplicableTier(ctx context.Context, daysUntilTour int) (*models.ExpirationTier, error)

	// Bulk operations
	GetMultipleConfigs(ctx context.Context, keys []string) (map[string]*models.SystemConfig, error)
	BulkUpdateConfigs(ctx context.Context, updates []dto.ConfigUpdateRequest) error

	// Cache and performance
	RefreshConfigCache(ctx context.Context) error
	GetCachedConfig(ctx context.Context, configKey string) (*models.SystemConfig, error)

	// Agent booking limit operations
	GetAgentBookingLimits(ctx context.Context) (*models.AgentBookingLimitConfig, error)
	UpdateAgentBookingLimits(ctx context.Context, limits *models.AgentBookingLimitConfig, updatedBy string) error
	ValidateAgentBookingLimit(ctx context.Context, agentID uuid.UUID, periodID uuid.UUID, requestedQuantity int) error
}