package ports

import "context"

// EmailPort — ส่งอีเมล, เปลี่ยน provider ได้ (SMTP, SendGrid, SES, etc.)
type EmailPort interface {
	// SendEmail ส่งอีเมลพร้อม HTML body
	SendEmail(ctx context.Context, req EmailRequest) error
	// IsConfigured เช็คว่า config ครบไหม
	IsConfigured() bool
}

// EmailRequest ข้อมูลสำหรับส่งอีเมล
type EmailRequest struct {
	To      []string // ผู้รับ (อาจหลายคน)
	Subject string   // หัวข้อ
	HTML    string   // เนื้อหา HTML
}
