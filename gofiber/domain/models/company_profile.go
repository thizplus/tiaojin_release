package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyProfile struct {
	ID       uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Name     string    `gorm:"not null;size:200" json:"name"`
	TaxID    string    `gorm:"not null;size:13;column:tax_id" json:"tax_id"`
	Phone    string    `gorm:"not null;size:20" json:"phone"`
	Address  string    `gorm:"type:text;not null" json:"address"`
	Province string    `gorm:"not null;size:100" json:"province"`
	Zipcode  string    `gorm:"not null;size:5" json:"zipcode"`
	ATTA     string    `gorm:"size:500;column:atta" json:"atta"`
	DBD      string    `gorm:"size:500;column:dbd" json:"dbd"`
	PP20     string    `gorm:"size:500;column:pp20" json:"pp20"`
	Logo           *string `gorm:"size:500" json:"logo"`
	PdfHeaderImage *string `gorm:"size:500;column:pdf_header_image" json:"pdf_header_image"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User *User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}

func (CompanyProfile) TableName() string {
	return "company_profiles"
}