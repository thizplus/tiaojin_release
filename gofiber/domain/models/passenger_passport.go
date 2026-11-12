package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PassportStatus string

const (
	PassportStatusPending   PassportStatus = "pending"
	PassportStatusProcessed PassportStatus = "processed"
	PassportStatusError     PassportStatus = "error"
)

type PassengerPassport struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	BookingID uuid.UUID `gorm:"type:uuid;not null;index" json:"booking_id"`
	//PassportID      *string        `gorm:"size:50;column:passport_id" json:"passport_id"`
	Image          *string        `gorm:"size:500" json:"image"`
	PassportNumber string         `gorm:"not null;size:20;column:passport_number" json:"passport_number"`
	Status         PassportStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ProcessedDate  *time.Time     `gorm:"column:processed_date" json:"processed_date"`
	ErrorMessage   *string        `gorm:"type:text;column:error_message" json:"error_message"`

	TitleTH          *string    `gorm:"size:20;column:title_th" json:"title_th"`
	NameTH           *string    `gorm:"size:100;column:name_th" json:"name_th"`
	SurnameTH        *string    `gorm:"size:100;column:surname_th" json:"surname_th"`
	TitleEN          *string    `gorm:"size:20;column:title_en" json:"title_en"`
	NameEN           *string    `gorm:"size:100;column:name_en" json:"name_en"`
	SurnameEN        *string    `gorm:"size:100;column:surname_en" json:"surname_en"`
	Gender           *string    `gorm:"size:1" json:"gender"`
	Nationality      *string    `gorm:"size:50" json:"nationality"`
	CountryCode      *string    `gorm:"size:3;column:country_code" json:"country_code"`
	DateOfBirth      *time.Time `gorm:"type:date;column:date_of_birth" json:"date_of_birth"`
	IDNumber         *string    `gorm:"size:20;column:id_number" json:"id_number"`
	PlaceOfBirth     *string    `gorm:"size:100;column:place_of_birth" json:"place_of_birth"`
	DateOfExpiry     *time.Time `gorm:"type:date;column:date_of_expiry" json:"date_of_expiry"`
	DateOfIssue      *time.Time `gorm:"type:date;column:date_of_issue" json:"date_of_issue"`
	IssuingAuthority *string    `gorm:"size:100;column:issuing_authority" json:"issuing_authority"`

	// Vision API OCR fields
	OCRText       *string  `gorm:"type:text;column:ocr_text" json:"ocr_text"`
	OCRConfidence *float32 `gorm:"column:ocr_confidence" json:"ocr_confidence"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Booking *Booking `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"booking,omitempty"`
}

func (PassengerPassport) TableName() string {
	return "passenger_passports"
}

func (pp *PassengerPassport) IsProcessed() bool {
	return pp.Status == PassportStatusProcessed
}

func (pp *PassengerPassport) HasError() bool {
	return pp.Status == PassportStatusError
}

func (pp *PassengerPassport) IsPending() bool {
	return pp.Status == PassportStatusPending
}

func (pp *PassengerPassport) GetFullNameTH() string {
	if pp.NameTH != nil && pp.SurnameTH != nil {
		return *pp.NameTH + " " + *pp.SurnameTH
	}
	return ""
}

func (pp *PassengerPassport) GetFullNameEN() string {
	if pp.NameEN != nil && pp.SurnameEN != nil {
		return *pp.NameEN + " " + *pp.SurnameEN
	}
	return ""
}

func (pp *PassengerPassport) IsPassportExpired() bool {
	if pp.DateOfExpiry == nil {
		return false
	}
	return pp.DateOfExpiry.Before(time.Now())
}
