package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CityStatus string

const (
	CityStatusPublish CityStatus = "publish"
	CityStatusDraft   CityStatus = "draft"
)

type City struct {
	ID     uuid.UUID  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title  string     `gorm:"not null;size:200" json:"title"`
	Detail *string    `gorm:"type:text" json:"detail"`
	Image  string     `gorm:"not null;size:500;column:img" json:"image"`
	Status CityStatus `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Tours []Tour `gorm:"foreignKey:CityID" json:"tours,omitempty"`
}

func (City) TableName() string {
	return "cities"
}