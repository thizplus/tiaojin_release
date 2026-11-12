package ports

import (
	"context"
	"time"
)

// SlipVerificationPort ตรวจสอบสลิปโอนเงิน — เปลี่ยน provider ได้ (SlipOK, Vision AI, etc.)
type SlipVerificationPort interface {
	// VerifySlip ส่งรูปสลิปไปตรวจสอบ, amount เป็น optional สำหรับเช็คยอด
	VerifySlip(ctx context.Context, imageData []byte, amount *float64) (*SlipVerificationResult, error)
	// CheckQuota เช็คโควต้าคงเหลือ
	CheckQuota(ctx context.Context) (*SlipQuotaResult, error)
}

// SlipVerificationResult ผลการตรวจสอบสลิป (mapping จาก SlipOK API response)
type SlipVerificationResult struct {
	Success bool   // สลิปถูกต้องหรือไม่
	Message string // ข้อความ (เช่น "✅" หรือ error message)

	// ข้อมูลธุรกรรม
	TransRef       string    // เลขอ้างอิง
	TransDate      string    // วันที่ (yyyyMMdd)
	TransTime      string    // เวลา (HH:mm:ss)
	TransTimestamp time.Time // Timestamp เต็ม (ISO 8601)
	Amount         float64   // จำนวนเงิน

	// ธนาคาร (รหัส 3 หลัก เช่น "004" = KBANK)
	SendingBank   string // รหัสธนาคารผู้โอน
	ReceivingBank string // รหัสธนาคารผู้รับ

	// ผู้โอน
	SenderName        string // ชื่อผู้โอน (อาจถูกตัดทอน)
	SenderDisplayName string // ชื่อแสดงผู้โอน
	SenderAccount     string // เลขบัญชีผู้โอน (masked)

	// ผู้รับ
	ReceiverName        string // ชื่อผู้รับ (อาจถูกตัดทอน)
	ReceiverDisplayName string // ชื่อแสดงผู้รับ
	ReceiverAccount     string // เลขบัญชีผู้รับ (masked)

	// ข้อมูลเพิ่มเติม
	CountryCode    string  // รหัสประเทศ (TH)
	TransFeeAmount float64 // ค่าธรรมเนียม
	Ref1           string  // Reference 1
	Ref2           string  // Reference 2
	Ref3           string  // Reference 3

	// Error info (กรณีล้มเหลว)
	ErrorCode int    // SlipOK error code (1000-1015)
	RawJSON   string // response ดิบจาก provider
}

// SlipQuotaResult ข้อมูลโควต้า
type SlipQuotaResult struct {
	Quota        int    // โควต้าคงเหลือ
	OverQuota    int    // โควต้าที่ใช้เกิน
	SpecialQuota int    // โควต้าพิเศษ
	EndDate      string // วันหมดอายุ
}

// BankCodeToName แปลงรหัสธนาคาร → ชื่อย่อ
var BankCodeToName = map[string]string{
	"002": "BBL",
	"004": "KBANK",
	"006": "KTB",
	"011": "TTB",
	"014": "SCB",
	"025": "BAY",
	"069": "KKP",
	"022": "CIMBT",
	"067": "TISCO",
	"024": "UOBT",
	"071": "TCD",
	"073": "LHFG",
	"070": "ICBCT",
	"098": "SME",
	"034": "BAAC",
	"035": "EXIM",
	"030": "GSB",
	"033": "GHB",
}

// BankCodeToFullName แปลงรหัสธนาคาร → ชื่อเต็ม
var BankCodeToFullName = map[string]string{
	"002": "ธนาคารกรุงเทพ",
	"004": "ธนาคารกสิกรไทย",
	"006": "ธนาคารกรุงไทย",
	"011": "ธนาคารทหารไทยธนชาต",
	"014": "ธนาคารไทยพาณิชย์",
	"025": "ธนาคารกรุงศรีอยุธยา",
	"069": "ธนาคารเกียรตินาคินภัทร",
	"022": "ธนาคารซีไอเอ็มบีไทย",
	"067": "ธนาคารทิสโก้",
	"024": "ธนาคารยูโอบี",
	"071": "ธนาคารไทยเครดิตเพื่อรายย่อย",
	"073": "ธนาคารแลนด์ แอนด์ เฮ้าส์",
	"070": "ธนาคารไอซีบีซี (ไทย)",
	"098": "ธนาคารพัฒนาวิสาหกิจฯ",
	"034": "ธนาคารเพื่อการเกษตรฯ",
	"035": "ธนาคารเพื่อการส่งออกฯ",
	"030": "ธนาคารออมสิน",
	"033": "ธนาคารอาคารสงเคราะห์",
}
