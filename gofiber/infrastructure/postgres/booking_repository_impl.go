package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
)

type BookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) repositories.BookingRepository {
	return &BookingRepositoryImpl{db: db}
}

func (r *BookingRepositoryImpl) Create(ctx context.Context, entity *models.Booking) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}
	return nil
}

func (r *BookingRepositoryImpl) GetByID(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&booking).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get booking by ID: %w", err)
	}

	return &booking, nil
}

func (r *BookingRepositoryImpl) GetByIDWithPaymentSlips(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.WithContext(ctx).
		Preload("PaymentSlips").
		Where("id = ?", id).
		First(&booking).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get booking by ID with payment slips: %w", err)
	}

	return &booking, nil
}

func (r *BookingRepositoryImpl) Update(ctx context.Context, id uuid.UUID, entity *models.Booking) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(entity)

	if result.Error != nil {
		return fmt.Errorf("failed to update booking: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("booking with ID %s not found", id)
	}

	return nil
}

func (r *BookingRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Booking{}, "id = ?", id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete booking: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("booking with ID %s not found", id)
	}

	return nil
}

func (r *BookingRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking

	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list bookings: %w", err)
	}

	return bookings, nil
}

func (r *BookingRepositoryImpl) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&models.Booking{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count bookings: %w", err)
	}

	return count, nil
}

func (r *BookingRepositoryImpl) GetByInvoiceID(ctx context.Context, invoiceID string) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.WithContext(ctx).Where("invoice_id = ?", invoiceID).First(&booking).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get booking by invoice ID: %w", err)
	}

	return &booking, nil
}

func (r *BookingRepositoryImpl) GetByUserID(ctx context.Context, userID uuid.UUID, offset, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Sales").
		Preload("Tour").
		Preload("Period").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get bookings by user ID: %w", err)
	}

	// Manually load company profiles and sales for each booking
	for _, booking := range bookings {
		// Load company profile
		if booking.UserID != uuid.Nil {
			var companyProfile models.CompanyProfile
			err := r.db.WithContext(ctx).Where("user_id = ?", booking.UserID).First(&companyProfile).Error
			if err == nil {
				booking.CompanyProfile = &companyProfile
				fmt.Printf("Loaded company profile for booking %s, user %s, name: %s\n", booking.ID, booking.UserID, companyProfile.Name)
			} else {
				fmt.Printf("No company profile found for booking %s, user %s: %v\n", booking.ID, booking.UserID, err)
			}
		}

		// Debug check for Sales data loaded by GORM
		if booking.Sales != nil {
			fmt.Printf("Sales already loaded by GORM for booking %s: %s (%s)\n", booking.ID, booking.Sales.Code, booking.Sales.FirstName)
		} else {
			fmt.Printf("No sales data loaded for booking %s\n", booking.ID)
		}

		// Load sales data
		if booking.SalesID != uuid.Nil {
			// Check if Sales is already loaded by GORM preload
			if booking.Sales != nil {
				fmt.Printf("Sales already loaded by GORM for booking %s: %s\n", booking.ID, booking.Sales.Code)
			} else {
				var sales models.Sales
				err := r.db.WithContext(ctx).Where("id = ?", booking.SalesID).First(&sales).Error
				if err == nil {
					booking.Sales = &sales
					fmt.Printf("Manually loaded sales for booking %s: %s\n", booking.ID, sales.Code)
				} else {
					fmt.Printf("No sales found for booking %s, sales_id %s: %v\n", booking.ID, booking.SalesID, err)
				}
			}
		}
	}

	return bookings, nil
}

func (r *BookingRepositoryImpl) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Booking{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count bookings by user ID: %w", err)
	}

	return count, nil
}

func (r *BookingRepositoryImpl) GetByPeriodID(ctx context.Context, periodID uuid.UUID, offset, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking

	err := r.db.WithContext(ctx).
		Where("period_id = ?", periodID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get bookings by period ID: %w", err)
	}

	return bookings, nil
}

func (r *BookingRepositoryImpl) GetBySalesID(ctx context.Context, salesID uuid.UUID, offset, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking

	err := r.db.WithContext(ctx).
		Where("sales_id = ?", salesID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get bookings by sales ID: %w", err)
	}

	return bookings, nil
}

func (r *BookingRepositoryImpl) GetByStatus(ctx context.Context, status models.BookingStatus, offset, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking

	err := r.db.WithContext(ctx).
		Where("status = ?", status).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get bookings by status: %w", err)
	}

	return bookings, nil
}

func (r *BookingRepositoryImpl) GetByDateRange(ctx context.Context, startDate, endDate time.Time, offset, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking

	err := r.db.WithContext(ctx).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get bookings by date range: %w", err)
	}

	return bookings, nil
}

func (r *BookingRepositoryImpl) GetWithDetails(ctx context.Context, id uuid.UUID) (*models.Booking, error) {
	var booking models.Booking
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Period").
		Preload("Period.Tour").
		Preload("Sales").
		Preload("RoomTypes").
		Preload("RateTypes").
		Preload("PaymentSlips").
		Preload("PassengerPassports").
		Where("id = ?", id).
		First(&booking).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Let service decide if this is an error
		}
		return nil, fmt.Errorf("failed to get booking with details: %w", err)
	}

	// Manually load company profile and sales
	if booking.UserID != uuid.Nil {
		var companyProfile models.CompanyProfile
		err := r.db.WithContext(ctx).Where("user_id = ?", booking.UserID).First(&companyProfile).Error
		if err == nil {
			booking.CompanyProfile = &companyProfile
		}
	}

	// Load sales data
	if booking.SalesID != uuid.Nil {
		var sales models.Sales
		err := r.db.WithContext(ctx).Where("id = ?", booking.SalesID).First(&sales).Error
		if err == nil {
			booking.Sales = &sales
		}
	}

	return &booking, nil
}

func (r *BookingRepositoryImpl) GetPendingPayments(ctx context.Context, offset, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking

	err := r.db.WithContext(ctx).
		Where("status = ?", models.BookingStatusPending).
		Where("payment_due_date >= ?", time.Now()).
		Order("payment_due_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending payment bookings: %w", err)
	}

	return bookings, nil
}

func (r *BookingRepositoryImpl) GetOverduePayments(ctx context.Context, offset, limit int) ([]*models.Booking, error) {
	var bookings []*models.Booking

	err := r.db.WithContext(ctx).
		Where("status = ?", models.BookingStatusPending).
		Where("payment_due_date < ?", time.Now()).
		Order("payment_due_date ASC").
		Offset(offset).
		Limit(limit).
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get overdue payment bookings: %w", err)
	}

	return bookings, nil
}

func (r *BookingRepositoryImpl) UpdateStatus(ctx context.Context, id uuid.UUID, status models.BookingStatus) error {
	result := r.db.WithContext(ctx).
		Model(&models.Booking{}).
		Where("id = ?", id).
		Update("status", status)

	if result.Error != nil {
		return fmt.Errorf("failed to update booking status: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("booking with ID %s not found", id)
	}

	return nil
}

func (r *BookingRepositoryImpl) UpdatePaymentAmount(ctx context.Context, id uuid.UUID, amount float64) error {
	result := r.db.WithContext(ctx).
		Model(&models.Booking{}).
		Where("id = ?", id).
		Update("paid_amount", amount)

	if result.Error != nil {
		return fmt.Errorf("failed to update booking payment amount: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("booking with ID %s not found", id)
	}

	return nil
}

func (r *BookingRepositoryImpl) GenerateInvoiceID(ctx context.Context, year int) (string, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&models.Booking{}).
		Where("EXTRACT(YEAR FROM created_at) = ?", year).
		Count(&count).Error

	if err != nil {
		return "", fmt.Errorf("failed to count bookings for invoice ID generation: %w", err)
	}

	// Format: TJ-YYYY-NNNN (e.g., TJ-2024-0001)
	invoiceID := fmt.Sprintf("TJ-%d-%04d", year, count+1)
	return invoiceID, nil
}

// CountBookingsByPeriodAndStatus counts bookings and passengers by period and status
func (r *BookingRepositoryImpl) CountBookingsByPeriodAndStatus(ctx context.Context, periodID uuid.UUID, status string) (int, int, error) {
	log.Printf("[DEBUG] CountBookingsByPeriodAndStatus called - Period: %s, Status: %s", periodID.String(), status)

	// Validate inputs
	if periodID == uuid.Nil {
		log.Printf("[ERROR] Invalid period ID provided: %s", periodID.String())
		return 0, 0, fmt.Errorf("invalid period ID: cannot be nil")
	}

	if status == "" {
		log.Printf("[ERROR] Empty status provided")
		return 0, 0, fmt.Errorf("status cannot be empty")
	}

	// Validate database connection
	if r.db == nil {
		log.Printf("[ERROR] Database connection is nil")
		return 0, 0, fmt.Errorf("database connection not initialized")
	}

	var result struct {
		BookingCount int `gorm:"column:booking_count"`
		PassengerSum int `gorm:"column:passenger_sum"`
	}

	// Build and execute the query
	query := r.db.WithContext(ctx).
		Model(&models.Booking{}).
		Select("COUNT(*) as booking_count, COALESCE(SUM(sum_qty), 0) as passenger_sum").
		Where("period_id = ? AND status = ?", periodID, status)

	log.Printf("[DEBUG] Executing booking count query for period_id: %s, status: %s", periodID.String(), status)

	err := query.Scan(&result).Error

	if err != nil {
		log.Printf("[ERROR] Database query failed - Period: %s, Status: %s, Error: %v",
			periodID.String(), status, err)
		return 0, 0, fmt.Errorf("failed to count bookings by period and status: %w", err)
	}

	log.Printf("[DEBUG] Query results - Period: %s, Status: %s, BookingCount: %d, PassengerSum: %d",
		periodID.String(), status, result.BookingCount, result.PassengerSum)

	// Validate results
	if result.BookingCount < 0 {
		log.Printf("[WARNING] Negative booking count detected: %d", result.BookingCount)
		result.BookingCount = 0
	}

	if result.PassengerSum < 0 {
		log.Printf("[WARNING] Negative passenger sum detected: %d", result.PassengerSum)
		result.PassengerSum = 0
	}

	log.Printf("[DEBUG] CountBookingsByPeriodAndStatus completed - Period: %s, Status: %s, Count: %d, Passengers: %d",
		periodID.String(), status, result.BookingCount, result.PassengerSum)

	return result.BookingCount, result.PassengerSum, nil
}

// FindByPeriodIDAndStatuses finds bookings by period ID and multiple statuses
func (r *BookingRepositoryImpl) FindByPeriodIDAndStatuses(periodID uuid.UUID, statuses []string) ([]*models.Booking, error) {
	log.Printf("🔍 Finding bookings for period %s with statuses: %v", periodID, statuses)

	var bookings []*models.Booking

	err := r.db.Where("period_id = ? AND status IN ?", periodID, statuses).
		Find(&bookings).Error

	if err != nil {
		log.Printf("❌ Error finding bookings by period and statuses: %v", err)
		return nil, fmt.Errorf("failed to find bookings by period and statuses: %w", err)
	}

	log.Printf("✅ Found %d bookings for period %s", len(bookings), periodID)
	return bookings, nil
}

// ListWithFilters lists bookings with advanced filtering, search, and sorting
func (r *BookingRepositoryImpl) ListWithFilters(ctx context.Context, req dto.BookingListRequest) ([]*models.Booking, error) {
	var bookings []*models.Booking

	query := r.db.WithContext(ctx)

	// Apply filters
	if req.UserID != nil {
		query = query.Where("bookings.user_id = ?", *req.UserID)
	}
	if req.SalesID != nil {
		query = query.Where("bookings.sales_id = ?", *req.SalesID)
	}
	if req.TourID != nil {
		query = query.Where("bookings.tour_id = ?", *req.TourID)
	}
	if req.PeriodID != nil {
		query = query.Where("bookings.period_id = ?", *req.PeriodID)
	}
	// Handle status filtering (backward compatible)
	if req.Status != nil {
		query = query.Where("bookings.status = ?", *req.Status)
	} else if len(req.Statuses) > 0 {
		query = query.Where("bookings.status IN ?", req.Statuses)
	}
	if req.StartDate != nil {
		query = query.Where("bookings.created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("bookings.created_at <= ?", *req.EndDate)
	}

	// Apply search (invoice_id, customer_name, customer_phone)
	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + strings.ToLower(*req.Search) + "%"
		query = query.Where(
			"LOWER(bookings.invoice_id) LIKE ? OR LOWER(bookings.customer_name) LIKE ? OR LOWER(bookings.customer_phone) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	// Apply payment slip filtering using JOIN
	if req.PaymentSlipStatus != nil ||
	   req.HasPendingPayments != nil ||
	   req.HasUnconfirmedPayments != nil ||
	   req.RequiresPaymentApproval != nil {

		// Join with payment_slips table
		query = query.Joins("JOIN payment_slips ON payment_slips.booking_id = bookings.id")

		// Only include bookings that are not cancelled or expired when filtering by payment slips
		query = query.Where("bookings.status NOT IN ('cancel', 'expired')")

		if req.PaymentSlipStatus != nil {
			query = query.Where("payment_slips.status = ?", *req.PaymentSlipStatus)
		}

		if req.HasPendingPayments != nil && *req.HasPendingPayments {
			query = query.Where("payment_slips.status = 'pending'")
		}

		if req.HasUnconfirmedPayments != nil && *req.HasUnconfirmedPayments {
			query = query.Where("payment_slips.status IN ('pending', 'reject')")
		}

		if req.RequiresPaymentApproval != nil && *req.RequiresPaymentApproval {
			query = query.Where("payment_slips.status = 'pending'")
		}

		// Use GROUP BY to avoid duplicate bookings when multiple payment slips match
		query = query.Group("bookings.id")
	}

	// Apply sorting with column mapping
	orderClause := "bookings.created_at DESC" // default
	if req.Sort != nil && *req.Sort != "" {
		// Map API field names to database column names
		var sortField string
		switch *req.Sort {
		case "created_at":
			sortField = "bookings.created_at"
		case "updated_at":
			sortField = "bookings.updated_at"
		case "invoice_id":
			sortField = "bookings.invoice_id"
		case "customer_name":
			sortField = "bookings.customer_name"
		case "total_amount":
			sortField = "bookings.sum_total" // sum_total is the actual column name
		default:
			sortField = "bookings.created_at" // fallback to default
		}

		direction := "DESC"
		if req.Order != nil && strings.ToUpper(*req.Order) == "ASC" {
			direction = "ASC"
		}
		orderClause = fmt.Sprintf("%s %s", sortField, direction)
	}

	err := query.
		Preload("User").
		Preload("Sales").
		Preload("Tour").
		Preload("Period").
		Preload("PaymentSlips").
		Preload("PassengerPassports").
		Order(orderClause).
		Offset(req.Offset).
		Limit(req.Limit).
		Find(&bookings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to list bookings with filters: %w", err)
	}

	// Manually load company profiles and sales for each booking
	for _, booking := range bookings {
		// Load company profile
		if booking.UserID != uuid.Nil {
			var companyProfile models.CompanyProfile
			err := r.db.WithContext(ctx).Where("user_id = ?", booking.UserID).First(&companyProfile).Error
			if err == nil {
				booking.CompanyProfile = &companyProfile
				fmt.Printf("Loaded company profile for booking %s, user %s, name: %s\n", booking.ID, booking.UserID, companyProfile.Name)
			} else {
				fmt.Printf("No company profile found for booking %s, user %s: %v\n", booking.ID, booking.UserID, err)
			}
		}

		// Debug check for Sales data loaded by GORM
		if booking.Sales != nil {
			fmt.Printf("Sales already loaded by GORM for booking %s: %s (%s)\n", booking.ID, booking.Sales.Code, booking.Sales.FirstName)
		} else {
			fmt.Printf("No sales data loaded for booking %s\n", booking.ID)
		}

		// Load sales data
		if booking.SalesID != uuid.Nil {
			// Check if Sales is already loaded by GORM preload
			if booking.Sales != nil {
				fmt.Printf("Sales already loaded by GORM for booking %s: %s\n", booking.ID, booking.Sales.Code)
			} else {
				var sales models.Sales
				err := r.db.WithContext(ctx).Where("id = ?", booking.SalesID).First(&sales).Error
				if err == nil {
					booking.Sales = &sales
					fmt.Printf("Manually loaded sales for booking %s, sales %s\n", booking.ID, booking.SalesID)
				} else {
					fmt.Printf("No sales found for booking %s, sales %s: %v\n", booking.ID, booking.SalesID, err)
				}
			}
		}
	}

	return bookings, nil
}

// CountWithFilters counts bookings with the same filters as ListWithFilters
func (r *BookingRepositoryImpl) CountWithFilters(ctx context.Context, req dto.BookingListRequest) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&models.Booking{})

	// Apply same filters as ListWithFilters
	if req.UserID != nil {
		query = query.Where("bookings.user_id = ?", *req.UserID)
	}
	if req.SalesID != nil {
		query = query.Where("bookings.sales_id = ?", *req.SalesID)
	}
	if req.TourID != nil {
		query = query.Where("bookings.tour_id = ?", *req.TourID)
	}
	if req.PeriodID != nil {
		query = query.Where("bookings.period_id = ?", *req.PeriodID)
	}
	// Handle status filtering (backward compatible)
	if req.Status != nil {
		query = query.Where("bookings.status = ?", *req.Status)
	} else if len(req.Statuses) > 0 {
		query = query.Where("bookings.status IN ?", req.Statuses)
	}
	if req.StartDate != nil {
		query = query.Where("bookings.created_at >= ?", *req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("bookings.created_at <= ?", *req.EndDate)
	}

	// Apply same search
	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + strings.ToLower(*req.Search) + "%"
		query = query.Where(
			"LOWER(bookings.invoice_id) LIKE ? OR LOWER(bookings.customer_name) LIKE ? OR LOWER(bookings.customer_phone) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	// Apply same payment slip filtering using JOIN
	if req.PaymentSlipStatus != nil ||
	   req.HasPendingPayments != nil ||
	   req.HasUnconfirmedPayments != nil ||
	   req.RequiresPaymentApproval != nil {

		// Join with payment_slips table
		query = query.Joins("JOIN payment_slips ON payment_slips.booking_id = bookings.id")

		// Only include bookings that are not cancelled or expired when filtering by payment slips
		query = query.Where("bookings.status NOT IN ('cancel', 'expired')")

		if req.PaymentSlipStatus != nil {
			query = query.Where("payment_slips.status = ?", *req.PaymentSlipStatus)
		}

		if req.HasPendingPayments != nil && *req.HasPendingPayments {
			query = query.Where("payment_slips.status = 'pending'")
		}

		if req.HasUnconfirmedPayments != nil && *req.HasUnconfirmedPayments {
			query = query.Where("payment_slips.status IN ('pending', 'reject')")
		}

		if req.RequiresPaymentApproval != nil && *req.RequiresPaymentApproval {
			query = query.Where("payment_slips.status = 'pending'")
		}

		// Use GROUP BY to avoid counting duplicates
		query = query.Group("bookings.id")
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count bookings with filters: %w", err)
	}

	return count, nil
}

// WithTx returns a new repository instance with the given transaction
func (r *BookingRepositoryImpl) WithTx(tx *gorm.DB) repositories.BookingRepository {
	return &BookingRepositoryImpl{db: tx}
}