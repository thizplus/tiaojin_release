package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompanyBankAccount struct {
	ID            uuid.UUID      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	BankCode      string         `gorm:"not null;size:10;column:bank_code" json:"bank_code"`           // รหัสธนาคาร (เช่น 004, 014)
	BankName      string         `gorm:"not null;size:100;column:bank_name" json:"bank_name"`          // ชื่อธนาคาร (เช่น KBANK, SCB)
	AccountNumber string         `gorm:"not null;size:20;column:account_number" json:"account_number"` // เลขบัญชี
	AccountName   string         `gorm:"not null;size:200;column:account_name" json:"account_name"`    // ชื่อบัญชี
	BranchName    *string        `gorm:"size:200;column:branch_name" json:"branch_name"`               // ชื่อสาขา (optional)
	IsActive      bool           `gorm:"not null;default:true;column:is_active" json:"is_active"`      // เปิด/ปิดใช้งาน
	SortOrder     int            `gorm:"not null;default:0;column:sort_order" json:"sort_order"`        // ลำดับการแสดง
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CompanyBankAccount) TableName() string {
	return "company_bank_accounts"
}
