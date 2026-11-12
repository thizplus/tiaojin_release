package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
)


type BookingService interface {
	CreateBooking(ctx context.Context, userID uuid.UUID, req dto.CreateBookingRequest) (*dto.BookingResponse, error)
	GetBookingByID(ctx context.Context, bookingID uuid.UUID) (*dto.BookingResponse, error)
	GetBookingModelByID(ctx context.Context, bookingID uuid.UUID) (*models.Booking, error)
	GetBookingByInvoiceID(ctx context.Context, invoiceID string) (*dto.BookingResponse, error)
	GetBookingWithDetails(ctx context.Context, bookingID uuid.UUID) (*dto.BookingResponse, error)
	UpdateBooking(ctx context.Context, bookingID uuid.UUID, req dto.UpdateBookingRequest) (*dto.BookingResponse, error)
	UpdateBookingStatus(ctx context.Context, req dto.UpdateBookingStatusRequest) (*dto.BookingResponse, error)
	CancelBooking(ctx context.Context, bookingID uuid.UUID, reason string) error
	ListBookings(ctx context.Context, req dto.BookingListRequest) (*dto.BookingListResponse, error)
	GetUserBookings(ctx context.Context, userID uuid.UUID, offset, limit int) (*dto.BookingListResponse, error)
	GetSalesBookings(ctx context.Context, salesID uuid.UUID, offset, limit int) (*dto.BookingListResponse, error)

	// Payment Management
	UploadPaymentSlip(ctx context.Context, req dto.UploadPaymentSlipRequest) (*dto.PaymentSlipResponse, error)
	VerifyPaymentSlip(ctx context.Context, req dto.VerifyPaymentSlipRequest) (*dto.PaymentSlipResponse, error)
	GetPendingPayments(ctx context.Context, offset, limit int) ([]*dto.PaymentSlipResponse, error)
	CreateAdminPayment(ctx context.Context, bookingID uuid.UUID, req dto.CreateAdminPaymentRequest, adminUserID uuid.UUID) (*dto.AdminPaymentResponse, error)

	// Passport Management
	UploadPassport(ctx context.Context, req dto.UploadPassportRequest) (*dto.PassengerPassportResponse, error)
	ProcessPassportOCR(ctx context.Context, passportID uuid.UUID) (*dto.PassengerPassportResponse, error)
	UpdatePassportInfo(ctx context.Context, passportID uuid.UUID, passport *dto.PassengerPassportResponse) (*dto.PassengerPassportResponse, error)

	// Document Management
	GetBookingDocuments(ctx context.Context, bookingID uuid.UUID) (*dto.BookingDocumentsResponse, error)

	// Business Logic
	// CalculateBookingTotal calculates total booking costs and returns:
	// (sumRoom, sumRate, sumAdditionalServices, sumNet, error)
	CalculateBookingTotal(ctx context.Context, req dto.CreateBookingRequest) (float64, float64, float64, float64, error)
	CheckPeriodAvailability(ctx context.Context, periodID uuid.UUID, quantity int) (bool, error)
	GenerateInvoiceID(ctx context.Context) (string, error)
	GetBookingStats(ctx context.Context, dateFrom, dateTo *time.Time) (*dto.BookingStatsResponse, error)
	GetBookingStatistics(ctx context.Context, req dto.BookingListRequest) (*dto.BookingStatisticsResponse, error)
	GetAdminReport(ctx context.Context, year int, month int) (*dto.AdminReportResponse, error)
	GetSalesReport(ctx context.Context, userID uuid.UUID, year int, month int) (*dto.AdminReportResponse, error)

	// Tiered Expiration Management
	SetBookingExpiration(ctx context.Context, bookingID uuid.UUID, periodID uuid.UUID) error
	CalculateBookingExpiration(ctx context.Context, periodID uuid.UUID, createdAt time.Time) (*time.Time, error)
	GetExpiringBookings(ctx context.Context, withinHours int, offset, limit int) ([]*dto.BookingResponse, error)
	ProcessExpiredBookings(ctx context.Context, dryRun bool) (*dto.ExpirationProcessResult, error)
	RefreshBookingExpirations(ctx context.Context, periodID *uuid.UUID) error
	GetBookingExpirationStats(ctx context.Context) (*dto.ExpirationStatsResponse, error)

	// Reports
	GetAgentReport(ctx context.Context, userID uuid.UUID, year int, month int) (*dto.AdminReportResponse, error)

	// Dynamic Job-based Expiration Management
	CreateBookingExpiryJob(ctx context.Context, bookingID uuid.UUID, expiresAt time.Time) error
	DeleteBookingExpiryJob(ctx context.Context, bookingID uuid.UUID) error
	ExpireBooking(ctx context.Context, bookingID uuid.UUID) error

	// Payment Slip Reversal System
	ReversePaymentSlip(ctx context.Context, req dto.ReversePaymentSlipRequest, adminUserID uuid.UUID) (*dto.PaymentSlipResponse, error)
	GetPaymentSlipHistory(ctx context.Context, paymentSlipID uuid.UUID) (*dto.PaymentSlipHistoryResponse, error)
}