package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateSuperAdmin creates a super admin user
func CreateSuperAdmin(ctx context.Context, username, email, password string) error {
	fmt.Println("Creating super admin user...")

	// Normalize username and email
	username = strings.TrimSpace(strings.ToLower(username))
	email = strings.TrimSpace(strings.ToLower(email))

	// Check if user already exists
	var existingID string
	checkQuery := `SELECT id FROM users WHERE username = $1 OR email = $2 LIMIT 1`
	err := GetPostgresDB().QueryRowContext(ctx, checkQuery, username, email).Scan(&existingID)
	if err == nil {
		fmt.Printf("  User already exists with ID: %s\n", existingID)

		// Update to super admin if not already
		updateQuery := `UPDATE users SET role = 'super', status = 'active', updated_at = NOW() WHERE id = $1`
		_, err := GetPostgresDB().ExecContext(ctx, updateQuery, existingID)
		if err != nil {
			return fmt.Errorf("failed to update user role: %w", err)
		}
		fmt.Println("  ✓ Updated existing user to super admin")
		return nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Create new user
	newUserID := uuid.New()
	insertQuery := `
		INSERT INTO users (id, username, email, password, role, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, 'super', 'active', NOW(), NOW())
	`

	_, err = GetPostgresDB().ExecContext(ctx, insertQuery,
		newUserID,
		username,
		email,
		string(hashedPassword),
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("  ✓ Created super admin user\n")
	fmt.Printf("    ID: %s\n", newUserID)
	fmt.Printf("    Username: %s\n", username)
	fmt.Printf("    Email: %s\n", email)
	fmt.Printf("    Role: super\n")
	fmt.Printf("    Status: active\n")

	return nil
}
