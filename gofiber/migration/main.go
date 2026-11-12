package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables - try local .env first, then parent directory
	if err := godotenv.Load(".env"); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("Warning: .env file not found, using environment variables")
		}
	}

	fmt.Println("============================================")
	fmt.Println("  MongoDB to PostgreSQL Migration Tool")
	fmt.Println("============================================")
	fmt.Println()

	// Initialize context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
	defer cancel()

	// Initialize database connections
	fmt.Println("[1/3] Connecting to databases...")
	if err := InitDatabases(); err != nil {
		log.Fatalf("Failed to initialize databases: %v", err)
	}
	defer CloseDatabases()
	fmt.Println("      - MongoDB connected")
	fmt.Println("      - PostgreSQL connected")
	fmt.Println()

	// Initialize ID mapping storage
	fmt.Println("[2/3] Initializing ID mapping storage...")
	InitMapping()
	fmt.Println("      - Mapping storage ready")

	// Initialize migration logger
	if err := InitMigrationLogger(); err != nil {
		log.Printf("Warning: Failed to initialize logger: %v", err)
	}
	defer CloseMigrationLogger()
	fmt.Println()

	// Check command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "verify":
			fmt.Println("[3/3] Running verification only...")
			if err := VerifyMigration(ctx); err != nil {
				log.Fatalf("Verification failed: %v", err)
			}
			return
		case "check-slips":
			fmt.Println("[3/3] Checking Slips and Passports...")
			if err := CheckSlipsAndPassports(ctx); err != nil {
				log.Fatalf("Check failed: %v", err)
			}
			return
		case "fix-slip-amounts":
			fmt.Println("[3/3] Fixing Payment Slip Amounts...")
			if err := FixSlipAmounts(ctx); err != nil {
				log.Fatalf("Fix failed: %v", err)
			}
			return
		case "verify-payments":
			fmt.Println("[3/3] Verifying Payment Data...")
			if err := VerifyPayments(ctx); err != nil {
				log.Fatalf("Verify failed: %v", err)
			}
			return
		case "list-passports":
			fmt.Println("[3/3] Listing Bookings with Passports...")
			if err := ListInvoicesWithPassports(ctx); err != nil {
				log.Fatalf("List failed: %v", err)
			}
			return
		case "show-invoice":
			if len(os.Args) < 3 {
				fmt.Println("Usage: go run . show-invoice <invoice_id>")
				fmt.Println("Example: go run . show-invoice INV-2024-00001")
				return
			}
			if err := ShowInvoice(ctx, os.Args[2]); err != nil {
				log.Fatalf("Show invoice failed: %v", err)
			}
			return
		case "debug-tours":
			fmt.Println("[3/3] Debug Tours...")
			if err := CheckDuplicateTours(ctx); err != nil {
				log.Fatalf("Debug failed: %v", err)
			}
			return
		case "debug-period":
			if len(os.Args) < 3 {
				fmt.Println("Usage: go run . debug-period <tour_code>")
				return
			}
			fmt.Printf("[3/3] Debug Period for Tour: %s\n", os.Args[2])
			if err := CheckPeriodDetails(ctx, os.Args[2]); err != nil {
				log.Fatalf("Debug failed: %v", err)
			}
			return
		case "check-periods":
			if len(os.Args) > 2 {
				// Check specific tour
				fmt.Printf("[3/3] Checking Periods for Tour: %s\n", os.Args[2])
				if err := CheckTourPeriods(ctx, os.Args[2]); err != nil {
					log.Fatalf("Check tour periods failed: %v", err)
				}
			} else {
				fmt.Println("[3/3] Checking Period Dates...")
				if err := CheckPeriodDates(ctx); err != nil {
					log.Fatalf("Check periods failed: %v", err)
				}
			}
			return
		case "clear":
			fmt.Println("[3/3] Clearing PostgreSQL data...")
			if err := ClearPostgresData(ctx); err != nil {
				log.Fatalf("Failed to clear data: %v", err)
			}
			fmt.Println("      - PostgreSQL data cleared")
			return
		case "help":
			printHelp()
			return
		case "files":
			fmt.Println("[3/3] Migrating Files from DO Spaces to iDrive e2...")
			fmt.Println()
			if err := MigrateFiles(ctx); err != nil {
				log.Fatalf("File migration failed: %v", err)
			}
			return
		case "users":
			fmt.Println("[3/3] Migrating Users only...")
			stats := &MigrationStats{}
			if err := MigrateUsers(ctx, stats); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			fmt.Printf("      - Migrated %d users\n", stats.Users)
			return
		case "cities":
			fmt.Println("[3/3] Migrating Cities only...")
			stats := &MigrationStats{}
			if err := MigrateCities(ctx, stats); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			fmt.Printf("      - Migrated %d cities\n", stats.Cities)
			return
		case "banners":
			fmt.Println("[3/3] Migrating Banners only...")
			stats := &MigrationStats{}
			if err := MigrateBanners(ctx, stats); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			fmt.Printf("      - Migrated %d banners\n", stats.Banners)
			return
		case "sales":
			fmt.Println("[3/3] Migrating Sales only...")
			stats := &MigrationStats{}
			if err := MigrateSales(ctx, stats); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			fmt.Printf("      - Migrated %d sales\n", stats.Sales)
			return
		case "tours":
			fmt.Println("[3/3] Migrating Tours only...")
			stats := &MigrationStats{}
			if err := MigrateTours(ctx, stats); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			fmt.Printf("      - Migrated %d tours\n", stats.Tours)
			return
		case "periods":
			fmt.Println("[3/3] Migrating Periods only...")
			stats := &MigrationStats{}
			if err := MigratePeriods(ctx, stats); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			fmt.Printf("      - Migrated %d periods\n", stats.Periods)
			return
		case "bookings":
			fmt.Println("[3/3] Migrating Bookings only...")
			stats := &MigrationStats{}
			if err := MigrateBookings(ctx, stats); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
			fmt.Printf("      - Migrated %d bookings\n", stats.Bookings)
			return
		case "create-admin":
			if len(os.Args) < 4 {
				fmt.Println("Usage: go run . create-admin <username> <email> <password>")
				fmt.Println("Example: go run . create-admin admin@example.com admin@example.com MyPassword123")
				return
			}
			username := os.Args[2]
			email := os.Args[2] // Use same as username if only 3 args
			password := os.Args[3]
			if len(os.Args) >= 5 {
				email = os.Args[3]
				password = os.Args[4]
			}
			if err := CreateSuperAdmin(ctx, username, email, password); err != nil {
				log.Fatalf("Failed to create admin: %v", err)
			}
			return
		case "fix-tour-cities":
			fmt.Println("[3/3] Fixing Tour-City relationships...")
			if err := FixTourCities(ctx); err != nil {
				log.Fatalf("Failed to fix tour cities: %v", err)
			}
			return
		case "fix-user-company":
			fmt.Println("[3/3] Fixing User-Company relationships...")
			if err := FixUserCompany(ctx); err != nil {
				log.Fatalf("Failed to fix user company: %v", err)
			}
			return
		case "fix-booking-users":
			fmt.Println("[3/3] Fixing Booking-User relationships...")
			if err := FixBookingUsers(ctx); err != nil {
				log.Fatalf("Failed to fix booking users: %v", err)
			}
			return
		case "fix-user-status":
			fmt.Println("[3/3] Fixing User Status...")
			if err := FixUserStatus(ctx); err != nil {
				log.Fatalf("Failed to fix user status: %v", err)
			}
			return
		case "fix-all":
			fmt.Println("[3/3] Running All Fixes...")
			if err := FixAll(ctx); err != nil {
				log.Fatalf("Failed to run all fixes: %v", err)
			}
			return
		case "clean-periods":
			fmt.Println("[3/3] Cleaning Duplicate Periods...")
			if err := CleanDuplicatePeriods(ctx); err != nil {
				log.Fatalf("Failed to clean periods: %v", err)
			}
			return
		case "sync-full":
			fmt.Println("[3/3] Running Full Sync...")
			if err := SyncFull(ctx); err != nil {
				log.Fatalf("Full sync failed: %v", err)
			}

			// Run all fixes after sync
			fmt.Println("\n[4/4] Running All Fixes...")
			if err := FixAll(ctx); err != nil {
				log.Fatalf("Fix all failed: %v", err)
			}
			return
		case "sync-latest":
			// Get last sync time or default to 24 hours ago
			lastSync, _ := GetLastSyncTime(ctx)

			// Check for --since flag
			if len(os.Args) > 2 && os.Args[2] == "--since" {
				if len(os.Args) > 3 {
					parsedTime, err := time.Parse("2006-01-02", os.Args[3])
					if err != nil {
						log.Fatalf("Invalid date format. Use: YYYY-MM-DD")
					}
					lastSync = parsedTime
				}
			}

			// Check for --hours flag
			if len(os.Args) > 2 && os.Args[2] == "--hours" {
				if len(os.Args) > 3 {
					hours := 24
					fmt.Sscanf(os.Args[3], "%d", &hours)
					lastSync = time.Now().Add(-time.Duration(hours) * time.Hour)
				}
			}

			fmt.Println("[3/3] Syncing Latest Data with Fixes...")
			if err := SyncLatest(ctx, lastSync); err != nil {
				log.Fatalf("Sync latest failed: %v", err)
			}

			// Update last sync time
			if err := UpdateLastSyncTime(ctx); err != nil {
				log.Printf("Warning: Failed to update last sync time: %v", err)
			}
			return
		case "sync":
			// Create tracking table if not exists
			if err := CreateSyncTrackingTable(ctx); err != nil {
				log.Fatalf("Failed to create sync tracking table: %v", err)
			}

			// Get last sync time
			lastSync, err := GetLastSyncTime(ctx)
			if err != nil {
				log.Fatalf("Failed to get last sync time: %v", err)
			}

			// Check for --since flag
			if len(os.Args) > 2 && os.Args[2] == "--since" {
				if len(os.Args) > 3 {
					// Parse custom date
					parsedTime, err := time.Parse("2006-01-02", os.Args[3])
					if err != nil {
						log.Fatalf("Invalid date format. Use: YYYY-MM-DD")
					}
					lastSync = parsedTime
				}
			}

			// Check for --hours flag
			if len(os.Args) > 2 && os.Args[2] == "--hours" {
				if len(os.Args) > 3 {
					hours := 24
					fmt.Sscanf(os.Args[3], "%d", &hours)
					lastSync = time.Now().Add(-time.Duration(hours) * time.Hour)
				}
			}

			fmt.Println("[3/3] Running Daily Sync...")
			if err := SyncDaily(ctx, lastSync); err != nil {
				log.Fatalf("Sync failed: %v", err)
			}

			// Update last sync time
			if err := UpdateLastSyncTime(ctx); err != nil {
				log.Printf("Warning: Failed to update last sync time: %v", err)
			}
			return
		}
	}

	// Run full migration
	fmt.Println("[3/3] Starting migration...")
	fmt.Println()

	startTime := time.Now()
	stats := &MigrationStats{}

	// Migration order based on foreign key dependencies
	migrations := []struct {
		name string
		fn   func(context.Context, *MigrationStats) error
	}{
		{"Cities", MigrateCities},
		{"Banners", MigrateBanners},
		{"Sales", MigrateSales},
		{"Users", MigrateUsers},
		{"Tours", MigrateTours},
		{"Periods", MigratePeriods},
		{"Bookings", MigrateBookings},
	}

	for i, m := range migrations {
		fmt.Printf("--- [%d/%d] Migrating %s ---\n", i+1, len(migrations), m.name)
		mStart := time.Now()

		if err := m.fn(ctx, stats); err != nil {
			log.Fatalf("Migration failed at %s: %v", m.name, err)
		}

		fmt.Printf("    Completed in %v\n\n", time.Since(mStart).Round(time.Millisecond))
	}

	// Print summary
	fmt.Println("============================================")
	fmt.Println("  Migration Summary")
	fmt.Println("============================================")
	fmt.Printf("  Cities:    %d migrated\n", stats.Cities)
	fmt.Printf("  Banners:   %d migrated\n", stats.Banners)
	fmt.Printf("  Sales:     %d migrated\n", stats.Sales)
	fmt.Printf("  Users:     %d migrated\n", stats.Users)
	fmt.Printf("  Tours:     %d migrated\n", stats.Tours)
	fmt.Printf("  Periods:   %d migrated\n", stats.Periods)
	fmt.Printf("  Bookings:  %d migrated\n", stats.Bookings)
	fmt.Println("--------------------------------------------")
	fmt.Printf("  Total time: %v\n", time.Since(startTime).Round(time.Second))
	fmt.Println("============================================")

	// Run verification
	fmt.Println()
	fmt.Println("Running post-migration verification...")
	if err := VerifyMigration(ctx); err != nil {
		log.Printf("Warning: Verification found issues: %v", err)
	}
}

func printHelp() {
	fmt.Println("Usage: go run . [command]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  (none)    Run full migration")
	fmt.Println("  verify    Run verification only")
	fmt.Println("  clear     Clear PostgreSQL data")
	fmt.Println("  help      Show this help message")
	fmt.Println()
	fmt.Println("Single table migration:")
	fmt.Println("  users     Migrate users only")
	fmt.Println("  cities    Migrate cities only")
	fmt.Println("  banners   Migrate banners only")
	fmt.Println("  sales     Migrate sales only")
	fmt.Println("  tours     Migrate tours only")
	fmt.Println("  periods   Migrate periods only")
	fmt.Println("  bookings  Migrate bookings only")
	fmt.Println()
	fmt.Println("Fix commands:")
	fmt.Println("  fix-tour-cities    Fix tour-city relationships (match by city title)")
	fmt.Println("  fix-user-company   Fix user-company relationships (create missing company profiles)")
	fmt.Println("  fix-booking-users  Fix booking-user relationships (match by profile ID)")
	fmt.Println("  fix-user-status    Fix user status (approved -> active)")
	fmt.Println("  fix-all            Run ALL fixes in correct order")
	fmt.Println("  clean-periods      Delete duplicate periods (keep only first occurrence)")
	fmt.Println()
	fmt.Println("Recommended sync commands:")
	fmt.Println("  sync-full                Sync ALL data + fix all (recommended)")
	fmt.Println("  sync-latest              Sync incremental + fix all (since last run)")
	fmt.Println("  sync-latest --hours 24   Sync incremental + fix all (last 24 hours)")
	fmt.Println("  sync-latest --since 2025-01-01  Sync incremental + fix all (since date)")
	fmt.Println()
	fmt.Println("Daily sync:")
	fmt.Println("  sync              Sync new/updated data since last run")
	fmt.Println("  sync --hours 24   Sync data from last 24 hours")
	fmt.Println("  sync --since 2025-01-01  Sync data since specific date")
	fmt.Println()
	fmt.Println("File migration (DO Spaces → iDrive e2):")
	fmt.Println("  files companies   Company files (logo, atta, dbd, pp20)")
	fmt.Println("  files banners     Banner images")
	fmt.Println("  files cities      City images")
	fmt.Println("  files tours       Tour files (thumbnail, pdf, word)")
	fmt.Println("  files periods     Period attachment files")
	fmt.Println("  files slips       Payment slip images")
	fmt.Println("  files passports   Passport images")
	fmt.Println("  files all         Migrate all files")
}

// MigrationStats tracks migration progress
type MigrationStats struct {
	Cities   int
	Banners  int
	Sales    int
	Users    int
	Tours    int
	Periods  int
	Bookings int
}
