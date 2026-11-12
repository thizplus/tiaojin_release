package di

import (
	"context"
	"encoding/json"
	"fmt"
	"gofiber-template/application/serviceimpl"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/ports"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/infrastructure/postgres"
	"gofiber-template/infrastructure/redis"
	visionservice "gofiber-template/infrastructure/services"
	emailpkg "gofiber-template/infrastructure/email"
	"gofiber-template/infrastructure/slipverification"
	"gofiber-template/infrastructure/storage"
	"gofiber-template/interfaces/api/handlers"
	"gofiber-template/pkg/config"
	"gofiber-template/pkg/scheduler"
	"log"
	"time"

	"gorm.io/gorm"
)

type Container struct {
	// Configuration
	Config *config.Config

	// Infrastructure
	DB             *gorm.DB
	RedisClient    *redis.RedisClient
	Storage          ports.StoragePort            // S3-compatible storage (iDrive e2, R2, AWS S3, etc.)
	SlipVerification ports.SlipVerificationPort  // SlipOK verification (optional)
	EmailPort        ports.EmailPort             // SMTP email (optional)
	EventScheduler scheduler.EventScheduler

	// Repositories
	UserRepository                     repositories.UserRepository
	CompanyProfileRepository           repositories.CompanyProfileRepository
	CityRepository                     repositories.CityRepository
	BannerRepository                   repositories.BannerRepository
	SalesRepository                    repositories.SalesRepository
	BookingRepository                  repositories.BookingRepository
	BookingRateTypeRepository          repositories.BookingRateTypeRepository
	BookingRoomTypeRepository          repositories.BookingRoomTypeRepository
	BookingAdditionalServiceRepository repositories.BookingAdditionalServiceRepository
	PaymentSlipRepository              repositories.PaymentSlipRepository
	PassengerPassportRepository        repositories.PassengerPassportRepository
	TourRepository                     repositories.TourRepository
	PeriodRepository                   repositories.PeriodRepository
	PeriodRateItemRepository           repositories.PeriodRateItemRepository
	PeriodRoomTypeRepository           repositories.PeriodRoomTypeRepository
	PeriodAdditionalServiceRepository  repositories.PeriodAdditionalServiceRepository
	TaskRepository                     repositories.TaskRepository
	FileRepository                     repositories.FileRepository
	JobRepository                      repositories.JobRepository
	SystemConfigRepository             repositories.SystemConfigRepository
	SeatHoldRepository                 repositories.SeatHoldRepository
	APIKeyRepository                   repositories.APIKeyRepository
	CompanyBankAccountRepository       repositories.CompanyBankAccountRepository

	// Services
	AuthService     services.AuthService
	UserService     services.UserService
	CompanyService  services.CompanyService
	CityService     services.CityService
	BannerService   services.BannerService
	SalesService    services.SalesService
	TourService     services.TourService
	PeriodService   services.PeriodService
	BookingService  services.BookingService
	TaskService     services.TaskService
	FileService     services.FileService
	JobService      services.JobService
	ConfigService   services.ConfigService
	CompanyBankAccountService services.CompanyBankAccountService
	SeatHoldService          services.SeatHoldService
	SeatAvailabilityService  services.SeatAvailabilityService
	VisionService            services.VisionService
	LLMService               services.LLMService
	APIKeyService            services.APIKeyService
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) Initialize() error {
	if err := c.initConfig(); err != nil {
		return err
	}

	if err := c.initInfrastructure(); err != nil {
		return err
	}

	if err := c.initRepositories(); err != nil {
		return err
	}

	if err := c.initServices(); err != nil {
		return err
	}

	if err := c.initScheduler(); err != nil {
		return err
	}

	return nil
}

func (c *Container) initConfig() error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	c.Config = cfg
	log.Println("✓ Configuration loaded")
	return nil
}

func (c *Container) initInfrastructure() error {
	// Initialize Database
	dbConfig := postgres.DatabaseConfig{
		Host:     c.Config.Database.Host,
		Port:     c.Config.Database.Port,
		User:     c.Config.Database.User,
		Password: c.Config.Database.Password,
		DBName:   c.Config.Database.DBName,
		SSLMode:  c.Config.Database.SSLMode,
	}

	db, err := postgres.NewDatabase(dbConfig)
	if err != nil {
		return err
	}
	c.DB = db
	log.Println("✓ Database connected")

	// Run migrations
	if err := postgres.Migrate(db); err != nil {
		return err
	}
	log.Println("✓ Database migrated")

	// Initialize Redis
	redisConfig := redis.RedisConfig{
		Host:     c.Config.Redis.Host,
		Port:     c.Config.Redis.Port,
		Password: c.Config.Redis.Password,
		DB:       c.Config.Redis.DB,
	}
	c.RedisClient = redis.NewRedisClient(redisConfig)

	// Test Redis connection
	if err := c.RedisClient.Ping(context.Background()); err != nil {
		log.Printf("Warning: Redis connection failed: %v", err)
	} else {
		log.Println("✓ Redis connected")
	}

	// Initialize S3-compatible Storage (iDrive e2, Cloudflare R2, AWS S3, etc.)
	s3Config := storage.S3Config{
		Endpoint:        c.Config.S3.Endpoint,
		Region:          c.Config.S3.Region,
		AccessKeyID:     c.Config.S3.AccessKeyID,
		SecretAccessKey: c.Config.S3.SecretAccessKey,
		Bucket:          c.Config.S3.Bucket,
		CDNUrl:          c.Config.S3.CDNUrl,
	}
	s3Storage, err := storage.NewS3Storage(s3Config)
	if err != nil {
		log.Printf("Warning: Failed to initialize S3 storage: %v", err)
		log.Println("Storage features will be disabled")
	} else {
		c.Storage = s3Storage
		log.Println("✓ S3 Storage initialized")
	}

	return nil
}

func (c *Container) initRepositories() error {
	c.UserRepository = postgres.NewUserRepository(c.DB)
	c.CompanyProfileRepository = postgres.NewCompanyProfileRepository(c.DB)
	c.CityRepository = postgres.NewCityRepository(c.DB)
	c.BannerRepository = postgres.NewBannerRepository(c.DB)
	c.SalesRepository = postgres.NewSalesRepository(c.DB)
	c.BookingRepository = postgres.NewBookingRepository(c.DB)
	c.BookingRateTypeRepository = postgres.NewBookingRateTypeRepository(c.DB)
	c.BookingRoomTypeRepository = postgres.NewBookingRoomTypeRepository(c.DB)
	c.BookingAdditionalServiceRepository = postgres.NewBookingAdditionalServiceRepository(c.DB)
	c.PaymentSlipRepository = postgres.NewPaymentSlipRepository(c.DB)
	c.PassengerPassportRepository = postgres.NewPassengerPassportRepository(c.DB)
	c.TourRepository = postgres.NewTourRepository(c.DB)
	c.PeriodRepository = postgres.NewPeriodRepository(c.DB)
	c.PeriodRateItemRepository = postgres.NewPeriodRateItemRepository(c.DB)
	c.PeriodRoomTypeRepository = postgres.NewPeriodRoomTypeRepository(c.DB)
	c.PeriodAdditionalServiceRepository = postgres.NewPeriodAdditionalServiceRepository(c.DB)
	c.TaskRepository = postgres.NewTaskRepository(c.DB)
	c.FileRepository = postgres.NewFileRepository(c.DB)
	c.JobRepository = postgres.NewJobRepository(c.DB)
	c.SystemConfigRepository = postgres.NewSystemConfigRepository(c.DB)
	c.SeatHoldRepository = postgres.NewSeatHoldRepository(c.DB)
	c.APIKeyRepository = postgres.NewAPIKeyRepository(c.DB)
	c.CompanyBankAccountRepository = postgres.NewCompanyBankAccountRepository(c.DB)
	log.Println("✓ Repositories initialized")
	return nil
}

func (c *Container) initServices() error {
	c.AuthService = serviceimpl.NewAuthService(c.UserRepository, c.CompanyProfileRepository, c.Config.JWT.Secret)
	c.UserService = serviceimpl.NewUserService(c.UserRepository, c.Config.JWT.Secret)
	c.CompanyService = serviceimpl.NewCompanyService(c.CompanyProfileRepository, c.UserRepository)
	c.CompanyBankAccountService = serviceimpl.NewCompanyBankAccountService(c.CompanyBankAccountRepository)
	c.CityService = serviceimpl.NewCityService(c.CityRepository, c.TourRepository, c.PeriodRepository)
	c.BannerService = serviceimpl.NewBannerService(c.BannerRepository)
	c.SalesService = serviceimpl.NewSalesService(c.SalesRepository, c.BookingRepository, c.UserRepository)
	c.TourService = serviceimpl.NewTourService(c.TourRepository, c.CityRepository, c.PeriodRepository, c.BookingRepository, c.SeatHoldRepository)
	c.PeriodService = serviceimpl.NewPeriodService(
		c.PeriodRepository,
		c.TourRepository,
		c.PeriodRateItemRepository,
		c.PeriodRoomTypeRepository,
		c.PeriodAdditionalServiceRepository,
		c.BookingRepository,
		c.SeatHoldRepository,
	)

	// Initialize ConfigService first
	c.ConfigService = serviceimpl.NewConfigService(c.SystemConfigRepository, c.BookingRepository)

	// Initialize Vision API service if enabled
	if c.Config.GoogleCloud.VisionEnabled {
		visionService, err := visionservice.NewVisionService(c.Config)
		if err != nil {
			log.Printf("Warning: Failed to initialize Vision service: %v", err)
			log.Println("Vision API features will be disabled")
			c.VisionService = nil
		} else {
			c.VisionService = visionService
			log.Println("✓ Vision API service initialized")
		}
	} else {
		log.Println("Vision API disabled by configuration")
		c.VisionService = nil
	}

	// Initialize LLM service if enabled
	if c.Config.LLM.Enabled {
		llmService, err := visionservice.NewLLMService(c.Config)
		if err != nil {
			log.Printf("Warning: Failed to initialize LLM service: %v", err)
			log.Println("LLM features will be disabled")
			c.LLMService = nil
		} else {
			c.LLMService = llmService
			log.Println("✓ LLM service initialized")
		}
	} else {
		log.Println("LLM service disabled by configuration")
		c.LLMService = nil
	}

	// Initialize SlipOK verification service (optional)
	slipOKAdapter := slipverification.NewSlipOKAdapter()
	if slipOKAdapter.IsConfigured() {
		c.SlipVerification = slipOKAdapter
		log.Println("✓ SlipOK verification service initialized")
	} else {
		log.Println("SlipOK not configured (missing SLIPOK_BRANCH_ID or SLIPOK_API_KEY) — slip verification disabled")
		c.SlipVerification = nil
	}

	// Initialize SMTP email service (optional)
	smtpAdapter := emailpkg.NewSMTPAdapter()
	if smtpAdapter.IsConfigured() {
		c.EmailPort = smtpAdapter
		log.Println("✓ SMTP email service initialized")
	} else {
		log.Println("SMTP not configured (missing SMTP_HOST/PORT/USER/PASS/FROM) — email notifications disabled")
		c.EmailPort = nil
	}

	// Initialize EventScheduler and JobService early to resolve dependencies
	c.EventScheduler = scheduler.NewEventScheduler()
	c.JobService = serviceimpl.NewJobService(c.JobRepository, c.EventScheduler)
	log.Println("✓ JobService initialized early for dependency injection")

	// Initialize SeatAvailabilityService first to avoid dependency cycle
	c.SeatAvailabilityService = serviceimpl.NewSeatAvailabilityService(c.PeriodRepository, c.BookingRepository, c.SeatHoldRepository)

	// Initialize BookingService with all dependencies including JobService
	c.BookingService = serviceimpl.NewBookingService(
		c.BookingRepository,
		c.BookingRateTypeRepository,
		c.BookingRoomTypeRepository,
		c.BookingAdditionalServiceRepository,
		c.TourRepository,
		c.PeriodRepository,
		c.UserRepository,
		c.PaymentSlipRepository,
		c.PassengerPassportRepository,
		c.PeriodService,
		c.ConfigService,
		c.JobService, // Now properly injected
		c.VisionService, // Add Vision API service
		c.LLMService,    // Add LLM service
		c.PeriodRateItemRepository,
		c.PeriodRoomTypeRepository,
		c.PeriodAdditionalServiceRepository,
		c.SeatAvailabilityService,
		c.SalesRepository, // Add SalesRepository for sales report
		c.SlipVerification, // SlipOK verification (optional, nil if not configured)
		c.EmailPort,        // SMTP email (optional, nil if not configured)
	)

	// Inject BookingService into JobService for cross-dependency
	c.JobService.SetBookingService(c.BookingService)

	c.TaskService = serviceimpl.NewTaskService(c.TaskRepository, c.UserRepository)
	c.FileService = serviceimpl.NewFileService(c.FileRepository, c.UserRepository, c.Storage)
	c.SeatHoldService = serviceimpl.NewSeatHoldService(c.SeatHoldRepository, c.PeriodRepository, c.UserRepository, c.BookingRepository, c.PeriodService, c.JobService, c.SeatAvailabilityService)
	c.APIKeyService = serviceimpl.NewAPIKeyService(c.APIKeyRepository, c.UserRepository)

	// Inject SeatHoldService into JobService for cross-dependency
	c.JobService.SetSeatHoldService(c.SeatHoldService)
	log.Println("✓ Services initialized")
	return nil
}

func (c *Container) initScheduler() error {
	// EventScheduler and JobService are already initialized in initServices
	// Just start the scheduler
	c.EventScheduler.Start()
	log.Println("✓ Event scheduler started")

	// Load and schedule existing active jobs
	ctx := context.Background()
	jobs, _, err := c.JobService.ListJobs(ctx, 0, 1000)
	if err != nil {
		log.Printf("Warning: Failed to load existing jobs: %v", err)
		return nil
	}

	activeJobCount := 0
	for _, job := range jobs {
		if job.IsActive {
			// Recalculate nextRun with correct timezone when loading jobs
			nextRun, err := scheduler.GetNextRunTime(job.CronExpr)
			if err != nil {
				log.Printf("Warning: Failed to calculate next run time for job %s: %v", job.Name, err)
				continue
			}

			// Update the job in database with corrected nextRun
			updatedJob := *job
			updatedJob.NextRun = nextRun
			if updateErr := c.JobRepository.Update(ctx, job.ID, &updatedJob); updateErr != nil {
				log.Printf("Warning: Failed to update nextRun for job %s: %v", job.Name, updateErr)
			}

			err = c.EventScheduler.AddJob(job.ID.String(), job.CronExpr, func() {
				c.JobService.ExecuteJob(ctx, job)
			})
			if err != nil {
				log.Printf("Warning: Failed to schedule job %s: %v", job.Name, err)
			} else {
				activeJobCount++
				log.Printf("Scheduled job %s, next run: %s", job.Name, nextRun.Format(time.RFC3339))
			}
		}
	}

	if activeJobCount > 0 {
		log.Printf("✓ Scheduled %d active jobs", activeJobCount)
	}

	// Create backup cleanup job for expired bookings
	if err := c.createBackupCleanupJob(ctx); err != nil {
		log.Printf("Warning: Failed to create backup cleanup job: %v", err)
	}

	return nil
}

func (c *Container) createBackupCleanupJob(ctx context.Context) error {
	// Create a backup cleanup job that runs every hour to catch any missed expiry jobs
	jobPayload := map[string]interface{}{
		"type":        "booking_expiry_cleanup",
		"description": "Backup cleanup job to expire bookings that might have been missed by individual expiry jobs",
	}

	payloadBytes, err := json.Marshal(jobPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal job payload: %v", err)
	}

	// Check if backup cleanup job already exists
	jobs, _, err := c.JobService.ListJobs(ctx, 0, 1000)
	if err != nil {
		return fmt.Errorf("failed to list existing jobs: %v", err)
	}

	for _, job := range jobs {
		if job.Name == "backup-booking-expiry-cleanup" {
			log.Println("✓ Backup cleanup job already exists")
			return nil
		}
	}

	createJobReq := &dto.CreateJobRequest{
		Name:     "backup-booking-expiry-cleanup",
		CronExpr: "0 */1 * * *", // Run every hour
		Payload:  string(payloadBytes),
	}

	job, err := c.JobService.CreateJob(ctx, createJobReq)
	if err != nil {
		return fmt.Errorf("failed to create backup cleanup job: %v", err)
	}

	log.Printf("✓ Backup cleanup job created: %s", job.ID)
	return nil
}

func (c *Container) Cleanup() error {
	log.Println("Starting cleanup...")

	// Stop scheduler
	if c.EventScheduler != nil {
		if c.EventScheduler.IsRunning() {
			c.EventScheduler.Stop()
			log.Println("✓ Event scheduler stopped")
		} else {
			log.Println("✓ Event scheduler was already stopped")
		}
	}

	// Close Redis connection
	if c.RedisClient != nil {
		if err := c.RedisClient.Close(); err != nil {
			log.Printf("Warning: Failed to close Redis connection: %v", err)
		} else {
			log.Println("✓ Redis connection closed")
		}
	}

	// Close database connection
	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err == nil {
			if err := sqlDB.Close(); err != nil {
				log.Printf("Warning: Failed to close database connection: %v", err)
			} else {
				log.Println("✓ Database connection closed")
			}
		}
	}

	log.Println("✓ Cleanup completed")
	return nil
}

func (c *Container) GetServices() (services.UserService, services.TaskService, services.FileService, services.JobService) {
	return c.UserService, c.TaskService, c.FileService, c.JobService
}

func (c *Container) GetConfig() *config.Config {
	return c.Config
}

func (c *Container) GetHandlerServices() *handlers.Services {
	return &handlers.Services{
		AuthService:     c.AuthService,
		UserService:     c.UserService,
		CompanyService:  c.CompanyService,
		CityService:     c.CityService,
		BannerService:   c.BannerService,
		SalesService:    c.SalesService,
		TourService:     c.TourService,
		PeriodService:   c.PeriodService,
		BookingService:  c.BookingService,
		TaskService:     c.TaskService,
		FileService:     c.FileService,
		JobService:      c.JobService,
		ConfigService:           c.ConfigService,
		SeatHoldService:         c.SeatHoldService,
		SeatAvailabilityService: c.SeatAvailabilityService,
		VisionService:           c.VisionService,
		LLMService:              c.LLMService,
		APIKeyService:           c.APIKeyService,
		Storage:                    c.Storage,
		Config:                     c.Config,
		CompanyBankAccountService:  c.CompanyBankAccountService,
	}
}
