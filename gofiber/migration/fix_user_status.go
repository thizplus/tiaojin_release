package main

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// FixUserStatus updates user status based on MongoDB data
func FixUserStatus(ctx context.Context) error {
	fmt.Println("\n=== Fixing User Status ===")

	// Step 1: Get users with approved status from MongoDB
	fmt.Println("[1/3] Loading MongoDB users with status...")

	mongoStatusMap, err := getMongoUserStatusMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to get MongoDB user status: %w", err)
	}
	fmt.Printf("      Found %d users in MongoDB\n", len(mongoStatusMap))

	// Step 2: Get pending users from PostgreSQL
	fmt.Println("[2/3] Finding pending users in PostgreSQL...")

	pendingUsers, err := getPendingUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to get pending users: %w", err)
	}
	fmt.Printf("      Found %d pending users in PostgreSQL\n", len(pendingUsers))

	// Step 3: Update status based on MongoDB
	fmt.Println("[3/3] Updating user status...")

	updated := 0
	notFound := 0
	alreadyCorrect := 0

	for _, pgUser := range pendingUsers {
		// Find matching MongoDB user by username
		mongoStatus, ok := mongoStatusMap[normalizeForStatusMatch(pgUser.Username)]
		if !ok {
			// Try email
			if pgUser.Email != nil {
				mongoStatus, ok = mongoStatusMap[normalizeForStatusMatch(*pgUser.Email)]
			}
			if !ok {
				notFound++
				continue
			}
		}

		// Map MongoDB status to PostgreSQL status
		newStatus := mapUserStatus(mongoStatus)

		// Skip if already correct
		if newStatus == "pending" {
			alreadyCorrect++
			continue
		}

		// Update user status
		err := updateUserStatus(ctx, pgUser.ID, newStatus)
		if err != nil {
			fmt.Printf("      ERROR: Failed to update user %s: %v\n", pgUser.Username, err)
			continue
		}
		updated++
	}

	// Summary
	fmt.Printf("\n=== Fix User Status Summary ===\n")
	fmt.Printf("    Updated to active: %d users\n", updated)
	fmt.Printf("    Already correct: %d users\n", alreadyCorrect)
	fmt.Printf("    Not found in MongoDB: %d users\n", notFound)

	// Verification
	fmt.Println("\n=== Verification ===")

	var activeCount, pendingCount, inactiveCount int
	GetPostgresDB().QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE status = 'active' AND deleted_at IS NULL").Scan(&activeCount)
	GetPostgresDB().QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE status = 'pending' AND deleted_at IS NULL").Scan(&pendingCount)
	GetPostgresDB().QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE status = 'inactive' AND deleted_at IS NULL").Scan(&inactiveCount)

	fmt.Printf("    Active users: %d\n", activeCount)
	fmt.Printf("    Pending users: %d\n", pendingCount)
	fmt.Printf("    Inactive users: %d\n", inactiveCount)

	return nil
}

// MongoUserStatus represents user with status from MongoDB
type MongoUserStatus struct {
	Email  string `bson:"email"`
	Status string `bson:"status"`
}

// getMongoUserStatusMap returns map of email -> status from MongoDB
func getMongoUserStatusMap(ctx context.Context) (map[string]string, error) {
	collection := GetMongoDB().Collection("tbl_users")

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result := make(map[string]string)
	for cursor.Next(ctx) {
		var user MongoUserStatus
		if err := cursor.Decode(&user); err != nil {
			continue
		}
		result[normalizeForStatusMatch(user.Email)] = user.Status
	}

	return result, cursor.Err()
}

// PendingUser represents a pending user from PostgreSQL
type PendingUser struct {
	ID       string
	Username string
	Email    *string
}

// getPendingUsers returns users with pending status from PostgreSQL
func getPendingUsers(ctx context.Context) ([]PendingUser, error) {
	query := `SELECT id, username, email FROM users WHERE status = 'pending' AND deleted_at IS NULL`

	rows, err := GetPostgresDB().QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []PendingUser
	for rows.Next() {
		var u PendingUser
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			continue
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

// updateUserStatus updates user status in PostgreSQL
func updateUserStatus(ctx context.Context, userID string, status string) error {
	query := `UPDATE users SET status = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`
	_, err := GetPostgresDB().ExecContext(ctx, query, status, userID)
	return err
}

func normalizeForStatusMatch(s string) string {
	return strings.ToLower(strings.TrimSpace(strings.ReplaceAll(s, " ", "")))
}
