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

type BookingAdditionalServiceRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingAdditionalServiceRepository(db *gorm.DB) repositories.BookingAdditionalServiceRepository {
	return &BookingAdditionalServiceRepositoryImpl{db: db}
}

func (r *BookingAdditionalServiceRepositoryImpl) Create(ctx context.Context, entity *models.BookingAdditionalService) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create booking additional service: %w", err)
	}
	return nil
}

func (r *BookingAdditionalServiceRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.BookingAdditionalService, error) {
	var service models.BookingAdditionalService
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&service).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get booking additional service: %w", err)
	}
	return &service, nil
}

func (r *BookingAdditionalServiceRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.BookingAdditionalService) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity).Error; err != nil {
		return fmt.Errorf("failed to update booking additional service: %w", err)
	}
	return nil
}

func (r *BookingAdditionalServiceRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.BookingAdditionalService{}).Error; err != nil {
		return fmt.Errorf("failed to delete booking additional service: %w", err)
	}
	return nil
}

func (r *BookingAdditionalServiceRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.BookingAdditionalService, error) {
	var services []*models.BookingAdditionalService
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&services).Error; err != nil {
		return nil, fmt.Errorf("failed to list booking additional services: %w", err)
	}
	return services, nil
}

func (r *BookingAdditionalServiceRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.BookingAdditionalService{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count booking additional services: %w", err)
	}
	return count, nil
}

func (r *BookingAdditionalServiceRepositoryImpl) GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.BookingAdditionalService, error) {
	var services []*models.BookingAdditionalService
	if err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).Find(&services).Error; err != nil {
		return nil, fmt.Errorf("failed to get booking additional services by booking ID: %w", err)
	}
	return services, nil
}

func (r *BookingAdditionalServiceRepositoryImpl) DeleteByBookingID(ctx context.Context, bookingID uuid.UUID) error {
	if err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).Delete(&models.BookingAdditionalService{}).Error; err != nil {
		return fmt.Errorf("failed to delete booking additional services by booking ID: %w", err)
	}
	return nil
}

func (r *BookingAdditionalServiceRepositoryImpl) CreateBatch(ctx context.Context, entities []*models.BookingAdditionalService) error {
	if len(entities) == 0 {
		return nil
	}
	if err := r.db.WithContext(ctx).CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("failed to batch create booking additional services: %w", err)
	}
	return nil
}

func (r *BookingAdditionalServiceRepositoryImpl) GetTotalByBookingID(ctx context.Context, bookingID uuid.UUID) (float64, error) {
	var total float64
	if err := r.db.WithContext(ctx).Model(&models.BookingAdditionalService{}).
		Where("booking_id = ?", bookingID).
		Select("COALESCE(SUM(price * quantity), 0)").
		Scan(&total).Error; err != nil {
		return 0, fmt.Errorf("failed to calculate total for booking additional services: %w", err)
	}
	return total, nil
}

func (r *BookingAdditionalServiceRepositoryImpl) GetTotalByBookingIDAndPassengerCount(ctx context.Context, bookingID uuid.UUID, passengerCount int) (float64, error) {
	var services []*models.BookingAdditionalService
	if err := r.db.WithContext(ctx).Where("booking_id = ?", bookingID).Find(&services).Error; err != nil {
		return 0, fmt.Errorf("failed to get booking additional services for total calculation: %w", err)
	}

	var total float64
	for _, service := range services {
		total += service.CalculateTotal(passengerCount)
	}

	return total, nil
}