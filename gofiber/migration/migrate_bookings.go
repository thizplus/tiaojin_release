package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoRoomType represents room type in booking from MongoDB
type MongoRoomType struct {
	Title string      `bson:"title"`
	Price interface{} `bson:"price"` // Can be string or number
	Qty   interface{} `bson:"qty"`   // Can be string or number
}

// MongoRateType represents rate type in booking from MongoDB
type MongoRateType struct {
	Title   string      `bson:"title"`
	Price   interface{} `bson:"price"` // Can be string or number
	Qty     interface{} `bson:"qty"`   // Can be string or number
	Exclude bool        `bson:"exclude"`
}

// MongoSlipType represents payment slip from MongoDB
type MongoSlipType struct {
	ID         string      `bson:"id"`
	Image      string      `bson:"img"`
	CreateDate string      `bson:"createDate"`
	Amount     interface{} `bson:"amount"` // Can be string or number
	Status     string      `bson:"status"`
	UpdateDate string      `bson:"updateDate"`
}

// MongoPassportPersonal represents personal info in passport from MongoDB
type MongoPassportPersonal struct {
	TitleTH      string `bson:"titleTH"`
	NameTH       string `bson:"nameTH"`
	SurnameTH    string `bson:"surnameTH"`
	TitleEN      string `bson:"titleEN"`
	NameEN       string `bson:"nameEN"`
	SurnameEN    string `bson:"surnameEN"`
	Gender       string `bson:"gender"`
	Nationality  string `bson:"nationality"`
	DateOfBirth  string `bson:"dateOfBirth"`
	IDNumber     string `bson:"idNumber"`
	PlaceOfBirth string `bson:"placeOfBirth"`
	DateOfExpiry string `bson:"dateOfExpiry"`
}

// MongoPassportType represents passport from MongoDB
type MongoPassportType struct {
	ID             string                 `bson:"id"`
	Image          string                 `bson:"img"`
	CreateDate     string                 `bson:"createDate"`
	PassportNumber string                 `bson:"passportNumber"`
	Personal       *MongoPassportPersonal `bson:"personal"`
	Status         string                 `bson:"status"`
	ProcessingDate string                 `bson:"processingDate"`
	ErrorMessage   string                 `bson:"errorMessage"`
}

// MongoBooking represents booking document from MongoDB
type MongoBooking struct {
	ID            primitive.ObjectID  `bson:"_id"`
	InvoiceID     string              `bson:"invoiceId"`
	ProfileID     primitive.ObjectID  `bson:"profileId"`
	TourID        primitive.ObjectID  `bson:"tour"`
	TourName      string              `bson:"tourName"`
	TourCode      string              `bson:"tourCode"`
	PeriodID      primitive.ObjectID  `bson:"periodCode"`
	PeriodStart   string              `bson:"periodStart"`
	PeriodEnd     string              `bson:"periodEnd"`
	CustomerName  string              `bson:"cusName"`
	CustomerPhone string              `bson:"cusTel"`
	SalesID       primitive.ObjectID  `bson:"salesCode"`
	SalesBooking  string              `bson:"salesBooking"`
	Note          string              `bson:"note"`
	RoomTypes     []MongoRoomType     `bson:"roomTypes"`
	RateTypes     []MongoRateType     `bson:"rateTypes"`
	SumRoom       interface{}         `bson:"sumRoom"`       // Can be string or number
	SumRate       interface{}         `bson:"sumRate"`       // Can be string or number
	SumQty        interface{}         `bson:"sumQty"`        // Can be string or number
	SumNet        interface{}         `bson:"sumNet"`        // Can be string or number
	SumCommission interface{}         `bson:"sumCommission"` // Can be string or number
	Discount      interface{}         `bson:"discount"`      // Can be string or number
	SumDeposit    interface{}         `bson:"sumDeposit"`    // Can be string or number
	DateDeposit   string              `bson:"dateDeposit"`
	DateFinal     string              `bson:"dateFinal"`
	SumPayment    interface{}         `bson:"sumPayment"` // Can be string or number
	Status        string              `bson:"status"`
	Slips         []MongoSlipType     `bson:"slip"`
	Passports     []MongoPassportType `bson:"passport"`
	CreateDate    string              `bson:"createDate"`
	UpdateDate    string              `bson:"updateDate"`
}

// MigrateBookings migrates bookings from MongoDB to PostgreSQL
func MigrateBookings(ctx context.Context, stats *MigrationStats) error {
	collection := GetMongoDB().Collection("tbl_booking")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	var bookings []MongoBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return fmt.Errorf("failed to decode bookings: %w", err)
	}

	fmt.Printf("    Found %d bookings in MongoDB\n", len(bookings))

	// Get default sales UUID (first one in PostgreSQL)
	var defaultSalesUUID uuid.UUID
	err = GetPostgresDB().QueryRowContext(ctx, "SELECT id FROM sales LIMIT 1").Scan(&defaultSalesUUID)
	if err != nil {
		return fmt.Errorf("failed to get default sales: %w", err)
	}

	// Get default user UUID (first one in PostgreSQL)
	var defaultUserUUID uuid.UUID
	err = GetPostgresDB().QueryRowContext(ctx, "SELECT id FROM users LIMIT 1").Scan(&defaultUserUUID)
	if err != nil {
		return fmt.Errorf("failed to get default user: %w", err)
	}

	// Track used invoice IDs to handle duplicates
	usedInvoiceIDs := make(map[string]int)

	skipped := 0
	skippedExisting := 0
	usedDefaultSales := 0
	usedDefaultUser := 0
	renamedInvoices := 0
	for _, booking := range bookings {
		// Check if booking already exists by invoice_id
		var existingID string
		checkQuery := `SELECT id FROM bookings WHERE invoice_id = $1 AND deleted_at IS NULL LIMIT 1`
		err := GetPostgresDB().QueryRowContext(ctx, checkQuery, booking.InvoiceID).Scan(&existingID)
		if err == nil {
			// Booking already exists, set mapping and skip
			existingUUID, _ := uuid.Parse(existingID)
			SetMapping(CollectionBookings, booking.ID, existingUUID)
			skippedExisting++
			continue
		}
		newBookingID := uuid.New()
		SetMapping(CollectionBookings, booking.ID, newBookingID)

		// Get user UUID from profile mapping - use default if not found
		userUUID, ok := GetMapping(CollectionUsers, booking.ProfileID)
		if !ok {
			userUUID = defaultUserUUID
			usedDefaultUser++
		}

		// Get tour UUID
		tourUUID, ok := GetMapping(CollectionTours, booking.TourID)
		if !ok {
			fmt.Printf("    WARNING: Tour not found for booking %s (tour_id: %s) - skipping\n", booking.InvoiceID, booking.TourID.Hex())
			skipped++
			continue
		}

		// Get period UUID
		periodUUID, ok := GetMapping(CollectionPeriods, booking.PeriodID)
		if !ok {
			fmt.Printf("    WARNING: Period not found for booking %s (period_id: %s) - skipping\n", booking.InvoiceID, booking.PeriodID.Hex())
			skipped++
			continue
		}

		// Get sales UUID - use default if not found
		salesUUID, ok := GetMapping(CollectionSales, booking.SalesID)
		if !ok {
			salesUUID = defaultSalesUUID
			usedDefaultSales++
		}

		// Parse dates
		createdAt := parseDate(booking.CreateDate)
		updatedAt := parseDate(booking.UpdateDate)
		if updatedAt.IsZero() {
			updatedAt = createdAt
		}

		dateDeposit := parseBookingDate(booking.DateDeposit)
		dateFinal := parseBookingDate(booking.DateFinal)

		// Map status
		status := mapBookingStatus(booking.Status)

		note := stringOrNil(booking.Note)

		// Calculate sum_commission_qty (count of rate items that have commission)
		sumCommissionQty := 0
		for _, rt := range booking.RateTypes {
			if !rt.Exclude {
				sumCommissionQty += toInt(rt.Qty)
			}
		}

		// Insert booking
		bookingQuery := `
			INSERT INTO bookings (
				id, invoice_id, user_id, tour_id, period_id, sales_id,
				customer_name, customer_phone, note,
				sum_room, sum_rate, sum_additional_services, sum_qty, sum_commission_qty, sum_net, sum_commission,
				discount, sum_deposit, sum_payment,
				date_deposit, date_final,
				status, created_at, updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
		`

		// Ensure sum_net > 0 (check constraint)
		sumNet := toFloat64(booking.SumNet)
		if sumNet <= 0 {
			sumNet = 1 // Set minimum value to satisfy constraint
		}

		// Handle duplicate invoice IDs
		invoiceID := booking.InvoiceID
		if count, exists := usedInvoiceIDs[invoiceID]; exists {
			// Rename duplicate invoice ID
			invoiceID = fmt.Sprintf("%s-%d", booking.InvoiceID, count+1)
			renamedInvoices++
		}
		usedInvoiceIDs[booking.InvoiceID]++

		_, err = GetPostgresDB().ExecContext(ctx, bookingQuery,
			newBookingID,
			invoiceID,
			userUUID,
			tourUUID,
			periodUUID,
			salesUUID,
			sanitizeUTF8(booking.CustomerName),
			truncateString(sanitizeUTF8(booking.CustomerPhone), 20),
			note,
			toFloat64(booking.SumRoom),
			toFloat64(booking.SumRate),
			0.0, // sum_additional_services - not in old system
			toInt(booking.SumQty),
			sumCommissionQty,
			sumNet,
			toFloat64(booking.SumCommission),
			toFloat64(booking.Discount),
			toFloat64(booking.SumDeposit),
			toFloat64(booking.SumPayment),
			dateDeposit,
			dateFinal,
			status,
			createdAt,
			updatedAt,
		)

		if err != nil {
			fmt.Printf("    WARNING: Failed to insert booking %s: %v - skipping\n", booking.InvoiceID, err)
			skipped++
			continue
		}

		// Insert room types
		if err := insertBookingRoomTypes(ctx, newBookingID, booking.RoomTypes, createdAt); err != nil {
			fmt.Printf("    WARNING: Failed to insert room types for booking %s: %v\n", booking.InvoiceID, err)
		}

		// Insert rate types
		if err := insertBookingRateTypes(ctx, newBookingID, booking.RateTypes, createdAt); err != nil {
			fmt.Printf("    WARNING: Failed to insert rate types for booking %s: %v\n", booking.InvoiceID, err)
		}

		// Insert payment slips
		if err := insertPaymentSlips(ctx, newBookingID, booking.Slips, createdAt); err != nil {
			fmt.Printf("    WARNING: Failed to insert payment slips for booking %s: %v\n", booking.InvoiceID, err)
		}

		// Insert passports
		if err := insertPassports(ctx, newBookingID, booking.Passports, createdAt); err != nil {
			fmt.Printf("    WARNING: Failed to insert passports for booking %s: %v\n", booking.InvoiceID, err)
		}

		stats.Bookings++
	}

	if skippedExisting > 0 {
		fmt.Printf("    INFO: %d bookings already exist (skipped)\n", skippedExisting)
	}
	if skipped > 0 {
		fmt.Printf("    WARNING: %d bookings were skipped due to missing references\n", skipped)
	}
	if usedDefaultSales > 0 {
		fmt.Printf("    INFO: %d bookings used default sales (original sales not found)\n", usedDefaultSales)
	}
	if usedDefaultUser > 0 {
		fmt.Printf("    INFO: %d bookings used default user (original user not found)\n", usedDefaultUser)
	}
	if renamedInvoices > 0 {
		fmt.Printf("    INFO: %d bookings had duplicate invoice IDs (renamed)\n", renamedInvoices)
	}
	fmt.Printf("    Migrated %d bookings\n", stats.Bookings)
	return nil
}

func insertBookingRoomTypes(ctx context.Context, bookingID uuid.UUID, roomTypes []MongoRoomType, createdAt time.Time) error {
	if roomTypes == nil || len(roomTypes) == 0 {
		return nil
	}

	query := `
		INSERT INTO booking_room_types (id, booking_id, title, price, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	for _, rt := range roomTypes {
		if toInt(rt.Qty) <= 0 {
			continue
		}
		_, err := GetPostgresDB().ExecContext(ctx, query,
			uuid.New(),
			bookingID,
			sanitizeUTF8(rt.Title),
			toFloat64(rt.Price),
			toInt(rt.Qty),
			createdAt,
			createdAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertBookingRateTypes(ctx context.Context, bookingID uuid.UUID, rateTypes []MongoRateType, createdAt time.Time) error {
	if rateTypes == nil || len(rateTypes) == 0 {
		return nil
	}

	query := `
		INSERT INTO booking_rate_types (id, booking_id, title, price, quantity, exclude, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	for _, rt := range rateTypes {
		if toInt(rt.Qty) <= 0 {
			continue
		}
		_, err := GetPostgresDB().ExecContext(ctx, query,
			uuid.New(),
			bookingID,
			sanitizeUTF8(rt.Title),
			toFloat64(rt.Price),
			toInt(rt.Qty),
			rt.Exclude,
			createdAt,
			createdAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertPaymentSlips(ctx context.Context, bookingID uuid.UUID, slips []MongoSlipType, defaultCreatedAt time.Time) error {
	if slips == nil || len(slips) == 0 {
		return nil
	}

	query := `
		INSERT INTO payment_slips (
			id, booking_id, slip_id, image, amount, status,
			payment_method, is_admin_created, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	for _, slip := range slips {
		amount := toFloat64(slip.Amount)
		// Skip slips with invalid amount (constraint requires amount > 0)
		if amount <= 0 {
			continue
		}
		slipID := slip.ID
		if slipID == "" {
			slipID = uuid.New().String()[:20]
		}

		createdAt := parseDate(slip.CreateDate)
		if createdAt.IsZero() {
			createdAt = defaultCreatedAt
		}

		updatedAt := parseDate(slip.UpdateDate)
		if updatedAt.IsZero() {
			updatedAt = createdAt
		}

		status := mapSlipStatus(slip.Status)
		image := stringOrNil(slip.Image)

		_, err := GetPostgresDB().ExecContext(ctx, query,
			uuid.New(),
			bookingID,
			slipID,
			image,
			amount,
			status,
			"slip",
			false,
			createdAt,
			updatedAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertPassports(ctx context.Context, bookingID uuid.UUID, passports []MongoPassportType, defaultCreatedAt time.Time) error {
	if passports == nil || len(passports) == 0 {
		return nil
	}

	query := `
		INSERT INTO passenger_passports (
			id, booking_id, image, passport_number, status,
			processed_date, error_message,
			title_th, name_th, surname_th,
			title_en, name_en, surname_en,
			gender, nationality, date_of_birth, id_number, place_of_birth, date_of_expiry,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
	`

	for _, passport := range passports {
		createdAt := parseDate(passport.CreateDate)
		if createdAt.IsZero() {
			createdAt = defaultCreatedAt
		}

		status := mapPassportStatus(passport.Status)
		image := stringOrNil(passport.Image)
		errorMessage := stringOrNil(passport.ErrorMessage)

		var processingDate *time.Time
		if passport.ProcessingDate != "" {
			pd := parseDate(passport.ProcessingDate)
			if !pd.IsZero() {
				processingDate = &pd
			}
		}

		// Personal info
		var titleTH, nameTH, surnameTH, titleEN, nameEN, surnameEN *string
		var gender, nationality, idNumber, placeOfBirth *string
		var dateOfBirth, dateOfExpiry *time.Time

		if passport.Personal != nil {
			titleTH = truncateStringPtr(stringOrNil(passport.Personal.TitleTH), 20)
			nameTH = truncateStringPtr(stringOrNil(passport.Personal.NameTH), 100)
			surnameTH = truncateStringPtr(stringOrNil(passport.Personal.SurnameTH), 100)
			titleEN = truncateStringPtr(stringOrNil(passport.Personal.TitleEN), 20)
			nameEN = truncateStringPtr(stringOrNil(passport.Personal.NameEN), 100)
			surnameEN = truncateStringPtr(stringOrNil(passport.Personal.SurnameEN), 100)
			gender = truncateStringPtr(stringOrNil(passport.Personal.Gender), 1)
			nationality = truncateStringPtr(stringOrNil(passport.Personal.Nationality), 50)
			idNumber = truncateStringPtr(stringOrNil(passport.Personal.IDNumber), 20)
			placeOfBirth = truncateStringPtr(stringOrNil(passport.Personal.PlaceOfBirth), 100)

			if passport.Personal.DateOfBirth != "" {
				dob := parseDate(passport.Personal.DateOfBirth)
				if !dob.IsZero() {
					dateOfBirth = &dob
				}
			}
			if passport.Personal.DateOfExpiry != "" {
				doe := parseDate(passport.Personal.DateOfExpiry)
				if !doe.IsZero() {
					dateOfExpiry = &doe
				}
			}
		}

		passportNumber := truncateString(passport.PassportNumber, 20)
		if passportNumber == "" {
			passportNumber = "UNKNOWN"
		}

		_, err := GetPostgresDB().ExecContext(ctx, query,
			uuid.New(),
			bookingID,
			image,
			passportNumber,
			status,
			processingDate,
			errorMessage,
			titleTH,
			nameTH,
			surnameTH,
			titleEN,
			nameEN,
			surnameEN,
			gender,
			nationality,
			dateOfBirth,
			idNumber,
			placeOfBirth,
			dateOfExpiry,
			createdAt,
			createdAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseBookingDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}

	// Try parsing with various formats
	formats := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	return time.Now()
}

func mapBookingStatus(status string) string {
	switch status {
	case "pending":
		return "pending"
	case "processing":
		return "processing"
	case "confirm", "confirmed":
		return "confirm"
	case "cancel", "cancelled":
		return "cancel"
	case "complete":
		return "complete"
	case "expired":
		return "expired"
	default:
		return "pending"
	}
}

func mapSlipStatus(status string) string {
	switch status {
	case "pending":
		return "pending"
	case "confirm", "confirmed", "approved":
		return "confirm"
	case "reject", "rejected":
		return "reject"
	case "reviewed":
		return "reviewed"
	default:
		return "pending"
	}
}

func mapPassportStatus(status string) string {
	switch status {
	case "pending":
		return "pending"
	case "processing":
		return "processing"
	case "done":
		return "done"
	default:
		return "pending"
	}
}
