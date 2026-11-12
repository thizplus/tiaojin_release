package dto

import "github.com/google/uuid"

type CompanyProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	TaxID     string    `json:"tax_id"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Province  string    `json:"province"`
	Zipcode   string    `json:"zipcode"`
	ATTA      string    `json:"atta"`
	DBD       string    `json:"dbd"`
	PP20      string    `json:"pp20"`
	Logo           *string `json:"logo"`
	PdfHeaderImage *string `json:"pdfHeaderImage"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type CreateCompanyProfileRequest struct {
	UserID   uuid.UUID `json:"user_id" validate:"required"`
	Name     string    `json:"name" validate:"required,max=200"`
	TaxID    string    `json:"tax_id" validate:"required,len=13"`
	Phone    string    `json:"phone" validate:"required"`
	Address  string    `json:"address" validate:"required"`
	Province string    `json:"province" validate:"required"`
	Zipcode  string    `json:"zipcode" validate:"required,len=5"`
	ATTA     string    `json:"atta,omitempty"`
	DBD      string    `json:"dbd,omitempty"`
	PP20     string    `json:"pp20,omitempty"`
	Logo           *string `json:"logo,omitempty"`
	PdfHeaderImage *string `json:"pdfHeaderImage,omitempty"`
}

type UpdateCompanyProfileRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,max=200"`
	TaxID    *string `json:"tax_id,omitempty" validate:"omitempty,len=13"`
	Phone    *string `json:"phone,omitempty"`
	Address  *string `json:"address,omitempty"`
	Province *string `json:"province,omitempty"`
	Zipcode  *string `json:"zipcode,omitempty" validate:"omitempty,len=5"`
	ATTA     *string `json:"atta,omitempty"`
	DBD      *string `json:"dbd,omitempty"`
	PP20     *string `json:"pp20,omitempty"`
	Logo           *string `json:"logo,omitempty"`
	PdfHeaderImage *string `json:"pdfHeaderImage,omitempty"`
}

type UploadLogoRequest struct {
	UserID   uuid.UUID `json:"user_id" validate:"required"`
	LogoData []byte    `json:"logo_data" validate:"required"`
	FileName string    `json:"file_name" validate:"required"`
	FileType string    `json:"file_type" validate:"required,oneof=image/jpeg image/png image/webp"`
}

type UploadPdfHeaderImageRequest struct {
	UserID    uuid.UUID `json:"user_id" validate:"required"`
	ImageData []byte    `json:"image_data" validate:"required"`
	FileName  string    `json:"file_name" validate:"required"`
	FileType  string    `json:"file_type" validate:"required,oneof=image/jpeg image/png image/webp"`
}

type CompanyProfileListRequest struct {
	Search      *string `json:"search,omitempty"`
	Offset      int     `json:"offset" validate:"min=0"`
	Limit       int     `json:"limit" validate:"min=1,max=100"`
	PendingOnly bool    `json:"pending_only,omitempty"`
}

type CompanyProfileListResponse struct {
	Profiles   []*CompanyProfileResponse `json:"profiles"`
	Pagination PaginationMeta            `json:"pagination"`
}

type ProfileApprovalResponse struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	UpdatedAt string    `json:"updated_at"`
}