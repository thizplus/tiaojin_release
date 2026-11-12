package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookingResponse struct {
	ID            uuid.UUID                    `json:"id"`
	InvoiceID     string                       `json:"invoice_id"`
	UserID        uuid.UUID                    `json:"user_id"`
	TourID        uuid.UUID                    `json:"tour_id"`
	PeriodID      uuid.UUID                    `json:"period_id"`
	SalesID       uuid.UUID                    `json:"sales_id"`
	CustomerName  string                       `json:"customer_name"`
	CustomerPhone string                       `json:"customer_phone"`
	Note          *string                      `json:"note"`
	SumRoom               float64                      `json:"sum_room"`
	SumRate               float64                      `json:"sum_rate"`
	SumAdditionalServices float64                      `json:"sum_additional_services"`
	SumQty                int                          `json:"sum_qty"`
	SumCommissionQty      int                          `json:"sum_commission_qty"`
	SumNet                float64                      `json:"sum_net"`
	SumCommission         float64                      `json:"sum_commission"`
	Discount              float64                      `json:"discount"`
	SumDeposit            float64                      `json:"sum_deposit"`
	SumPayment            float64                      `json:"sum_payment"`

	// ค่าบริการ
	ServiceFeePerPax float64 `json:"service_fee_per_pax"`
	ServiceFeeTotal  float64 `json:"service_fee_total"`

	// VAT fields (ภาษีมูลค่าเพิ่ม — คิดจากค่าบริการเท่านั้น)
	IncludeVAT   bool    `json:"include_vat"`
	VATRate      float64 `json:"vat_rate"`
	VATAmount    float64 `json:"vat_amount"`

	// Withholding Tax fields (หัก ณ ที่จ่าย)
	IncludeWHT   bool    `json:"include_wht"`
	WHTRate      float64 `json:"wht_rate"`
	WHTAmount    float64 `json:"wht_amount"`

	// Calculated totals
	NetBeforeTax float64 `json:"net_before_tax"`
	GrandTotal   float64 `json:"grand_total"`

	DateDeposit   string                       `json:"date_deposit"`
	DateFinal     string                       `json:"date_final"`
	Status        string                       `json:"status"`
	ExpiresAt     *string                      `json:"expires_at,omitempty"`
	CreatedAt     string                       `json:"created_at"`
	UpdatedAt     string                       `json:"updated_at"`

	// Relations
	TourDetails        *TourResponse                      `json:"tour_details,omitempty"`
	PeriodDetails      *PeriodSummaryResponse             `json:"period_details,omitempty"`
	SalesDetails       *SalesResponse                     `json:"sales_details,omitempty"`
	CompanyProfile     *CompanyProfileResponse            `json:"company_profile,omitempty"`
	UserDetails        *UserResponse                      `json:"user_details,omitempty"`
	RoomTypes          []*BookingRoomTypeResponse         `json:"room_types,omitempty"`
	RateTypes          []*BookingRateTypeResponse         `json:"rate_types,omitempty"`
	AdditionalServices []*BookingAdditionalServiceResponse `json:"additional_services,omitempty"`
	PaymentSlips       []*PaymentSlipResponse             `json:"payment_slips,omitempty"`
	Passports          []*PassengerPassportResponse       `json:"passports,omitempty"`
}

type CreateBookingRequest struct {
	TourID             uuid.UUID                                `json:"tour_id" validate:"required"`
	PeriodID           uuid.UUID                                `json:"period_id" validate:"required"`
	SalesID            uuid.UUID                                `json:"sales_id" validate:"required"`
	UserID             *uuid.UUID                               `json:"user_id,omitempty"` // Optional: for admin/super to book for others
	CustomerName       string                                   `json:"customer_name" validate:"required,max=200"`
	CustomerPhone      string                                   `json:"customer_phone" validate:"required"`
	Note               *string                                  `json:"note,omitempty"`
	RoomTypes          []*CreateBookingRoomTypeRequest          `json:"room_types"`
	RateTypes          []*CreateBookingRateTypeRequest          `json:"rate_types" validate:"required,min=1"`
	AdditionalServices []*CreateBookingAdditionalServiceRequest `json:"additional_services"`

	// VAT and Withholding Tax options (Agent เลือกตอนจอง)
	IncludeVAT bool `json:"include_vat"` // รวม VAT 7% (default: false)
	IncludeWHT bool `json:"include_wht"` // หัก ณ ที่จ่าย 3% (default: false)
}

type CreateBookingRoomTypeRequest struct {
	ID       string `json:"id" validate:"required"`
	Quantity int    `json:"quantity" validate:"min=1"`
}

type CreateBookingRateTypeRequest struct {
	ID       string `json:"id" validate:"required"`
	Quantity int    `json:"quantity" validate:"min=1"`
}

type CreateBookingAdditionalServiceRequest struct {
	ID       string `json:"id" validate:"required"`
	Quantity int    `json:"quantity" validate:"min=1"`
}

type UpdateBookingRequest struct {
	CustomerName  *string `json:"customer_name,omitempty" validate:"omitempty,max=200"`
	CustomerPhone *string `json:"customer_phone,omitempty"`
	Note          *string `json:"note,omitempty"`
}

type UpdateBookingStatusRequest struct {
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
	Status    string    `json:"status" validate:"required,oneof=pending processing confirm cancel completed"`
	Reason    *string   `json:"reason,omitempty"`
}

type BookingRoomTypeResponse struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Price    float64   `json:"price"`
	Quantity int       `json:"quantity"`
	Total    float64   `json:"total"`
}

type BookingRateTypeResponse struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Price    float64   `json:"price"`
	Quantity int       `json:"quantity"`
	Exclude  bool      `json:"exclude"`
	Total    float64   `json:"total"`
}

type BookingAdditionalServiceResponse struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Price    float64   `json:"price"`
	Quantity int       `json:"quantity"`
	RateType string    `json:"rate_type"`
	UnitName string    `json:"unit_name"`
	Total    float64   `json:"total"`
}


type PassengerPassportResponse struct {
	ID               uuid.UUID `json:"id"`
	BookingID        uuid.UUID `json:"booking_id"`
	PassportNumber   string    `json:"passport_number"`
	Image            *string   `json:"image"`
	Status           string    `json:"status"`
	TitleTH          *string   `json:"title_th"`
	NameTH           *string   `json:"name_th"`
	SurnameTH        *string   `json:"surname_th"`
	TitleEN          *string   `json:"title_en"`
	NameEN           *string   `json:"name_en"`
	SurnameEN        *string   `json:"surname_en"`
	Gender           *string   `json:"gender"`
	Nationality      *string   `json:"nationality"`
	CountryCode      *string   `json:"country_code"`
	DateOfBirth      *string   `json:"date_of_birth"`
	DateOfExpiry     *string   `json:"date_of_expiry"`
	DateOfIssue      *string   `json:"date_of_issue"`
	IDNumber         *string   `json:"id_number"`
	PlaceOfBirth     *string   `json:"place_of_birth"`
	IssuingAuthority *string   `json:"issuing_authority"`
	OCRText          *string   `json:"ocr_text"`
	OCRConfidence    *float32  `json:"ocr_confidence"`
	ProcessedDate    *string   `json:"processed_date"`
	ErrorMessage     *string   `json:"error_message"`
	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`
}



type UploadPassportRequest struct {
	BookingID      uuid.UUID `json:"booking_id" validate:"required"`
	PassportNumber string    `json:"passport_number" validate:"required"`
	ImageData      []byte    `json:"image_data" validate:"required"`
	FileName       string    `json:"file_name" validate:"required"`
	ImageURL       string    `json:"image_url,omitempty"` // S3 Storage URL
}

type BookingDocument struct {
	ID         uuid.UUID `json:"id"`
	Filename   string    `json:"filename"`
	Type       string    `json:"type"` // payment_slip, passport, other
	UploadedAt string    `json:"uploaded_at"`
	FileURL    string    `json:"file_url"`
	FileSize   int64     `json:"file_size"`
	Status     string    `json:"status"` // pending, confirm, processed, etc.
	Amount     *float64  `json:"amount,omitempty"` // for payment slips
}

type BookingDocumentsResponse struct {
	BookingID uuid.UUID         `json:"booking_id"`
	Documents []BookingDocument `json:"documents"`
}

type BookingListRequest struct {
	UserID    *uuid.UUID `json:"user_id,omitempty"`
	SalesID   *uuid.UUID `json:"sales_id,omitempty"`
	TourID    *uuid.UUID `json:"tour_id,omitempty"`
	PeriodID  *uuid.UUID `json:"period_id,omitempty"`
	Status    *string    `json:"status,omitempty" validate:"omitempty"`           // Single status (backward compatibility)
	Statuses  []string   `json:"statuses,omitempty" validate:"omitempty,dive,oneof=pending processing confirm cancel completed"` // Multiple statuses
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	Search    *string    `json:"search,omitempty"`
	Sort      *string    `json:"sort,omitempty" validate:"omitempty,oneof=created_at updated_at invoice_id customer_name total_amount"`
	Order     *string    `json:"order,omitempty" validate:"omitempty,oneof=asc desc"`

	// Payment slip filtering parameters
	PaymentSlipStatus        *string `json:"payment_slip_status,omitempty" validate:"omitempty,oneof=pending confirm reject"`
	HasPendingPayments       *bool   `json:"has_pending_payments,omitempty"`
	HasUnconfirmedPayments   *bool   `json:"has_unconfirmed_payments,omitempty"`
	RequiresPaymentApproval  *bool   `json:"requires_payment_approval,omitempty"`

	Offset    int        `json:"offset" validate:"min=0"`
	Limit     int        `json:"limit" validate:"min=1,max=100"`
}

type BookingListResponse struct {
	Bookings   []*BookingResponse `json:"bookings"`
	Pagination PaginationMeta     `json:"pagination"`
}

type BookingStatsResponse struct {
	TotalBookings     int64   `json:"total_bookings"`
	PendingBookings   int64   `json:"pending_bookings"`
	ConfirmedBookings int64   `json:"confirmed_bookings"`
	CancelledBookings int64   `json:"cancelled_bookings"`
	TotalRevenue      float64 `json:"total_revenue"`
	PendingRevenue    float64 `json:"pending_revenue"`
}

type BookingStatisticsResponse struct {
	TotalRecords      int64   `json:"total_records"`
	TotalRevenue      float64 `json:"total_revenue"`
	TotalDeposit      float64 `json:"total_deposit"`
	TotalPayments     float64 `json:"total_payments"`
	TotalCommission   float64 `json:"total_commission"`
	StatusCounts      map[string]int64 `json:"status_counts"`
	AvgBookingValue   float64 `json:"avg_booking_value"`
}

// Admin Report DTOs
type AdminReportResponse struct {
	TotalBookings     int64                    `json:"total_bookings"`
	TotalSeats        int64                    `json:"total_seats"`
	TotalSales        float64                  `json:"total_sales"`        // ยอดขายทั้งหมด (SumNet)
	ActualRevenue     float64                  `json:"actual_revenue"`     // เงินที่เก็บได้จริง (จาก PaymentSlips ที่ confirm)
	DailyStatusReport []*DailyBookingStatusData `json:"daily_status_report"`
}

type DailyBookingStatusData struct {
	Date       string `json:"date"`       // "2024-04-01"
	Pending    int64  `json:"pending"`
	Processing int64  `json:"processing"`
	Confirm    int64  `json:"confirm"`
	Cancel     int64  `json:"cancel"`
	Completed  int64  `json:"completed"`
	Expired    int64  `json:"expired"`
	TotalSeats int64  `json:"total_seats"` // จำนวนที่นั่งทั้งหมดวันนั้น
}


// Payment Slip DTOs
type UploadPaymentSlipRequest struct {
	BookingID uuid.UUID `json:"booking_id" validate:"required"`
	Amount    float64   `json:"amount" validate:"required,gt=0"`
	Image     string    `json:"image" validate:"required"`
	ImageData []byte    `json:"-"` // raw image bytes สำหรับส่งไป SlipOK (ไม่ต้อง serialize)
}

type VerifyPaymentSlipRequest struct {
	PaymentSlipID      uuid.UUID `json:"payment_slip_id" validate:"required"`
	Status             string    `json:"status" validate:"required,oneof=confirm reject"`
	Note               *string   `json:"note,omitempty"`
	ActualAmount       *float64  `json:"actual_amount,omitempty" validate:"omitempty,gt=0"`
	AdjustmentReason   *string   `json:"adjustment_reason,omitempty"`
}

type PaymentSlipResponse struct {
	ID               uuid.UUID `json:"id"`
	BookingID        uuid.UUID `json:"booking_id"`
	SlipID           string    `json:"slip_id"`
	Image            *string   `json:"image,omitempty"`                      // Optional for admin direct payments
	Amount           float64   `json:"amount"`
	OriginalAmount   *float64  `json:"original_amount,omitempty"`
	AdjustmentReason *string   `json:"adjustment_reason,omitempty"`
	Status           string    `json:"status"`
	Note             *string   `json:"note,omitempty"`

	// Payment method and admin creation tracking
	PaymentMethod      string  `json:"payment_method"`
	IsAdminCreated     bool    `json:"is_admin_created"`
	CreatedByAdmin     *string `json:"created_by_admin,omitempty"`
	AdminCreatedNote   *string `json:"admin_created_note,omitempty"`

	// SlipOK verification fields
	PaymentDate     *string `json:"payment_date,omitempty"`
	SendingBank     *string `json:"sending_bank,omitempty"`
	ReceivingBank   *string `json:"receiving_bank,omitempty"`
	ReceiverAccount *string `json:"receiver_account,omitempty"`
	ReceiverName    *string `json:"receiver_name,omitempty"`
	SenderName      *string `json:"sender_name,omitempty"`
	TransRef        *string `json:"trans_ref,omitempty"`
	SlipVerified    bool    `json:"slip_verified"`
	SlipVerifyError *string `json:"slip_verify_error,omitempty"`

	// Reversal tracking fields
	IsReversed     bool    `json:"is_reversed"`
	ReversedAt     *string `json:"reversed_at,omitempty"`
	ReversedBy     *string `json:"reversed_by,omitempty"`
	ReversalReason *string `json:"reversal_reason,omitempty"`

	CreatedAt        string    `json:"created_at"`
	UpdatedAt        string    `json:"updated_at"`

	// Optional booking details
	Booking *BookingResponse `json:"booking,omitempty"`
}

// Reversal System DTOs
type ReversePaymentSlipRequest struct {
	PaymentSlipID  uuid.UUID `json:"payment_slip_id" validate:"required"`
	ReversalReason string    `json:"reversal_reason" validate:"required,min=10,max=500"`
}

type PaymentSlipHistoryResponse struct {
	PaymentSlip *PaymentSlipResponse              `json:"payment_slip"`
	Timeline    []*PaymentSlipTimelineEntry       `json:"timeline"`
	BookingInfo *PaymentSlipBookingInfoResponse   `json:"booking_info"`
}

type PaymentSlipTimelineEntry struct {
	Action      string    `json:"action"`
	Timestamp   string    `json:"timestamp"`
	PerformedBy *string   `json:"performed_by,omitempty"`
	Details     string    `json:"details"`
	Amount      *float64  `json:"amount,omitempty"`
}

type PaymentSlipBookingInfoResponse struct {
	BookingID           uuid.UUID `json:"booking_id"`
	InvoiceID           string    `json:"invoice_id"`
	CustomerName        string    `json:"customer_name"`
	TotalAmount         float64   `json:"total_amount"`
	DepositAmount       float64   `json:"deposit_amount"`
	CurrentPayments     float64   `json:"current_payments"`
	RemainingBalance    float64   `json:"remaining_balance"`
	BookingStatus       string    `json:"booking_status"`
}

// Admin Direct Payment DTOs
type CreateAdminPaymentRequest struct {
	Amount        float64  `json:"amount" validate:"required,gt=0"`
	PaymentMethod string   `json:"payment_method" validate:"required,oneof=cash bank_transfer credit_card adjustment"`
	Note          *string  `json:"note,omitempty"`
	SlipImage     *string  `json:"slip_image,omitempty"` // Optional slip image URL
}

type AdminPaymentResponse struct {
	PaymentSlip   *PaymentSlipResponse `json:"payment_slip"`
	BookingStatus string              `json:"booking_status"`
	Message       string              `json:"message"`
}

