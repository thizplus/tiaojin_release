package errors

import "errors"

// Common Entity Not Found Errors
var (
	ErrUserNotFound           = errors.New("user not found")
	ErrCityNotFound          = errors.New("city not found")
	ErrTourNotFound          = errors.New("tour not found")
	ErrBannerNotFound        = errors.New("banner not found")
	ErrBookingNotFound       = errors.New("booking not found")
	ErrPeriodNotFound        = errors.New("period not found")
	ErrSalesNotFound         = errors.New("sales not found")
	ErrCompanyProfileNotFound = errors.New("company profile not found")
	ErrPaymentSlipNotFound   = errors.New("payment slip not found")
	ErrPassportNotFound      = errors.New("passport not found")
	ErrNotFound              = errors.New("resource not found")
	ErrSeatHoldNotFound      = errors.New("seat hold not found")
)

// Authentication & Authorization Errors
var (
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrAccountInactive    = errors.New("account is not active")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
)

// Status Validation Errors
var (
	ErrInvalidBannerStatus = errors.New("invalid banner status")
	ErrInvalidCityStatus   = errors.New("invalid city status")
	ErrInvalidTourStatus   = errors.New("invalid tour status")
	ErrInvalidBookingStatus = errors.New("invalid booking status")
	ErrInvalidPeriodStatus  = errors.New("invalid period status")
	ErrInvalidSalesStatus   = errors.New("invalid sales status")
)

// Business Logic Errors
var (
	ErrEmptyBannerList          = errors.New("banner list cannot be empty")
	ErrDuplicateBannerID        = errors.New("duplicate banner ID in order request")
	ErrTourCodeExists           = errors.New("tour code already exists")
	ErrSalesCodeExists          = errors.New("sales code already exists")
	ErrCompanyProfileExists     = errors.New("company profile already exists for this user")
	ErrTaxIDAlreadyExists       = errors.New("tax ID already exists")
	ErrInsufficientCapacity     = errors.New("insufficient period capacity")
	ErrInsufficientSeats        = errors.New("insufficient available seats")
	ErrBookingCannotBeCancelled = errors.New("booking cannot be cancelled")
	ErrPeriodNotAvailable       = errors.New("period not available for booking")
	ErrTourHasBookings          = errors.New("cannot delete tour with existing bookings")
	ErrPeriodHasBookings        = errors.New("cannot delete period with existing bookings")
	ErrSalesHasBookings         = errors.New("cannot delete sales with existing bookings")
	ErrSalesAlreadyActive       = errors.New("sales is already active")
	ErrSalesAlreadyInactive     = errors.New("sales is already inactive")
)

// Validation Errors
var (
	ErrInvalidTourCode       = errors.New("invalid tour code format")
	ErrInvalidSalesCode      = errors.New("invalid sales code format")
	ErrInvalidTaxIDFormat    = errors.New("invalid tax ID format")
	ErrInvalidZipcodeFormat  = errors.New("invalid zipcode format")
	ErrInvalidPeriodDates    = errors.New("invalid period dates")
	ErrInvalidPaymentAmount  = errors.New("invalid payment amount")
	ErrInvalidRateType       = errors.New("invalid rate type")
	ErrPeriodInPast          = errors.New("period start date is in the past")
)

// File Upload Errors
var (
	ErrInvalidLogoFormat        = errors.New("invalid logo format")
	ErrInvalidImageFormat       = errors.New("invalid image format")
	ErrInvalidTourFileType      = errors.New("invalid file type")
	ErrInvalidPeriodFileType    = errors.New("invalid file type")
	ErrLogoTooLarge            = errors.New("logo file too large")
	ErrTourFileTooLarge        = errors.New("file too large")
)

// Notification Service Errors
var (
	ErrNotificationNotSent         = errors.New("notification not sent")
	ErrInvalidTelegramConfig       = errors.New("invalid telegram configuration")
	ErrInvalidEmailConfig          = errors.New("invalid email configuration")
	ErrInvalidSMSConfig            = errors.New("invalid SMS configuration")
	ErrNotificationServiceDisabled = errors.New("notification service is disabled")
	ErrInvalidTemplate             = errors.New("invalid notification template")
	ErrRateLimitExceeded           = errors.New("rate limit exceeded")
)