package serviceimpl

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
)

type CompanyBankAccountServiceImpl struct {
	repo repositories.CompanyBankAccountRepository
}

func NewCompanyBankAccountService(repo repositories.CompanyBankAccountRepository) services.CompanyBankAccountService {
	return &CompanyBankAccountServiceImpl{repo: repo}
}

func (s *CompanyBankAccountServiceImpl) Create(ctx context.Context, req dto.CreateCompanyBankAccountRequest) (*dto.CompanyBankAccountResponse, error) {
	entity := &models.CompanyBankAccount{
		BankCode:      req.BankCode,
		BankName:      req.BankName,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
		BranchName:    req.BranchName,
		IsActive:      true,
		SortOrder:     0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		entity.SortOrder = *req.SortOrder
	}

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, fmt.Errorf("failed to create bank account: %w", err)
	}

	return s.toResponse(entity), nil
}

func (s *CompanyBankAccountServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*dto.CompanyBankAccountResponse, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, fmt.Errorf("bank account not found")
	}
	return s.toResponse(entity), nil
}

func (s *CompanyBankAccountServiceImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCompanyBankAccountRequest) (*dto.CompanyBankAccountResponse, error) {
	entity, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, fmt.Errorf("bank account not found")
	}

	if req.BankCode != nil {
		entity.BankCode = *req.BankCode
	}
	if req.BankName != nil {
		entity.BankName = *req.BankName
	}
	if req.AccountNumber != nil {
		entity.AccountNumber = *req.AccountNumber
	}
	if req.AccountName != nil {
		entity.AccountName = *req.AccountName
	}
	if req.BranchName != nil {
		entity.BranchName = req.BranchName
	}
	if req.IsActive != nil {
		entity.IsActive = *req.IsActive
	}
	if req.SortOrder != nil {
		entity.SortOrder = *req.SortOrder
	}

	entity.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, id, entity); err != nil {
		return nil, err
	}

	return s.toResponse(entity), nil
}

func (s *CompanyBankAccountServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *CompanyBankAccountServiceImpl) List(ctx context.Context) ([]*dto.CompanyBankAccountResponse, error) {
	entities, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	return s.toResponseList(entities), nil
}

func (s *CompanyBankAccountServiceImpl) ListActive(ctx context.Context) ([]*dto.CompanyBankAccountResponse, error) {
	entities, err := s.repo.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	return s.toResponseList(entities), nil
}

func (s *CompanyBankAccountServiceImpl) toResponse(entity *models.CompanyBankAccount) *dto.CompanyBankAccountResponse {
	return &dto.CompanyBankAccountResponse{
		ID:            entity.ID.String(),
		BankCode:      entity.BankCode,
		BankName:      entity.BankName,
		AccountNumber: entity.AccountNumber,
		AccountName:   entity.AccountName,
		BranchName:    entity.BranchName,
		IsActive:      entity.IsActive,
		SortOrder:     entity.SortOrder,
		CreatedAt:     entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     entity.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *CompanyBankAccountServiceImpl) toResponseList(entities []*models.CompanyBankAccount) []*dto.CompanyBankAccountResponse {
	responses := make([]*dto.CompanyBankAccountResponse, len(entities))
	for i, entity := range entities {
		responses[i] = s.toResponse(entity)
	}
	return responses
}
