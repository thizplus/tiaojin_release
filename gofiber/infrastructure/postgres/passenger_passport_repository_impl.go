package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type PassengerPassportRepositoryImpl struct {
	db *gorm.DB
}

func NewPassengerPassportRepository(db *gorm.DB) repositories.PassengerPassportRepository {
	return &PassengerPassportRepositoryImpl{db: db}
}

func (r *PassengerPassportRepositoryImpl) Create(ctx context.Context, entity *models.PassengerPassport) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create passenger passport: %w", err)
	}
	return nil
}

func (r *PassengerPassportRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.PassengerPassport, error) {
	var passport models.PassengerPassport
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&passport).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get passenger passport by ID: %w", err)
	}

	return &passport, nil
}

func (r *PassengerPassportRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.PassengerPassport) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update passenger passport: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("passenger passport with ID %s not found", id)
	}

	return nil
}

func (r *PassengerPassportRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.PassengerPassport{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete passenger passport: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("passenger passport with ID %s not found", id)
	}

	return nil
}

func (r *PassengerPassportRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.PassengerPassport, error) {
	var passports []*models.PassengerPassport

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&passports).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list passenger passports: %w", err)
	}

	return passports, nil
}

func (r *PassengerPassportRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.PassengerPassport{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count passenger passports: %w", err)
	}

	return count, nil
}

func (r *PassengerPassportRepositoryImpl) GetByBookingID(ctx context.Context, bookingID uuid.UUID) ([]*models.PassengerPassport, error) {
	var passports []*models.PassengerPassport

	err := r.db.WithContext(ctx).
		Where("booking_id = ?", bookingID).
		Order("created_at ASC").
		Find(&passports).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get passenger passports by booking ID: %w", err)
	}

	return passports, nil
}

func (r *PassengerPassportRepositoryImpl) GetByPassportNumber(ctx context.Context, passportNumber string) (*models.PassengerPassport, error) {
	var passport models.PassengerPassport
	err := r.db.WithContext(ctx).Where("passport_number = ?", passportNumber).First(&passport).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get passenger passport by passport number: %w", err)
	}

	return &passport, nil
}

func (r *PassengerPassportRepositoryImpl) GetByStatus(ctx context.Context, status models.PassportStatus, offset, limit int) ([]*models.PassengerPassport, error) {
	var passports []*models.PassengerPassport

	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&passports).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get passenger passports by status: %w", err)
	}

	return passports, nil
}

func (r *PassengerPassportRepositoryImpl) GetPendingProcessing(ctx context.Context, offset, limit int) ([]*models.PassengerPassport, error) {
	var passports []*models.PassengerPassport

	err := r.db.WithContext(ctx).
		Where("status = ?", models.PassportStatusPending).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&passports).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending processing passports: %w", err)
	}

	return passports, nil
}

func (r *PassengerPassportRepositoryImpl) GetExpiring(ctx context.Context, daysBeforeExpiry int, offset, limit int) ([]*models.PassengerPassport, error) {
	var passports []*models.PassengerPassport
	expiryDate := time.Now().AddDate(0, 0, daysBeforeExpiry)

	err := r.db.WithContext(ctx).
		Where("date_of_expiry <= ? AND date_of_expiry > ?", expiryDate, time.Now()).
		Order("date_of_expiry ASC").
		Offset(offset).
		Limit(limit).
		Find(&passports).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get expiring passports: %w", err)
	}

	return passports, nil
}

func (r *PassengerPassportRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.PassportStatus) error {
	result := r.db.WithContext(ctx).
		Model(&models.PassengerPassport{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update passenger passport status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("passenger passport with ID %s not found", id)
	}

	return nil
}

func (r *PassengerPassportRepositoryImpl) UpdatePersonalInfo(ctx context.Context, id uuid.UUID, passport *models.PassengerPassport) error {
	updates := map[string]interface{}{
		"name_en":                passport.NameEN,
		"surname_en":             passport.SurnameEN,
		"name_th":                passport.NameTH,
		"surname_th":             passport.SurnameTH,
		"passport_number":        passport.PassportNumber,
		"date_of_expiry":         passport.DateOfExpiry,
		"nationality":            passport.Nationality,
		"date_of_birth":          passport.DateOfBirth,
		"gender":                 passport.Gender,
	}

	result := r.db.WithContext(ctx).
		Model(&models.PassengerPassport{}).
		Where("id = ?", id).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update passenger passport personal info: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("passenger passport with ID %s not found", id)
	}

	return nil
}

func (r *PassengerPassportRepositoryImpl) CreateBatch(ctx context.Context, entities []*models.PassengerPassport) error {
	if len(entities) == 0 {
		return nil
	}

	if err := r.db.WithContext(ctx).CreateInBatches(entities, 100).Error; err != nil {
		return fmt.Errorf("failed to create passenger passports in batch: %w", err)
	}

	return nil
}

func (r *PassengerPassportRepositoryImpl) DeleteByBookingID(ctx context.Context, bookingID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.PassengerPassport{}, "booking_id = ?", bookingID)

	if result.Error != nil {
		return fmt.Errorf("failed to delete passenger passports by booking ID: %w", result.Error)
	}

	return nil
}

// WithTx returns a new repository instance with the given transaction
func (r *PassengerPassportRepositoryImpl) WithTx(tx *gorm.DB) repositories.PassengerPassportRepository {
	return &PassengerPassportRepositoryImpl{db: tx}
}