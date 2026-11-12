package handlers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/repositories"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/utils"
)

type JobHandler struct {
	jobService services.JobService
}

func NewJobHandler(jobService services.JobService) *JobHandler {
	return &JobHandler{
		jobService: jobService,
	}
}

func (h *JobHandler) CreateJob(c *fiber.Ctx) error {
	var req dto.CreateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	if err := utils.ValidateStruct(&req); err != nil {
		errors := utils.GetValidationErrors(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  errors,
		})
	}

	job, err := h.jobService.CreateJob(c.Context(), &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Job creation failed", err)
	}

	jobResponse := dto.JobToJobResponse(job)
	return utils.SuccessResponse(c, "Job created successfully", jobResponse)
}

func (h *JobHandler) GetJob(c *fiber.Ctx) error {
	jobIDStr := c.Params("id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid job ID")
	}

	job, err := h.jobService.GetJob(c.Context(), jobID)
	if err != nil {
		return utils.NotFoundResponse(c, "Job not found")
	}

	jobResponse := dto.JobToJobResponse(job)
	return utils.SuccessResponse(c, "Job retrieved successfully", jobResponse)
}

func (h *JobHandler) UpdateJob(c *fiber.Ctx) error {
	jobIDStr := c.Params("id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid job ID")
	}

	var req dto.UpdateJobRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ValidationErrorResponse(c, "Invalid request body")
	}

	job, err := h.jobService.UpdateJob(c.Context(), jobID, &req)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Job update failed", err)
	}

	jobResponse := dto.JobToJobResponse(job)
	return utils.SuccessResponse(c, "Job updated successfully", jobResponse)
}

func (h *JobHandler) DeleteJob(c *fiber.Ctx) error {
	jobIDStr := c.Params("id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid job ID")
	}

	err = h.jobService.DeleteJob(c.Context(), jobID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Job deletion failed", err)
	}

	return utils.SuccessResponse(c, "Job deleted successfully", nil)
}

func (h *JobHandler) StartJob(c *fiber.Ctx) error {
	jobIDStr := c.Params("id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid job ID")
	}

	err = h.jobService.StartJob(c.Context(), jobID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to start job", err)
	}

	return utils.SuccessResponse(c, "Job started successfully", nil)
}

func (h *JobHandler) StopJob(c *fiber.Ctx) error {
	jobIDStr := c.Params("id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid job ID")
	}

	err = h.jobService.StopJob(c.Context(), jobID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Failed to stop job", err)
	}

	return utils.SuccessResponse(c, "Job stopped successfully", nil)
}

func (h *JobHandler) ListJobs(c *fiber.Ctx) error {
	pagination := utils.ParsePaginationParams(c, 10)

	// Check if filters are provided
	status := c.Query("status")
	isActiveStr := c.Query("is_active")
	jobType := c.Query("job_type")

	// If any filter is provided, use filtered list
	if status != "" || isActiveStr != "" || jobType != "" {
		return h.listJobsWithFilter(c, pagination.Offset, pagination.Limit, status, isActiveStr, jobType)
	}

	// Default behavior - list all jobs
	jobs, total, err := h.jobService.ListJobs(c.Context(), pagination.Offset, pagination.Limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve jobs", err)
	}

	jobResponses := make([]dto.JobResponse, len(jobs))
	for i, job := range jobs {
		jobResponses[i] = *dto.JobToJobResponse(job)
	}

	response := &dto.JobListResponse{
		Jobs: jobResponses,
		Meta: utils.CalculatePagination(total, pagination.Offset, pagination.Limit),
	}

	return utils.SuccessResponse(c, "Jobs retrieved successfully", response)
}

func (h *JobHandler) listJobsWithFilter(c *fiber.Ctx, offset, limit int, status, isActiveStr, jobType string) error {
	// Import repositories package for JobListFilter
	filter := struct {
		Status   *string
		IsActive *bool
		JobType  *string
	}{}

	// Parse filters
	if status != "" {
		filter.Status = &status
	}

	if isActiveStr != "" {
		isActive, err := strconv.ParseBool(isActiveStr)
		if err != nil {
			return utils.ValidationErrorResponse(c, "Invalid is_active parameter")
		}
		filter.IsActive = &isActive
	}

	if jobType != "" {
		filter.JobType = &jobType
	}

	// Convert to repositories.JobListFilter
	repoFilter := repositories.JobListFilter{
		Status:   filter.Status,
		IsActive: filter.IsActive,
		JobType:  filter.JobType,
	}

	jobs, total, err := h.jobService.ListJobsWithFilter(c.Context(), offset, limit, repoFilter)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to retrieve filtered jobs", err)
	}

	jobResponses := make([]dto.JobResponse, len(jobs))
	for i, job := range jobs {
		jobResponses[i] = *dto.JobToJobResponse(job)
	}

	response := &dto.JobListResponse{
		Jobs: jobResponses,
		Meta: dto.PaginationMeta{
			Total:       total,
			HasMore:     offset+limit < int(total),
			CurrentPage: (offset / limit) + 1,
			TotalPages:  int((total + int64(limit) - 1) / int64(limit)),
			PageSize:    limit,
			Limit:       limit,
			Offset:      offset,
		},
	}

	return utils.SuccessResponse(c, "Filtered jobs retrieved successfully", response)
}

// ExecuteJob executes a job immediately (Super admin only)
func (h *JobHandler) ExecuteJob(c *fiber.Ctx) error {
	jobIDStr := c.Params("id")
	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		return utils.ValidationErrorResponse(c, "Invalid job ID")
	}

	// Get the job first
	job, err := h.jobService.GetJob(c.Context(), jobID)
	if err != nil {
		return utils.NotFoundResponse(c, "Job not found")
	}

	// Execute the job immediately
	err = h.jobService.ExecuteJob(c.Context(), job)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to execute job", err)
	}

	return utils.SuccessResponse(c, "Job executed successfully", nil)
}