package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingAdditionalService struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	BookingID uuid.UUID `gorm:"type:uuid;not null;index" json:"booking_id"`
	Title     string    `gorm:"not null;size:100" json:"title"`
	Price     float64   `gorm:"type:decimal(10,2);not null;check:price >= 0" json:"price"`
	Quantity  int       `gorm:"not null;check:quantity > 0" json:"quantity"`
	RateType  string    `gorm:"not null;size:20" json:"rate_type"` // per_unit, per_booking
	UnitName  string    `gorm:"not null;size:50" json:"unit_name"` // คน, วัน, ชิ้น, เหมา

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Booking *Booking `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"booking,omitempty"`
}

func (BookingAdditionalService) TableName() string {
	return "booking_additional_services"
}

func (bas *BookingAdditionalService) TotalPrice() float64 {
	return bas.Price * float64(bas.Quantity)
}

func (bas *BookingAdditionalService) IsPerUnitRate() bool {
	return bas.RateType == "per_unit"
}

func (bas *BookingAdditionalService) IsPerBookingRate() bool {
	return bas.RateType == "per_booking"
}

func (bas *BookingAdditionalService) IsPersonUnit() bool {
	return bas.UnitName == "คน"
}

func (bas *BookingAdditionalService) IsDayUnit() bool {
	return bas.UnitName == "วัน"
}

func (bas *BookingAdditionalService) IsPieceUnit() bool {
	return bas.UnitName == "ชิ้น"
}

func (bas *BookingAdditionalService) IsPackageUnit() bool {
	return bas.UnitName == "เหมา"
}

func (bas *BookingAdditionalService) CalculateTotal(passengerCount int) float64 {
	switch bas.RateType {
	case "per_unit":
		// For per_unit, quantity represents the number of units
		// For person units, we might multiply by passenger count if it's meant to be per person
		if bas.UnitName == "คน" {
			// For person units, use quantity as the number of people (not necessarily all passengers)
			return bas.Price * float64(bas.Quantity)
		}
		// For other units (วัน, ชิ้น), just use quantity
		return bas.Price * float64(bas.Quantity)
	case "per_booking":
		// For per_booking, always flat rate regardless of passenger count
		return bas.Price * float64(bas.Quantity)
	default:
		return bas.TotalPrice()
	}
}

func (bas *BookingAdditionalService) GetUnitDisplayName() string {
	return bas.UnitName
}

func (bas *BookingAdditionalService) GetRateDescription() string {
	switch bas.RateType {
	case "per_unit":
		return fmt.Sprintf("ราคาต่อ%s", bas.UnitName)
	case "per_booking":
		return "ราคาเหมา"
	default:
		return "ราคาพิเศษ"
	}
}