package repositories

import (
	"context"

	"gofiber-template/domain/models"
)

type SystemConfigRepository interface {
	// Basic CRUD operations
	Create(ctx context.Context, entity *models.SystemConfig) error
	GetByID(ctx context.Context, id string) (*models.SystemConfig, error)
	GetByKey(ctx context.Context, configKey string) (*models.SystemConfig, error)
	Update(ctx context.Context, id string, entity *models.SystemConfig) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*models.SystemConfig, error)
	Count(ctx context.Context) (int64, error)

	// Configuration-specific operations
	GetByType(ctx context.Context, configType models.ConfigType, offset, limit int) ([]*models.SystemConfig, error)
	GetActiveConfigs(ctx context.Context, configType models.ConfigType) ([]*models.SystemConfig, error)
	UpdateByKey(ctx context.Context, configKey string, configValue models.ConfigValueJSON) error
	GetMultipleByKeys(ctx context.Context, keys []string) ([]*models.SystemConfig, error)

	// Batch operations
	BulkUpdate(ctx context.Context, configs []*models.SystemConfig) error
	BulkCreate(ctx context.Context, configs []*models.SystemConfig) error

	// Search and filter operations
	SearchByKeyPattern(ctx context.Context, pattern string, offset, limit int) ([]*models.SystemConfig, error)
	GetConfigsUpdatedAfter(ctx context.Context, after string, offset, limit int) ([]*models.SystemConfig, error)
}