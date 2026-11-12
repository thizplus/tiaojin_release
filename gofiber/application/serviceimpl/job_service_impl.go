package serviceimpl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/models"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/scheduler"
	"time"

	"github.com/google/uuid"
)

type JobServiceImpl struct {
	jobRepo         repositories.JobRepository
	scheduler       scheduler.EventScheduler
	bookingService  services.BookingService
	seatHoldService services.SeatHoldService
}

func NewJobService(jobRepo repositories.JobRepository, scheduler scheduler.EventScheduler) services.JobService {
	return &JobServiceImpl{
		jobRepo:   jobRepo,
		scheduler: scheduler,
	}
}

func (s *JobServiceImpl) SetBookingService(bookingService services.BookingService) {
	s.bookingService = bookingService
}

func (s *JobServiceImpl) SetSeatHoldService(seatHoldService services.SeatHoldService) {
	s.seatHoldService = seatHoldService
}

func (s *JobServiceImpl) CreateJob(ctx context.Context, req *dto.CreateJobRequest) (*models.Job, error) {
	if err := scheduler.ValidateCronExpression(req.CronExpr); err != nil {
		return nil, fmt.Errorf("invalid cron expression: %v", err)
	}

	existingJob, _ := s.jobRepo.GetByName(ctx, req.Name)
	if existingJob != nil {
		return nil, errors.New("job with this name already exists")
	}

	nextRun, err := scheduler.GetNextRunTime(req.CronExpr)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate next run time: %v", err)
	}

	job := &models.Job{
		ID:        uuid.New(),
		Name:      req.Name,
		CronExpr:  req.CronExpr,
		Payload:   req.Payload,
		Status:    "active",
		NextRun:   nextRun,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.jobRepo.Create(ctx, job)
	if err != nil {
		return nil, err
	}

	err = s.scheduler.AddJob(job.ID.String(), req.CronExpr, func() {
		s.ExecuteJob(context.Background(), job)
	})
	if err != nil {
		s.jobRepo.Delete(ctx, job.ID)
		return nil, fmt.Errorf("failed to schedule job: %v", err)
	}

	return job, nil
}

func (s *JobServiceImpl) GetJob(ctx context.Context, jobID uuid.UUID) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, errors.New("job not found")
	}
	return job, nil
}

func (s *JobServiceImpl) UpdateJob(ctx context.Context, jobID uuid.UUID, req *dto.UpdateJobRequest) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return nil, errors.New("job not found")
	}

	needsReschedule := false

	if req.Name != "" {
		job.Name = req.Name
	}
	if req.CronExpr != "" {
		if err := scheduler.ValidateCronExpression(req.CronExpr); err != nil {
			return nil, fmt.Errorf("invalid cron expression: %v", err)
		}
		job.CronExpr = req.CronExpr
		needsReschedule = true
	}
	if req.Payload != "" {
		job.Payload = req.Payload
	}
	if req.IsActive != job.IsActive {
		job.IsActive = req.IsActive
		needsReschedule = true
	}

	if needsReschedule {
		s.scheduler.RemoveJob(jobID.String())
		if job.IsActive {
			nextRun, err := scheduler.GetNextRunTime(job.CronExpr)
			if err != nil {
				return nil, fmt.Errorf("failed to calculate next run time: %v", err)
			}
			job.NextRun = nextRun

			err = s.scheduler.AddJob(jobID.String(), job.CronExpr, func() {
				s.ExecuteJob(context.Background(), job)
			})
			if err != nil {
				return nil, fmt.Errorf("failed to reschedule job: %v", err)
			}
		}
	}

	job.UpdatedAt = time.Now()

	err = s.jobRepo.Update(ctx, jobID, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (s *JobServiceImpl) DeleteJob(ctx context.Context, jobID uuid.UUID) error {
	_, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return errors.New("job not found")
	}

	s.scheduler.RemoveJob(jobID.String())

	return s.jobRepo.Delete(ctx, jobID)
}

func (s *JobServiceImpl) DeleteJobByName(ctx context.Context, jobName string) error {
	job, err := s.jobRepo.GetByName(ctx, jobName)
	if err != nil {
		return errors.New("job not found")
	}

	s.scheduler.RemoveJob(job.ID.String())

	return s.jobRepo.Delete(ctx, job.ID)
}

func (s *JobServiceImpl) ListJobs(ctx context.Context, offset, limit int) ([]*models.Job, int64, error) {
	jobs, err := s.jobRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.jobRepo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	return jobs, count, nil
}

func (s *JobServiceImpl) ListJobsWithFilter(ctx context.Context, offset, limit int, filter repositories.JobListFilter) ([]*models.Job, int64, error) {
	jobs, err := s.jobRepo.ListWithFilter(ctx, offset, limit, filter)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.jobRepo.CountWithFilter(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return jobs, count, nil
}

func (s *JobServiceImpl) StartJob(ctx context.Context, jobID uuid.UUID) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return errors.New("job not found")
	}

	if job.IsActive {
		return errors.New("job is already active")
	}

	job.IsActive = true
	job.UpdatedAt = time.Now()

	nextRun, err := scheduler.GetNextRunTime(job.CronExpr)
	if err != nil {
		return fmt.Errorf("failed to calculate next run time: %v", err)
	}
	job.NextRun = nextRun

	err = s.scheduler.AddJob(jobID.String(), job.CronExpr, func() {
		s.ExecuteJob(context.Background(), job)
	})
	if err != nil {
		return fmt.Errorf("failed to start job: %v", err)
	}

	return s.jobRepo.Update(ctx, jobID, job)
}

func (s *JobServiceImpl) StopJob(ctx context.Context, jobID uuid.UUID) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return errors.New("job not found")
	}

	if !job.IsActive {
		return errors.New("job is already inactive")
	}

	job.IsActive = false
	job.UpdatedAt = time.Now()

	s.scheduler.RemoveJob(jobID.String())

	return s.jobRepo.Update(ctx, jobID, job)
}

func (s *JobServiceImpl) ExecuteJob(ctx context.Context, job *models.Job) error {
	now := time.Now()

	fmt.Printf("Executing job: %s at %s\n", job.Name, now.Format(time.RFC3339))

	job.LastRun = &now
	job.Status = "running"

	nextRun, err := scheduler.GetNextRunTime(job.CronExpr)
	if err == nil {
		job.NextRun = nextRun
	}

	s.jobRepo.Update(ctx, job.ID, job)

	// Parse job payload to determine job type and execute appropriate handler
	var payload map[string]interface{}
	if job.Payload != "" {
		if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
			fmt.Printf("Failed to parse job payload: %v\n", err)
			job.Status = "failed"
			job.UpdatedAt = time.Now()
			s.jobRepo.Update(ctx, job.ID, job)
			return fmt.Errorf("failed to parse job payload: %v", err)
		}
	}

	// Execute job based on type
	jobType, exists := payload["type"]
	if exists {
		switch jobType {
		case "booking_expiry":
			if bookingIDStr, ok := payload["booking_id"].(string); ok {
				if bookingID, err := uuid.Parse(bookingIDStr); err == nil {
					fmt.Printf("Processing booking expiry for booking ID: %s\n", bookingID)
					if s.bookingService != nil {
						if expireErr := s.bookingService.ExpireBooking(ctx, bookingID); expireErr != nil {
							fmt.Printf("Failed to expire booking %s: %v\n", bookingID, expireErr)
							job.Status = "failed"
							job.UpdatedAt = time.Now()
							s.jobRepo.Update(ctx, job.ID, job)
							return fmt.Errorf("failed to expire booking %s: %v", bookingID, expireErr)
						}
						fmt.Printf("Successfully expired booking %s\n", bookingID)
						// For expiry jobs, the job is deleted by the service, so we return early
						return nil
					} else {
						fmt.Printf("BookingService not available for booking expiry job\n")
						job.Status = "failed"
						job.UpdatedAt = time.Now()
						s.jobRepo.Update(ctx, job.ID, job)
						return fmt.Errorf("booking service not available")
					}
				} else {
					fmt.Printf("Invalid booking ID format in job payload: %s\n", bookingIDStr)
					job.Status = "failed"
					job.UpdatedAt = time.Now()
					s.jobRepo.Update(ctx, job.ID, job)
					return fmt.Errorf("invalid booking ID format: %s", bookingIDStr)
				}
			} else {
				fmt.Printf("Missing booking_id in job payload\n")
				job.Status = "failed"
				job.UpdatedAt = time.Now()
				s.jobRepo.Update(ctx, job.ID, job)
				return fmt.Errorf("missing booking_id in job payload")
			}
		case "seat_hold_expiry":
			if seatHoldIDStr, ok := payload["seat_hold_id"].(string); ok {
				if seatHoldID, err := uuid.Parse(seatHoldIDStr); err == nil {
					fmt.Printf("Processing seat hold expiry for seat hold ID: %s\n", seatHoldID)
					if s.seatHoldService != nil {
						if expireErr := s.seatHoldService.ExpireSeatHold(ctx, seatHoldID); expireErr != nil {
							fmt.Printf("Failed to expire seat hold %s: %v\n", seatHoldID, expireErr)
							job.Status = "failed"
							job.UpdatedAt = time.Now()
							s.jobRepo.Update(ctx, job.ID, job)
							return fmt.Errorf("failed to expire seat hold %s: %v", seatHoldID, expireErr)
						}
						fmt.Printf("Successfully expired seat hold %s\n", seatHoldID)
						// For expiry jobs, the job is deleted by the service, so we return early
						return nil
					} else {
						fmt.Printf("SeatHoldService not available for seat hold expiry job\n")
						job.Status = "failed"
						job.UpdatedAt = time.Now()
						s.jobRepo.Update(ctx, job.ID, job)
						return fmt.Errorf("seat hold service not available")
					}
				} else {
					fmt.Printf("Invalid seat hold ID format in job payload: %s\n", seatHoldIDStr)
					job.Status = "failed"
					job.UpdatedAt = time.Now()
					s.jobRepo.Update(ctx, job.ID, job)
					return fmt.Errorf("invalid seat hold ID format: %s", seatHoldIDStr)
				}
			} else {
				fmt.Printf("Missing seat_hold_id in job payload\n")
				job.Status = "failed"
				job.UpdatedAt = time.Now()
				s.jobRepo.Update(ctx, job.ID, job)
				return fmt.Errorf("missing seat_hold_id in job payload")
			}
		case "booking_expiry_cleanup":
			fmt.Printf("Processing backup booking expiry cleanup job\n")
			if s.bookingService != nil {
				// Use existing ProcessExpiredBookings method for backup cleanup
				result, err := s.bookingService.ProcessExpiredBookings(ctx, false) // dryRun = false
				if err != nil {
					fmt.Printf("Failed to process expired bookings cleanup: %v\n", err)
					job.Status = "failed"
					job.UpdatedAt = time.Now()
					s.jobRepo.Update(ctx, job.ID, job)
					return fmt.Errorf("failed to process expired bookings cleanup: %v", err)
				}
				fmt.Printf("Backup cleanup processed %d expired bookings\n", result.ProcessedCount)
			} else {
				fmt.Printf("BookingService not available for cleanup job\n")
				job.Status = "failed"
				job.UpdatedAt = time.Now()
				s.jobRepo.Update(ctx, job.ID, job)
				return fmt.Errorf("booking service not available")
			}
		default:
			fmt.Printf("Unknown job type: %s\n", jobType)
		}
	} else {
		fmt.Printf("No job type specified in payload, executing as generic job\n")
	}

	job.Status = "completed"
	job.UpdatedAt = time.Now()

	return s.jobRepo.Update(ctx, job.ID, job)
}
