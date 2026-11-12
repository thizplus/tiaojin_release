package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

// Rate limiting config
const (
	DelayBetweenFiles   = 200 * time.Millisecond // Delay ระหว่างไฟล์
	DelayBetweenRecords = 100 * time.Millisecond // Delay ระหว่าง record
)

// S3ClientConfig holds S3 connection config
type S3ClientConfig struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
	PublicURL string
}

// DigitalOcean Spaces (Source)
var doSpacesConfig = S3ClientConfig{
	Endpoint:  "https://sgp1.digitaloceanspaces.com",
	Region:    "sgp1",
	AccessKey: "DO00FVR27QFDGU8FAKZL",
	SecretKey: "1ZXXPzGFHvdsqyfroyFc9PImZyT43GP9SRLZ1uJeWjE",
	Bucket:    "tiaojin",
	PublicURL: "https://tiaojin.sgp1.cdn.digitaloceanspaces.com",
}

// Cloudflare R2 (Destination)
var r2Config = S3ClientConfig{
	Endpoint:  "https://fcc0e164ed5f9fcf121a73f8f111ccd1.r2.cloudflarestorage.com",
	Region:    "auto",
	AccessKey: "ec1293bc10ab6bff0c09788a73b839dd",
	SecretKey: "b325b8272364c52e5b49b53b47aa7e4d747357595d903e72450f98ee583b5e89",
	Bucket:    "tiaojin",
	PublicURL: "https://pub-1fc850a8ec4e4e1395125c2ce7811e5d.r2.dev",
}

var sourceClient *s3.Client
var destClient *s3.Client

func createS3Client(cfg S3ClientConfig) (*s3.Client, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               cfg.Endpoint,
			SigningRegion:     cfg.Region,
			HostnameImmutable: true,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKey,
			cfg.SecretKey,
			"",
		)),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return client, nil
}

// initS3Clients initializes both source and destination S3 clients
func initS3Clients() error {
	var err error

	fmt.Println("Connecting to DigitalOcean Spaces...")
	sourceClient, err = createS3Client(doSpacesConfig)
	if err != nil {
		return fmt.Errorf("failed to create source client: %w", err)
	}
	fmt.Println("  ✓ Connected to DO Spaces")

	fmt.Println("Connecting to Cloudflare R2...")
	destClient, err = createS3Client(r2Config)
	if err != nil {
		return fmt.Errorf("failed to create destination client: %w", err)
	}
	fmt.Println("  ✓ Connected to Cloudflare R2")

	return nil
}

// extractPathFromURL extracts the file path from a full URL or returns the path as-is
func extractPathFromURL(urlOrPath string) string {
	if urlOrPath == "" {
		return ""
	}

	// Remove DO Spaces CDN URL prefix
	path := strings.TrimPrefix(urlOrPath, doSpacesConfig.PublicURL)
	path = strings.TrimPrefix(path, "https://tiaojin.sgp1.digitaloceanspaces.com")
	path = strings.TrimPrefix(path, "/")

	return path
}

// isAlreadyMigrated checks if the path is already in the new format
// Old format: uploads/banner/xxx_yyy.jpg (no UUID structure)
// New format: uploads/banners/{entity-uuid}/image/{file-uuid}.jpg
func isAlreadyMigrated(urlOrPath string) bool {
	if urlOrPath == "" {
		return true // empty = nothing to migrate
	}
	// Already migrated if it contains R2 URL
	if strings.Contains(urlOrPath, "r2.dev") || strings.Contains(urlOrPath, "r2.cloudflarestorage.com") {
		return true
	}
	// Check if it's the NEW format with UUID structure
	// New format has pattern: uploads/{category}/{uuid}/{type}/{uuid}.ext
	// Old format has pattern: uploads/{category}/{timestamp}_{random}.ext
	if strings.HasPrefix(urlOrPath, "uploads/") {
		// Count path segments - new format has 5 segments (uploads/category/entityId/type/filename)
		segments := strings.Split(urlOrPath, "/")
		if len(segments) >= 5 {
			// Check if 3rd segment looks like a UUID (36 chars with dashes)
			if len(segments[2]) == 36 && strings.Contains(segments[2], "-") {
				return true // New format - already migrated
			}
		}
		// Old format like uploads/banner/xxx.jpg - needs migration
		return false
	}
	return false
}

// downloadFile downloads a file from DO Spaces and returns its content
func downloadFile(ctx context.Context, path string) ([]byte, string, error) {
	output, err := sourceClient.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(doSpacesConfig.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, "", err
	}
	defer output.Body.Close()

	content, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, "", err
	}

	contentType := "application/octet-stream"
	if output.ContentType != nil {
		contentType = *output.ContentType
	}

	return content, contentType, nil
}

// uploadFile uploads a file to Cloudflare R2 with a new UUID-based path
func uploadFile(ctx context.Context, content []byte, newPath string, contentType string) error {
	_, err := destClient.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r2Config.Bucket),
		Key:         aws.String(newPath),
		Body:        strings.NewReader(string(content)),
		ContentType: aws.String(contentType),
	})
	return err
}

// generateNewPath generates a new path with UUID filename
func generateNewPath(category, entityID, fileType, originalPath string) string {
	ext := filepath.Ext(originalPath)
	if ext == "" {
		ext = ".jpg" // default
	}
	newFileName := uuid.New().String() + ext
	return fmt.Sprintf("uploads/%s/%s/%s/%s", category, entityID, fileType, newFileName)
}

// generateNewURL returns just the path (not full URL)
// Full URL will be constructed from S3_PUBLIC_URL in .env when displaying
func generateNewURL(newPath string) string {
	// เก็บแค่ path ไม่ต้องต่อ PublicURL
	return newPath
}

// migrateFileField migrates a single file field and returns the new URL
func migrateFileField(ctx context.Context, oldURL, category, entityID, fileType string) (string, error) {
	if oldURL == "" {
		return "", nil
	}

	oldPath := extractPathFromURL(oldURL)
	if oldPath == "" {
		return "", nil
	}

	// Download from DO Spaces
	content, contentType, err := downloadFile(ctx, oldPath)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}

	// Rate limiting - delay after download
	time.Sleep(DelayBetweenFiles)

	// Generate new path with UUID
	newPath := generateNewPath(category, entityID, fileType, oldPath)

	// Upload to Cloudflare R2
	if err := uploadFile(ctx, content, newPath, contentType); err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	// Rate limiting - delay after upload
	time.Sleep(DelayBetweenFiles)

	return generateNewURL(newPath), nil
}

// MigrateFiles is the main entry point for file migration
func MigrateFiles(ctx context.Context) error {
	if len(os.Args) < 3 {
		printFileMigrationHelp()
		return nil
	}

	if err := initS3Clients(); err != nil {
		return err
	}

	target := os.Args[2]
	switch target {
	case "companies":
		return migrateCompanyFiles(ctx)
	case "banners":
		return migrateBannerFiles(ctx)
	case "cities":
		return migrateCityFiles(ctx)
	case "tours":
		return migrateTourFiles(ctx)
	case "periods":
		return migratePeriodFiles(ctx)
	case "slips":
		return migratePaymentSlipFiles(ctx)
	case "passports":
		return migratePassportFiles(ctx)
	case "all":
		return migrateAllFiles(ctx)
	default:
		fmt.Printf("Unknown target: %s\n", target)
		printFileMigrationHelp()
		return nil
	}
}

func printFileMigrationHelp() {
	fmt.Println("Usage: go run . files <target>")
	fmt.Println()
	fmt.Println("Targets:")
	fmt.Println("  companies   Migrate company profile files (logo, atta, dbd, pp20)")
	fmt.Println("  banners     Migrate banner images")
	fmt.Println("  cities      Migrate city images")
	fmt.Println("  tours       Migrate tour files (thumbnail, pdf, word)")
	fmt.Println("  periods     Migrate period attachment files")
	fmt.Println("  slips       Migrate payment slip images")
	fmt.Println("  passports   Migrate passport images")
	fmt.Println("  all         Migrate all files")
}

func migrateAllFiles(ctx context.Context) error {
	targets := []struct {
		name string
		fn   func(context.Context) error
	}{
		{"Companies", migrateCompanyFiles},
		{"Banners", migrateBannerFiles},
		{"Cities", migrateCityFiles},
		{"Tours", migrateTourFiles},
		{"Periods", migratePeriodFiles},
		{"Payment Slips", migratePaymentSlipFiles},
		{"Passports", migratePassportFiles},
	}

	for _, t := range targets {
		fmt.Printf("\n=== Migrating %s ===\n", t.name)
		if err := t.fn(ctx); err != nil {
			fmt.Printf("ERROR migrating %s: %v\n", t.name, err)
			// Continue with other targets
		}
	}

	fmt.Println("\n=== All file migrations completed ===")
	return nil
}

// migrateCompanyFiles migrates company profile files (logo, atta, dbd, pp20)
func migrateCompanyFiles(ctx context.Context) error {
	fmt.Println("Migrating company profile files...")

	// Count total records first
	var total int
	if err := GetPostgresDB().QueryRowContext(ctx, `SELECT COUNT(*) FROM company_profiles WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return fmt.Errorf("count failed: %w", err)
	}
	fmt.Printf("  Found %d company profiles\n", total)

	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT id, user_id, logo, atta, dbd, pp20
		FROM company_profiles
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var migrated, skipped, errors, current int

	for rows.Next() {
		current++
		var id, userID string
		var logo, atta, dbd, pp20 sql.NullString

		if err := rows.Scan(&id, &userID, &logo, &atta, &dbd, &pp20); err != nil {
			fmt.Printf("  ERROR scanning row: %v\n", err)
			errors++
			continue
		}

		var newLogo, newAtta, newDBD, newPP20 sql.NullString
		hasChanges := false
		filesMigrated := 0

		// Migrate logo
		if logo.Valid && logo.String != "" && !isAlreadyMigrated(logo.String) {
			fmt.Printf("  [%d/%d] Migrating logo for %s...\n", current, total, id[:8])
			newPath, err := migrateFileField(ctx, logo.String, "users", userID, "logo")
			if err != nil {
				fmt.Printf("\n  ERROR migrating logo for %s: %v\n", id, err)
				errors++
			} else if newPath != "" {
				newLogo = sql.NullString{String: newPath, Valid: true}
				hasChanges = true
				filesMigrated++
			}
		}

		// Migrate atta
		if atta.Valid && atta.String != "" && !isAlreadyMigrated(atta.String) {
			fmt.Printf("\r  [%d/%d] Migrating atta for %s...", current, total, id[:8])
			newPath, err := migrateFileField(ctx, atta.String, "users", userID, "atta")
			if err != nil {
				fmt.Printf("\n  ERROR migrating atta for %s: %v\n", id, err)
				errors++
			} else if newPath != "" {
				newAtta = sql.NullString{String: newPath, Valid: true}
				hasChanges = true
				filesMigrated++
			}
		}

		// Migrate dbd
		if dbd.Valid && dbd.String != "" && !isAlreadyMigrated(dbd.String) {
			fmt.Printf("\r  [%d/%d] Migrating dbd for %s...", current, total, id[:8])
			newPath, err := migrateFileField(ctx, dbd.String, "users", userID, "dbd")
			if err != nil {
				fmt.Printf("\n  ERROR migrating dbd for %s: %v\n", id, err)
				errors++
			} else if newPath != "" {
				newDBD = sql.NullString{String: newPath, Valid: true}
				hasChanges = true
				filesMigrated++
			}
		}

		// Migrate pp20
		if pp20.Valid && pp20.String != "" && !isAlreadyMigrated(pp20.String) {
			fmt.Printf("\r  [%d/%d] Migrating pp20 for %s...", current, total, id[:8])
			newPath, err := migrateFileField(ctx, pp20.String, "users", userID, "pp20")
			if err != nil {
				fmt.Printf("\n  ERROR migrating pp20 for %s: %v\n", id, err)
				errors++
			} else if newPath != "" {
				newPP20 = sql.NullString{String: newPath, Valid: true}
				hasChanges = true
				filesMigrated++
			}
		}

		// Update database if there are changes
		if hasChanges {
			_, err := GetPostgresDB().ExecContext(ctx, `
				UPDATE company_profiles
				SET logo = COALESCE($1, logo),
				    atta = COALESCE($2, atta),
				    dbd = COALESCE($3, dbd),
				    pp20 = COALESCE($4, pp20),
				    updated_at = NOW()
				WHERE id = $5
			`, newLogo, newAtta, newDBD, newPP20, id)
			if err != nil {
				fmt.Printf("\n  ERROR updating %s: %v\n", id, err)
				errors++
			} else {
				migrated++
				fmt.Printf("\r  [%d/%d] ✓ Migrated %d files for %s\n", current, total, filesMigrated, id[:8])
			}
		} else {
			skipped++
		}

		// Delay between records
		time.Sleep(DelayBetweenRecords)
	}

	fmt.Printf("\n  Summary: Migrated: %d, Skipped: %d, Errors: %d\n", migrated, skipped, errors)
	return nil
}

// migrateBannerFiles migrates banner images
func migrateBannerFiles(ctx context.Context) error {
	fmt.Println("Migrating banner files...")

	var total int
	if err := GetPostgresDB().QueryRowContext(ctx, `SELECT COUNT(*) FROM banners WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return fmt.Errorf("count failed: %w", err)
	}
	fmt.Printf("  Found %d banners\n", total)

	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT id, img FROM banners WHERE deleted_at IS NULL
	`)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var migrated, skipped, errors, current int

	for rows.Next() {
		current++
		var id, img string

		if err := rows.Scan(&id, &img); err != nil {
			fmt.Printf("  ERROR scanning row: %v\n", err)
			errors++
			continue
		}

		if isAlreadyMigrated(img) {
			skipped++
			continue
		}

		fmt.Printf("\r  [%d/%d] Migrating banner %s...", current, total, id[:8])
		newURL, err := migrateFileField(ctx, img, "banners", id, "image")
		if err != nil {
			fmt.Printf("\n  ERROR migrating banner %s: %v\n", id, err)
			errors++
			continue
		}

		if newURL != "" {
			_, err := GetPostgresDB().ExecContext(ctx, `
				UPDATE banners SET img = $1, updated_at = NOW() WHERE id = $2
			`, newURL, id)
			if err != nil {
				fmt.Printf("\n  ERROR updating banner %s: %v\n", id, err)
				errors++
			} else {
				migrated++
				fmt.Printf("\r  [%d/%d] ✓ Migrated banner %s\n", current, total, id[:8])
			}
		}

		time.Sleep(DelayBetweenRecords)
	}

	fmt.Printf("\n  Summary: Migrated: %d, Skipped: %d, Errors: %d\n", migrated, skipped, errors)
	return nil
}

// migrateCityFiles migrates city images
func migrateCityFiles(ctx context.Context) error {
	fmt.Println("Migrating city files...")

	var total int
	if err := GetPostgresDB().QueryRowContext(ctx, `SELECT COUNT(*) FROM cities WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return fmt.Errorf("count failed: %w", err)
	}
	fmt.Printf("  Found %d cities\n", total)

	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT id, img FROM cities WHERE deleted_at IS NULL
	`)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var migrated, skipped, errors, current int

	for rows.Next() {
		current++
		var id, img string

		if err := rows.Scan(&id, &img); err != nil {
			fmt.Printf("  ERROR scanning row: %v\n", err)
			errors++
			continue
		}

		if isAlreadyMigrated(img) {
			skipped++
			continue
		}

		fmt.Printf("\r  [%d/%d] Migrating city %s...", current, total, id[:8])
		newURL, err := migrateFileField(ctx, img, "cities", id, "image")
		if err != nil {
			fmt.Printf("\n  ERROR migrating city %s: %v\n", id, err)
			errors++
			continue
		}

		if newURL != "" {
			_, err := GetPostgresDB().ExecContext(ctx, `
				UPDATE cities SET img = $1, updated_at = NOW() WHERE id = $2
			`, newURL, id)
			if err != nil {
				fmt.Printf("\n  ERROR updating city %s: %v\n", id, err)
				errors++
			} else {
				migrated++
				fmt.Printf("\r  [%d/%d] ✓ Migrated city %s\n", current, total, id[:8])
			}
		}

		time.Sleep(DelayBetweenRecords)
	}

	fmt.Printf("\n  Summary: Migrated: %d, Skipped: %d, Errors: %d\n", migrated, skipped, errors)
	return nil
}

// migrateTourFiles migrates tour files (thumbnail, pdf, word)
func migrateTourFiles(ctx context.Context) error {
	fmt.Println("Migrating tour files...")

	var total int
	if err := GetPostgresDB().QueryRowContext(ctx, `SELECT COUNT(*) FROM tours WHERE deleted_at IS NULL`).Scan(&total); err != nil {
		return fmt.Errorf("count failed: %w", err)
	}
	fmt.Printf("  Found %d tours\n", total)

	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT id, thumbnail, file_pdf, file_word FROM tours WHERE deleted_at IS NULL
	`)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var migrated, skipped, errors, current int

	for rows.Next() {
		current++
		var id string
		var thumbnail, filePDF, fileWord sql.NullString

		if err := rows.Scan(&id, &thumbnail, &filePDF, &fileWord); err != nil {
			fmt.Printf("  ERROR scanning row: %v\n", err)
			errors++
			continue
		}

		var newThumbnail, newPDF, newWord sql.NullString
		hasChanges := false
		filesMigrated := 0

		// Migrate thumbnail
		if thumbnail.Valid && !isAlreadyMigrated(thumbnail.String) {
			fmt.Printf("\r  [%d/%d] Migrating thumbnail for tour %s...", current, total, id[:8])
			newURL, err := migrateFileField(ctx, thumbnail.String, "tours", id, "thumbnail")
			if err != nil {
				fmt.Printf("\n  ERROR migrating thumbnail for %s: %v\n", id, err)
				errors++
			} else if newURL != "" {
				newThumbnail = sql.NullString{String: newURL, Valid: true}
				hasChanges = true
				filesMigrated++
			}
		}

		// Migrate PDF
		if filePDF.Valid && !isAlreadyMigrated(filePDF.String) {
			fmt.Printf("\r  [%d/%d] Migrating PDF for tour %s...", current, total, id[:8])
			newURL, err := migrateFileField(ctx, filePDF.String, "tours", id, "pdf")
			if err != nil {
				fmt.Printf("\n  ERROR migrating pdf for %s: %v\n", id, err)
				errors++
			} else if newURL != "" {
				newPDF = sql.NullString{String: newURL, Valid: true}
				hasChanges = true
				filesMigrated++
			}
		}

		// Migrate Word
		if fileWord.Valid && !isAlreadyMigrated(fileWord.String) {
			fmt.Printf("\r  [%d/%d] Migrating Word for tour %s...", current, total, id[:8])
			newURL, err := migrateFileField(ctx, fileWord.String, "tours", id, "word")
			if err != nil {
				fmt.Printf("\n  ERROR migrating word for %s: %v\n", id, err)
				errors++
			} else if newURL != "" {
				newWord = sql.NullString{String: newURL, Valid: true}
				hasChanges = true
				filesMigrated++
			}
		}

		if hasChanges {
			_, err := GetPostgresDB().ExecContext(ctx, `
				UPDATE tours
				SET thumbnail = COALESCE($1, thumbnail),
				    file_pdf = COALESCE($2, file_pdf),
				    file_word = COALESCE($3, file_word),
				    updated_at = NOW()
				WHERE id = $4
			`, newThumbnail, newPDF, newWord, id)
			if err != nil {
				fmt.Printf("\n  ERROR updating tour %s: %v\n", id, err)
				errors++
			} else {
				migrated++
				fmt.Printf("\r  [%d/%d] ✓ Migrated %d files for tour %s\n", current, total, filesMigrated, id[:8])
			}
		} else {
			skipped++
		}

		time.Sleep(DelayBetweenRecords)
	}

	fmt.Printf("\n  Summary: Migrated: %d, Skipped: %d, Errors: %d\n", migrated, skipped, errors)
	return nil
}

// FileAttachment represents a file attachment in period
type FileAttachment struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

// migratePeriodFiles migrates period attachment files (JSONB array)
func migratePeriodFiles(ctx context.Context) error {
	fmt.Println("Migrating period files...")

	var total int
	if err := GetPostgresDB().QueryRowContext(ctx, `SELECT COUNT(*) FROM periods WHERE deleted_at IS NULL AND files IS NOT NULL`).Scan(&total); err != nil {
		return fmt.Errorf("count failed: %w", err)
	}
	fmt.Printf("  Found %d periods with files\n", total)

	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT id, files FROM periods WHERE deleted_at IS NULL AND files IS NOT NULL
	`)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var migrated, skipped, errors, current int

	for rows.Next() {
		current++
		var id string
		var filesJSON []byte

		if err := rows.Scan(&id, &filesJSON); err != nil {
			fmt.Printf("  ERROR scanning row: %v\n", err)
			errors++
			continue
		}

		if len(filesJSON) == 0 || string(filesJSON) == "null" || string(filesJSON) == "[]" {
			skipped++
			continue
		}

		var files []FileAttachment
		if err := json.Unmarshal(filesJSON, &files); err != nil {
			fmt.Printf("  ERROR parsing files JSON for %s: %v\n", id, err)
			errors++
			continue
		}

		if len(files) == 0 {
			skipped++
			continue
		}

		hasChanges := false
		filesMigrated := 0
		for i, file := range files {
			if isAlreadyMigrated(file.URL) {
				continue
			}

			fmt.Printf("\r  [%d/%d] Migrating attachment %d for period %s...", current, total, i+1, id[:8])
			newURL, err := migrateFileField(ctx, file.URL, "periods", id, "attachment")
			if err != nil {
				fmt.Printf("\n  ERROR migrating period file %s: %v\n", id, err)
				errors++
				continue
			}

			if newURL != "" {
				files[i].URL = newURL
				hasChanges = true
				filesMigrated++
			}
		}

		if hasChanges {
			newFilesJSON, err := json.Marshal(files)
			if err != nil {
				fmt.Printf("\n  ERROR marshaling files for %s: %v\n", id, err)
				errors++
				continue
			}

			_, err = GetPostgresDB().ExecContext(ctx, `
				UPDATE periods SET files = $1, updated_at = NOW() WHERE id = $2
			`, newFilesJSON, id)
			if err != nil {
				fmt.Printf("\n  ERROR updating period %s: %v\n", id, err)
				errors++
			} else {
				migrated++
				fmt.Printf("\r  [%d/%d] ✓ Migrated %d files for period %s\n", current, total, filesMigrated, id[:8])
			}
		} else {
			skipped++
		}

		time.Sleep(DelayBetweenRecords)
	}

	fmt.Printf("\n  Summary: Migrated: %d, Skipped: %d, Errors: %d\n", migrated, skipped, errors)
	return nil
}

// migratePaymentSlipFiles migrates payment slip images
func migratePaymentSlipFiles(ctx context.Context) error {
	fmt.Println("Migrating payment slip files...")

	var total int
	if err := GetPostgresDB().QueryRowContext(ctx, `SELECT COUNT(*) FROM payment_slips WHERE deleted_at IS NULL AND image IS NOT NULL`).Scan(&total); err != nil {
		return fmt.Errorf("count failed: %w", err)
	}
	fmt.Printf("  Found %d payment slips with images\n", total)

	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT id, booking_id, image FROM payment_slips WHERE deleted_at IS NULL AND image IS NOT NULL
	`)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var migrated, skipped, errors, current int

	for rows.Next() {
		current++
		var id, bookingID string
		var image sql.NullString

		if err := rows.Scan(&id, &bookingID, &image); err != nil {
			fmt.Printf("  ERROR scanning row: %v\n", err)
			errors++
			continue
		}

		if !image.Valid || isAlreadyMigrated(image.String) {
			skipped++
			continue
		}

		fmt.Printf("\r  [%d/%d] Migrating slip %s...", current, total, id[:8])
		newURL, err := migrateFileField(ctx, image.String, "bookings", bookingID, "slip")
		if err != nil {
			fmt.Printf("\n  ERROR migrating slip %s: %v\n", id, err)
			errors++
			continue
		}

		if newURL != "" {
			_, err := GetPostgresDB().ExecContext(ctx, `
				UPDATE payment_slips SET image = $1, updated_at = NOW() WHERE id = $2
			`, newURL, id)
			if err != nil {
				fmt.Printf("\n  ERROR updating slip %s: %v\n", id, err)
				errors++
			} else {
				migrated++
				fmt.Printf("\r  [%d/%d] ✓ Migrated slip %s\n", current, total, id[:8])
			}
		}

		time.Sleep(DelayBetweenRecords)
	}

	fmt.Printf("\n  Summary: Migrated: %d, Skipped: %d, Errors: %d\n", migrated, skipped, errors)
	return nil
}

// migratePassportFiles migrates passport images
func migratePassportFiles(ctx context.Context) error {
	fmt.Println("Migrating passport files...")

	var total int
	if err := GetPostgresDB().QueryRowContext(ctx, `SELECT COUNT(*) FROM passenger_passports WHERE deleted_at IS NULL AND image IS NOT NULL`).Scan(&total); err != nil {
		return fmt.Errorf("count failed: %w", err)
	}
	fmt.Printf("  Found %d passports with images\n", total)

	rows, err := GetPostgresDB().QueryContext(ctx, `
		SELECT id, booking_id, image FROM passenger_passports WHERE deleted_at IS NULL AND image IS NOT NULL
	`)
	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var migrated, skipped, errors, current int

	for rows.Next() {
		current++
		var id, bookingID string
		var image sql.NullString

		if err := rows.Scan(&id, &bookingID, &image); err != nil {
			fmt.Printf("  ERROR scanning row: %v\n", err)
			errors++
			continue
		}

		if !image.Valid || isAlreadyMigrated(image.String) {
			skipped++
			continue
		}

		fmt.Printf("\r  [%d/%d] Migrating passport %s...", current, total, id[:8])
		newURL, err := migrateFileField(ctx, image.String, "bookings", bookingID, "passport")
		if err != nil {
			fmt.Printf("\n  ERROR migrating passport %s: %v\n", id, err)
			errors++
			continue
		}

		if newURL != "" {
			_, err := GetPostgresDB().ExecContext(ctx, `
				UPDATE passenger_passports SET image = $1, updated_at = NOW() WHERE id = $2
			`, newURL, id)
			if err != nil {
				fmt.Printf("\n  ERROR updating passport %s: %v\n", id, err)
				errors++
			} else {
				migrated++
				fmt.Printf("\r  [%d/%d] ✓ Migrated passport %s\n", current, total, id[:8])
			}
		}

		time.Sleep(DelayBetweenRecords)
	}

	fmt.Printf("\n  Summary: Migrated: %d, Skipped: %d, Errors: %d\n", migrated, skipped, errors)
	return nil
}

