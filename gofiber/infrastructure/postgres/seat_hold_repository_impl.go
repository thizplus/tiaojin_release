package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SeatHoldRepositoryImpl struct {
	db *gorm.DB
}

// NewSeatHoldRepository creates a new seat hold repository implementation
func NewSeatHoldRepository(db *gorm.DB) repositories.SeatHoldRepository {
	return &SeatHoldRepositoryImpl{
		db: db,
	}
}

// Create creates a new seat hold
func (r *SeatHoldRepositoryImpl) Create(ctx context.Context, seatHold *models.SeatHold) error {
	if err := r.db.WithContext(ctx).Create(seatHold).Error; err != nil {
		log.Printf("Error creating seat hold: %v", err)
		return fmt.Errorf("failed to create seat hold: %w", err)
	}
	return nil
}

// GetByID retrieves a seat hold by ID
func (r *SeatHoldRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.SeatHold, error) {
	var seatHold models.SeatHold
	err := r.db.WithContext(ctx).
		Preload("Period").
		Preload("Admin").
		First(&seatHold, "id = ?", id).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Printf("Error getting seat hold by ID: %v", err)
		return nil, fmt.Errorf("failed to get seat hold: %w", err)
	}
	return &seatHold, nil
}

// Update updates a seat hold
func (r *SeatHoldRepositoryImpl) Update(ctx context.Context, seatHold *models.SeatHold) error {
	if err := r.db.WithContext(ctx).Save(seatHold).Error; err != nil {
		log.Printf("Error updating seat hold: %v", err)
		return fmt.Errorf("failed to update seat hold: %w", err)
	}
	return nil
}

// Delete deletes a seat hold
func (r *SeatHoldRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&models.SeatHold{}, "id = ?", id).Error; err != nil {
		log.Printf("Error deleting seat hold: %v", err)
		return fmt.Errorf("failed to delete seat hold: %w", err)
	}
	return nil
}

// GetByPeriodID retrieves seat holds by period ID
func (r *SeatHoldRepositoryImpl) GetByPeriodID(ctx context.Context, periodID uuid.UUID, includeInactive bool) ([]*models.SeatHold, error) {
	var seatHolds []*models.SeatHold
	query := r.db.WithContext(ctx).
		Preload("Admin").
		Where("period_id = ?", periodID)

	if !includeInactive {
		query = query.Where("status = ?", models.SeatHoldStatusActive)
	}

	err := query.Order("created_at DESC").Find(&seatHolds).Error
	if err != nil {
		log.Printf("Error getting seat holds by period ID: %v", err)
		return nil, fmt.Errorf("failed to get seat holds by period: %w", err)
	}
	return seatHolds, nil
}

// GetByAdminID retrieves seat holds by admin ID with pagination
func (r *SeatHoldRepositoryImpl) GetByAdminID(ctx context.Context, adminID uuid.UUID, offset, limit int) ([]*models.SeatHold, error) {
	var seatHolds []*models.SeatHold
	err := r.db.WithContext(ctx).
		Preload("Period").
		Where("admin_id = ?", adminID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&seatHolds).Error

	if err != nil {
		log.Printf("Error getting seat holds by admin ID: %v", err)
		return nil, fmt.Errorf("failed to get seat holds by admin: %w", err)
	}
	return seatHolds, nil
}

// GetActiveHolds retrieves all active seat holds with pagination
func (r *SeatHoldRepositoryImpl) GetActiveHolds(ctx context.Context, offset, limit int) ([]*models.SeatHold, error) {
	var seatHolds []*models.SeatHold
	err := r.db.WithContext(ctx).
		Preload("Period").
		Preload("Admin").
		Where("status = ?", models.SeatHoldStatusActive).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&seatHolds).Error

	if err != nil {
		log.Printf("Error getting active seat holds: %v", err)
		return nil, fmt.Errorf("failed to get active seat holds: %w", err)
	}
	return seatHolds, nil
}

// GetExpiredHolds retrieves expired holds that should be cleaned up
func (r *SeatHoldRepositoryImpl) GetExpiredHolds(ctx context.Context, beforeTime time.Time) ([]*models.SeatHold, error) {
	var seatHolds []*models.SeatHold
	err := r.db.WithContext(ctx).
		Where("status = ? AND expires_at IS NOT NULL AND expires_at <= ?",
			models.SeatHoldStatusActive, beforeTime).
		Find(&seatHolds).Error

	if err != nil {
		log.Printf("Error getting expired seat holds: %v", err)
		return nil, fmt.Errorf("failed to get expired seat holds: %w", err)
	}
	return seatHolds, nil
}

// GetTotalHeldSeats gets total held seats for a period (including all statuses)
func (r *SeatHoldRepositoryImpl) GetTotalHeldSeats(ctx context.Context, periodID uuid.UUID) (int, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.SeatHold{}).
		Where("period_id = ?", periodID).
		Select("COALESCE(SUM(held_seats), 0)").
		Scan(&total).Error

	if err != nil {
		log.Printf("Error getting total held seats: %v", err)
		return 0, fmt.Errorf("failed to get total held seats: %w", err)
	}
	return int(total), nil
}

// GetActiveHeldSeats gets currently active held seats for a period
func (r *SeatHoldRepositoryImpl) GetActiveHeldSeats(ctx context.Context, periodID uuid.UUID) (int, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.SeatHold{}).
		Where("period_id = ? AND status = ?", periodID, models.SeatHoldStatusActive).
		Select("COALESCE(SUM(held_seats), 0)").
		Scan(&total).Error

	if err != nil {
		log.Printf("Error getting active held seats: %v", err)
		return 0, fmt.Errorf("failed to get active held seats: %w", err)
	}
	return int(total), nil
}

// GetSeatHoldSummary gets aggregated information about seat holds for a period
func (r *SeatHoldRepositoryImpl) GetSeatHoldSummary(ctx context.Context, periodID uuid.UUID) (*models.SeatHoldSummary, error) {
	summary := &models.SeatHoldSummary{
		PeriodID: periodID,
	}

	// Get counts and totals
	type summaryData struct {
		ActiveHolds    int64      `json:"active_holds"`
		HeldSeats      int64      `json:"held_seats"`
		ExpiredHolds   int64      `json:"expired_holds"`
		EarliestExpiry *time.Time `json:"earliest_expiry"`
		LatestExpiry   *time.Time `json:"latest_expiry"`
	}

	var data summaryData
	err := r.db.WithContext(ctx).
		Model(&models.SeatHold{}).
		Select(`
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_holds,
			COALESCE(SUM(CASE WHEN status = 'active' THEN held_seats ELSE 0 END), 0) as held_seats,
			COUNT(CASE WHEN status = 'expired' THEN 1 END) as expired_holds,
			MIN(CASE WHEN status = 'active' AND expires_at IS NOT NULL THEN expires_at END) as earliest_expiry,
			MAX(CASE WHEN status = 'active' AND expires_at IS NOT NULL THEN expires_at END) as latest_expiry
		`).
		Where("period_id = ?", periodID).
		Scan(&data).Error

	if err != nil {
		log.Printf("Error getting seat hold summary: %v", err)
		return nil, fmt.Errorf("failed to get seat hold summary: %w", err)
	}

	summary.TotalActiveHolds = int(data.ActiveHolds)
	summary.TotalHeldSeats = int(data.HeldSeats)
	summary.TotalExpiredHolds = int(data.ExpiredHolds)
	summary.EarliestExpiry = data.EarliestExpiry
	summary.LatestExpiry = data.LatestExpiry

	return summary, nil
}

// ExpireOldHolds marks old holds as expired
func (r *SeatHoldRepositoryImpl) ExpireOldHolds(ctx context.Context, beforeTime time.Time) (int, error) {
	result := r.db.WithContext(ctx).
		Model(&models.SeatHold{}).
		Where("status = ? AND expires_at IS NOT NULL AND expires_at <= ?",
			models.SeatHoldStatusActive, beforeTime).
		Updates(map[string]interface{}{
			"status":     models.SeatHoldStatusExpired,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		log.Printf("Error expiring old holds: %v", result.Error)
		return 0, fmt.Errorf("failed to expire old holds: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}

// ReleaseHoldsByAdmin releases all active holds by an admin
func (r *SeatHoldRepositoryImpl) ReleaseHoldsByAdmin(ctx context.Context, adminID uuid.UUID) (int, error) {
	result := r.db.WithContext(ctx).
		Model(&models.SeatHold{}).
		Where("admin_id = ? AND status = ?", adminID, models.SeatHoldStatusActive).
		Updates(map[string]interface{}{
			"status":     models.SeatHoldStatusReleased,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		log.Printf("Error releasing holds by admin: %v", result.Error)
		return 0, fmt.Errorf("failed to release holds by admin: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}

// ReleaseHoldsByPeriod releases all active holds for a period
func (r *SeatHoldRepositoryImpl) ReleaseHoldsByPeriod(ctx context.Context, periodID uuid.UUID) (int, error) {
	result := r.db.WithContext(ctx).
		Model(&models.SeatHold{}).
		Where("period_id = ? AND status = ?", periodID, models.SeatHoldStatusActive).
		Updates(map[string]interface{}{
			"status":     models.SeatHoldStatusReleased,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		log.Printf("Error releasing holds by period: %v", result.Error)
		return 0, fmt.Errorf("failed to release holds by period: %w", result.Error)
	}

	return int(result.RowsAffected), nil
}

// CheckAvailableSeatsForHold checks if there are enough available seats for a new hold
func (r *SeatHoldRepositoryImpl) CheckAvailableSeatsForHold(ctx context.Context, periodID uuid.UUID, requestedSeats int) (bool, error) {
	// Get period info
	var period models.Period
	err := r.db.WithContext(ctx).First(&period, "id = ?", periodID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, fmt.Errorf("period not found")
		}
		return false, fmt.Errorf("failed to get period: %w", err)
	}

	// Get currently held seats
	heldSeats, err := r.GetActiveHeldSeats(ctx, periodID)
	if err != nil {
		return false, fmt.Errorf("failed to get held seats: %w", err)
	}

	// Calculate real-time booked passengers
	var bookedPassengers int
	activeStatuses := []string{"pending", "processing", "confirm"}
	err = r.db.WithContext(ctx).Model(&models.Booking{}).
		Where("period_id = ? AND status IN ?", periodID, activeStatuses).
		Select("COALESCE(SUM(sum_qty), 0)").
		Scan(&bookedPassengers).Error
	if err != nil {
		log.Printf("Error calculating booked passengers: %v", err)
		// Conservative fallback: assume no seats available
		return false, fmt.Errorf("failed to calculate booked passengers: %w", err)
	}

	// Calculate available seats using real-time data
	availableSeats := period.Seat - bookedPassengers - heldSeats

	return availableSeats >= requestedSeats, nil
}

// GetConflictingHolds gets existing active holds that might conflict
func (r *SeatHoldRepositoryImpl) GetConflictingHolds(ctx context.Context, periodID uuid.UUID, adminID uuid.UUID) ([]*models.SeatHold, error) {
	var seatHolds []*models.SeatHold
	err := r.db.WithContext(ctx).
		Where("period_id = ? AND admin_id = ? AND status = ?",
			periodID, adminID, models.SeatHoldStatusActive).
		Find(&seatHolds).Error

	if err != nil {
		log.Printf("Error getting conflicting holds: %v", err)
		return nil, fmt.Errorf("failed to get conflicting holds: %w", err)
	}
	return seatHolds, nil
}

// GetHoldsExpiringBefore gets holds that expire before a certain time
func (r *SeatHoldRepositoryImpl) GetHoldsExpiringBefore(ctx context.Context, beforeTime time.Time, limit int) ([]*models.SeatHold, error) {
	var seatHolds []*models.SeatHold
	err := r.db.WithContext(ctx).
		Preload("Period").
		Preload("Admin").
		Where("status = ? AND expires_at IS NOT NULL AND expires_at <= ?",
			models.SeatHoldStatusActive, beforeTime).
		Order("expires_at ASC").
		Limit(limit).
		Find(&seatHolds).Error

	if err != nil {
		log.Printf("Error getting holds expiring before: %v", err)
		return nil, fmt.Errorf("failed to get expiring holds: %w", err)
	}
	return seatHolds, nil
}

// GetHoldsByStatus gets holds by status with pagination
func (r *SeatHoldRepositoryImpl) GetHoldsByStatus(ctx context.Context, status models.SeatHoldStatus, offset, limit int) ([]*models.SeatHold, error) {
	var seatHolds []*models.SeatHold
	err := r.db.WithContext(ctx).
		Preload("Period").
		Preload("Admin").
		Where("status = ?", status).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&seatHolds).Error

	if err != nil {
		log.Printf("Error getting holds by status: %v", err)
		return nil, fmt.Errorf("failed to get holds by status: %w", err)
	}
	return seatHolds, nil
}

// GetAdminHoldStatistics gets statistics for an admin's holds
func (r *SeatHoldRepositoryImpl) GetAdminHoldStatistics(ctx context.Context, adminID uuid.UUID) (map[string]interface{}, error) {
	type stats struct {
		TotalHolds     int64 `json:"total_holds"`
		ActiveHolds    int64 `json:"active_holds"`
		ReleasedHolds  int64 `json:"released_holds"`
		ExpiredHolds   int64 `json:"expired_holds"`
		ConvertedHolds int64 `json:"converted_holds"`
		TotalSeatsHeld int64 `json:"total_seats_held"`
	}

	var data stats
	err := r.db.WithContext(ctx).
		Model(&models.SeatHold{}).
		Select(`
			COUNT(*) as total_holds,
			COUNT(CASE WHEN status = 'active' THEN 1 END) as active_holds,
			COUNT(CASE WHEN status = 'released' THEN 1 END) as released_holds,
			COUNT(CASE WHEN status = 'expired' THEN 1 END) as expired_holds,
			COUNT(CASE WHEN status = 'converted' THEN 1 END) as converted_holds,
			COALESCE(SUM(CASE WHEN status = 'active' THEN held_seats ELSE 0 END), 0) as total_seats_held
		`).
		Where("admin_id = ?", adminID).
		Scan(&data).Error

	if err != nil {
		log.Printf("Error getting admin hold statistics: %v", err)
		return nil, fmt.Errorf("failed to get admin statistics: %w", err)
	}

	return map[string]interface{}{
		"total_holds":      data.TotalHolds,
		"active_holds":     data.ActiveHolds,
		"released_holds":   data.ReleasedHolds,
		"expired_holds":    data.ExpiredHolds,
		"converted_holds":  data.ConvertedHolds,
		"total_seats_held": data.TotalSeatsHeld,
	}, nil
}

// CountByPeriodID counts seat holds by period ID
func (r *SeatHoldRepositoryImpl) CountByPeriodID(ctx context.Context, periodID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.SeatHold{}).
		Where("period_id = ?", periodID).
		Count(&count).Error

	if err != nil {
		log.Printf("Error counting holds by period: %v", err)
		return 0, fmt.Errorf("failed to count holds by period: %w", err)
	}
	return count, nil
}

// CountByAdminID counts seat holds by admin ID
func (r *SeatHoldRepositoryImpl) CountByAdminID(ctx context.Context, adminID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.SeatHold{}).
		Where("admin_id = ?", adminID).
		Count(&count).Error

	if err != nil {
		log.Printf("Error counting holds by admin: %v", err)
		return 0, fmt.Errorf("failed to count holds by admin: %w", err)
	}
	return count, nil
}

// CountByStatus counts seat holds by status
func (r *SeatHoldRepositoryImpl) CountByStatus(ctx context.Context, status models.SeatHoldStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.SeatHold{}).
		Where("status = ?", status).
		Count(&count).Error

	if err != nil {
		log.Printf("Error counting holds by status: %v", err)
		return 0, fmt.Errorf("failed to count holds by status: %w", err)
	}
	return count, nil
}

// FindActiveBYPeriodID finds active seat holds by period ID
func (r *SeatHoldRepositoryImpl) FindActiveBYPeriodID(periodID uuid.UUID) ([]*models.SeatHold, error) {
	log.Printf("🔍 Finding active seat holds for period: %s", periodID)

	var seatHolds []*models.SeatHold

	err := r.db.Where("period_id = ? AND status = ?", periodID, models.SeatHoldStatusActive).
		Find(&seatHolds).Error

	if err != nil {
		log.Printf("❌ Error finding active seat holds: %v", err)
		return nil, fmt.Errorf("failed to find active seat holds: %w", err)
	}

	log.Printf("✅ Found %d active seat holds for period %s", len(seatHolds), periodID)
	return seatHolds, nil
}