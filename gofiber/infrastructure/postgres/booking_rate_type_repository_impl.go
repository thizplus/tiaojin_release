package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type BookingRateTypeRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRateTypeRepository(db *gorm.DB) repositories.BookingRateTypeRepository {
	return &BookingRateTypeRepositoryImpl{db: db}
}

func (r *BookingRateTypeRepositoryImpl) Create(ctx context.Context, entity *models.BookingRateType) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create booking rate type: %w", err)
	}
	return nil
}

func (r *BookingRateTypeRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.BookingRateType, error) {
	var rateType models.BookingRateType
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rateType).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get booking rate type by ID: %w", err)
	}

	return &rateType, nil
}

func (r *BookingRateTypeRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.BookingRateType) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update booking rate type: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("booking rate type with ID %s not found", id)
	}

	return nil
}

func (r *BookingRateTypeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.BookingRateType{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete booking rate type: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("booking rate type with ID %s not found", id)
	}

	return nil
}

func (r *BookingRateTypeRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.BookingRateType, error) {
	var rateTypes []*models.BookingRateType

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&rateTypes).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list booking rate types: %w", err)
	}

	return rateTypes, nil
}

func (r *BookingRateTypeRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.BookingRateType{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count booking rate types: %w", err)
	}

	return count, nil
}

func (r *BookingRateTypeRepositoryImpl) GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.BookingRateType, error) {
	var rateTypes []*models.BookingRateType

	err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Order("created_at ASC").
		Find(&rateTypes).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get booking rate types by booking ID: %w", err)
	}

	return rateTypes, nil
}

func (r *BookingRateTypeRepositoryImpl) DeleteByBookingID(ctx context.Context, bookingID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.BookingRateType{}, "booking_id = ?", bookingID)

	if result.Error != nil {
		return fmt.Errorf("failed to delete booking rate types by booking ID: %w", result.Error)
	}

	return nil
}

func (r *BookingRateTypeRepositoryImpl) CreateBatch(ctx context.Context, entities []*models.BookingRateType) error {
	if len(entities) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("failed to create booking rate types in batch: %w", err)
	}

	return nil
}

func (r *BookingRateTypeRepositoryImpl) GetTotalByBookingID(ctx context.Context, bookingID uuid.UUID) (float64, error) {
	var total float64

	err := r.db.WithContext(ctx).
		Model(&models.BookingRateType{}).
		Where("booking_id = ?", bookingID).
		Select("COALESCE(SUM(quantity * price), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, fmt.Errorf("failed to calculate total by booking ID: %w", err)
	}

	return total, nil
}

func (r *BookingRateTypeRepositoryImpl) GetIncludedByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.BookingRateType, error) {
	var rateTypes []*models.BookingRateType

	err := r.db.WithContext(ctx).
		Where("booking_id = ? AND rate_type = ?", bookingID, "included").
		Order("created_at ASC").
		Find(&rateTypes).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get included rate types by booking ID: %w", err)
	}

	return rateTypes, nil
}

func (r *BookingRateTypeRepositoryImpl) GetExcludedByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.BookingRateType, error) {
	var rateTypes []*models.BookingRateType

	err := r.db.WithContext(ctx).
		Where("booking_id = ? AND rate_type = ?", bookingID, "excluded").
		Order("created_at ASC").
		Find(&rateTypes).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get excluded rate types by booking ID: %w", err)
	}

	return rateTypes, nil
}

// WithTx returns a new repository instance with the given transaction
func (r *BookingRateTypeRepositoryImpl) WithTx(tx *gorm.DB) repositories.BookingRateTypeRepository {
	return &BookingRateTypeRepositoryImpl{db: tx}
}