package postgres

import (
	"context"
	"time"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodRateItemRepositoryImpl struct {
	db *gorm.DB
}

func NewPeriodRateItemRepository(db *gorm.DB) repositories.PeriodRateItemRepository {
	return &PeriodRateItemRepositoryImpl{db: db}
}

func (r *PeriodRateItemRepositoryImpl) Create(ctx context.Context, item *models.PeriodRateItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PeriodRateItemRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.PeriodRateItem, error) {
	var item models.PeriodRateItem
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PeriodRateItemRepositoryImpl) GetByPeriodID(ctx context.Context, periodID uuid.UUID) ([]*models.PeriodRateItem, error) {
	var items []*models.PeriodRateItem
	err := r.db.WithContext(ctx).
		Where("period_id = ?", periodID).
		Order("sort_order ASC, created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *PeriodRateItemRepositoryImpl) Update(ctx context.Context, id uuid.UUID, item *models.PeriodRateItem) error {
	item.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(item).Error
}

func (r *PeriodRateItemRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PeriodRateItem{}).Error
}

func (r *PeriodRateItemRepositoryImpl) DeleteByPeriodID(ctx context.Context, periodID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("period_id = ?", periodID).Delete(&models.PeriodRateItem{}).Error
}

func (r *PeriodRateItemRepositoryImpl) CreateBatch(ctx context.Context, items []*models.PeriodRateItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(items, 100).Error
}

func (r *PeriodRateItemRepositoryImpl) UpdateBatch(ctx context.Context, periodID uuid.UUID, items []*models.PeriodRateItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Hard delete existing items to avoid primary key conflicts
		if err := tx.Unscoped().Where("period_id = ?", periodID).Delete(&models.PeriodRateItem{}).Error; err != nil {
			return err
		}

		// Create new items
		if len(items) > 0 {
			for _, item := range items {
				item.PeriodID = periodID
				// Only set timestamps for items without ID (new items)
				if item.ID == uuid.Nil {
					item.ID = uuid.New()
				}
				item.CreatedAt = time.Now()
				item.UpdatedAt = time.Now()
			}
			if err := tx.CreateInBatches(items, 100).Error; err != nil {
				return err
			}
		}

		return nil
	})
}