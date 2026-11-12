package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MongoUserWithProfile represents user with embedded profile from MongoDB
type MongoUserWithProfile struct {
	ID      primitive.ObjectID `bson:"_id"`
	Email   string             `bson:"email"`
	Profile *MongoUserProfile  `bson:"profile"`
}

// FixUserCompany creates missing company_profiles for existing users
func FixUserCompany(ctx context.Context) error {
	fmt.Println("\n=== Fixing User-Company Relationships ===")

	// Step 1: Find users WITHOUT company_profile in PostgreSQL
	fmt.Println("[1/4] Finding users without company profile...")

	usersWithoutCompany, err := getUsersWithoutCompanyProfile(ctx)
	if err != nil {
		return fmt.Errorf("failed to get users without company: %w", err)
	}
	fmt.Printf("      Found %d users without company profile\n", len(usersWithoutCompany))

	if len(usersWithoutCompany) == 0 {
		fmt.Println("      All users have company profiles. Nothing to fix.")
		return nil
	}

	// Step 2: Get all users with profiles from MongoDB
	fmt.Println("[2/4] Loading MongoDB users with profiles...")

	mongoUsersMap, err := getMongoUsersWithProfiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MongoDB users: %w", err)
	}
	fmt.Printf("      Found %d users with profiles in MongoDB\n", len(mongoUsersMap))

	// Step 3: Create missing company_profiles
	fmt.Println("[3/4] Creating missing company profiles...")

	created := 0
	notFound := 0
	noProfile := 0
	errors := 0

	for _, pgUser := range usersWithoutCompany {
		// Find matching MongoDB user by username/email
		mongoUser, ok := mongoUsersMap[normalizeForMatch(pgUser.Username)]
		if !ok {
			// Try email if username doesn't match
			if pgUser.Email != nil {
				mongoUser, ok = mongoUsersMap[normalizeForMatch(*pgUser.Email)]
			}
			if !ok {
				notFound++
				continue
			}
		}

		// Skip if no profile in MongoDB
		if mongoUser.Profile == nil {
			noProfile++
			continue
		}

		// Create company_profile
		err := createCompanyProfile(ctx, pgUser.ID, mongoUser.Profile)
		if err != nil {
			fmt.Printf("      ERROR: Failed to create company for user %s: %v\n", pgUser.Username, err)
			errors++
			continue
		}
		created++
	}

	// Step 4: Summary
	fmt.Printf("\n=== Fix User-Company Summary ===\n")
	fmt.Printf("    Created: %d company profiles\n", created)
	fmt.Printf("    Not found in MongoDB: %d users\n", notFound)
	fmt.Printf("    No profile in MongoDB: %d users\n", noProfile)
	fmt.Printf("    Errors: %d\n", errors)

	// Verification
	fmt.Println("\n=== Verification ===")

	var totalUsers int
	GetPostgresDB().QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&totalUsers)

	var withCompany int
	GetPostgresDB().QueryRowContext(ctx, `
		SELECT COUNT(*) FROM users u
		WHERE u.deleted_at IS NULL
		AND EXISTS (SELECT 1 FROM company_profiles cp WHERE cp.user_id = u.id AND cp.deleted_at IS NULL)
	`).Scan(&withCompany)

	fmt.Printf("    Total users: %d\n", totalUsers)
	fmt.Printf("    Users with company: %d\n", withCompany)
	fmt.Printf("    Users without company: %d\n", totalUsers-withCompany)

	return nil
}

// PostgresUserSimple represents minimal user data from PostgreSQL
type PostgresUserSimple struct {
	ID       uuid.UUID
	Username string
	Email    *string
}

// getUsersWithoutCompanyProfile returns users that don't have company_profile
func getUsersWithoutCompanyProfile(ctx context.Context) ([]PostgresUserSimple, error) {
	query := `
		SELECT u.id, u.username, u.email
		FROM users u
		LEFT JOIN company_profiles cp ON u.id = cp.user_id AND cp.deleted_at IS NULL
		WHERE u.deleted_at IS NULL AND cp.id IS NULL
	`

	rows, err := GetPostgresDB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []PostgresUserSimple
	for rows.Next() {
		var u PostgresUserSimple
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			continue
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

// getMongoUsersWithProfiles returns map of email -> MongoUserWithProfile
func getMongoUsersWithProfiles(ctx context.Context) (map[string]*MongoUserWithProfile, error) {
	collection := GetMongoDB().Collection("tbl_users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]*MongoUserWithProfile)
	for cursor.Next(ctx) {
		var user MongoUserWithProfile
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		// Index by normalized email
		result[normalizeForMatch(user.Email)] = &user
	}

	return result, cursor.Err()
}

// createCompanyProfile creates a company_profile for a PostgreSQL user
func createCompanyProfile(ctx context.Context, userID uuid.UUID, profile *MongoUserProfile) error {
	profileID := uuid.New()

	// Handle empty or null values
	logo := stringOrNil(profile.Logo)

	// Truncate fields to fit database constraints
	taxID := truncateString(profile.TaxID, 13)
	phone := truncateString(profile.Phone, 20)
	province := truncateString(profile.Province, 100)
	zipcode := truncateString(profile.Zipcode, 5)
	name := truncateString(profile.Name, 200)

	query := `
		INSERT INTO company_profiles (id, user_id, name, tax_id, phone, address, province, zipcode, atta, dbd, pp20, logo, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, NOW(), NOW())
	`

	_, err := GetPostgresDB().ExecContext(ctx, query,
		profileID,
		userID,
		name,
		taxID,
		phone,
		sanitizeUTF8(profile.Address),
		province,
		zipcode,
		sanitizeUTF8(profile.ATTA),
		sanitizeUTF8(profile.DBD),
		sanitizeUTF8(profile.PP20),
		logo,
	)

	return err
}

func normalizeForMatch(s string) string {
	return strings.ToLower(strings.TrimSpace(strings.ReplaceAll(s, " ", "")))
}
