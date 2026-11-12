package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)

type CompanyBankAccountService interface {
	Create(ctx context.Context, req dto.CreateCompanyBankAccountRequest) (*dto.CompanyBankAccountResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.CompanyBankAccountResponse, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateCompanyBankAccountRequest) (*dto.CompanyBankAccountResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]*dto.CompanyBankAccountResponse, error)
	ListActive(ctx context.Context) ([]*dto.CompanyBankAccountResponse, error)
}
