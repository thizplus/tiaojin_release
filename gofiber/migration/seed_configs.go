package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// DefaultExpirationTiers returns the default expiration tier configuration
func DefaultExpirationTiers() map[string]interface{} {
	return map[string]interface{}{
		"tiers": []map[string]interface{}{
			{
				"name":           "Critical",
				"description":    "ทัวร์ที่ใกล้จะเริ่มมาก ต้องยืนยันด่วน",
				"threshold_days": 3,
				"expire_minutes": 30,
				"sort_order":     0,
			},
			{
				"name":           "Urgent",
				"description":    "ทัวร์ที่ใกล้จะเริ่ม",
				"threshold_days": 7,
				"expire_hours":   6,
				"sort_order":     1,
			},
			{
				"name":           "Normal",
				"description":    "ทัวร์ปกติ",
				"threshold_days": 14,
				"expire_hours":   24,
				"sort_order":     2,
			},
			{
				"name":           "Standard",
				"description":    "ทัวร์ล่วงหน้า",
				"threshold_days": 30,
				"expire_days":    2,
				"sort_order":     3,
			},
			{
				"name":           "Default",
				"description":    "ทัวร์ล่วงหน้ามาก",
				"threshold_days": nil,
				"expire_days":    3,
				"sort_order":     4,
			},
		},
		"metadata": map[string]interface{}{
			"version":      "1.0",
			"last_updated": time.Now().Format(time.RFC3339),
			"updated_by":   "system",
		},
	}
}

// DefaultAgentBookingLimits returns the default agent booking limit configuration
func DefaultAgentBookingLimits() map[string]interface{} {
	return map[string]interface{}{
		"global_settings": map[string]interface{}{
			"enabled":                  false,
			"default_limit_per_period": 10,
			"count_all_statuses":       []string{"pending", "confirmed", "paid"},
		},
		"period_specific_limits": []interface{}{},
		"agent_global_bypasses":  []interface{}{},
		"metadata": map[string]interface{}{
			"version":      "1.0",
			"last_updated": time.Now().Format(time.RFC3339),
			"updated_by":   "system",
		},
	}
}

// SeedDefaultConfigs seeds default system configurations
func SeedDefaultConfigs(ctx context.Context) error {
	fmt.Println("\n--- Seeding Default System Configurations ---")

	configs := []struct {
		ConfigKey   string
		ConfigValue map[string]interface{}
		ConfigType  string
		Description string
	}{
		{
			ConfigKey:   "booking.expiration_tiers",
			ConfigValue: DefaultExpirationTiers(),
			ConfigType:  "expiration",
			Description: "Booking expiration tier configuration - กำหนดระยะเวลาหมดอายุการจองตามความใกล้ของวันเริ่มทัวร์",
		},
		{
			ConfigKey:   "booking.agent_period_limits",
			ConfigValue: DefaultAgentBookingLimits(),
			ConfigType:  "booking",
			Description: "Agent booking limit configuration per period - จำกัดจำนวนการจองของ Agent ต่อ Period",
		},
	}

	for _, cfg := range configs {
		// Check if config already exists
		var count int
		err := GetPostgresDB().QueryRowContext(ctx,
			"SELECT COUNT(*) FROM system_configs WHERE config_key = $1",
			cfg.ConfigKey,
		).Scan(&count)

		if err != nil {
			fmt.Printf("    ERROR checking config '%s': %v\n", cfg.ConfigKey, err)
			continue
		}

		if count > 0 {
			fmt.Printf("    SKIP: Config '%s' already exists\n", cfg.ConfigKey)
			continue
		}

		// Convert to JSON
		jsonValue, err := json.Marshal(cfg.ConfigValue)
		if err != nil {
			fmt.Printf("    ERROR marshaling config '%s': %v\n", cfg.ConfigKey, err)
			continue
		}

		// Insert config
		_, err = GetPostgresDB().ExecContext(ctx, `
			INSERT INTO system_configs (config_key, config_value, config_type, description, is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, true, NOW(), NOW())
		`, cfg.ConfigKey, string(jsonValue), cfg.ConfigType, cfg.Description)

		if err != nil {
			fmt.Printf("    ERROR creating config '%s': %v\n", cfg.ConfigKey, err)
			return fmt.Errorf("failed to seed config '%s': %w", cfg.ConfigKey, err)
		}

		fmt.Printf("    CREATED: %s\n", cfg.ConfigKey)
	}

	fmt.Println("    Default system configurations seeded successfully")
	return nil
}

// SeedDefaultJobs seeds default scheduled jobs
func SeedDefaultJobs(ctx context.Context) error {
	fmt.Println("\n--- Seeding Default Scheduled Jobs ---")

	// Check if backup-booking-expiry-cleanup job exists
	var count int
	err := GetPostgresDB().QueryRowContext(ctx,
		"SELECT COUNT(*) FROM jobs WHERE name = $1",
		"backup-booking-expiry-cleanup",
	).Scan(&count)

	if err != nil {
		fmt.Printf("    ERROR checking job: %v\n", err)
		return nil // Don't fail the migration
	}

	if count > 0 {
		fmt.Println("    SKIP: Job 'backup-booking-expiry-cleanup' already exists")
		return nil
	}

	// Create the backup cleanup job
	payload := `{"type": "booking_expiry_cleanup", "description": "Backup cleanup job to expire bookings that might have been missed by individual expiry jobs"}`

	_, err = GetPostgresDB().ExecContext(ctx, `
		INSERT INTO jobs (name, cron_expr, payload, status, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, 'active', true, NOW(), NOW())
	`, "backup-booking-expiry-cleanup", "0 */1 * * *", payload)

	if err != nil {
		fmt.Printf("    ERROR creating job: %v\n", err)
		return nil // Don't fail the migration
	}

	fmt.Println("    CREATED: backup-booking-expiry-cleanup (runs every hour)")
	return nil
}
