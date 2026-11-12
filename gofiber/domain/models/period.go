package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PeriodStatus string

const (
	PeriodStatusPublish PeriodStatus = "publish"
	PeriodStatusDraft   PeriodStatus = "draft"
)

type RateType string

const (
	RateTypeRegular RateType = "regular"
	RateTypeSale    RateType = "sale"
)

type FlightSector struct {
	Date   string `json:"date"`
	Flight string `json:"flight"`
	Route  string `json:"route"`
}

type FlightSectors []FlightSector

func (fs FlightSectors) Value() (driver.Value, error) {
	return json.Marshal(fs)
}

func (fs *FlightSectors) Scan(value interface{}) error {
	if value == nil {
		*fs = FlightSectors{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, fs)
}

type RateItem struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Regular float64 `json:"regular"`
	Sale    float64 `json:"sale"`
	Exclude bool    `json:"exclude,omitempty"`
}

type RateItems []RateItem

func (ri RateItems) Value() (driver.Value, error) {
	return json.Marshal(ri)
}

func (ri *RateItems) Scan(value interface{}) error {
	if value == nil {
		*ri = RateItems{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, ri)
}

type CommissionItem struct {
	Regular float64 `json:"regular"`
	Sale    float64 `json:"sale"`
}

type CommissionItems []CommissionItem

func (ci CommissionItems) Value() (driver.Value, error) {
	return json.Marshal(ci)
}

func (ci *CommissionItems) Scan(value interface{}) error {
	if value == nil {
		*ci = CommissionItems{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, ci)
}

type FileAttachment struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type FileAttachments []FileAttachment

func (fa FileAttachments) Value() (driver.Value, error) {
	return json.Marshal(fa)
}

func (fa *FileAttachments) Scan(value interface{}) error {
	if value == nil {
		*fa = FileAttachments{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, fa)
}

type RoomType struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

type RoomTypes []RoomType

func (rt RoomTypes) Value() (driver.Value, error) {
	return json.Marshal(rt)
}

func (rt *RoomTypes) Scan(value interface{}) error {
	if value == nil {
		*rt = RoomTypes{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, rt)
}

type AdditionalService struct {
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	RateType string  `json:"rate_type"`
	UnitName string  `json:"unit_name"`
}

type AdditionalServices []AdditionalService

func (as AdditionalServices) Value() (driver.Value, error) {
	return json.Marshal(as)
}

func (as *AdditionalServices) Scan(value interface{}) error {
	if value == nil {
		*as = AdditionalServices{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, as)
}

type Period struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	TourID      uuid.UUID `gorm:"type:uuid;not null;index" json:"tour_id"`
	StartDate   time.Time `gorm:"type:date;not null;column:start_date" json:"start_date"`
	EndDate     time.Time `gorm:"type:date;not null;column:end_date" json:"end_date"`
	Seat        int       `gorm:"not null;check:seat > 0 AND seat <= 999" json:"seat"`
	Deposit     float64   `gorm:"type:decimal(10,2);not null;check:deposit >= 0" json:"deposit"`
	DepositType string    `gorm:"type:varchar(20);not null;default:'percent';column:deposit_type" json:"deposit_type"`
	//SingleStay         float64            `gorm:"type:decimal(10,2);not null;default:0;column:single_stay" json:"single_stay"`
	Airline            string             `gorm:"not null;size:100" json:"airline"`
	Airport            string             `gorm:"not null;size:10" json:"airport"`
	Note               *string            `gorm:"type:text" json:"note"`
	Sector1            FlightSectors      `gorm:"type:jsonb;column:sector_1" json:"sector_1"`
	Sector2            FlightSectors      `gorm:"type:jsonb;column:sector_2" json:"sector_2"`
	RateType           RateType           `gorm:"type:varchar(20);not null;column:rate_type" json:"rate_type"`
	RateItems          RateItems          `gorm:"type:jsonb;column:rate_items" json:"rate_items"`
	CommissionType     RateType           `gorm:"type:varchar(20);not null;column:commission_type" json:"commission_type"`
	CommissionItems    CommissionItems    `gorm:"type:jsonb;column:commission_items" json:"commission_items"`
	RoomTypes          RoomTypes          `gorm:"type:jsonb;column:room_types" json:"room_types"`
	AdditionalServices AdditionalServices `gorm:"type:jsonb;column:additional_services" json:"additional_services"`
	Files              FileAttachments    `gorm:"type:jsonb" json:"files"`
	Status             PeriodStatus       `gorm:"type:varchar(20);not null;default:'draft'" json:"status"`

	// ค่าบริการ/คน (บาท, รวม VAT) — ใช้คำนวณ VAT ตอนสร้าง Booking
	ServiceFeePerPax float64 `gorm:"type:decimal(10,2);not null;default:0;column:service_fee_per_pax" json:"service_fee_per_pax"`

	// Expiration Override - สำหรับทัวร์ไฟไหม้หรือ period พิเศษ
	ExpirationOverride bool    `gorm:"default:false;column:expiration_override" json:"expiration_override"`
	ExpirationMinutes  *int    `gorm:"column:expiration_minutes" json:"expiration_minutes,omitempty"`
	ExpirationHours    *int    `gorm:"column:expiration_hours" json:"expiration_hours,omitempty"`
	ExpirationDays     *int    `gorm:"column:expiration_days" json:"expiration_days,omitempty"`
	ExpirationMessage  *string `gorm:"type:varchar(255);column:expiration_message" json:"expiration_message,omitempty"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Tour     *Tour     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tour,omitempty"`
	Bookings []Booking `gorm:"foreignKey:PeriodID" json:"bookings,omitempty"`

	// New stable reference tables (will replace JSONB arrays)
	PeriodRateItems          []PeriodRateItem          `gorm:"foreignKey:PeriodID" json:"period_rate_items,omitempty"`
	PeriodRoomTypes          []PeriodRoomType          `gorm:"foreignKey:PeriodID" json:"period_room_types,omitempty"`
	PeriodAdditionalServices []PeriodAdditionalService `gorm:"foreignKey:PeriodID" json:"period_additional_services,omitempty"`
}

func (Period) TableName() string {
	return "periods"
}

func (p *Period) IsPublished() bool {
	return p.Status == PeriodStatusPublish
}

// HasExpirationOverride checks if period has custom expiration settings
func (p *Period) HasExpirationOverride() bool {
	return p.ExpirationOverride && (p.ExpirationMinutes != nil || p.ExpirationHours != nil || p.ExpirationDays != nil)
}

// CalculateOverrideExpiration calculates expiration time from period override settings
func (p *Period) CalculateOverrideExpiration(createdAt time.Time) *time.Time {
	if !p.HasExpirationOverride() {
		return nil
	}

	if p.ExpirationMinutes != nil && *p.ExpirationMinutes > 0 {
		t := createdAt.Add(time.Duration(*p.ExpirationMinutes) * time.Minute)
		return &t
	}
	if p.ExpirationHours != nil && *p.ExpirationHours > 0 {
		t := createdAt.Add(time.Duration(*p.ExpirationHours) * time.Hour)
		return &t
	}
	if p.ExpirationDays != nil && *p.ExpirationDays > 0 {
		t := createdAt.Add(time.Duration(*p.ExpirationDays) * 24 * time.Hour)
		return &t
	}

	return nil
}

// GetExpirationDisplay returns human-readable expiration duration
func (p *Period) GetExpirationDisplay() string {
	if !p.HasExpirationOverride() {
		return ""
	}

	if p.ExpirationMinutes != nil && *p.ExpirationMinutes > 0 {
		return fmt.Sprintf("%d นาที", *p.ExpirationMinutes)
	}
	if p.ExpirationHours != nil && *p.ExpirationHours > 0 {
		return fmt.Sprintf("%d ชั่วโมง", *p.ExpirationHours)
	}
	if p.ExpirationDays != nil && *p.ExpirationDays > 0 {
		return fmt.Sprintf("%d วัน", *p.ExpirationDays)
	}

	return ""
}
