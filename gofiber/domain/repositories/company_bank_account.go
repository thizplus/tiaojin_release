package repositories

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/models"
)

type CompanyBankAccountRepository interface {
	Create(ctx context.Context, entity *models.CompanyBankAccount) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.CompanyBankAccount, error)
	Update(ctx context.Context, id uuid.UUID, entity *models.CompanyBankAccount) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*models.CompanyBankAccount, error)        // ทั้งหมด (admin)
	ListActive(ctx context.Context) ([]*models.CompanyBankAccount, error)  // เฉพาะ active (agent/public)
	Count(ctx context.Context) (int64, error)
}
