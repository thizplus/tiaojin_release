package services

import (
	"context"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
)


type NotificationService interface {
	// Telegram Notifications
	SendTelegramNotification(ctx context.Context, req dto.TelegramNotificationRequest) (*dto.NotificationResponse, error)
	SendBookingNotification(ctx context.Context, notifyType dto.NotificationType, data dto.BookingNotificationData) (*dto.NotificationResponse, error)
	SendPaymentNotification(ctx context.Context, notifyType dto.NotificationType, data dto.PaymentNotificationData) (*dto.NotificationResponse, error)
	SendUserNotification(ctx context.Context, notifyType dto.NotificationType, data dto.UserNotificationData) (*dto.NotificationResponse, error)
	SendPeriodNotification(ctx context.Context, notifyType dto.NotificationType, data dto.PeriodNotificationData) (*dto.NotificationResponse, error)

	// Email Notifications
	SendEmailNotification(ctx context.Context, req dto.EmailNotificationRequest) (*dto.NotificationResponse, error)
	SendBookingConfirmationEmail(ctx context.Context, bookingID uuid.UUID, userEmail string) (*dto.NotificationResponse, error)
	SendPaymentConfirmationEmail(ctx context.Context, bookingID uuid.UUID, userEmail string) (*dto.NotificationResponse, error)
	SendWelcomeEmail(ctx context.Context, userEmail, userName string) (*dto.NotificationResponse, error)
	SendPasswordResetEmail(ctx context.Context, userEmail, resetToken string) (*dto.NotificationResponse, error)

	// SMS Notifications
	SendSMSNotification(ctx context.Context, req dto.SMSNotificationRequest) (*dto.NotificationResponse, error)
	SendBookingConfirmationSMS(ctx context.Context, phone string, bookingInfo dto.BookingNotificationData) (*dto.NotificationResponse, error)

	// System Alerts
	SendSystemAlert(ctx context.Context, title, message string, priority dto.NotificationPriority) (*dto.NotificationResponse, error)
	SendErrorAlert(ctx context.Context, service, errorMsg string, data map[string]interface{}) (*dto.NotificationResponse, error)

	// Template Management
	FormatBookingMessage(data dto.BookingNotificationData, notifyType dto.NotificationType) string
	FormatPaymentMessage(data dto.PaymentNotificationData, notifyType dto.NotificationType) string
	FormatUserMessage(data dto.UserNotificationData, notifyType dto.NotificationType) string
	FormatPeriodMessage(data dto.PeriodNotificationData, notifyType dto.NotificationType) string

	// Settings Management
	GetNotificationSettings(ctx context.Context) (*dto.NotificationSettings, error)
	UpdateNotificationSettings(ctx context.Context, settings dto.NotificationSettings) error
	TestTelegramConnection(ctx context.Context) error
	TestEmailConnection(ctx context.Context) error

	// Batch Operations
	SendBulkNotifications(ctx context.Context, notifications []dto.TelegramNotificationRequest) ([]*dto.NotificationResponse, error)
	SendBulkEmails(ctx context.Context, emails []dto.EmailNotificationRequest) ([]*dto.NotificationResponse, error)
}