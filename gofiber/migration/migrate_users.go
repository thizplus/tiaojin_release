package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// emailRegex สำหรับ validate email format
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// MongoUserProfile represents the embedded profile in user document
type MongoUserProfile struct {
	Name    string `bson:"name"`
	TaxID   string `bson:"taxId"`
	Phone   string `bson:"tel"`
	Address string `bson:"address"`
	Province string `bson:"province"`
	Zipcode string `bson:"zipcode"`
	ATTA    string `bson:"atta"`
	DBD     string `bson:"dbd"`
	PP20    string `bson:"pp20"`
	Logo    string `bson:"logo"`
}

// MongoUser represents user document from MongoDB
type MongoUser struct {
	ID         primitive.ObjectID `bson:"_id"`
	Email      string             `bson:"email"`
	Password   string             `bson:"password"`
	Role       string             `bson:"role"`
	Status     string             `bson:"status"`
	Profile    *MongoUserProfile  `bson:"profile"`
	CreateDate string             `bson:"createDate"`
	UpdateDate string             `bson:"updateDate"`
}

// MigrateUsers migrates users and company profiles from MongoDB to PostgreSQL
func MigrateUsers(ctx context.Context, stats *MigrationStats) error {
	collection := GetMongoDB().Collection("tbl_users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return fmt.Errorf("failed to query MongoDB: %w", err)
	}
	defer cursor.Close(ctx)

	var users []MongoUser
	if err := cursor.All(ctx, &users); err != nil {
		return fmt.Errorf("failed to decode users: %w", err)
	}

	fmt.Printf("    Found %d users in MongoDB\n", len(users))

	// Track usernames to detect duplicates
	normalizedUsernames := make(map[string]bool)
	duplicateCount := 0
	skippedCount := 0
	emailCount := 0

	for _, user := range users {
		// MongoDB email → PostgreSQL username (ใช้เป็น login identifier)
		username := normalizeUsername(user.Email)

		// Check for duplicate usernames in this batch
		if normalizedUsernames[username] {
			fmt.Printf("    WARNING: Duplicate username found after normalization: %s (original: %s) - skipping\n", username, user.Email)
			duplicateCount++
			continue
		}
		normalizedUsernames[username] = true

		// Check if user already exists in PostgreSQL
		var existingID string
		checkQuery := `SELECT id FROM users WHERE username = $1 LIMIT 1`
		err := GetPostgresDB().QueryRowContext(ctx, checkQuery, username).Scan(&existingID)
		if err == nil {
			// User already exists, skip
			skippedCount++
			continue
		}

		newUserID := uuid.New()
		SetMapping(CollectionUsers, user.ID, newUserID)

		// Parse dates
		createdAt := parseDate(user.CreateDate)
		updatedAt := parseDate(user.UpdateDate)
		if updatedAt.IsZero() {
			updatedAt = createdAt
		}

		// Map role
		role := mapUserRole(user.Role)

		// Map status
		status := mapUserStatus(user.Status)

		// ตรวจสอบว่า MongoDB email เป็น email format จริงหรือไม่
		// ถ้าใช่ → copy ไปใส่ PostgreSQL email ด้วย
		// ถ้าไม่ใช่ → PostgreSQL email = NULL
		var email *string
		if isValidEmail(user.Email) {
			normalizedEmail := normalizeEmail(user.Email)
			email = &normalizedEmail
			emailCount++
		}

		// Insert user into PostgreSQL with username and optional email
		userQuery := `
			INSERT INTO users (id, username, email, password, role, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		_, execErr := GetPostgresDB().ExecContext(ctx, userQuery,
			newUserID,
			username,
			email, // อาจเป็น NULL ถ้าไม่ใช่ email format
			user.Password,
			role,
			status,
			createdAt,
			updatedAt,
		)

		if execErr != nil {
			return fmt.Errorf("failed to insert user %s: %w", username, execErr)
		}

		// Insert company profile if exists
		if user.Profile != nil {
			profileID := uuid.New()

			// Handle empty or null values
			logo := stringOrNil(user.Profile.Logo)

			// Truncate fields to fit database constraints
			taxID := truncateString(user.Profile.TaxID, 13)
			phone := truncateString(user.Profile.Phone, 20)
			province := truncateString(user.Profile.Province, 100)
			zipcode := truncateString(user.Profile.Zipcode, 5)
			name := truncateString(user.Profile.Name, 200)

			profileQuery := `
				INSERT INTO company_profiles (id, user_id, name, tax_id, phone, address, province, zipcode, atta, dbd, pp20, logo, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
			`

			_, profileErr := GetPostgresDB().ExecContext(ctx, profileQuery,
				profileID,
				newUserID,
				name,
				taxID,
				phone,
				sanitizeUTF8(user.Profile.Address),
				province,
				zipcode,
				sanitizeUTF8(user.Profile.ATTA),
				sanitizeUTF8(user.Profile.DBD),
				sanitizeUTF8(user.Profile.PP20),
				logo,
				createdAt,
				updatedAt,
			)

			if profileErr != nil {
				// Log error and skip this profile instead of failing
				fmt.Printf("    WARNING: Failed to insert company profile for user %s (skipping profile): %v\n", username, profileErr)
			}
		}

		stats.Users++
	}

	if duplicateCount > 0 {
		fmt.Printf("    WARNING: %d duplicate usernames were skipped\n", duplicateCount)
	}
	if skippedCount > 0 {
		fmt.Printf("    INFO: %d users already exist in PostgreSQL (skipped)\n", skippedCount)
	}
	fmt.Printf("    Migrated %d users (%d with valid email)\n", stats.Users, emailCount)
	return nil
}

// normalizeUsername normalizes username by trimming spaces and converting to lowercase
func normalizeUsername(username string) string {
	// Trim leading/trailing spaces
	username = strings.TrimSpace(username)
	// Convert to lowercase
	username = strings.ToLower(username)
	// Remove any internal spaces
	username = strings.ReplaceAll(username, " ", "")
	// Truncate to max 100 characters
	if len(username) > 100 {
		username = username[:100]
	}
	return username
}

// normalizeEmail normalizes email by trimming spaces and converting to lowercase
func normalizeEmail(email string) string {
	// Trim leading/trailing spaces
	email = strings.TrimSpace(email)
	// Convert to lowercase
	email = strings.ToLower(email)
	// Remove any internal spaces (shouldn't happen but just in case)
	email = strings.ReplaceAll(email, " ", "")
	return email
}

// isValidEmail checks if string is a valid email format
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	return emailRegex.MatchString(email)
}

func mapUserRole(role string) string {
	switch role {
	case "super", "superadmin":
		return "super"
	case "admin":
		return "admin"
	case "agent", "user", "guest": // guest จากระบบเก่า → agent
		return "agent"
	case "sales":
		return "sales"
	default:
		return "agent"
	}
}

func mapUserStatus(status string) string {
	switch status {
	case "active", "publish", "approve", "approved":
		return "active"
	case "pending", "":
		return "pending"
	case "inactive", "draft":
		return "inactive"
	default:
		return "pending"
	}
}

func stringOrNil(s string) *string {
	s = sanitizeUTF8(s)
	if s == "" {
		return nil
	}
	return &s
}

// truncateString truncates a string to the specified max length (in runes, not bytes)
func truncateString(s string, maxLen int) string {
	s = sanitizeUTF8(s) // Sanitize first
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}

// sanitizeUTF8 removes invalid UTF-8 bytes from a string
func sanitizeUTF8(s string) string {
	// Use ToValidUTF8 to replace invalid sequences with empty string
	return strings.ToValidUTF8(s, "")
}

// truncateStringPtr truncates a string pointer to maxLen characters
func truncateStringPtr(s *string, maxLen int) *string {
	if s == nil {
		return nil
	}
	truncated := truncateString(*s, maxLen)
	return &truncated
}
