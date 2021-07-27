package main

import (
	"context"
	"github.com/google/uuid"
	"job/db"
	"time"
)

type JobService interface {
	GetJob(string) (Job, error)
	ListJob() ([]Job, error)
	CreateJob(request createJobRequest) (Job, error)
	UpdateJob(request updateJobRequest) (Job, error)
}

type jobService struct {
	DB *db.Queries
}

func (j jobService) GetJob(id string) (Job, error) {
	ctx := context.Background()
	result, err := j.DB.GetJob(ctx, uuid.MustParse(id))
	if err != nil {
		return Job{}, err
	}
	return Job{
		Uuid:    result.Uuid.String(),
		Name:    result.Name,
		Command: result.Command,
	}, nil
}

func (j jobService) ListJob() ([]Job, error) {
	ctx := context.Background()
	var jobs []Job
	result, err := j.DB.ListJob(ctx, db.ListJobParams{
		CreatedAt: time.Time{},
		Limit:     5,
	})
	if err != nil {
		return jobs, err
	}
	for _, job := range result {
		jobs = append(jobs, Job{
			Uuid:      job.Uuid.String(),
			Name:      job.Name,
			Command:   job.Command,
			CreatedAt: job.CreatedAt,
			UpdatedAt: job.UpdatedAt,
		})
	}

	return jobs, nil
}

func (j jobService) CreateJob(payload createJobRequest) (Job, error) {
	ctx := context.Background()
	result, err := j.DB.CreateJob(ctx, db.CreateJobParams{
		Uuid:      uuid.New(),
		Name:      payload.Name,
		Command:   payload.Command,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return Job{}, err
	}
	return Job{
		Uuid:      result.Uuid.String(),
		Name:      result.Name,
		Command:   result.Command,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

func (j jobService) UpdateJob(payload updateJobRequest) (Job, error) {
	ctx := context.Background()
	result, err := j.DB.UpdateJob(ctx, db.UpdateJobParams{
		Uuid:      uuid.MustParse(payload.Id),
		Name:      payload.Name,
		Command:   payload.Command,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return Job{}, err
	}
	return Job{
		Uuid:      result.Uuid.String(),
		Name:      result.Name,
		Command:   result.Command,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}, nil
}
