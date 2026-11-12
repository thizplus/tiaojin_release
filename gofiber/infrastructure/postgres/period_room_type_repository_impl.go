package postgres

import (
	"context"
	"time"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodRoomTypeRepositoryImpl struct {
	db *gorm.DB
}

func NewPeriodRoomTypeRepository(db *gorm.DB) repositories.PeriodRoomTypeRepository {
	return &PeriodRoomTypeRepositoryImpl{db: db}
}

func (r *PeriodRoomTypeRepositoryImpl) Create(ctx context.Context, item *models.PeriodRoomType) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PeriodRoomTypeRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.PeriodRoomType, error) {
	var item models.PeriodRoomType
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PeriodRoomTypeRepositoryImpl) GetByPeriodID(ctx context.Context, periodID uuid.UUID) ([]*models.PeriodRoomType, error) {
	var items []*models.PeriodRoomType
	err := r.db.WithContext(ctx).
		Where("period_id = ?", periodID).
		Order("sort_order ASC, created_at ASC").
		Find(&items).Error
	return items, err
}

func (r *PeriodRoomTypeRepositoryImpl) Update(ctx context.Context, id uuid.UUID, item *models.PeriodRoomType) error {
	item.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Where("id = ?", id).Updates(item).Error
}

func (r *PeriodRoomTypeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.PeriodRoomType{}).Error
}

func (r *PeriodRoomTypeRepositoryImpl) DeleteByPeriodID(ctx context.Context, periodID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("period_id = ?", periodID).Delete(&models.PeriodRoomType{}).Error
}

func (r *PeriodRoomTypeRepositoryImpl) CreateBatch(ctx context.Context, items []*models.PeriodRoomType) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).CreateInBatches(items, 100).Error
}

func (r *PeriodRoomTypeRepositoryImpl) UpdateBatch(ctx context.Context, periodID uuid.UUID, items []*models.PeriodRoomType) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Hard delete existing items to avoid primary key conflicts
		if err := tx.Unscoped().Where("period_id = ?", periodID).Delete(&models.PeriodRoomType{}).Error; err != nil {
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