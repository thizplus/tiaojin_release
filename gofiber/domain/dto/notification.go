package dto

import "github.com/google/uuid"

type NotificationType string

const (
	NotificationTypeNewBooking       NotificationType = "new_booking"
	NotificationTypeBookingCancelled NotificationType = "booking_cancelled"
	NotificationTypePaymentReceived  NotificationType = "payment_received"
	NotificationTypePaymentVerified  NotificationType = "payment_verified"
	NotificationTypePassportUploaded NotificationType = "passport_uploaded"
	NotificationTypeUserRegistered   NotificationType = "user_registered"
	NotificationTypePeriodFull       NotificationType = "period_full"
	NotificationTypeSystemAlert      NotificationType = "system_alert"
)

type NotificationPriority string

const (
	NotificationPriorityLow    NotificationPriority = "low"
	NotificationPriorityNormal NotificationPriority = "normal"
	NotificationPriorityHigh   NotificationPriority = "high"
	NotificationPriorityUrgent NotificationPriority = "urgent"
)

type TelegramNotificationRequest struct {
	Type     NotificationType       `json:"type" validate:"required"`
	Priority NotificationPriority   `json:"priority" validate:"required"`
	Title    string                 `json:"title" validate:"required,max=200"`
	Message  string                 `json:"message" validate:"required,max=4000"`
	Data     map[string]interface{} `json:"data,omitempty"`
	ChatID   *string                `json:"chat_id,omitempty"`
}

type BookingNotificationData struct {
	BookingID    uuid.UUID `json:"booking_id"`
	InvoiceID    string    `json:"invoice_id"`
	CustomerName string    `json:"customer_name"`
	TourName     string    `json:"tour_name"`
	PeriodStart  string    `json:"period_start"`
	Amount       float64   `json:"amount"`
	SalesName    string    `json:"sales_name"`
	CompanyName  string    `json:"company_name"`
}

type PaymentNotificationData struct {
	BookingID    uuid.UUID `json:"booking_id"`
	InvoiceID    string    `json:"invoice_id"`
	SlipID       string    `json:"slip_id"`
	Amount       float64   `json:"amount"`
	CustomerName string    `json:"customer_name"`
	CompanyName  string    `json:"company_name"`
	PaymentImage string    `json:"payment_image,omitempty"`
}

type UserNotificationData struct {
	UserID      uuid.UUID `json:"user_id"`
	Email       string    `json:"email"`
	CompanyName string    `json:"company_name"`
	Status      string    `json:"status"`
}

type PeriodNotificationData struct {
	PeriodID    uuid.UUID `json:"period_id"`
	TourName    string    `json:"tour_name"`
	TourCode    string    `json:"tour_code"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
	TotalSeats  int       `json:"total_seats"`
	BookedSeats int       `json:"booked_seats"`
}

type EmailNotificationRequest struct {
	To       []string               `json:"to" validate:"required,min=1"`
	CC       []string               `json:"cc,omitempty"`
	BCC      []string               `json:"bcc,omitempty"`
	Subject  string                 `json:"subject" validate:"required,max=200"`
	Template string                 `json:"template" validate:"required"`
	Data     map[string]interface{} `json:"data"`
	Priority NotificationPriority   `json:"priority" validate:"required"`
}

type SMSNotificationRequest struct {
	Phone    string                 `json:"phone" validate:"required"`
	Message  string                 `json:"message" validate:"required,max=160"`
	Data     map[string]interface{} `json:"data,omitempty"`
	Priority NotificationPriority   `json:"priority" validate:"required"`
}

type NotificationResponse struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Status  string  `json:"status"`
	Message string  `json:"message,omitempty"`
	SentAt  string  `json:"sent_at"`
	Error   *string `json:"error,omitempty"`
}

type NotificationSettings struct {
	TelegramEnabled   bool     `json:"telegram_enabled"`
	TelegramChatID    string   `json:"telegram_chat_id"`
	TelegramToken     string   `json:"telegram_token"`
	EmailEnabled      bool     `json:"email_enabled"`
	SMSEnabled        bool     `json:"sms_enabled"`
	NotificationTypes []string `json:"notification_types"`
}