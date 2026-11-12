package handlers

import (
	"gofiber-template/domain/ports"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/config"
)

// Services contains all the services needed for handlers
type Services struct {
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
	SeatHoldService          services.SeatHoldService
	SeatAvailabilityService  services.SeatAvailabilityService
	VisionService            services.VisionService
	LLMService               services.LLMService
	APIKeyService            services.APIKeyService
	Storage                    ports.StoragePort // S3-compatible storage for presigned URLs
	Config                     *config.Config
	CompanyBankAccountService  services.CompanyBankAccountService
}

// Handlers contains all HTTP handlers
type Handlers struct {
	AuthHandler      *AuthHandler
	UserHandler      *UserHandler
	CompanyHandler   *CompanyHandler
	CityHandler      *CityHandler
	BannerHandler    *BannerHandler
	SalesHandler     *SalesHandler
	TourHandler      *TourHandler
	PeriodHandler    *PeriodHandler
	BookingHandler   *BookingHandler
	TaskHandler      *TaskHandler
	FileHandler      *FileHandler
	JobHandler       *JobHandler
	AgentHandler     *AgentHandler     // NEW
	AnalyticsHandler *AnalyticsHandler // NEW
	ConfigHandler            *ConfigHandler            // NEW
	SeatHoldHandler          *SeatHoldHandler          // NEW
	SeatAvailabilityHandler  *SeatAvailabilityHandler  // NEW
	DebugHandler             *DebugHandler             // NEW - Vision API Debug
	APIKeyHandler            *APIKeyHandler            // NEW - API Key Management
	PublicTourHandler        *PublicTourHandler        // NEW - Public API for external access
	StorageHandler           *StorageHandler           // Storage file access with presigned URLs
	BankAccountHandler       *CompanyBankAccountHandler // บัญชีธนาคารบริษัท
}

// NewHandlers creates a new instance of Handlers with all dependencies
func NewHandlers(services *Services) *Handlers {
	return &Handlers{
		AuthHandler:      NewAuthHandler(services.AuthService, services.UserService),
		UserHandler:      NewUserHandler(services.UserService),
		CompanyHandler:   NewCompanyHandler(services.CompanyService),
		CityHandler:      NewCityHandler(services.CityService),
		BannerHandler:    NewBannerHandler(services.BannerService),
		SalesHandler:     NewSalesHandler(services.SalesService, services.BookingService, services.FileService),
		TourHandler:      NewTourHandler(services.TourService),
		PeriodHandler:    NewPeriodHandler(services.PeriodService),
		BookingHandler:   NewBookingHandler(services.BookingService, services.FileService),
		TaskHandler:      NewTaskHandler(services.TaskService),
		FileHandler:      NewFileHandler(services.FileService),
		JobHandler:       NewJobHandler(services.JobService),
		AgentHandler:     NewAgentHandler(services.TourService, services.PeriodService, services.BookingService, services.UserService, services.FileService, services.CompanyService, services.SalesService),
		AnalyticsHandler: NewAnalyticsHandler(services.BookingService, services.SalesService, services.TourService, services.UserService),
		ConfigHandler:           NewConfigHandler(services.ConfigService, services.BookingService),
		SeatHoldHandler:         NewSeatHoldHandler(services.SeatHoldService),
		SeatAvailabilityHandler: NewSeatAvailabilityHandler(services.SeatAvailabilityService),
		DebugHandler:            NewDebugHandler(services.VisionService, services.LLMService, services.Config),
		APIKeyHandler:           NewAPIKeyHandler(services.APIKeyService),
		PublicTourHandler:       NewPublicTourHandler(services.TourService),
		StorageHandler:          NewStorageHandler(services.Storage),
		BankAccountHandler:      NewCompanyBankAccountHandler(services.CompanyBankAccountService),
	}
}