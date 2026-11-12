package dto

// Request DTOs

type CreateCompanyBankAccountRequest struct {
	BankCode      string  `json:"bankCode" validate:"required,max=10"`
	BankName      string  `json:"bankName" validate:"required,max=100"`
	AccountNumber string  `json:"accountNumber" validate:"required,max=20"`
	AccountName   string  `json:"accountName" validate:"required,max=200"`
	BranchName    *string `json:"branchName,omitempty"`
	IsActive      *bool   `json:"isActive,omitempty"`
	SortOrder     *int    `json:"sortOrder,omitempty"`
}

type UpdateCompanyBankAccountRequest struct {
	BankCode      *string `json:"bankCode,omitempty" validate:"omitempty,max=10"`
	BankName      *string `json:"bankName,omitempty" validate:"omitempty,max=100"`
	AccountNumber *string `json:"accountNumber,omitempty" validate:"omitempty,max=20"`
	AccountName   *string `json:"accountName,omitempty" validate:"omitempty,max=200"`
	BranchName    *string `json:"branchName,omitempty"`
	IsActive      *bool   `json:"isActive,omitempty"`
	SortOrder     *int    `json:"sortOrder,omitempty"`
}

// Response DTO

type CompanyBankAccountResponse struct {
	ID            string  `json:"id"`
	BankCode      string  `json:"bankCode"`
	BankName      string  `json:"bankName"`
	AccountNumber string  `json:"accountNumber"`
	AccountName   string  `json:"accountName"`
	BranchName    *string `json:"branchName,omitempty"`
	IsActive      bool    `json:"isActive"`
	SortOrder     int     `json:"sortOrder"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}
