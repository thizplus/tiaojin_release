package config

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
)

type Config struct {
	App         AppConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	S3          S3Config
	GoogleCloud GoogleCloudConfig
	LLM         LLMConfig
}

type AppConfig struct {
	Name string
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret string
}

// S3Config for S3-compatible storage (Cloudflare R2, iDrive e2, AWS S3, etc.)
type S3Config struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	CDNUrl          string
}

type GoogleCloudConfig struct {
	ProjectID       string
	CredentialsPath string
	VisionEnabled   bool
}

type LLMConfig struct {
	Provider string // gemini, openai, claude
	APIKey   string
	Model    string // gemini-1.5-flash, gpt-4o-mini, claude-3-5-haiku
	Enabled  bool
}

func LoadConfig() (*Config, error) {
	// Try to load .env file, but don't fail if it doesn't exist (for production)
	_ = godotenv.Load()

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	visionEnabled, _ := strconv.ParseBool(getEnv("GOOGLE_VISION_ENABLED", "false"))
	llmEnabled, _ := strconv.ParseBool(getEnv("LLM_ENABLED", "false"))

	config := &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "GoFiber Template"),
			Port: getEnv("APP_PORT", "3000"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "gofiber_template"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "your-secret-key"),
		},
		S3: S3Config{
			Endpoint:        getEnv("S3_ENDPOINT", ""),
			Region:          getEnv("S3_REGION", "us-east-1"),
			AccessKeyID:     getEnv("S3_ACCESS_KEY", ""),
			SecretAccessKey: getEnv("S3_SECRET_KEY", ""),
			Bucket:          getEnv("S3_BUCKET", ""),
			CDNUrl:          getEnv("S3_PUBLIC_URL", ""),
		},
		GoogleCloud: GoogleCloudConfig{
			ProjectID:       getEnv("GOOGLE_CLOUD_PROJECT_ID", ""),
			CredentialsPath: getEnv("GOOGLE_APPLICATION_CREDENTIALS", ""),
			VisionEnabled:   visionEnabled,
		},
		LLM: LLMConfig{
			Provider: getEnv("LLM_PROVIDER", "gemini"),
			APIKey:   getEnv("LLM_API_KEY", ""),
			Model:    getEnv("LLM_MODEL", "gemini-1.5-flash"),
			Enabled:  llmEnabled,
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}