package postgres

import (
	"fmt"
	"log"
	"time"

	"gofiber-template/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewDatabase(config DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %v", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(1 * time.Minute)

	return db, nil
}

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		// Independent models (no foreign keys or primary entities)
		&models.User{},
		&models.City{},
		&models.Sales{},
		&models.Banner{},
		&models.Job{},
		&models.SystemConfig{},      // System configuration storage
		&models.CompanyBankAccount{}, // บัญชีธนาคารบริษัท

		// Level 1 dependencies (depend on User/City)
		&models.Task{},
		&models.File{},
		&models.CompanyProfile{},
		&models.Tour{},
		&models.APIKey{},

		// Level 2 dependencies (depend on Tour)
		&models.Period{},

		// Level 2.5 dependencies (depend on Period) - New stable reference tables
		&models.PeriodRateItem{},
		&models.PeriodRoomType{},
		&models.PeriodAdditionalService{},
		&models.SeatHold{},

		// Level 3 dependencies (depend on User, Tour, Period, Sales)
		&models.Booking{},

		// Level 4 dependencies (child tables depend on Booking)
		&models.BookingRoomType{},
		&models.BookingRateType{},
		&models.BookingAdditionalService{},
		&models.PaymentSlip{},
		&models.PassengerPassport{},
	)
	if err != nil {
		return err
	}

	// Seed default configurations
	return SeedDefaultConfigs(db)
}

// SeedDefaultConfigs creates default system configurations if they don't exist
func SeedDefaultConfigs(db *gorm.DB) error {
	log.Println("[SEED] Checking default system configurations...")

	// Define default configs
	defaultConfigs := []struct {
		ConfigKey   string
		ConfigValue map[string]interface{}
		ConfigType  models.ConfigType
		Description string
	}{
		{
			ConfigKey: "booking.expiration_tiers",
			ConfigValue: map[string]interface{}{
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
			},
			ConfigType:  models.ConfigTypeExpiration,
			Description: "Booking expiration tier configuration - กำหนดระยะเวลาหมดอายุการจองตามความใกล้ของวันเริ่มทัวร์",
		},
		{
			ConfigKey: "booking.agent_period_limits",
			ConfigValue: map[string]interface{}{
				"global_settings": map[string]interface{}{
					"enabled":                false,
					"default_limit_per_period": 10,
					"count_all_statuses":     []string{"pending", "confirmed", "paid"},
				},
				"period_specific_limits": []interface{}{},
				"agent_global_bypasses":  []interface{}{},
				"metadata": map[string]interface{}{
					"version":      "1.0",
					"last_updated": time.Now().Format(time.RFC3339),
					"updated_by":   "system",
				},
			},
			ConfigType:  models.ConfigTypeBooking,
			Description: "Agent booking limit configuration per period - จำกัดจำนวนการจองของ Agent ต่อ Period",
		},
	}

	for _, cfg := range defaultConfigs {
		// Check if config already exists
		var existingConfig models.SystemConfig
		result := db.Where("config_key = ?", cfg.ConfigKey).First(&existingConfig)

		if result.Error == nil {
			log.Printf("[SEED] Config '%s' already exists, skipping", cfg.ConfigKey)
			continue
		}

		// Create new config (ConfigValueJSON is map[string]interface{})
		newConfig := models.SystemConfig{
			ConfigKey:   cfg.ConfigKey,
			ConfigValue: models.ConfigValueJSON(cfg.ConfigValue),
			ConfigType:  cfg.ConfigType,
			Description: &cfg.Description,
			IsActive:    true,
		}

		if err := db.Create(&newConfig).Error; err != nil {
			log.Printf("[SEED] Error creating config '%s': %v", cfg.ConfigKey, err)
			return fmt.Errorf("failed to seed config '%s': %w", cfg.ConfigKey, err)
		}

		log.Printf("[SEED] Created default config: %s", cfg.ConfigKey)
	}

	log.Println("[SEED] Default system configurations check completed")
	return nil
}
