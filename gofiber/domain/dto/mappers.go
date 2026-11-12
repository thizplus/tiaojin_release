package dto

import (
	"gofiber-template/domain/models"
	"gofiber-template/pkg/cdn"
)

func UserToUserResponse(user *models.User) *UserResponse {
	if user == nil {
		return nil
	}

	var companyProfile *CompanyProfileResponse
	if user.CompanyProfile != nil {
		companyProfile = CompanyProfileToResponse(user.CompanyProfile)
	}

	return &UserResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		Role:           string(user.Role),
		Status:         string(user.Status),
		CreatedAt:      user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:      user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		CompanyProfile: companyProfile,
	}
}

func RegisterRequestToUser(req *RegisterRequest) *models.User {
	return &models.User{
		Username: req.Username,
		Email:    req.Email,
		Role:     models.UserRoleAgent,
		Status:   models.UserStatusPending,
	}
}

func CompanyProfileToResponse(profile *models.CompanyProfile) *CompanyProfileResponse {
	if profile == nil {
		return nil
	}
	return &CompanyProfileResponse{
		ID:        profile.ID,
		UserID:    profile.UserID,
		Name:      profile.Name,
		TaxID:     profile.TaxID,
		Phone:     profile.Phone,
		Address:   profile.Address,
		Province:  profile.Province,
		Zipcode:   profile.Zipcode,
		ATTA:      profile.ATTA,
		DBD:       profile.DBD,
		PP20:      profile.PP20,
		Logo:      profile.Logo,
		CreatedAt: profile.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: profile.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func UpdateUserRequestToUser(req *UpdateUserRequest) *models.User {
	return &models.User{
		// Note: Only basic fields like Email, Role, Status can be updated on User
		// CompanyProfile fields are handled separately
	}
}

func TaskToTaskResponse(task *models.Task, user *models.User) *TaskResponse {
	if task == nil {
		return nil
	}
	taskResp := &TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		UserID:      task.UserID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
	if user != nil {
		taskResp.User = *UserToUserResponse(user)
	}
	return taskResp
}

func CreateTaskRequestToTask(req *CreateTaskRequest) *models.Task {
	return &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}
}

func UpdateTaskRequestToTask(req *UpdateTaskRequest) *models.Task {
	return &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}
}

func JobToJobResponse(job *models.Job) *JobResponse {
	if job == nil {
		return nil
	}
	return &JobResponse{
		ID:        job.ID,
		Name:      job.Name,
		CronExpr:  job.CronExpr,
		Payload:   job.Payload,
		Status:    job.Status,
		LastRun:   job.LastRun,
		NextRun:   job.NextRun,
		IsActive:  job.IsActive,
		CreatedAt: job.CreatedAt,
		UpdatedAt: job.UpdatedAt,
	}
}

func CreateJobRequestToJob(req *CreateJobRequest) *models.Job {
	return &models.Job{
		Name:     req.Name,
		CronExpr: req.CronExpr,
		Payload:  req.Payload,
	}
}

func UpdateJobRequestToJob(req *UpdateJobRequest) *models.Job {
	return &models.Job{
		Name:     req.Name,
		CronExpr: req.CronExpr,
		Payload:  req.Payload,
		IsActive: req.IsActive,
	}
}

func FileToFileResponse(file *models.File) *FileResponse {
	if file == nil {
		return nil
	}
	return &FileResponse{
		ID:        file.ID,
		FileName:  file.FileName,
		FileSize:  file.FileSize,
		MimeType:  file.MimeType,
		URL:       cdn.GetFullFileURL(file.URL), // Construct full URL for frontend
		CDNPath:   file.CDNPath,
		UserID:    file.UserID,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}
}