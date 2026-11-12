package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type PeriodRepositoryImpl struct {
	db *gorm.DB
}

func NewPeriodRepository(db *gorm.DB) repositories.PeriodRepository {
	return &PeriodRepositoryImpl{db: db}
}

func (r *PeriodRepositoryImpl) Create(ctx context.Context, entity *models.Period) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create period: %w", err)
	}
	return nil
}

func (r *PeriodRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Period, error) {
	var period models.Period
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&period).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get period by ID: %w", err)
	}

	return &period, nil
}

func (r *PeriodRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.Period) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update period: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("period with ID %s not found", id)
	}

	return nil
}

func (r *PeriodRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Period{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete period: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("period with ID %s not found", id)
	}

	return nil
}

func (r *PeriodRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Period, error) {
	var periods []*models.Period

	err := r.db.WithContext(ctx).
		Order("start_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&periods).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list periods: %w", err)
	}

	return periods, nil
}

func (r *PeriodRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Period{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count periods: %w", err)
	}

	return count, nil
}

func (r *PeriodRepositoryImpl) GetByTour(ctx context.Context, tourID uuid.UUID, offset, limit int) ([]*models.Period, error) {
	var periods []*models.Period

	err := r.db.WithContext(ctx).
		Where("tour_id = ?", tourID).
		Order("start_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&periods).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get periods by tour: %w", err)
	}

	return periods, nil
}

func (r *PeriodRepositoryImpl) GetByTourID(ctx context.Context, tourID uuid.UUID) ([]*models.Period, error) {
	var periods []*models.Period

	err := r.db.WithContext(ctx).
		Where("tour_id = ?", tourID).
		Order("start_date ASC").
		Find(&periods).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get periods by tour ID: %w", err)
	}

	return periods, nil
}

func (r *PeriodRepositoryImpl) GetAvailable(ctx context.Context, offset, limit int) ([]*models.Period, error) {
	var periods []*models.Period

	err := r.db.WithContext(ctx).
		Where("status = ? AND start_date > ? AND max_participants > current_participants",
			  models.PeriodStatusPublish, time.Now()).
		Order("start_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&periods).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get available periods: %w", err)
	}

	return periods, nil
}

func (r *PeriodRepositoryImpl) GetAvailableByTour(ctx context.Context, tourID uuid.UUID, offset, limit int) ([]*models.Period, error) {
	var periods []*models.Period

	err := r.db.WithContext(ctx).
		Where("tour_id = ? AND status = ? AND start_date > ? AND max_participants > current_participants",
			  tourID, models.PeriodStatusPublish, time.Now()).
		Order("start_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&periods).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get available periods by tour: %w", err)
	}

	return periods, nil
}

func (r *PeriodRepositoryImpl) GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Period, error) {
	var periods []*models.Period

	err := r.db.WithContext(ctx).
		Where("start_date >= ? AND end_date <= ?", startDate, endDate).
		Order("start_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&periods).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get periods by date range: %w", err)
	}

	return periods, nil
}

func (r *PeriodRepositoryImpl) GetUpcoming(ctx context.Context, offset, limit int) ([]*models.Period, error) {
	var periods []*models.Period

	err := r.db.WithContext(ctx).
		Where("start_date > ?", time.Now()).
		Order("start_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&periods).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming periods: %w", err)
	}

	return periods, nil
}

func (r *PeriodRepositoryImpl) UpdateBookingCount(ctx context.Context, id uuid.UUID, bookingCount int) error {
	result := r.db.WithContext(ctx).
		Model(&models.Period{}).
		Where("id = ?", id).
		Update("current_participants", bookingCount)

	if result.Error != nil {
		return fmt.Errorf("failed to update period booking count: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("period with ID %s not found", id)
	}

	return nil
}

func (r *PeriodRepositoryImpl) IncrementBookingCount(ctx context.Context, id uuid.UUID, increment int) error {
	result := r.db.WithContext(ctx).
		Model(&models.Period{}).
		Where("id = ?", id).
		UpdateColumn("current_participants", gorm.Expr("current_participants + ?", increment))

	if result.Error != nil {
		return fmt.Errorf("failed to increment period booking count: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("period with ID %s not found", id)
	}

	return nil
}

func (r *PeriodRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PeriodStatus) error {
	result := r.db.WithContext(ctx).
		Model(&models.Period{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update period status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("period with ID %s not found", id)
	}

	return nil
}

func (r *PeriodRepositoryImpl) GetByTourIDWithPagination(ctx context.Context, tourID uuid.UUID, offset, limit int) ([]*models.Period, error) {
	var periods []*models.Period

	err := r.db.WithContext(ctx).
		Preload("PeriodRateItems").
		Preload("PeriodRoomTypes").
		Preload("PeriodAdditionalServices").
		Where("tour_id = ?", tourID).
		Order("start_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&periods).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get periods by tour ID with pagination: %w", err)
	}

	return periods, nil
}

func (r *PeriodRepositoryImpl) CountByTourID(ctx context.Context, tourID uuid.UUID) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Period{}).
		Where("tour_id = ?", tourID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count periods by tour ID: %w", err)
	}

	return count, nil
}

// ListWithFilters returns periods with multiple filters applied
func (r *PeriodRepositoryImpl) ListWithFilters(ctx context.Context, req dto.PeriodListRequest) ([]*models.Period, int64, error) {
	var periods []*models.Period
	var total int64

	// Build base query
	query := r.db.WithContext(ctx).Model(&models.Period{})

	// Apply filters
	if req.TourID != nil {
		query = query.Where("tour_id = ?", *req.TourID)
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if req.DateFrom != nil {
		query = query.Where("start_date >= ?", *req.DateFrom)
	}

	if req.DateTo != nil {
		query = query.Where("start_date <= ?", *req.DateTo)
	}

	if req.Available != nil && *req.Available {
		// Available means: published and has available seats
		query = query.Where("status = ?", "publish").
			Where("(seat - booking - COALESCE(held_seats, 0)) > 0")
	}

	if req.Upcoming != nil && *req.Upcoming {
		// Upcoming means start_date is in the future
		now := time.Now()
		query = query.Where("start_date > ?", now)
	}

	// Get total count before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count periods: %w", err)
	}

	// Apply pagination and sorting
	err := query.
		Order("start_date ASC, created_at DESC").
		Offset(req.Offset).
		Limit(req.Limit).
		Find(&periods).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to list periods with filters: %w", err)
	}

	return periods, total, nil
}

// WithTx returns a new repository instance with the given transaction
func (r *PeriodRepositoryImpl) WithTx(tx *gorm.DB) repositories.PeriodRepository {
	return &PeriodRepositoryImpl{db: tx}
}