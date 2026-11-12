package email

import (
	"context"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	"gofiber-template/domain/ports"
)

// SMTPAdapter — implement EmailPort ด้วย SMTP
type SMTPAdapter struct {
	host     string
	port     string
	user     string
	password string
	from     string
	fromName string
}

// NewSMTPAdapter สร้าง adapter จาก env
func NewSMTPAdapter() *SMTPAdapter {
	fromName := os.Getenv("SMTP_FROM_NAME")
	if fromName == "" {
		fromName = "TIAOJIN"
	}

	return &SMTPAdapter{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		user:     os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASS"),
		from:     os.Getenv("SMTP_FROM"),
		fromName: fromName,
	}
}

// IsConfigured เช็คว่า config ครบไหม
func (a *SMTPAdapter) IsConfigured() bool {
	return a.host != "" && a.port != "" && a.user != "" && a.password != "" && a.from != ""
}

// SendEmail ส่งอีเมลผ่าน SMTP
func (a *SMTPAdapter) SendEmail(ctx context.Context, req ports.EmailRequest) error {
	if !a.IsConfigured() {
		return fmt.Errorf("SMTP not configured")
	}

	if len(req.To) == 0 {
		return fmt.Errorf("no recipients")
	}

	// สร้าง email headers + body
	fromHeader := fmt.Sprintf("%s <%s>", a.fromName, a.from)
	toHeader := strings.Join(req.To, ", ")

	msg := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s",
		fromHeader, toHeader, req.Subject, req.HTML)

	// SMTP auth
	addr := fmt.Sprintf("%s:%s", a.host, a.port)
	auth := smtp.PlainAuth("", a.user, a.password, a.host)

	// ส่ง
	err := smtp.SendMail(addr, auth, a.from, req.To, []byte(msg))
	if err != nil {
		log.Printf("[EMAIL] Failed to send email to %s: %v", toHeader, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("[EMAIL] Sent to %s subject=%q", toHeader, req.Subject)
	return nil
}
