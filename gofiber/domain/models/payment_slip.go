package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentSlipStatus string

const (
	PaymentSlipStatusPending  PaymentSlipStatus = "pending"
	PaymentSlipStatusConfirm  PaymentSlipStatus = "confirm"
	PaymentSlipStatusReject   PaymentSlipStatus = "reject"
	PaymentSlipStatusReversed PaymentSlipStatus = "reversed"
)

type PaymentMethod string

const (
	PaymentMethodSlip        PaymentMethod = "slip"          // Traditional slip upload
	PaymentMethodCash        PaymentMethod = "cash"          // เงินสด
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer" // โอนธนาคาร
	PaymentMethodCreditCard  PaymentMethod = "credit_card"   // บัตรเครดิต
	PaymentMethodAdjustment  PaymentMethod = "adjustment"    // ปรับปรุงยอด
)

type PaymentSlip struct {
	ID               uuid.UUID         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	BookingID        uuid.UUID         `gorm:"type:uuid;not null;index" json:"booking_id"`
	SlipID           string            `gorm:"uniqueIndex;not null;size:50;column:slip_id" json:"slip_id"`
	Image            *string           `gorm:"size:500" json:"image"`                                                      // Optional for admin direct payments
	Amount           float64           `gorm:"type:decimal(10,2);not null;check:amount > 0" json:"amount"`
	OriginalAmount   *float64          `gorm:"type:decimal(10,2);column:original_amount" json:"original_amount"`
	AdjustmentReason *string           `gorm:"type:text;column:adjustment_reason" json:"adjustment_reason"`
	Status           PaymentSlipStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Note             *string           `gorm:"type:text;column:note" json:"note"`

	// Payment method and admin creation tracking
	PaymentMethod    PaymentMethod `gorm:"type:varchar(20);not null;default:'slip'" json:"payment_method"`
	IsAdminCreated   bool          `gorm:"not null;default:false;column:is_admin_created" json:"is_admin_created"`
	CreatedByAdmin   *uuid.UUID    `gorm:"type:uuid;column:created_by_admin" json:"created_by_admin,omitempty"`
	AdminCreatedNote *string       `gorm:"type:text;column:admin_created_note" json:"admin_created_note,omitempty"`

	// SlipOK verification fields
	PaymentDate     *time.Time `gorm:"column:payment_date" json:"payment_date"`                          // วันเวลาโอนจริง (จาก SlipOK transTimestamp)
	SendingBank     *string    `gorm:"size:100;column:sending_bank" json:"sending_bank"`                 // ธนาคารผู้โอน (ชื่อย่อ เช่น KBANK)
	ReceivingBank   *string    `gorm:"size:100;column:receiving_bank" json:"receiving_bank"`              // ธนาคารผู้รับ
	ReceiverAccount *string    `gorm:"size:50;column:receiver_account" json:"receiver_account"`           // เลขบัญชีผู้รับ (masked)
	ReceiverName    *string    `gorm:"size:100;column:receiver_name" json:"receiver_name"`                // ชื่อผู้รับ
	SenderName      *string    `gorm:"size:100;column:sender_name" json:"sender_name"`                   // ชื่อผู้โอน
	TransRef        *string    `gorm:"size:50;column:trans_ref" json:"trans_ref"`                         // เลขอ้างอิงธุรกรรม
	SlipVerified    bool       `gorm:"not null;default:false;column:slip_verified" json:"slip_verified"`  // ผ่าน SlipOK หรือยัง
	SlipVerifyError *string    `gorm:"size:255;column:slip_verify_error" json:"slip_verify_error"`        // error message จาก SlipOK (ถ้าล้มเหลว)

	// Reversal tracking fields
	IsReversed     bool       `gorm:"not null;default:false;column:is_reversed" json:"is_reversed"`
	ReversedAt     *time.Time `gorm:"column:reversed_at" json:"reversed_at,omitempty"`
	ReversedBy     *uuid.UUID `gorm:"type:uuid;column:reversed_by" json:"reversed_by,omitempty"`
	ReversalReason *string    `gorm:"type:text;column:reversal_reason" json:"reversal_reason,omitempty"`

	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Booking        *Booking `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"booking,omitempty"`
	ReversedByUser *User    `gorm:"foreignKey:ReversedBy;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;" json:"reversed_by_user,omitempty"`
	CreatedByUser  *User    `gorm:"foreignKey:CreatedByAdmin;constraint:OnUpdate:SET NULL,OnDelete:SET NULL;" json:"created_by_user,omitempty"`
}

func (PaymentSlip) TableName() string {
	return "payment_slips"
}

func (ps *PaymentSlip) IsVerified() bool {
	return ps.Status == PaymentSlipStatusConfirm && !ps.IsReversed
}

func (ps *PaymentSlip) IsRejected() bool {
	return ps.Status == PaymentSlipStatusReject
}

func (ps *PaymentSlip) IsPending() bool {
	return ps.Status == PaymentSlipStatusPending
}

func (ps *PaymentSlip) IsReversedStatus() bool {
	return ps.Status == PaymentSlipStatusReversed || ps.IsReversed
}

// Reversal helper methods
func (ps *PaymentSlip) IsReversible() bool {
	return ps.Status == PaymentSlipStatusConfirm && !ps.IsReversed
}

func (ps *PaymentSlip) MarkAsReversed(reversedBy uuid.UUID, reason string) {
	now := time.Now()
	ps.Status = PaymentSlipStatusReversed
	ps.IsReversed = true
	ps.ReversedAt = &now
	ps.ReversedBy = &reversedBy
	ps.ReversalReason = &reason
	ps.UpdatedAt = now
}

func (ps *PaymentSlip) CanBeIncludedInTotal() bool {
	return ps.Status == PaymentSlipStatusConfirm && !ps.IsReversed
}

// Admin Direct Payment helper methods
func (ps *PaymentSlip) IsAdminDirectPayment() bool {
	return ps.IsAdminCreated && ps.PaymentMethod != PaymentMethodSlip
}

func (ps *PaymentSlip) RequiresSlipImage() bool {
	return ps.PaymentMethod == PaymentMethodSlip
}

// PaymentMethod validation
func IsValidPaymentMethod(method string) bool {
	switch PaymentMethod(method) {
	case PaymentMethodSlip, PaymentMethodCash, PaymentMethodBankTransfer, PaymentMethodCreditCard, PaymentMethodAdjustment:
		return true
	default:
		return false
	}
}