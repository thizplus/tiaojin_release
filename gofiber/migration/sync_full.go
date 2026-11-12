package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SyncFullStats tracks full sync progress
type SyncFullStats struct {
	CitiesTotal      int
	CitiesInserted   int
	CitiesUpdated    int
	CitiesSkipped    int
	BannersTotal     int
	BannersInserted  int
	BannersUpdated   int
	BannersSkipped   int
	ToursTotal       int
	ToursInserted    int
	ToursUpdated     int
	ToursSkipped     int
	PeriodsTotal     int
	PeriodsInserted  int
	PeriodsUpdated   int
	PeriodsSkipped   int
	BookingsTotal    int
	BookingsInserted int
	BookingsUpdated  int
	BookingsSkipped  int
	UsersTotal       int
	UsersInserted    int
	UsersUpdated     int
	UsersSkipped     int
	SalesTotal       int
	SalesInserted    int
	SalesUpdated     int
	SalesSkipped     int
}

// SyncFull runs full sync from MongoDB to PostgreSQL
// Syncs ALL records regardless of date
func SyncFull(ctx context.Context) error {
	fmt.Println("\n============================================")
	fmt.Println("  Full Sync - All Data")
	fmt.Println("============================================")
	fmt.Println()

	LogSection("FULL SYNC - START")
	startTime := time.Now()
	stats := &SyncFullStats{}

	// Seed default system configurations first
	fmt.Println("--- [0] Seeding Default Configurations ---")
	if err := SeedDefaultConfigs(ctx); err != nil {
		LogError("Seed configs failed: %v", err)
		// Don't fail, continue with sync
	}
	if err := SeedDefaultJobs(ctx); err != nil {
		LogError("Seed jobs failed: %v", err)
		// Don't fail, continue with sync
	}
	fmt.Println()

	// Sync in order of dependencies
	syncTasks := []struct {
		name string
		fn   func(context.Context, *SyncFullStats) error
	}{
		{"Cities", SyncAllCities},
		{"Banners", SyncAllBanners},
		{"Sales", SyncAllSales},
		{"Users", SyncAllUsers},
		{"Tours", SyncAllTours},
		{"Periods", SyncAllPeriods},
		{"Bookings", SyncAllBookings},
	}

	for i, task := range syncTasks {
		fmt.Printf("--- [%d/%d] Syncing %s ---\n", i+1, len(syncTasks), task.name)
		taskStart := time.Now()

		if err := task.fn(ctx, stats); err != nil {
			LogError("Sync %s failed: %v", task.name, err)
			return fmt.Errorf("sync %s failed: %w", task.name, err)
		}

		fmt.Printf("    Completed in %v\n\n", time.Since(taskStart).Round(time.Millisecond))
	}

	// Print summary
	fmt.Println("============================================")
	fmt.Println("  Full Sync Summary")
	fmt.Println("============================================")
	fmt.Printf("  Cities:   %d total, %d inserted, %d updated, %d skipped\n",
		stats.CitiesTotal, stats.CitiesInserted, stats.CitiesUpdated, stats.CitiesSkipped)
	fmt.Printf("  Banners:  %d total, %d inserted, %d updated, %d skipped\n",
		stats.BannersTotal, stats.BannersInserted, stats.BannersUpdated, stats.BannersSkipped)
	fmt.Printf("  Sales:    %d total, %d inserted, %d updated, %d skipped\n",
		stats.SalesTotal, stats.SalesInserted, stats.SalesUpdated, stats.SalesSkipped)
	fmt.Printf("  Users:    %d total, %d inserted, %d updated, %d skipped\n",
		stats.UsersTotal, stats.UsersInserted, stats.UsersUpdated, stats.UsersSkipped)
	fmt.Printf("  Tours:    %d total, %d inserted, %d updated, %d skipped\n",
		stats.ToursTotal, stats.ToursInserted, stats.ToursUpdated, stats.ToursSkipped)
	fmt.Printf("  Periods:  %d total, %d inserted, %d updated, %d skipped\n",
		stats.PeriodsTotal, stats.PeriodsInserted, stats.PeriodsUpdated, stats.PeriodsSkipped)
	fmt.Printf("  Bookings: %d total, %d inserted, %d updated, %d skipped\n",
		stats.BookingsTotal, stats.BookingsInserted, stats.BookingsUpdated, stats.BookingsSkipped)
	fmt.Println("--------------------------------------------")
	fmt.Printf("  Total time: %v\n", time.Since(startTime).Round(time.Second))
	fmt.Println("============================================")

	LogSection("FULL SYNC - COMPLETE")
	LogSuccess("Full sync completed in %v", time.Since(startTime).Round(time.Second))

	return nil
}

// SyncAllCities syncs all cities from MongoDB
func SyncAllCities(ctx context.Context, stats *SyncFullStats) error {
	collection := GetMongoDB().Collection("tbl_city")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var cities []MongoCity
	if err := cursor.All(ctx, &cities); err != nil {
		return err
	}

	stats.CitiesTotal = len(cities)
	fmt.Printf("    Found %d cities in MongoDB\n", len(cities))

	for _, city := range cities {
		inserted, err := upsertCity(ctx, city)
		if err != nil {
			fmt.Printf("    ERROR: %s: %v\n", city.Title, err)
			stats.CitiesSkipped++
			continue
		}
		if inserted {
			stats.CitiesInserted++
		} else {
			stats.CitiesUpdated++
		}
	}

	return nil
}

// SyncAllBanners syncs all banners from MongoDB
func SyncAllBanners(ctx context.Context, stats *SyncFullStats) error {
	collection := GetMongoDB().Collection("tbl_banner")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var banners []MongoBanner
	if err := cursor.All(ctx, &banners); err != nil {
		return err
	}

	stats.BannersTotal = len(banners)
	fmt.Printf("    Found %d banners in MongoDB\n", len(banners))

	for _, banner := range banners {
		inserted, err := upsertBanner(ctx, banner)
		if err != nil {
			fmt.Printf("    ERROR: banner: %v\n", err)
			stats.BannersSkipped++
			continue
		}
		if inserted {
			stats.BannersInserted++
		} else {
			stats.BannersUpdated++
		}
	}

	return nil
}

// upsertBanner inserts or updates a banner
func upsertBanner(ctx context.Context, banner MongoBanner) (bool, error) {
	var existingID string
	var existingImg *string
	checkQuery := `SELECT id, img FROM banners WHERE img = $1 AND deleted_at IS NULL LIMIT 1`
	err := GetPostgresDB().QueryRowContext(ctx, checkQuery, banner.Image).Scan(&existingID, &existingImg)

	createdAt := parseDate(banner.CreateDate)
	updatedAt := parseDate(banner.UpdateDate)
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	status := mapBannerStatus(banner.Status)

	if err == nil {
		// Update existing - preserve existing image
		existingUUID, _ := uuid.Parse(existingID)
		SetMapping(CollectionBanners, banner.ID, existingUUID)

		updateQuery := `UPDATE banners SET status = $1, updated_at = $2 WHERE id = $3`
		_, err = GetPostgresDB().ExecContext(ctx, updateQuery, status, updatedAt, existingID)
		return false, err
	}

	// Insert new
	newID := uuid.New()
	SetMapping(CollectionBanners, banner.ID, newID)

	// Get current max display_order
	var maxOrder int
	GetPostgresDB().QueryRowContext(ctx, `SELECT COALESCE(MAX(display_order), 0) FROM banners WHERE deleted_at IS NULL`).Scan(&maxOrder)
	maxOrder++

	title := fmt.Sprintf("Banner %d", maxOrder)

	insertQuery := `
		INSERT INTO banners (id, title, img, link_url, alt, display_order, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err = GetPostgresDB().ExecContext(ctx, insertQuery,
		newID, title, banner.Image, "", "", maxOrder, status, createdAt, updatedAt)

	return true, err
}

// SyncAllSales syncs all sales from MongoDB
func SyncAllSales(ctx context.Context, stats *SyncFullStats) error {
	collection := GetMongoDB().Collection("tbl_sales")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var sales []MongoSales
	if err := cursor.All(ctx, &sales); err != nil {
		return err
	}

	stats.SalesTotal = len(sales)
	fmt.Printf("    Found %d sales in MongoDB\n", len(sales))

	for _, s := range sales {
		inserted, err := upsertSales(ctx, s)
		if err != nil {
			fmt.Printf("    ERROR: %s: %v\n", s.Code, err)
			stats.SalesSkipped++
			continue
		}
		if inserted {
			stats.SalesInserted++
		} else {
			stats.SalesUpdated++
		}
	}

	return nil
}

// SyncAllUsers syncs all users from MongoDB
func SyncAllUsers(ctx context.Context, stats *SyncFullStats) error {
	collection := GetMongoDB().Collection("tbl_users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var users []MongoUser
	if err := cursor.All(ctx, &users); err != nil {
		return err
	}

	stats.UsersTotal = len(users)
	fmt.Printf("    Found %d users in MongoDB\n", len(users))

	for _, user := range users {
		inserted, err := upsertUser(ctx, user)
		if err != nil {
			stats.UsersSkipped++
			continue
		}
		if inserted {
			stats.UsersInserted++
		} else {
			stats.UsersUpdated++
		}
	}

	return nil
}

// SyncAllTours syncs all tours from MongoDB
func SyncAllTours(ctx context.Context, stats *SyncFullStats) error {
	collection := GetMongoDB().Collection("tbl_tours")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var tours []MongoTour
	if err := cursor.All(ctx, &tours); err != nil {
		return err
	}

	stats.ToursTotal = len(tours)
	fmt.Printf("    Found %d tours in MongoDB\n", len(tours))

	cityMap, _ := buildCityMap(ctx)

	for _, tour := range tours {
		inserted, err := upsertTour(ctx, tour, cityMap)
		if err != nil {
			fmt.Printf("    ERROR: %s: %v\n", tour.Code, err)
			stats.ToursSkipped++
			continue
		}
		if inserted {
			stats.ToursInserted++
		} else {
			stats.ToursUpdated++
		}
	}

	return nil
}

// SyncAllPeriods syncs all periods from MongoDB
func SyncAllPeriods(ctx context.Context, stats *SyncFullStats) error {
	collection := GetMongoDB().Collection("tbl_period")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var periods []MongoPeriod
	if err := cursor.All(ctx, &periods); err != nil {
		return err
	}

	stats.PeriodsTotal = len(periods)
	fmt.Printf("    Found %d periods in MongoDB\n", len(periods))

	for _, period := range periods {
		inserted, err := upsertPeriodFull(ctx, period)
		if err != nil {
			stats.PeriodsSkipped++
			continue
		}
		if inserted {
			stats.PeriodsInserted++
		} else {
			stats.PeriodsUpdated++
		}
	}

	return nil
}

// SyncAllBookings syncs all bookings from MongoDB
func SyncAllBookings(ctx context.Context, stats *SyncFullStats) error {
	collection := GetMongoDB().Collection("tbl_booking")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)

	var bookings []MongoBooking
	if err := cursor.All(ctx, &bookings); err != nil {
		return err
	}

	stats.BookingsTotal = len(bookings)
	fmt.Printf("    Found %d bookings in MongoDB\n", len(bookings))

	// Track first few errors for debugging
	errorCount := 0
	for _, booking := range bookings {
		inserted, err := upsertBookingFull(ctx, booking)
		if err != nil {
			if errorCount < 5 {
				fmt.Printf("    ERROR booking %s: %v\n", booking.InvoiceID, err)
				errorCount++
			}
			stats.BookingsSkipped++
			continue
		}
		if inserted {
			stats.BookingsInserted++
		} else {
			stats.BookingsUpdated++
		}
	}

	return nil
}

// upsertSales inserts or updates a sales record
func upsertSales(ctx context.Context, sales MongoSales) (bool, error) {
	var existingID string
	checkQuery := `SELECT id FROM sales WHERE code = $1 AND deleted_at IS NULL LIMIT 1`
	err := GetPostgresDB().QueryRowContext(ctx, checkQuery, sales.Code).Scan(&existingID)

	createdAt := parseDate(sales.CreateDate)
	updatedAt := parseDate(sales.UpdateDate)
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}

	// Map status to active boolean
	active := sales.Status == "active" || sales.Status == "publish" || sales.Status == ""

	firstName := sales.FirstName
	lastName := sales.LastName
	nickName := sales.NickName

	if err == nil {
		// Update existing
		existingUUID, _ := uuid.Parse(existingID)
		SetMapping(CollectionSales, sales.ID, existingUUID)

		updateQuery := `
			UPDATE sales SET first_name = $1, last_name = $2, nick_name = $3, active = $4, updated_at = $5
			WHERE id = $6
		`
		_, err = GetPostgresDB().ExecContext(ctx, updateQuery,
			firstName, lastName, nickName, active, updatedAt, existingID)
		return false, err
	}

	// Insert new
	newID := uuid.New()
	SetMapping(CollectionSales, sales.ID, newID)

	insertQuery := `
		INSERT INTO sales (id, code, first_name, last_name, nick_name, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = GetPostgresDB().ExecContext(ctx, insertQuery,
		newID, sales.Code, firstName, lastName, nickName, active, createdAt, updatedAt)

	return true, err
}

// upsertUser inserts or updates a user record
func upsertUser(ctx context.Context, user MongoUser) (bool, error) {
	username := normalizeUsername(user.Email)

	var existingID string
	checkQuery := `SELECT id FROM users WHERE username = $1 LIMIT 1`
	err := GetPostgresDB().QueryRowContext(ctx, checkQuery, username).Scan(&existingID)

	createdAt := parseDate(user.CreateDate)
	updatedAt := parseDate(user.UpdateDate)
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	role := mapUserRole(user.Role)
	status := mapUserStatus(user.Status)

	if err == nil {
		// Update existing - only update role and status
		updateQuery := `
			UPDATE users SET role = $1, status = $2, updated_at = $3
			WHERE id = $4
		`
		_, err = GetPostgresDB().ExecContext(ctx, updateQuery,
			role, status, updatedAt, existingID)
		return false, err
	}

	// Insert new
	var email *string
	if isValidEmail(user.Email) {
		normalizedEmail := normalizeEmail(user.Email)
		email = &normalizedEmail
	}

	insertQuery := `
		INSERT INTO users (id, username, email, password, role, status, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
	`
	_, err = GetPostgresDB().ExecContext(ctx, insertQuery,
		username, email, user.Password, role, status, createdAt, updatedAt)

	return true, err
}

// upsertPeriodFull inserts or updates a period with full data
func upsertPeriodFull(ctx context.Context, period MongoPeriod) (bool, error) {
	// Get tour UUID by looking up tour ObjectID in mapping
	tourOID, err := primitive.ObjectIDFromHex(period.TourID)
	if err != nil {
		return false, fmt.Errorf("invalid tour ID: %s", period.TourID)
	}

	// Look up tour by MongoDB ID mapping
	pgTourUUID, ok := GetMapping(CollectionTours, tourOID)
	if !ok {
		// Tour not found, skip this period
		return false, fmt.Errorf("tour not found for period")
	}
	tourUUID := pgTourUUID.String()

	startDate := parsePeriodDate(period.StartDate)
	endDate := parsePeriodDate(period.EndDate)

	// Convert to Thailand timezone (UTC+7) before extracting date
	// MongoDB stores "27 Mar midnight Thailand" as "26 Mar 17:00 UTC"
	// We need to convert to Thailand time first, then extract the date
	thailandTZ := time.FixedZone("Asia/Bangkok", 7*60*60) // UTC+7
	startDateTH := startDate.In(thailandTZ)
	endDateTH := endDate.In(thailandTZ)

	// Normalize dates to DATE only (strip time component) for comparison
	// PostgreSQL stores as DATE type, so we must compare DATE to DATE
	startDateOnly := time.Date(startDateTH.Year(), startDateTH.Month(), startDateTH.Day(), 0, 0, 0, 0, time.UTC)
	endDateOnly := time.Date(endDateTH.Year(), endDateTH.Month(), endDateTH.Day(), 0, 0, 0, 0, time.UTC)

	// DEBUG: Log for specific tour
	if period.TourID == "69a128444dc9528d3f125eb9" {
		fmt.Printf("\n[DEBUG BCCKG25-SL] MongoDB Period ID: %s\n", period.ID.Hex())
		fmt.Printf("[DEBUG BCCKG25-SL] tour_id (mongo): %s -> pgTourUUID: %s\n", period.TourID, tourUUID)
		fmt.Printf("[DEBUG BCCKG25-SL] startDate UTC: %v -> Thailand: %v -> normalized: %v\n", startDate, startDateTH, startDateOnly)
		fmt.Printf("[DEBUG BCCKG25-SL] endDate UTC: %v -> Thailand: %v -> normalized: %v\n", endDate, endDateTH, endDateOnly)
		fmt.Printf("[DEBUG BCCKG25-SL] createDate: %v\n", period.CreateDate)
	}

	var existingID string
	checkQuery := `
		SELECT id FROM periods
		WHERE tour_id = $1 AND start_date = $2::date AND end_date = $3::date AND deleted_at IS NULL
		LIMIT 1
	`
	err = GetPostgresDB().QueryRowContext(ctx, checkQuery, tourUUID, startDateOnly, endDateOnly).Scan(&existingID)

	// DEBUG: Log check result
	if period.TourID == "69a128444dc9528d3f125eb9" {
		if err == nil {
			fmt.Printf("[DEBUG BCCKG25-SL] FOUND existing period: %s -> UPDATE\n", existingID)
		} else {
			fmt.Printf("[DEBUG BCCKG25-SL] NOT FOUND -> INSERT NEW (err: %v)\n", err)
		}
	}

	createdAt := parseDate(period.CreateDate)
	updatedAt := parseDate(period.UpdateDate)
	if updatedAt.IsZero() {
		updatedAt = createdAt
	}
	status := mapPeriodStatus(period.Status)

	if err == nil {
		// Update existing - only update seat and status, preserve other data
		existingUUID, _ := uuid.Parse(existingID)
		SetMapping(CollectionPeriods, period.ID, existingUUID)

		updateQuery := `
			UPDATE periods SET seat = $1, status = $2, updated_at = $3
			WHERE id = $4
		`
		_, err = GetPostgresDB().ExecContext(ctx, updateQuery,
			period.Seat, status, updatedAt, existingID)
		return false, err
	}

	// Insert new with full data
	newID := uuid.New()
	SetMapping(CollectionPeriods, period.ID, newID)

	// Convert sector arrays to JSON
	sector1JSON := convertSectorsToJSON(period.Sector1)
	sector2JSON := convertSectorsToJSON(period.Sector2)

	// Convert commission items to JSON
	commissionItemsJSON := convertCommissionItemsToJSON(period.Commission.Item)

	// Convert files to JSON
	filesJSON := convertFilesToJSON(period.Files)

	// Map rate type
	rateType := mapRateType(period.Rate.Type)
	commissionType := mapRateType(period.Commission.Type)

	// Default deposit type
	// Backend validation: oneof=percent per_unit
	depositType := "percent"
	if period.Deposit > 100 {
		depositType = "per_unit"
	}

	note := stringOrNil(period.Note)
	airline := truncateString(period.Airline, 100)
	airport := truncateString(period.Airport, 10)

	// JSONB columns (rate_items, room_types, additional_services) are empty
	// Actual data goes to separate tables: period_rate_items, period_room_types
	insertQuery := `
		INSERT INTO periods (
			id, tour_id, start_date, end_date, seat, deposit, deposit_type,
			airline, airport, note, sector_1, sector_2,
			rate_type, rate_items, commission_type, commission_items,
			room_types, additional_services, files, status, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`

	_, err = GetPostgresDB().ExecContext(ctx, insertQuery,
		newID,
		tourUUID,
		startDateOnly,
		endDateOnly,
		period.Seat,
		period.Deposit,
		depositType,
		airline,
		airport,
		note,
		sector1JSON,
		sector2JSON,
		rateType,
		"[]", // rate_items - empty JSONB, data goes to period_rate_items table
		commissionType,
		commissionItemsJSON,
		"[]", // room_types - empty JSONB, data goes to period_room_types table
		"[]", // additional_services - empty by default (not in MongoDB)
		filesJSON,
		status,
		createdAt,
		updatedAt,
	)

	if err != nil {
		return false, err
	}

	// Insert rate items into separate table
	if err := insertPeriodRateItems(ctx, newID, period.Rate.Item, createdAt); err != nil {
		fmt.Printf("    WARNING: Failed to insert rate items for period %s: %v\n", period.ID.Hex(), err)
	}

	// Insert room types into separate table
	if err := insertPeriodRoomTypes(ctx, newID, period.SingleStay, createdAt); err != nil {
		fmt.Printf("    WARNING: Failed to insert room types for period %s: %v\n", period.ID.Hex(), err)
	}

	return true, nil
}

// upsertBookingFull inserts or updates a booking with full data
func upsertBookingFull(ctx context.Context, booking MongoBooking) (bool, error) {
	var existingID string
	checkQuery := `SELECT id FROM bookings WHERE invoice_id = $1 AND deleted_at IS NULL LIMIT 1`
	err := GetPostgresDB().QueryRowContext(ctx, checkQuery, booking.InvoiceID).Scan(&existingID)

	if err == nil {
		// Booking exists - update status only
		status := mapBookingStatus(booking.Status)
		updateQuery := `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`
		_, err = GetPostgresDB().ExecContext(ctx, updateQuery, status, existingID)

		// Set mapping for existing booking
		existingUUID, _ := uuid.Parse(existingID)
		SetMapping(CollectionBookings, booking.ID, existingUUID)
		return false, err
	}

	// New booking - insert with full data
	newBookingID := uuid.New()
	SetMapping(CollectionBookings, booking.ID, newBookingID)

	// Get default user UUID
	var defaultUserUUID uuid.UUID
	err = GetPostgresDB().QueryRowContext(ctx, "SELECT id FROM users LIMIT 1").Scan(&defaultUserUUID)
	if err != nil {
		return false, fmt.Errorf("failed to get default user: %w", err)
	}

	// Get default sales UUID
	var defaultSalesUUID uuid.UUID
	err = GetPostgresDB().QueryRowContext(ctx, "SELECT id FROM sales LIMIT 1").Scan(&defaultSalesUUID)
	if err != nil {
		return false, fmt.Errorf("failed to get default sales: %w", err)
	}

	// Get user UUID from profile mapping - use default if not found
	userUUID, ok := GetMapping(CollectionUsers, booking.ProfileID)
	if !ok {
		userUUID = defaultUserUUID
	}

	// Get tour UUID - try mapping first, then lookup by code
	tourUUID, ok := GetMapping(CollectionTours, booking.TourID)
	if !ok {
		// Try to find by tour code
		var tourID string
		err = GetPostgresDB().QueryRowContext(ctx,
			"SELECT id FROM tours WHERE code = $1 AND deleted_at IS NULL LIMIT 1",
			booking.TourCode).Scan(&tourID)
		if err != nil {
			return false, fmt.Errorf("tour not found: %s", booking.TourCode)
		}
		tourUUID, _ = uuid.Parse(tourID)
	}

	// Get period UUID - try mapping first, then lookup by tour + dates
	periodUUID, ok := GetMapping(CollectionPeriods, booking.PeriodID)
	if !ok {
		// Try to find by tour_id + period dates
		startDate := parseDate(booking.PeriodStart)
		endDate := parseDate(booking.PeriodEnd)
		var periodID string
		err = GetPostgresDB().QueryRowContext(ctx,
			"SELECT id FROM periods WHERE tour_id = $1 AND start_date = $2 AND end_date = $3 AND deleted_at IS NULL LIMIT 1",
			tourUUID, startDate, endDate).Scan(&periodID)
		if err != nil {
			return false, fmt.Errorf("period not found for tour %s", tourUUID)
		}
		periodUUID, _ = uuid.Parse(periodID)
	}

	// Get sales UUID - use default if not found
	salesUUID, ok := GetMapping(CollectionSales, booking.SalesID)
	if !ok {
		salesUUID = defaultSalesUUID
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

	// Calculate sum_commission_qty
	sumCommissionQty := 0
	for _, rt := range booking.RateTypes {
		if !rt.Exclude {
			sumCommissionQty += toInt(rt.Qty)
		}
	}

	// Ensure sum_net > 0 (check constraint)
	sumNet := toFloat64(booking.SumNet)
	if sumNet <= 0 {
		sumNet = 1
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

	_, err = GetPostgresDB().ExecContext(ctx, bookingQuery,
		newBookingID,
		booking.InvoiceID,
		userUUID,
		tourUUID,
		periodUUID,
		salesUUID,
		sanitizeUTF8(booking.CustomerName),
		truncateString(sanitizeUTF8(booking.CustomerPhone), 20),
		note,
		toFloat64(booking.SumRoom),
		toFloat64(booking.SumRate),
		0.0, // sum_additional_services
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
		return false, fmt.Errorf("failed to insert booking: %w", err)
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

	return true, nil
}

func mapSalesStatus(status string) string {
	switch status {
	case "active", "publish":
		return "active"
	case "inactive", "draft":
		return "inactive"
	default:
		return "active"
	}
}
