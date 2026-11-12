package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoFlightSector represents flight sector from MongoDB
type MongoFlightSector struct {
	Date   string `bson:"date"`
	Flight string `bson:"flight"`
	Route  string `bson:"route"`
}

// MongoRateItem represents rate item from MongoDB (flexible types)
type MongoRateItem struct {
	Title   interface{} `bson:"title"`   // Can be string or int
	Regular interface{} `bson:"regular"` // Can be string or float
	Sale    interface{} `bson:"sale"`    // Can be string or float
	Exclude bool        `bson:"exclude"` // Exclude from commission calculation
}

// MongoRate represents rate structure from MongoDB
type MongoRate struct {
	Type string          `bson:"type"`
	Item []MongoRateItem `bson:"item"`
}

// MongoCommissionItem represents commission item from MongoDB
type MongoCommissionItem struct {
	Regular interface{} `bson:"regular"` // Can be string or float
	Sale    interface{} `bson:"sale"`    // Can be string or float
}

// MongoCommission represents commission structure from MongoDB
type MongoCommission struct {
	Type string                `bson:"type"`
	Item []MongoCommissionItem `bson:"item"`
}

// MongoFile represents file attachment from MongoDB
type MongoFile struct {
	Name string `bson:"name"`
	Type string `bson:"type"`
	URL  string `bson:"url"`
}

// MongoPeriod represents period document from MongoDB
type MongoPeriod struct {
	ID         primitive.ObjectID  `bson:"_id"`
	TourID     string              `bson:"tour"` // Can be string or ObjectID
	StartDate  interface{}         `bson:"start_date"`
	EndDate    interface{}         `bson:"end_date"`
	Seat       int                 `bson:"seat"`
	Booking    int                 `bson:"booking"`
	Deposit    float64             `bson:"deposit"`
	SingleStay float64             `bson:"single_stay"`
	Airline    string              `bson:"airline"`
	Airport    string              `bson:"airport"`
	Note       string              `bson:"note"`
	Sector1    []MongoFlightSector `bson:"sector_1"`
	Sector2    []MongoFlightSector `bson:"sector_2"`
	Rate       MongoRate           `bson:"rate"`
	Commission MongoCommission     `bson:"commission"`
	Files      []MongoFile         `bson:"files"`
	Status     string              `bson:"status"`
	CreateDate string              `bson:"createDate"`
	UpdateDate string              `bson:"updateDate"`
}

// MigratePeriods migrates periods from MongoDB to PostgreSQL
func MigratePeriods(ctx context.Context, stats *MigrationStats) error {
	collection := GetMongoDB().Collection("tbl_period")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	var periods []MongoPeriod
	if err := cursor.All(ctx, &periods); err != nil {
		return fmt.Errorf("failed to decode periods: %w", err)
	}

	fmt.Printf("    Found %d periods in MongoDB\n", len(periods))

	skippedCount := 0
	warningCount := 0

	for _, period := range periods {
		// Get tour UUID from mapping
		tourOID, err := primitive.ObjectIDFromHex(period.TourID)
		if err != nil {
			fmt.Printf("    WARNING: Invalid tour ID %s for period %s - skipping\n", period.TourID, period.ID.Hex())
			warningCount++
			continue
		}

		tourUUID, ok := GetMapping(CollectionTours, tourOID)
		if !ok {
			fmt.Printf("    WARNING: Tour not found for period %s (tour_id: %s) - skipping\n", period.ID.Hex(), period.TourID)
			warningCount++
			continue
		}

		// Parse dates
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

		// Check if period already exists by tour_id + start_date + end_date
		var existingID string
		checkQuery := `SELECT id FROM periods WHERE tour_id = $1 AND start_date = $2::date AND end_date = $3::date AND deleted_at IS NULL LIMIT 1`
		err = GetPostgresDB().QueryRowContext(ctx, checkQuery, tourUUID, startDateOnly, endDateOnly).Scan(&existingID)

		// DEBUG for BCCKG25-SL tour
		if period.TourID == "69a128444dc9528d3f125eb9" {
			fmt.Printf("\n[DEBUG MIGRATE] Period ID: %s\n", period.ID.Hex())
			fmt.Printf("[DEBUG MIGRATE] tourUUID: %s\n", tourUUID)
			fmt.Printf("[DEBUG MIGRATE] startDateOnly: %v\n", startDateOnly)
			fmt.Printf("[DEBUG MIGRATE] endDateOnly: %v\n", endDateOnly)
			if err == nil {
				fmt.Printf("[DEBUG MIGRATE] FOUND existing: %s -> SKIP\n", existingID)
			} else {
				fmt.Printf("[DEBUG MIGRATE] NOT FOUND (err: %v) -> INSERT\n", err)
			}
		}

		if err == nil {
			// Period already exists, set mapping and skip
			existingUUID, _ := uuid.Parse(existingID)
			SetMapping(CollectionPeriods, period.ID, existingUUID)
			skippedCount++
			continue
		}

		newID := uuid.New()
		SetMapping(CollectionPeriods, period.ID, newID)

		createdAt := parseDate(period.CreateDate)
		updatedAt := parseDate(period.UpdateDate)
		if updatedAt.IsZero() {
			updatedAt = createdAt
		}

		// Convert sector arrays to JSON
		sector1JSON := convertSectorsToJSON(period.Sector1)
		sector2JSON := convertSectorsToJSON(period.Sector2)

		// Convert commission items to JSON
		commissionItemsJSON := convertCommissionItemsToJSON(period.Commission.Item)

		// Convert files to JSON
		filesJSON := convertFilesToJSON(period.Files)

		// Map status
		status := mapPeriodStatus(period.Status)

		// Map rate type
		rateType := mapRateType(period.Rate.Type)
		commissionType := mapRateType(period.Commission.Type)

		// Default deposit type to percent if deposit <= 100, otherwise per_unit
		// Backend validation: oneof=percent per_unit
		depositType := "percent"
		if period.Deposit > 100 {
			depositType = "per_unit"
		}

		// Insert into PostgreSQL
		// JSONB columns (rate_items, room_types, additional_services) are empty
		// Actual data goes to separate tables: period_rate_items, period_room_types
		query := `
			INSERT INTO periods (
				id, tour_id, start_date, end_date, seat, deposit, deposit_type,
				airline, airport, note, sector_1, sector_2,
				rate_type, rate_items, commission_type, commission_items,
				room_types, additional_services, files, status, created_at, updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
		`

		note := stringOrNil(period.Note)
		airline := truncateString(period.Airline, 100)
		airport := truncateString(period.Airport, 10)

		_, err = GetPostgresDB().ExecContext(ctx, query,
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
			return fmt.Errorf("failed to insert period %s: %w", period.ID.Hex(), err)
		}

		// Insert rate items into separate table
		if err := insertPeriodRateItems(ctx, newID, period.Rate.Item, createdAt); err != nil {
			fmt.Printf("    WARNING: Failed to insert rate items for period %s: %v\n", period.ID.Hex(), err)
		}

		// Insert room types into separate table
		if err := insertPeriodRoomTypes(ctx, newID, period.SingleStay, createdAt); err != nil {
			fmt.Printf("    WARNING: Failed to insert room types for period %s: %v\n", period.ID.Hex(), err)
		}

		stats.Periods++
	}

	if skippedCount > 0 {
		fmt.Printf("    INFO: %d periods already exist (skipped)\n", skippedCount)
	}
	if warningCount > 0 {
		fmt.Printf("    WARNING: %d periods skipped due to missing tour reference\n", warningCount)
	}
	fmt.Printf("    Migrated %d periods\n", stats.Periods)
	return nil
}

func parsePeriodDate(dateValue interface{}) time.Time {
	if dateValue == nil {
		return time.Now()
	}

	switch v := dateValue.(type) {
	case string:
		return parseDate(v)
	case primitive.DateTime:
		return v.Time()
	case time.Time:
		return v
	case map[string]interface{}:
		// Handle {$date: "..."} format
		if dateStr, ok := v["$date"].(string); ok {
			return parseDate(dateStr)
		}
	case primitive.M:
		if dateVal, ok := v["$date"]; ok {
			switch dv := dateVal.(type) {
			case string:
				return parseDate(dv)
			case primitive.DateTime:
				return dv.Time()
			}
		}
	}

	return time.Now()
}

func convertSectorsToJSON(sectors []MongoFlightSector) string {
	if sectors == nil || len(sectors) == 0 {
		return "[]"
	}

	type FlightSector struct {
		Date   string `json:"date"`
		Flight string `json:"flight"`
		Route  string `json:"route"`
	}

	result := make([]FlightSector, len(sectors))
	for i, s := range sectors {
		result[i] = FlightSector{
			Date:   convertDateFormat(sanitizeUTF8(s.Date)),
			Flight: sanitizeUTF8(s.Flight),
			Route:  sanitizeUTF8(s.Route),
		}
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}

// convertDateFormat converts date from DD/MM/YYYY to YYYY-MM-DD
// MongoDB: "02/06/2026" -> PostgreSQL: "2026-06-02"
func convertDateFormat(date string) string {
	if date == "" {
		return ""
	}

	// Try to parse DD/MM/YYYY format
	parts := strings.Split(date, "/")
	if len(parts) == 3 {
		day := parts[0]
		month := parts[1]
		year := parts[2]

		// Validate and pad if needed
		if len(day) == 1 {
			day = "0" + day
		}
		if len(month) == 1 {
			month = "0" + month
		}

		return fmt.Sprintf("%s-%s-%s", year, month, day)
	}

	// If already in YYYY-MM-DD format or other format, return as-is
	return date
}

func convertRateItemsToJSON(items []MongoRateItem) string {
	if items == nil || len(items) == 0 {
		return "[]"
	}

	// Match model: ID, Title, Regular, Sale, Exclude (no SortOrder)
	type RateItem struct {
		ID      string  `json:"id"`
		Title   string  `json:"title"`
		Regular float64 `json:"regular"`
		Sale    float64 `json:"sale"`
		Exclude bool    `json:"exclude"`
	}

	result := make([]RateItem, len(items))
	for i, item := range items {
		result[i] = RateItem{
			ID:      uuid.New().String(),
			Title:   toString(item.Title),
			Regular: toFloat64(item.Regular),
			Sale:    toFloat64(item.Sale),
			Exclude: item.Exclude,
		}
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}

func convertCommissionItemsToJSON(items []MongoCommissionItem) string {
	if items == nil || len(items) == 0 {
		return "[]"
	}

	type CommissionItem struct {
		Regular float64 `json:"regular"`
		Sale    float64 `json:"sale"`
	}

	result := make([]CommissionItem, len(items))
	for i, item := range items {
		result[i] = CommissionItem{
			Regular: toFloat64(item.Regular),
			Sale:    toFloat64(item.Sale),
		}
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}

func convertFilesToJSON(files []MongoFile) string {
	if files == nil || len(files) == 0 {
		return "[]"
	}

	type FileItem struct {
		Name string `json:"name"`
		Type string `json:"type"`
		URL  string `json:"url"`
	}

	result := make([]FileItem, len(files))
	for i, f := range files {
		result[i] = FileItem{
			Name: sanitizeUTF8(f.Name),
			Type: sanitizeUTF8(f.Type),
			URL:  sanitizeUTF8(f.URL),
		}
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}

// convertSingleStayToRoomTypesJSON converts MongoDB single_stay to PostgreSQL room_types JSONB
// Creates standard 4 room types: TWIN, DOUBLE, TRIPLE, SINGLE
// MongoDB: single_stay: 5500
// PostgreSQL: room_types: [
//
//	{"id": "uuid", "title": "TWIN", "price": 0},
//	{"id": "uuid", "title": "DOUBLE", "price": 0},
//	{"id": "uuid", "title": "TRIPLE", "price": 0},
//	{"id": "uuid", "title": "SINGLE", "price": 5500}
//
// ]
func convertSingleStayToRoomTypesJSON(singleStay float64) string {
	// Match model: ID, Title, Price (no SortOrder)
	type RoomType struct {
		ID    string  `json:"id"`
		Title string  `json:"title"`
		Price float64 `json:"price"`
	}

	// Standard 4 room types for every period
	roomTypes := []RoomType{
		{
			ID:    uuid.New().String(),
			Title: "TWIN",
			Price: 0,
		},
		{
			ID:    uuid.New().String(),
			Title: "DOUBLE",
			Price: 0,
		},
		{
			ID:    uuid.New().String(),
			Title: "TRIPLE",
			Price: 0,
		},
		{
			ID:    uuid.New().String(),
			Title: "SINGLE",
			Price: singleStay,
		},
	}

	jsonBytes, err := json.Marshal(roomTypes)
	if err != nil {
		return "[]"
	}
	return string(jsonBytes)
}

// insertPeriodRateItems inserts rate items into period_rate_items table
func insertPeriodRateItems(ctx context.Context, periodID uuid.UUID, items []MongoRateItem, createdAt time.Time) error {
	if items == nil || len(items) == 0 {
		return nil
	}

	query := `
		INSERT INTO period_rate_items (id, period_id, title, regular, sale, exclude, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	for i, item := range items {
		title := toString(item.Title)
		if title == "" {
			continue
		}
		_, err := GetPostgresDB().ExecContext(ctx, query,
			uuid.New(),
			periodID,
			title,
			toFloat64(item.Regular),
			toFloat64(item.Sale),
			item.Exclude,
			i, // sort_order based on position
			createdAt,
			createdAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// insertPeriodRoomTypes inserts room types into period_room_types table
// Creates standard 4 room types: TWIN, DOUBLE, TRIPLE, SINGLE
func insertPeriodRoomTypes(ctx context.Context, periodID uuid.UUID, singleStay float64, createdAt time.Time) error {
	query := `
		INSERT INTO period_room_types (id, period_id, title, price, sort_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	// Standard 4 room types
	roomTypes := []struct {
		Title string
		Price float64
	}{
		{"TWIN", 0},
		{"DOUBLE", 0},
		{"TRIPLE", 0},
		{"SINGLE", singleStay},
	}

	for i, rt := range roomTypes {
		_, err := GetPostgresDB().ExecContext(ctx, query,
			uuid.New(),
			periodID,
			rt.Title,
			rt.Price,
			i, // sort_order based on position
			createdAt,
			createdAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func mapPeriodStatus(status string) string {
	switch status {
	case "publish", "active":
		return "publish"
	case "draft", "inactive", "":
		return "draft"
	default:
		return "draft"
	}
}

func mapRateType(rateType string) string {
	switch rateType {
	case "sale":
		return "sale"
	case "regular", "":
		return "regular"
	default:
		return "regular"
	}
}

// toString converts interface{} to string, handling various types
func toString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return sanitizeUTF8(val)
	case int:
		return fmt.Sprintf("%d", val)
	case int32:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%.0f", val)
	case float32:
		return fmt.Sprintf("%.0f", val)
	default:
		return sanitizeUTF8(fmt.Sprintf("%v", val))
	}
}

// toFloat64 converts interface{} to float64, handling strings and various number types
func toFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int32:
		return float64(val)
	case int64:
		return float64(val)
	case string:
		var f float64
		fmt.Sscanf(val, "%f", &f)
		return f
	default:
		return 0
	}
}

// toInt converts interface{} to int, handling strings and various number types
func toInt(v interface{}) int {
	if v == nil {
		return 0
	}

	switch val := v.(type) {
	case int:
		return val
	case int32:
		return int(val)
	case int64:
		return int(val)
	case float64:
		return int(val)
	case float32:
		return int(val)
	case string:
		var i int
		fmt.Sscanf(val, "%d", &i)
		return i
	default:
		return 0
	}
}
