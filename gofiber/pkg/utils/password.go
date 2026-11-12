package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt with default cost
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a plain text password with a hashed password
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsValidPassword validates if a password meets the minimum requirements
func IsValidPassword(password string) bool {
	// Minimum password requirements:
	// - At least 8 characters
	// - Maximum 72 characters (bcrypt limitation)
	return len(password) >= 8 && len(password) <= 72
}