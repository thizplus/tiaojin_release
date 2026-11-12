package serviceimpl

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	apperrors "gofiber-template/pkg/errors"
)


type NotificationServiceImpl struct {
	settings *dto.NotificationSettings
}

func NewNotificationService() services.NotificationService {
	// Initialize with default settings
	settings := &dto.NotificationSettings{
		TelegramEnabled:   false, // Default to disabled
		EmailEnabled:      false,
		SMSEnabled:        false,
		NotificationTypes: []string{},
	}

	return &NotificationServiceImpl{
		settings: settings,
	}
}

// Telegram Notifications
func (s *NotificationServiceImpl) SendTelegramNotification(ctx context.Context, req dto.TelegramNotificationRequest) (*dto.NotificationResponse, error) {
	// Business validation - check if telegram is enabled
	if !s.settings.TelegramEnabled {
		return nil, apperrors.ErrNotificationServiceDisabled
	}

	// Business validation - check configuration
	if s.settings.TelegramToken == "" {
		return nil, apperrors.ErrInvalidTelegramConfig
	}

	// Generate notification ID
	notificationID := s.generateNotificationID()

	// TODO: Implement actual Telegram API call
	// For now, just log the notification
	log.Printf("Sending Telegram notification: %s - %s", req.Title, req.Message)

	// Simulate sending
	response := &dto.NotificationResponse{
		ID:      notificationID,
		Type:    "telegram",
		Status:  "sent",
		Message: fmt.Sprintf("Notification sent to Telegram: %s", req.Title),
		SentAt:  time.Now().Format(time.RFC3339),
	}

	return response, nil
}

func (s *NotificationServiceImpl) SendBookingNotification(ctx context.Context, notifyType dto.NotificationType, data dto.BookingNotificationData) (*dto.NotificationResponse, error) {
	// Format message based on type
	message := s.FormatBookingMessage(data, notifyType)
	title := s.getBookingNotificationTitle(notifyType)

	// Create telegram request
	req := dto.TelegramNotificationRequest{
		Type:     notifyType,
		Priority: dto.NotificationPriorityNormal,
		Title:    title,
		Message:  message,
		Data: map[string]interface{}{
			"booking_id":    data.BookingID,
			"invoice_id":    data.InvoiceID,
			"customer_name": data.CustomerName,
			"tour_name":     data.TourName,
			"amount":        data.Amount,
		},
	}

	return s.SendTelegramNotification(ctx, req)
}

func (s *NotificationServiceImpl) SendPaymentNotification(ctx context.Context, notifyType dto.NotificationType, data dto.PaymentNotificationData) (*dto.NotificationResponse, error) {
	// Format message based on type
	message := s.FormatPaymentMessage(data, notifyType)
	title := s.getPaymentNotificationTitle(notifyType)

	// Create telegram request
	req := dto.TelegramNotificationRequest{
		Type:     notifyType,
		Priority: dto.NotificationPriorityHigh,
		Title:    title,
		Message:  message,
		Data: map[string]interface{}{
			"booking_id":    data.BookingID,
			"invoice_id":    data.InvoiceID,
			"slip_id":       data.SlipID,
			"amount":        data.Amount,
			"customer_name": data.CustomerName,
		},
	}

	return s.SendTelegramNotification(ctx, req)
}

func (s *NotificationServiceImpl) SendUserNotification(ctx context.Context, notifyType dto.NotificationType, data dto.UserNotificationData) (*dto.NotificationResponse, error) {
	// Format message based on type
	message := s.FormatUserMessage(data, notifyType)
	title := s.getUserNotificationTitle(notifyType)

	// Create telegram request
	req := dto.TelegramNotificationRequest{
		Type:     notifyType,
		Priority: dto.NotificationPriorityNormal,
		Title:    title,
		Message:  message,
		Data: map[string]interface{}{
			"user_id":      data.UserID,
			"email":        data.Email,
			"company_name": data.CompanyName,
			"status":       data.Status,
		},
	}

	return s.SendTelegramNotification(ctx, req)
}

func (s *NotificationServiceImpl) SendPeriodNotification(ctx context.Context, notifyType dto.NotificationType, data dto.PeriodNotificationData) (*dto.NotificationResponse, error) {
	// Format message based on type
	message := s.FormatPeriodMessage(data, notifyType)
	title := s.getPeriodNotificationTitle(notifyType)

	// Create telegram request
	req := dto.TelegramNotificationRequest{
		Type:     notifyType,
		Priority: dto.NotificationPriorityHigh,
		Title:    title,
		Message:  message,
		Data: map[string]interface{}{
			"period_id":    data.PeriodID,
			"tour_name":    data.TourName,
			"tour_code":    data.TourCode,
			"start_date":   data.StartDate,
			"total_seats":  data.TotalSeats,
			"booked_seats": data.BookedSeats,
		},
	}

	return s.SendTelegramNotification(ctx, req)
}

// Email Notifications
func (s *NotificationServiceImpl) SendEmailNotification(ctx context.Context, req dto.EmailNotificationRequest) (*dto.NotificationResponse, error) {
	// Business validation - check if email is enabled
	if !s.settings.EmailEnabled {
		return nil, apperrors.ErrNotificationServiceDisabled
	}

	// Generate notification ID
	notificationID := s.generateNotificationID()

	// TODO: Implement actual email sending
	// For now, just log the notification
	log.Printf("Sending email to %v: %s", req.To, req.Subject)

	// Simulate sending
	response := &dto.NotificationResponse{
		ID:      notificationID,
		Type:    "email",
		Status:  "sent",
		Message: fmt.Sprintf("Email sent to %v: %s", req.To, req.Subject),
		SentAt:  time.Now().Format(time.RFC3339),
	}

	return response, nil
}

func (s *NotificationServiceImpl) SendBookingConfirmationEmail(ctx context.Context, bookingID uuid.UUID, userEmail string) (*dto.NotificationResponse, error) {
	req := dto.EmailNotificationRequest{
		To:       []string{userEmail},
		Subject:  "Booking Confirmation",
		Template: "booking_confirmation",
		Data: map[string]interface{}{
			"booking_id": bookingID,
		},
		Priority: dto.NotificationPriorityHigh,
	}

	return s.SendEmailNotification(ctx, req)
}

func (s *NotificationServiceImpl) SendPaymentConfirmationEmail(ctx context.Context, bookingID uuid.UUID, userEmail string) (*dto.NotificationResponse, error) {
	req := dto.EmailNotificationRequest{
		To:       []string{userEmail},
		Subject:  "Payment Confirmation",
		Template: "payment_confirmation",
		Data: map[string]interface{}{
			"booking_id": bookingID,
		},
		Priority: dto.NotificationPriorityHigh,
	}

	return s.SendEmailNotification(ctx, req)
}

func (s *NotificationServiceImpl) SendWelcomeEmail(ctx context.Context, userEmail, userName string) (*dto.NotificationResponse, error) {
	req := dto.EmailNotificationRequest{
		To:       []string{userEmail},
		Subject:  "Welcome to Tiaojin Travel",
		Template: "welcome_email",
		Data: map[string]interface{}{
			"user_name": userName,
		},
		Priority: dto.NotificationPriorityNormal,
	}

	return s.SendEmailNotification(ctx, req)
}

func (s *NotificationServiceImpl) SendPasswordResetEmail(ctx context.Context, userEmail, resetToken string) (*dto.NotificationResponse, error) {
	req := dto.EmailNotificationRequest{
		To:       []string{userEmail},
		Subject:  "Password Reset Request",
		Template: "password_reset",
		Data: map[string]interface{}{
			"reset_token": resetToken,
		},
		Priority: dto.NotificationPriorityHigh,
	}

	return s.SendEmailNotification(ctx, req)
}

// SMS Notifications
func (s *NotificationServiceImpl) SendSMSNotification(ctx context.Context, req dto.SMSNotificationRequest) (*dto.NotificationResponse, error) {
	// Business validation - check if SMS is enabled
	if !s.settings.SMSEnabled {
		return nil, apperrors.ErrNotificationServiceDisabled
	}

	// Generate notification ID
	notificationID := s.generateNotificationID()

	// TODO: Implement actual SMS sending
	// For now, just log the notification
	log.Printf("Sending SMS to %s: %s", req.Phone, req.Message)

	// Simulate sending
	response := &dto.NotificationResponse{
		ID:      notificationID,
		Type:    "sms",
		Status:  "sent",
		Message: fmt.Sprintf("SMS sent to %s", req.Phone),
		SentAt:  time.Now().Format(time.RFC3339),
	}

	return response, nil
}

func (s *NotificationServiceImpl) SendBookingConfirmationSMS(ctx context.Context, phone string, bookingInfo dto.BookingNotificationData) (*dto.NotificationResponse, error) {
	message := fmt.Sprintf("Your booking %s for %s has been confirmed. Amount: %.2f THB",
		bookingInfo.InvoiceID, bookingInfo.TourName, bookingInfo.Amount)

	req := dto.SMSNotificationRequest{
		Phone:    phone,
		Message:  message,
		Priority: dto.NotificationPriorityHigh,
		Data: map[string]interface{}{
			"booking_id": bookingInfo.BookingID,
			"invoice_id": bookingInfo.InvoiceID,
		},
	}

	return s.SendSMSNotification(ctx, req)
}

// System Alerts
func (s *NotificationServiceImpl) SendSystemAlert(ctx context.Context, title, message string, priority dto.NotificationPriority) (*dto.NotificationResponse, error) {
	req := dto.TelegramNotificationRequest{
		Type:     dto.NotificationTypeSystemAlert,
		Priority: priority,
		Title:    title,
		Message:  message,
		Data: map[string]interface{}{
			"alert_type": "system",
			"timestamp":  time.Now().Unix(),
		},
	}

	return s.SendTelegramNotification(ctx, req)
}

func (s *NotificationServiceImpl) SendErrorAlert(ctx context.Context, service, errorMsg string, data map[string]interface{}) (*dto.NotificationResponse, error) {
	title := fmt.Sprintf("Error in %s", service)
	message := fmt.Sprintf("Service: %s\nError: %s", service, errorMsg)

	// Add data to message
	if data != nil {
		message += "\nAdditional Data:"
		for key, value := range data {
			message += fmt.Sprintf("\n- %s: %v", key, value)
		}
	}

	req := dto.TelegramNotificationRequest{
		Type:     dto.NotificationTypeSystemAlert,
		Priority: dto.NotificationPriorityUrgent,
		Title:    title,
		Message:  message,
		Data: map[string]interface{}{
			"service":    service,
			"error":      errorMsg,
			"data":       data,
			"timestamp":  time.Now().Unix(),
		},
	}

	return s.SendTelegramNotification(ctx, req)
}

// Template Management
func (s *NotificationServiceImpl) FormatBookingMessage(data dto.BookingNotificationData, notifyType dto.NotificationType) string {
	switch notifyType {
	case dto.NotificationTypeNewBooking:
		return fmt.Sprintf(`🎫 **New Booking Created**

📋 Invoice: %s
👤 Customer: %s
🏢 Company: %s
🌍 Tour: %s
📅 Date: %s
💰 Amount: %.2f THB
👨‍💼 Sales: %s

ℹ️ Booking ID: %s`,
			data.InvoiceID, data.CustomerName, data.CompanyName,
			data.TourName, data.PeriodStart, data.Amount,
			data.SalesName, data.BookingID.String())

	case dto.NotificationTypeBookingCancelled:
		return fmt.Sprintf(`❌ **Booking Cancelled**

📋 Invoice: %s
👤 Customer: %s
🏢 Company: %s
🌍 Tour: %s
💰 Amount: %.2f THB

ℹ️ Booking ID: %s`,
			data.InvoiceID, data.CustomerName, data.CompanyName,
			data.TourName, data.Amount, data.BookingID.String())

	default:
		return fmt.Sprintf("Booking notification: %s for %s", notifyType, data.InvoiceID)
	}
}

func (s *NotificationServiceImpl) FormatPaymentMessage(data dto.PaymentNotificationData, notifyType dto.NotificationType) string {
	switch notifyType {
	case dto.NotificationTypePaymentReceived:
		return fmt.Sprintf(`💳 **Payment Slip Uploaded**

📋 Invoice: %s
🧾 Slip ID: %s
👤 Customer: %s
🏢 Company: %s
💰 Amount: %.2f THB

⚠️ Waiting for verification
ℹ️ Booking ID: %s`,
			data.InvoiceID, data.SlipID, data.CustomerName,
			data.CompanyName, data.Amount, data.BookingID.String())

	case dto.NotificationTypePaymentVerified:
		return fmt.Sprintf(`✅ **Payment Verified**

📋 Invoice: %s
🧾 Slip ID: %s
👤 Customer: %s
🏢 Company: %s
💰 Amount: %.2f THB

✅ Payment confirmed
ℹ️ Booking ID: %s`,
			data.InvoiceID, data.SlipID, data.CustomerName,
			data.CompanyName, data.Amount, data.BookingID.String())

	default:
		return fmt.Sprintf("Payment notification: %s for %s", notifyType, data.InvoiceID)
	}
}

func (s *NotificationServiceImpl) FormatUserMessage(data dto.UserNotificationData, notifyType dto.NotificationType) string {
	switch notifyType {
	case dto.NotificationTypeUserRegistered:
		return fmt.Sprintf(`👤 **New User Registration**

📧 Email: %s
🏢 Company: %s
📊 Status: %s

⚠️ Waiting for approval
ℹ️ User ID: %s`,
			data.Email, data.CompanyName, data.Status, data.UserID.String())

	default:
		return fmt.Sprintf("User notification: %s for %s", notifyType, data.Email)
	}
}

func (s *NotificationServiceImpl) FormatPeriodMessage(data dto.PeriodNotificationData, notifyType dto.NotificationType) string {
	switch notifyType {
	case dto.NotificationTypePeriodFull:
		return fmt.Sprintf(`🎫 **Period Fully Booked**

🌍 Tour: %s (%s)
📅 Period: %s - %s
👥 Seats: %d/%d (FULL)

🚫 No more bookings accepted
ℹ️ Period ID: %s`,
			data.TourName, data.TourCode, data.StartDate, data.EndDate,
			data.BookedSeats, data.TotalSeats, data.PeriodID.String())

	default:
		return fmt.Sprintf("Period notification: %s for %s", notifyType, data.TourCode)
	}
}

// Settings Management
func (s *NotificationServiceImpl) GetNotificationSettings(ctx context.Context) (*dto.NotificationSettings, error) {
	return s.settings, nil
}

func (s *NotificationServiceImpl) UpdateNotificationSettings(ctx context.Context, settings dto.NotificationSettings) error {
	// Business validation - validate settings
	if settings.TelegramEnabled && settings.TelegramToken == "" {
		return apperrors.ErrInvalidTelegramConfig
	}

	// Update settings
	s.settings = &settings
	log.Printf("Notification settings updated")

	return nil
}

func (s *NotificationServiceImpl) TestTelegramConnection(ctx context.Context) error {
	if !s.settings.TelegramEnabled {
		return apperrors.ErrNotificationServiceDisabled
	}

	if s.settings.TelegramToken == "" {
		return apperrors.ErrInvalidTelegramConfig
	}

	// TODO: Implement actual Telegram connection test
	log.Printf("Testing Telegram connection...")

	// Simulate test
	return nil
}

func (s *NotificationServiceImpl) TestEmailConnection(ctx context.Context) error {
	if !s.settings.EmailEnabled {
		return apperrors.ErrNotificationServiceDisabled
	}

	// TODO: Implement actual email connection test
	log.Printf("Testing email connection...")

	// Simulate test
	return nil
}

// Batch Operations
func (s *NotificationServiceImpl) SendBulkNotifications(ctx context.Context, notifications []dto.TelegramNotificationRequest) ([]*dto.NotificationResponse, error) {
	responses := make([]*dto.NotificationResponse, len(notifications))

	for i, notification := range notifications {
		response, err := s.SendTelegramNotification(ctx, notification)
		if err != nil {
			// Continue with other notifications even if one fails
			responses[i] = &dto.NotificationResponse{
				ID:     s.generateNotificationID(),
				Type:   "telegram",
				Status: "failed",
				SentAt: time.Now().Format(time.RFC3339),
				Error:  func() *string { s := err.Error(); return &s }(),
			}
		} else {
			responses[i] = response
		}
	}

	return responses, nil
}

func (s *NotificationServiceImpl) SendBulkEmails(ctx context.Context, emails []dto.EmailNotificationRequest) ([]*dto.NotificationResponse, error) {
	responses := make([]*dto.NotificationResponse, len(emails))

	for i, email := range emails {
		response, err := s.SendEmailNotification(ctx, email)
		if err != nil {
			// Continue with other emails even if one fails
			responses[i] = &dto.NotificationResponse{
				ID:     s.generateNotificationID(),
				Type:   "email",
				Status: "failed",
				SentAt: time.Now().Format(time.RFC3339),
				Error:  func() *string { s := err.Error(); return &s }(),
			}
		} else {
			responses[i] = response
		}
	}

	return responses, nil
}

// Helper methods
func (s *NotificationServiceImpl) generateNotificationID() string {
	return fmt.Sprintf("ntf_%d_%s", time.Now().Unix(), uuid.New().String()[:8])
}

func (s *NotificationServiceImpl) getBookingNotificationTitle(notifyType dto.NotificationType) string {
	switch notifyType {
	case dto.NotificationTypeNewBooking:
		return "New Booking"
	case dto.NotificationTypeBookingCancelled:
		return "Booking Cancelled"
	default:
		return "Booking Update"
	}
}

func (s *NotificationServiceImpl) getPaymentNotificationTitle(notifyType dto.NotificationType) string {
	switch notifyType {
	case dto.NotificationTypePaymentReceived:
		return "Payment Slip Uploaded"
	case dto.NotificationTypePaymentVerified:
		return "Payment Verified"
	default:
		return "Payment Update"
	}
}

func (s *NotificationServiceImpl) getUserNotificationTitle(notifyType dto.NotificationType) string {
	switch notifyType {
	case dto.NotificationTypeUserRegistered:
		return "New User Registration"
	default:
		return "User Update"
	}
}

func (s *NotificationServiceImpl) getPeriodNotificationTitle(notifyType dto.NotificationType) string {
	switch notifyType {
	case dto.NotificationTypePeriodFull:
		return "Period Fully Booked"
	default:
		return "Period Update"
	}
}