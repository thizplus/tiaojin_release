package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodRoomType struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	PeriodID  uuid.UUID `gorm:"type:uuid;not null;index" json:"period_id"`
	Title     string    `gorm:"not null;size:100" json:"title"`
	Price     float64   `gorm:"type:decimal(10,2);not null;check:price >= 0" json:"price"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Period *Period `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"period,omitempty"`
}

func (PeriodRoomType) TableName() string {
	return "period_room_types"
}

func (prt *PeriodRoomType) CalculateTotal(quantity int) float64 {
	return prt.Price * float64(quantity)
}

func (prt *PeriodRoomType) IsAvailable() bool {
	return prt.Price > 0
}