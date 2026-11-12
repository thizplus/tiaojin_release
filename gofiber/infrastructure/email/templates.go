package email

import (
	"fmt"
	"strings"
)

// BookingEmailData ข้อมูลสำหรับ email template
type BookingEmailData struct {
	InvoiceID     string
	CustomerName  string
	TourName      string
	TourCode      string
	PeriodStart   string
	PeriodEnd     string
	TotalAmount   string
	DepositAmount string
	DateDeposit   string
	DateFinal     string
	CompanyName   string
	AgentName     string
	Reason        string // สำหรับ cancel
}

func baseTemplate(title, content string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"></head>
<body style="font-family: 'Google Sans', 'Noto Sans Thai', sans-serif; background: #f5f5f5; padding: 20px; margin: 0;">
  <div style="max-width: 600px; margin: 0 auto; background: #fff; border-radius: 12px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08);">
    <div style="background: #e54b2a; padding: 24px; text-align: center;">
      <h1 style="color: #fff; margin: 0; font-size: 20px;">TIAOJIN</h1>
      <p style="color: rgba(255,255,255,0.8); margin: 4px 0 0; font-size: 12px;">ทัวร์จีนต้องเที่ยวจีน</p>
    </div>
    <div style="padding: 32px 24px;">
      <h2 style="color: #333; margin: 0 0 20px; font-size: 18px;">%s</h2>
      %s
    </div>
    <div style="background: #fafafa; padding: 16px 24px; text-align: center; font-size: 12px; color: #999;">
      <p style="margin: 0;">© TIAOJIN - ทัวร์จีนต้องเที่ยวจีน</p>
    </div>
  </div>
</body>
</html>`, title, content)
}

func infoRow(label, value string) string {
	return fmt.Sprintf(`<tr><td style="padding:6px 0;color:#666;width:140px;">%s</td><td style="padding:6px 0;color:#333;font-weight:500;">%s</td></tr>`, label, value)
}

func bookingTable(d BookingEmailData) string {
	rows := []string{
		infoRow("เลขที่ใบจอง", d.InvoiceID),
		infoRow("ชื่อลูกค้า", d.CustomerName),
		infoRow("โปรแกรมทัวร์", fmt.Sprintf("%s %s", d.TourCode, d.TourName)),
	}
	if d.PeriodStart != "" {
		rows = append(rows, infoRow("วันเดินทาง", fmt.Sprintf("%s - %s", d.PeriodStart, d.PeriodEnd)))
	}
	if d.TotalAmount != "" {
		rows = append(rows, infoRow("ยอดรวม", d.TotalAmount+" บาท"))
	}
	if d.DepositAmount != "" {
		rows = append(rows, infoRow("มัดจำ", d.DepositAmount+" บาท"))
	}
	if d.DateDeposit != "" {
		rows = append(rows, infoRow("กำหนดมัดจำ", d.DateDeposit))
	}
	if d.DateFinal != "" {
		rows = append(rows, infoRow("กำหนดชำระเต็ม", d.DateFinal))
	}

	return fmt.Sprintf(`<table style="width:100%%;border-collapse:collapse;">%s</table>`, strings.Join(rows, ""))
}

// NewBookingEmail — สร้างการจองสำเร็จ
func NewBookingEmail(d BookingEmailData) (subject, html string) {
	subject = fmt.Sprintf("การจองใหม่ #%s - %s %s", d.InvoiceID, d.TourCode, d.TourName)
	content := fmt.Sprintf(`
		<p style="color:#333;line-height:1.6;">สวัสดีคุณ %s,</p>
		<p style="color:#333;line-height:1.6;">การจองของคุณถูกสร้างเรียบร้อยแล้ว กรุณาชำระเงินภายในกำหนด</p>
		<div style="background:#f8f9fa;border-radius:8px;padding:16px;margin:16px 0;">
			%s
		</div>
		<div style="background:#fff3cd;border:1px solid #ffc107;border-radius:8px;padding:12px;margin:16px 0;">
			<p style="margin:0;color:#856404;font-size:14px;">⚠️ กรุณาชำระมัดจำภายใน <strong>%s</strong></p>
		</div>
	`, d.AgentName, bookingTable(d), d.DateDeposit)
	html = baseTemplate("🎉 สร้างการจองสำเร็จ", content)
	return
}

// PaymentConfirmedEmail — ชำระเงินครบ ยืนยันการจอง
func PaymentConfirmedEmail(d BookingEmailData) (subject, html string) {
	subject = fmt.Sprintf("ยืนยันการจอง #%s - ชำระเงินครบแล้ว", d.InvoiceID)
	content := fmt.Sprintf(`
		<p style="color:#333;line-height:1.6;">สวัสดีคุณ %s,</p>
		<p style="color:#333;line-height:1.6;">การจองของคุณได้รับการ<strong style="color:#28a745;">ยืนยันแล้ว</strong> เนื่องจากชำระเงินครบถ้วน</p>
		<div style="background:#d4edda;border:1px solid #28a745;border-radius:8px;padding:12px;margin:16px 0;text-align:center;">
			<p style="margin:0;color:#155724;font-size:16px;font-weight:bold;">✅ การจองยืนยันแล้ว</p>
		</div>
		<div style="background:#f8f9fa;border-radius:8px;padding:16px;margin:16px 0;">
			%s
		</div>
	`, d.AgentName, bookingTable(d))
	html = baseTemplate("✅ ยืนยันการจอง — ชำระเงินครบ", content)
	return
}

// BookingCancelledEmail — ยกเลิกการจอง
func BookingCancelledEmail(d BookingEmailData) (subject, html string) {
	subject = fmt.Sprintf("ยกเลิกการจอง #%s", d.InvoiceID)
	reason := d.Reason
	if reason == "" {
		reason = "ไม่ระบุ"
	}
	content := fmt.Sprintf(`
		<p style="color:#333;line-height:1.6;">สวัสดีคุณ %s,</p>
		<p style="color:#333;line-height:1.6;">การจองของคุณถูก<strong style="color:#dc3545;">ยกเลิก</strong>แล้ว</p>
		<div style="background:#f8d7da;border:1px solid #dc3545;border-radius:8px;padding:12px;margin:16px 0;">
			<p style="margin:0;color:#721c24;font-size:14px;">❌ เหตุผล: %s</p>
		</div>
		<div style="background:#f8f9fa;border-radius:8px;padding:16px;margin:16px 0;">
			%s
		</div>
	`, d.AgentName, reason, bookingTable(d))
	html = baseTemplate("❌ ยกเลิกการจอง", content)
	return
}

// BookingExpiredEmail — หมดเวลาชำระเงิน
func BookingExpiredEmail(d BookingEmailData) (subject, html string) {
	subject = fmt.Sprintf("การจอง #%s หมดอายุ — ไม่ได้ชำระเงินภายในกำหนด", d.InvoiceID)
	content := fmt.Sprintf(`
		<p style="color:#333;line-height:1.6;">สวัสดีคุณ %s,</p>
		<p style="color:#333;line-height:1.6;">การจองของคุณ<strong style="color:#dc3545;">หมดอายุ</strong>เนื่องจากไม่ได้ชำระเงินภายในเวลาที่กำหนด</p>
		<div style="background:#f8d7da;border:1px solid #dc3545;border-radius:8px;padding:12px;margin:16px 0;">
			<p style="margin:0;color:#721c24;font-size:14px;">⏰ การจองหมดอายุแล้ว กรุณาทำการจองใหม่หากต้องการ</p>
		</div>
		<div style="background:#f8f9fa;border-radius:8px;padding:16px;margin:16px 0;">
			%s
		</div>
	`, d.AgentName, bookingTable(d))
	html = baseTemplate("⏰ การจองหมดอายุ", content)
	return
}
