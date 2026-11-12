package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookingStatus string

const (
	BookingStatusPending    BookingStatus = "pending"
	BookingStatusProcessing BookingStatus = "processing"
	BookingStatusConfirm    BookingStatus = "confirm"
	BookingStatusCancel     BookingStatus = "cancel"
	BookingStatusCompleted  BookingStatus = "completed"
	BookingStatusExpired    BookingStatus = "expired"
)

type Booking struct {
	ID            uuid.UUID     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	InvoiceID     string        `gorm:"uniqueIndex;not null;size:20;column:invoice_id" json:"invoice_id"`
	UserID        uuid.UUID     `gorm:"type:uuid;not null;index" json:"user_id"`
	TourID        uuid.UUID     `gorm:"type:uuid;not null;index" json:"tour_id"`
	PeriodID      uuid.UUID     `gorm:"type:uuid;not null;index" json:"period_id"`
	SalesID       uuid.UUID     `gorm:"type:uuid;not null;index" json:"sales_id"`

	CustomerName  string  `gorm:"not null;size:200;column:customer_name" json:"customer_name"`
	CustomerPhone string  `gorm:"not null;size:20;column:customer_phone" json:"customer_phone"`
	Note          *string `gorm:"type:text" json:"note"`

	SumRoom               float64 `gorm:"type:decimal(10,2);not null;default:0;column:sum_room" json:"sum_room"`
	SumRate               float64 `gorm:"type:decimal(10,2);not null;default:0;column:sum_rate" json:"sum_rate"`
	SumAdditionalServices float64 `gorm:"type:decimal(10,2);not null;default:0;column:sum_additional_services" json:"sum_additional_services"`
	SumQty                int     `gorm:"not null;check:sum_qty > 0;column:sum_qty" json:"sum_qty"`
	SumCommissionQty      int     `gorm:"not null;default:0;column:sum_commission_qty" json:"sum_commission_qty"`
	SumNet                float64 `gorm:"type:decimal(10,2);not null;check:sum_net > 0;column:sum_net" json:"sum_net"`
	SumCommission         float64 `gorm:"type:decimal(10,2);not null;default:0;column:sum_commission" json:"sum_commission"`
	Discount              float64 `gorm:"type:decimal(10,2);not null;default:0;check:discount >= 0" json:"discount"`
	SumDeposit            float64 `gorm:"type:decimal(10,2);not null;check:sum_deposit >= 0;column:sum_deposit" json:"sum_deposit"`
	SumPayment            float64 `gorm:"type:decimal(10,2);not null;default:0;column:sum_payment" json:"sum_payment"`

	// ค่าบริการ (VAT คิดจากค่าบริการเท่านั้น)
	ServiceFeePerPax  float64 `gorm:"type:decimal(10,2);not null;default:0;column:service_fee_per_pax" json:"service_fee_per_pax"`   // ค่าบริการ/คน (รวม VAT เช่น 500)
	ServiceFeeTotal   float64 `gorm:"type:decimal(10,2);not null;default:0;column:service_fee_total" json:"service_fee_total"`       // ค่าบริการรวม = per_pax × qty

	// VAT fields (ภาษีมูลค่าเพิ่ม — คิดจากค่าบริการเท่านั้น)
	IncludeVAT bool    `gorm:"not null;default:false;column:include_vat" json:"include_vat"`
	VATRate    float64 `gorm:"type:decimal(5,2);not null;default:7.00;column:vat_rate" json:"vat_rate"`
	VATAmount  float64 `gorm:"type:decimal(10,2);not null;default:0;column:vat_amount" json:"vat_amount"`

	// Withholding Tax fields (หัก ณ ที่จ่าย)
	IncludeWHT bool    `gorm:"not null;default:false;column:include_wht" json:"include_wht"`
	WHTRate    float64 `gorm:"type:decimal(5,2);not null;default:3.00;column:wht_rate" json:"wht_rate"`
	WHTAmount  float64 `gorm:"type:decimal(10,2);not null;default:0;column:wht_amount" json:"wht_amount"`

	// Calculated totals
	NetBeforeTax float64 `gorm:"type:decimal(10,2);not null;default:0;column:net_before_tax" json:"net_before_tax"`
	GrandTotal   float64 `gorm:"type:decimal(10,2);not null;default:0;column:grand_total" json:"grand_total"`

	DateDeposit time.Time `gorm:"type:date;not null;column:date_deposit" json:"date_deposit"`
	DateFinal   time.Time `gorm:"type:date;not null;column:date_final" json:"date_final"`

	Status    BookingStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	ExpiresAt *time.Time    `gorm:"index" json:"expires_at,omitempty"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User                *User                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Tour                *Tour                 `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tour,omitempty"`
	Period              *Period               `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"period,omitempty"`
	Sales               *Sales                `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"sales,omitempty"`
	CompanyProfile      *CompanyProfile       `gorm:"-" json:"company_profile,omitempty"`
	BookingRoomTypes           []BookingRoomType           `gorm:"foreignKey:BookingID" json:"booking_room_types,omitempty"`
	BookingRateTypes           []BookingRateType           `gorm:"foreignKey:BookingID" json:"booking_rate_types,omitempty"`
	BookingAdditionalServices  []BookingAdditionalService  `gorm:"foreignKey:BookingID" json:"booking_additional_services,omitempty"`
	PaymentSlips               []PaymentSlip               `gorm:"foreignKey:BookingID" json:"payment_slips,omitempty"`
	PassengerPassports         []PassengerPassport         `gorm:"foreignKey:BookingID" json:"passenger_passports,omitempty"`
}

func (Booking) TableName() string {
	return "bookings"
}

// CalculateActualSumPayment คำนวณยอดรวมการชำระจาก PaymentSlips ที่ valid
func (b *Booking) CalculateActualSumPayment() float64 {
	var total float64
	for _, slip := range b.PaymentSlips {
		if slip.CanBeIncludedInTotal() {
			total += slip.Amount
		}
	}
	return total
}

// GetActualSumPayment ดึงยอดรวมการชำระที่แท้จริงจาก PaymentSlips
func (b *Booking) GetActualSumPayment() float64 {
	return b.CalculateActualSumPayment()
}

func (b *Booking) IsPaid() bool {
	actualPayment := b.GetActualSumPayment()
	// ใช้ GrandTotal ถ้ามี (รวม VAT - WHT) ไม่งั้นใช้ SumNet
	total := b.GrandTotal
	if total == 0 {
		total = b.SumNet
	}
	return actualPayment >= total
}

func (b *Booking) RemainingBalance() float64 {
	actualPayment := b.GetActualSumPayment()
	// ใช้ GrandTotal ถ้ามี (รวม VAT - WHT) ไม่งั้นใช้ SumNet
	total := b.GrandTotal
	if total == 0 {
		total = b.SumNet
	}
	return total - actualPayment
}

// GetPayableAmount returns the total amount customer needs to pay (GrandTotal or SumNet)
func (b *Booking) GetPayableAmount() float64 {
	if b.GrandTotal > 0 {
		return b.GrandTotal
	}
	return b.SumNet
}

// IsExpired checks if the booking has expired
func (b *Booking) IsExpired() bool {
	if b.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*b.ExpiresAt)
}

// IsNearExpiration checks if the booking will expire within the given duration
func (b *Booking) IsNearExpiration(threshold time.Duration) bool {
	if b.ExpiresAt == nil {
		return false
	}
	return time.Until(*b.ExpiresAt) <= threshold
}

// SetExpiration sets the expiration time for the booking
func (b *Booking) SetExpiration(expiration time.Time) {
	b.ExpiresAt = &expiration
}

// ClearExpiration removes the expiration time
func (b *Booking) ClearExpiration() {
	b.ExpiresAt = nil
}

// TimeUntilExpiration returns the duration until expiration
func (b *Booking) TimeUntilExpiration() *time.Duration {
	if b.ExpiresAt == nil {
		return nil
	}
	duration := time.Until(*b.ExpiresAt)
	return &duration
}

// CanExpire checks if the booking is in a state that allows expiration
func (b *Booking) CanExpire() bool {
	// Only pending bookings can expire
	// Processing bookings (with payment slips uploaded) are protected from expiry
	return b.Status == BookingStatusPending
}

// MarkAsExpired marks the booking as expired
func (b *Booking) MarkAsExpired() {
	b.Status = BookingStatusExpired
	now := time.Now()
	b.UpdatedAt = now
	if b.ExpiresAt == nil {
		b.ExpiresAt = &now
	}
}

// IsExpiredStatus checks if the booking status is expired
func (b *Booking) IsExpiredStatus() bool {
	return b.Status == BookingStatusExpired
}