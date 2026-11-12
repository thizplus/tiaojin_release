package postgres

import (
	"context"
	"errors"
	"fmt"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gorm.io/gorm"
)

type SystemConfigRepositoryImpl struct {
	db *gorm.DB
}

func NewSystemConfigRepository(db *gorm.DB) repositories.SystemConfigRepository {
	return &SystemConfigRepositoryImpl{db: db}
}

func (r *SystemConfigRepositoryImpl) Create(ctx context.Context, entity *models.SystemConfig) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create system config: %w", err)
	}
	return nil
}

func (r *SystemConfigRepositoryImpl) GetByID(ctx context.Context, id string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get system config by ID: %w", err)
	}

	return &config, nil
}

func (r *SystemConfigRepositoryImpl) GetByKey(ctx context.Context, configKey string) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := r.db.WithContext(ctx).Where("config_key = ?", configKey).First(&config).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get system config by key: %w", err)
	}

	return &config, nil
}

func (r *SystemConfigRepositoryImpl) Update(ctx context.Context, id string, entity *models.SystemConfig) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update system config: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("system config with ID %s not found", id)
	}

	return nil
}

func (r *SystemConfigRepositoryImpl) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.SystemConfig{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete system config: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("system config with ID %s not found", id)
	}

	return nil
}

func (r *SystemConfigRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig

	err := r.db.WithContext(ctx).
		Order("config_type, config_key").
		Offset(offset).
		Limit(limit).
		Find(&configs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list system configs: %w", err)
	}

	return configs, nil
}

func (r *SystemConfigRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.SystemConfig{}).Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count system configs: %w", err)
	}

	return count, nil
}

func (r *SystemConfigRepositoryImpl) GetByType(ctx context.Context, configType models.ConfigType, offset, limit int) ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig

	err := r.db.WithContext(ctx).
		Where("config_type = ?", configType).
		Order("config_key").
		Offset(offset).
		Limit(limit).
		Find(&configs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get system configs by type: %w", err)
	}

	return configs, nil
}

func (r *SystemConfigRepositoryImpl) GetActiveConfigs(ctx context.Context, configType models.ConfigType) ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig

	err := r.db.WithContext(ctx).
		Where("config_type = ? AND is_active = ?", configType, true).
		Order("config_key").
		Find(&configs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get active system configs: %w", err)
	}

	return configs, nil
}

func (r *SystemConfigRepositoryImpl) UpdateByKey(ctx context.Context, configKey string, configValue models.ConfigValueJSON) error {
	result := r.db.WithContext(ctx).
		Model(&models.SystemConfig{}).
		Where("config_key = ?", configKey).
		Updates(map[string]interface{}{
			"config_value": configValue,
			"updated_at":   "NOW()",
		})

	if result.Error != nil {
		return fmt.Errorf("failed to update system config by key: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("system config with key %s not found", configKey)
	}

	return nil
}

func (r *SystemConfigRepositoryImpl) GetMultipleByKeys(ctx context.Context, keys []string) ([]*models.SystemConfig, error) {
	if len(keys) == 0 {
		return []*models.SystemConfig{}, nil
	}

	var configs []*models.SystemConfig

	err := r.db.WithContext(ctx).
		Where("config_key IN ?", keys).
		Order("config_key").
		Find(&configs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get system configs by keys: %w", err)
	}

	return configs, nil
}

func (r *SystemConfigRepositoryImpl) BulkUpdate(ctx context.Context, configs []*models.SystemConfig) error {
	if len(configs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, config := range configs {
			if err := tx.Save(config).Error; err != nil {
				return fmt.Errorf("failed to update config %s: %w", config.ConfigKey, err)
			}
		}
		return nil
	})
}

func (r *SystemConfigRepositoryImpl) BulkCreate(ctx context.Context, configs []*models.SystemConfig) error {
	if len(configs) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(configs, 100).Error; err != nil {
		return fmt.Errorf("failed to bulk create system configs: %w", err)
	}

	return nil
}

func (r *SystemConfigRepositoryImpl) SearchByKeyPattern(ctx context.Context, pattern string, offset, limit int) ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig

	err := r.db.WithContext(ctx).
		Where("config_key ILIKE ?", pattern).
		Order("config_key").
		Offset(offset).
		Limit(limit).
		Find(&configs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to search system configs by pattern: %w", err)
	}

	return configs, nil
}

func (r *SystemConfigRepositoryImpl) GetConfigsUpdatedAfter(ctx context.Context, after string, offset, limit int) ([]*models.SystemConfig, error) {
	var configs []*models.SystemConfig

	err := r.db.WithContext(ctx).
		Where("updated_at > ?", after).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&configs).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get configs updated after: %w", err)
	}

	return configs, nil
}

// WithTx returns a new repository instance with the given transaction
func (r *SystemConfigRepositoryImpl) WithTx(tx *gorm.DB) repositories.SystemConfigRepository {
	return &SystemConfigRepositoryImpl{db: tx}
}