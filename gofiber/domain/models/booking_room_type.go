package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingRoomType struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	BookingID uuid.UUID `gorm:"type:uuid;not null;index" json:"booking_id"`
	Title     string    `gorm:"not null;size:100" json:"title"`
	Price     float64   `gorm:"type:decimal(10,2);not null;check:price >= 0" json:"price"`
	Quantity  int       `gorm:"not null;check:quantity > 0" json:"quantity"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Booking *Booking `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"booking,omitempty"`
}

func (BookingRoomType) TableName() string {
	return "booking_room_types"
}

func (brt *BookingRoomType) TotalPrice() float64 {
	return brt.Price * float64(brt.Quantity)
}