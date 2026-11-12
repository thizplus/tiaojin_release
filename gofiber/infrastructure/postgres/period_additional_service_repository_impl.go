package postgres

import (
	"context"
	"time"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodAdditionalServiceRepositoryImpl struct {
	db *gorm.DB
}

func NewPeriodAdditionalServiceRepository(db *gorm.DB) repositories.PeriodAdditionalServiceRepository {
	return &PeriodAdditionalServiceRepositoryImpl{db: db}
}

func (r *PeriodAdditionalServiceRepositoryImpl) Create(ctx context.Context, item *models.PeriodAdditionalService) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PeriodAdditionalServiceRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.PeriodAdditionalService, error) {
	var item models.PeriodAdditionalService
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PeriodAdditionalServiceRepositoryImpl) GetByPeriodID(ctx context.Context, periodID uuid.UUID) ([]*models.PeriodAdditionalService, error) {
	var items []*models.PeriodAdditionalService
	err := r.db.WithContext(ctx).
		Where("period_id = ?", periodID).
		Order("sort_order ASC, created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *PeriodAdditionalServiceRepositoryImpl) Update(ctx context.Context, id uuid.UUID, item *models.PeriodAdditionalService) error {
	item.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(item).Error
}

func (r *PeriodAdditionalServiceRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PeriodAdditionalService{}).Error
}

func (r *PeriodAdditionalServiceRepositoryImpl) DeleteByPeriodID(ctx context.Context, periodID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("period_id = ?", periodID).Delete(&models.PeriodAdditionalService{}).Error
}

func (r *PeriodAdditionalServiceRepositoryImpl) CreateBatch(ctx context.Context, items []*models.PeriodAdditionalService) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(items, 100).Error
}

func (r *PeriodAdditionalServiceRepositoryImpl) UpdateBatch(ctx context.Context, periodID uuid.UUID, items []*models.PeriodAdditionalService) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Hard delete existing items to avoid primary key conflicts
		if err := tx.Unscoped().Where("period_id = ?", periodID).Delete(&models.PeriodAdditionalService{}).Error; err != nil {
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