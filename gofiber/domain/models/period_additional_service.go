package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodAdditionalService struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	PeriodID  uuid.UUID `gorm:"type:uuid;not null;index" json:"period_id"`
	Title     string    `gorm:"not null;size:100" json:"title"`
	Price     float64   `gorm:"type:decimal(10,2);not null;check:price >= 0" json:"price"`
	RateType  string    `gorm:"not null;size:20;check:rate_type IN ('per_unit', 'per_booking')" json:"rate_type"`
	UnitName  string    `gorm:"not null;size:50" json:"unit_name"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Period *Period `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"period,omitempty"`
}

func (PeriodAdditionalService) TableName() string {
	return "period_additional_services"
}

func (pas *PeriodAdditionalService) CalculateTotal(quantity int) float64 {
	return pas.Price * float64(quantity)
}

func (pas *PeriodAdditionalService) IsPerUnitRate() bool {
	return pas.RateType == "per_unit"
}

func (pas *PeriodAdditionalService) IsPerBookingRate() bool {
	return pas.RateType == "per_booking"
}

func (pas *PeriodAdditionalService) IsPersonUnit() bool {
	return pas.UnitName == "คน"
}

func (pas *PeriodAdditionalService) IsDayUnit() bool {
	return pas.UnitName == "วัน"
}

func (pas *PeriodAdditionalService) IsPieceUnit() bool {
	return pas.UnitName == "ชิ้น"
}

func (pas *PeriodAdditionalService) IsPackageUnit() bool {
	return pas.UnitName == "เหมา"
}

func (pas *PeriodAdditionalService) CalculateTotalForBooking(quantity int, passengerCount int) float64 {
	switch pas.RateType {
	case "per_unit":
		if pas.UnitName == "คน" {
			return pas.Price * float64(quantity)
		}
		return pas.Price * float64(quantity)
	case "per_booking":
		return pas.Price * float64(quantity)
	default:
		return pas.CalculateTotal(quantity)
	}
}

func (pas *PeriodAdditionalService) GetUnitDisplayName() string {
	return pas.UnitName
}

func (pas *PeriodAdditionalService) GetRateDescription() string {
	switch pas.RateType {
	case "per_unit":
		return fmt.Sprintf("ราคาต่อ%s", pas.UnitName)
	case "per_booking":
		return "ราคาเหมา"
	default:
		return "ราคาพิเศษ"
	}
}