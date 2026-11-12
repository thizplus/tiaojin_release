package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TourStatus string

const (
	TourStatusPublish TourStatus = "publish"
	TourStatusDraft   TourStatus = "draft"
)

type Tour struct {
	ID        uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title     string     `gorm:"not null;size:300" json:"title"`
	Code      string     `gorm:"not null;uniqueIndex;size:20" json:"code"`
	CityID    *uuid.UUID `gorm:"type:uuid;index" json:"city_id,omitempty"`
	Detail    *string    `gorm:"type:text" json:"detail"`
	Thumbnail *string    `gorm:"size:500" json:"thumbnail"`
	FilePDF   *string    `gorm:"size:500;column:file_pdf" json:"file_pdf"`
	FileWord  *string    `gorm:"size:500;column:file_word" json:"file_word"`
	Status    TourStatus `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	City     *City     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"city,omitempty"`
	Periods  []Period  `gorm:"foreignKey:TourID" json:"periods,omitempty"`
	Bookings []Booking `gorm:"foreignKey:TourID" json:"bookings,omitempty"`
}

func (Tour) TableName() string {
	return "tours"
}