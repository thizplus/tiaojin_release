package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodRateItem struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	PeriodID  uuid.UUID `gorm:"type:uuid;not null;index" json:"period_id"`
	Title     string    `gorm:"not null;size:100" json:"title"`
	Regular   float64   `gorm:"type:decimal(10,2);not null;check:regular >= 0" json:"regular"`
	Sale      float64   `gorm:"type:decimal(10,2);not null;check:sale >= 0" json:"sale"`
	Exclude   bool      `gorm:"not null;default:false" json:"exclude"`
	SortOrder int       `gorm:"not null;default:0" json:"sort_order"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Period *Period `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"period,omitempty"`
}

func (PeriodRateItem) TableName() string {
	return "period_rate_items"
}

func (pri *PeriodRateItem) GetPrice(rateType RateType) float64 {
	switch rateType {
	case RateTypeRegular:
		return pri.Regular
	case RateTypeSale:
		return pri.Sale
	default:
		return pri.Regular
	}
}

func (pri *PeriodRateItem) IsIncludedInPackage() bool {
	return !pri.Exclude
}

func (pri *PeriodRateItem) IsExcludedFromPackage() bool {
	return pri.Exclude
}