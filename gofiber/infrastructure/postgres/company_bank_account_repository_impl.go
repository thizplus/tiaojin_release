package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type CompanyBankAccountRepositoryImpl struct {
	db *gorm.DB
}

func NewCompanyBankAccountRepository(db *gorm.DB) repositories.CompanyBankAccountRepository {
	return &CompanyBankAccountRepositoryImpl{db: db}
}

func (r *CompanyBankAccountRepositoryImpl) Create(ctx context.Context, entity *models.CompanyBankAccount) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create company bank account: %w", err)
	}
	return nil
}

func (r *CompanyBankAccountRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.CompanyBankAccount, error) {
	var account models.CompanyBankAccount
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get company bank account: %w", err)
	}
	return &account, nil
}

func (r *CompanyBankAccountRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.CompanyBankAccount) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)
	if result.Error != nil {
		return fmt.Errorf("failed to update company bank account: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("company bank account with ID %s not found", id)
	}
	return nil
}

func (r *CompanyBankAccountRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.CompanyBankAccount{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete company bank account: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("company bank account with ID %s not found", id)
	}
	return nil
}

func (r *CompanyBankAccountRepositoryImpl) List(ctx context.Context) ([]*models.CompanyBankAccount, error) {
	var accounts []*models.CompanyBankAccount
	err := r.db.WithContext(ctx).Order("sort_order ASC, created_at DESC").Find(&accounts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list company bank accounts: %w", err)
	}
	return accounts, nil
}

func (r *CompanyBankAccountRepositoryImpl) ListActive(ctx context.Context) ([]*models.CompanyBankAccount, error) {
	var accounts []*models.CompanyBankAccount
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Order("sort_order ASC, created_at DESC").Find(&accounts).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list active company bank accounts: %w", err)
	}
	return accounts, nil
}

func (r *CompanyBankAccountRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.CompanyBankAccount{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count company bank accounts: %w", err)
	}
	return count, nil
}
