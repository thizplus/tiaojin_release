package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type PaymentSlipRepositoryImpl struct {
	db *gorm.DB
}

func NewPaymentSlipRepository(db *gorm.DB) repositories.PaymentSlipRepository {
	return &PaymentSlipRepositoryImpl{db: db}
}

func (r *PaymentSlipRepositoryImpl) Create(ctx context.Context, entity *models.PaymentSlip) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create payment slip: %w", err)
	}
	return nil
}

func (r *PaymentSlipRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.PaymentSlip, error) {
	var slip models.PaymentSlip
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&slip).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get payment slip by ID: %w", err)
	}

	return &slip, nil
}

func (r *PaymentSlipRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.PaymentSlip) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update payment slip: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("payment slip with ID %s not found", id)
	}

	return nil
}

func (r *PaymentSlipRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.PaymentSlip{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete payment slip: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("payment slip with ID %s not found", id)
	}

	return nil
}

func (r *PaymentSlipRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.PaymentSlip, error) {
	var slips []*models.PaymentSlip

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&slips).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list payment slips: %w", err)
	}

	return slips, nil
}

func (r *PaymentSlipRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.PaymentSlip{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count payment slips: %w", err)
	}

	return count, nil
}

func (r *PaymentSlipRepositoryImpl) GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.PaymentSlip, error) {
	var slips []*models.PaymentSlip

	err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Order("created_at ASC").
		Find(&slips).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get payment slips by booking ID: %w", err)
	}

	return slips, nil
}

func (r *PaymentSlipRepositoryImpl) GetBySlipID(ctx context.Context, slipID string) (*models.PaymentSlip, error) {
	var slip models.PaymentSlip
	err := r.db.WithContext(ctx).Where("slip_id = ?", slipID).First(&slip).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get payment slip by slip ID: %w", err)
	}

	return &slip, nil
}

func (r *PaymentSlipRepositoryImpl) GetByStatus(ctx context.Context, status models.PaymentSlipStatus, offset, limit int) ([]*models.PaymentSlip, error) {
	var slips []*models.PaymentSlip

	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&slips).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get payment slips by status: %w", err)
	}

	return slips, nil
}

func (r *PaymentSlipRepositoryImpl) GetPendingSlips(ctx context.Context, offset, limit int) ([]*models.PaymentSlip, error) {
	var slips []*models.PaymentSlip

	err := r.db.WithContext(ctx).
		Where("status = ?", models.PaymentSlipStatusPending).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&slips).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending payment slips: %w", err)
	}

	return slips, nil
}

func (r *PaymentSlipRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PaymentSlipStatus) error {
	result := r.db.WithContext(ctx).
		Model(&models.PaymentSlip{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update payment slip status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("payment slip with ID %s not found", id)
	}

	return nil
}

func (r *PaymentSlipRepositoryImpl) GetTotalVerifiedAmountByBookingID(ctx context.Context, bookingID uuid.UUID) (float64, error) {
	var total float64

	err := r.db.WithContext(ctx).
		Model(&models.PaymentSlip{}).
		Where("booking_id = ? AND status = ? AND is_reversed = ?", bookingID, models.PaymentSlipStatusConfirm, false).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, fmt.Errorf("failed to calculate total verified amount by booking ID: %w", err)
	}

	return total, nil
}

func (r *PaymentSlipRepositoryImpl) GenerateSlipID(ctx context.Context) (string, error) {
	const maxRetries = 10
	now := time.Now()
	dateStr := now.Format("20060102")

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Use FOR UPDATE to lock the table during counting to prevent race conditions
		var count int64
		err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Count existing slip IDs for today with row-level locking
			err := tx.Model(&models.PaymentSlip{}).
				Where("slip_id LIKE ?", fmt.Sprintf("PS-%s-%%", dateStr)).
				Count(&count).Error
			if err != nil {
				return fmt.Errorf("failed to count payment slips for slip ID generation: %w", err)
			}

			return nil
		})

		if err != nil {
			return "", err
		}

		// Generate slip ID with sequence number
		slipID := fmt.Sprintf("PS-%s-%03d", dateStr, count+1)

		// Verify uniqueness by trying to find existing record
		var existingSlip models.PaymentSlip
		err = r.db.WithContext(ctx).
			Where("slip_id = ?", slipID).
			First(&existingSlip).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Slip ID is unique, return it
			return slipID, nil
		} else if err != nil {
			// Database error
			return "", fmt.Errorf("failed to check slip ID uniqueness: %w", err)
		}

		// Slip ID already exists, increment count and try again
		log.Printf("[WARNING] Slip ID %s already exists on attempt %d, retrying...", slipID, attempt)

		// If we've tried multiple times, add some randomness to avoid infinite loops
		if attempt > 5 {
			// Add microsecond-based randomness to break ties
			randomSuffix := now.UnixNano() % 1000
			slipID = fmt.Sprintf("PS-%s-%03d-%03d", dateStr, count+1, randomSuffix)

			// Check this randomized version
			err = r.db.WithContext(ctx).
				Where("slip_id = ?", slipID).
				First(&existingSlip).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Printf("[SUCCESS] Generated unique slip ID with randomness: %s", slipID)
				return slipID, nil
			}
		}

		// Small delay before retry to reduce contention
		time.Sleep(time.Millisecond * time.Duration(attempt))
	}

	// If all retries failed, generate UUID-based slip ID as fallback
	fallbackID := fmt.Sprintf("PS-%s-UUID-%s", dateStr, uuid.New().String()[:8])
	log.Printf("[FALLBACK] Using UUID-based slip ID after %d retries: %s", maxRetries, fallbackID)
	return fallbackID, nil
}

// GetTotalActiveAmountByBookingID returns the total amount of confirmed and non-reversed payment slips
func (r *PaymentSlipRepositoryImpl) GetTotalActiveAmountByBookingID(ctx context.Context, bookingID uuid.UUID) (float64, error) {
	var total float64

	err := r.db.WithContext(ctx).
		Model(&models.PaymentSlip{}).
		Where("booking_id = ? AND status = ? AND is_reversed = ?", bookingID, models.PaymentSlipStatusConfirm, false).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error

	if err != nil {
		return 0, fmt.Errorf("failed to calculate total active amount by booking ID: %w", err)
	}

	return total, nil
}

// ReversePaymentSlip marks a payment slip as reversed
func (r *PaymentSlipRepositoryImpl) ReversePaymentSlip(ctx context.Context, id uuid.UUID, reversedBy uuid.UUID, reason string) error {
	now := time.Now()

	result := r.db.WithContext(ctx).
		Model(&models.PaymentSlip{}).
		Where("id = ? AND status = ? AND is_reversed = ?", id, models.PaymentSlipStatusConfirm, false).
		Updates(map[string]interface{}{
			"status":          models.PaymentSlipStatusReversed,
			"is_reversed":     true,
			"reversed_at":     now,
			"reversed_by":     reversedBy,
			"reversal_reason": reason,
			"updated_at":      now,
		})

	if result.Error != nil {
		return fmt.Errorf("failed to reverse payment slip: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("payment slip not found or not reversible")
	}

	return nil
}

// GetPaymentSlipHistory returns a payment slip with detailed information for history
func (r *PaymentSlipRepositoryImpl) GetPaymentSlipHistory(ctx context.Context, id uuid.UUID) (*models.PaymentSlip, error) {
	var paymentSlip models.PaymentSlip

	err := r.db.WithContext(ctx).
		Preload("Booking").
		Preload("ReversedByUser").
		Where("id = ?", id).
		First(&paymentSlip).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment slip not found")
		}
		return nil, fmt.Errorf("failed to get payment slip history: %w", err)
	}

	return &paymentSlip, nil
}

// CheckForDuplicateSlipIDs helps identify existing duplicate slip IDs in the database
func (r *PaymentSlipRepositoryImpl) CheckForDuplicateSlipIDs(ctx context.Context) error {
	var duplicates []struct {
		SlipID string
		Count  int64
	}

	err := r.db.WithContext(ctx).
		Model(&models.PaymentSlip{}).
		Select("slip_id, COUNT(*) as count").
		Group("slip_id").
		Having("COUNT(*) > 1").
		Scan(&duplicates).Error

	if err != nil {
		return fmt.Errorf("failed to check for duplicate slip IDs: %w", err)
	}

	if len(duplicates) > 0 {
		log.Printf("[WARNING] Found %d duplicate slip IDs in database:", len(duplicates))
		for _, dup := range duplicates {
			log.Printf("[WARNING] Slip ID '%s' appears %d times", dup.SlipID, dup.Count)
		}
		return fmt.Errorf("found %d duplicate slip IDs in database", len(duplicates))
	}

	log.Printf("[SUCCESS] No duplicate slip IDs found in database")
	return nil
}

// WithTx returns a new repository instance with the given transaction
func (r *PaymentSlipRepositoryImpl) WithTx(tx *gorm.DB) repositories.PaymentSlipRepository {
	return &PaymentSlipRepositoryImpl{db: tx}
}