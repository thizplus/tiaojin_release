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

type BookingRoomTypeRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRoomTypeRepository(db *gorm.DB) repositories.BookingRoomTypeRepository {
	return &BookingRoomTypeRepositoryImpl{db: db}
}

func (r *BookingRoomTypeRepositoryImpl) Create(ctx context.Context, entity *models.BookingRoomType) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create booking room type: %w", err)
	}
	return nil
}

func (r *BookingRoomTypeRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.BookingRoomType, error) {
	var roomType models.BookingRoomType
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&roomType).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get booking room type by ID: %w", err)
	}

	return &roomType, nil
}

func (r *BookingRoomTypeRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.BookingRoomType) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update booking room type: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("booking room type with ID %s not found", id)
	}

	return nil
}

func (r *BookingRoomTypeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.BookingRoomType{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete booking room type: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("booking room type with ID %s not found", id)
	}

	return nil
}

func (r *BookingRoomTypeRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.BookingRoomType, error) {
	var roomTypes []*models.BookingRoomType

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&roomTypes).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list booking room types: %w", err)
	}

	return roomTypes, nil
}

func (r *BookingRoomTypeRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.BookingRoomType{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count booking room types: %w", err)
	}

	return count, nil
}

func (r *BookingRoomTypeRepositoryImpl) GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.BookingRoomType, error) {
	var roomTypes []*models.BookingRoomType

	err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Order("created_at ASC").
		Find(&roomTypes).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get booking room types by booking ID: %w", err)
	}

	return roomTypes, nil
}

func (r *BookingRoomTypeRepositoryImpl) DeleteByBookingID(ctx context.Context, bookingID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.BookingRoomType{}, "booking_id = ?", bookingID)

	if result.Error != nil {
		return fmt.Errorf("failed to delete booking room types by booking ID: %w", result.Error)
	}

	return nil
}

func (r *BookingRoomTypeRepositoryImpl) CreateBatch(ctx context.Context, entities []*models.BookingRoomType) error {
	if len(entities) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("failed to create booking room types in batch: %w", err)
	}

	return nil
}

func (r *BookingRoomTypeRepositoryImpl) GetTotalByBookingID(ctx context.Context, bookingID uuid.UUID) (float64, error) {
	var total float64

	err := r.db.WithContext(ctx).
		Model(&models.BookingRoomType{}).
		Where("booking_id = ?", bookingID).
		Select("COALESCE(SUM(quantity * price), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, fmt.Errorf("failed to calculate total by booking ID: %w", err)
	}

	return total, nil
}

// WithTx returns a new repository instance with the given transaction
func (r *BookingRoomTypeRepositoryImpl) WithTx(tx *gorm.DB) repositories.BookingRoomTypeRepository {
	return &BookingRoomTypeRepositoryImpl{db: tx}
}