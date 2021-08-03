package main

import (
	"context"
	"github.com/google/uuid"
	"log"
	"result/db"
)

type ResultService interface {
	CreateResult(payload createResultRequest) error
}

type resultService struct {
	DB *db.Queries
}

func (j resultService) CreateResult(payload createResultRequest) error {
	ctx := context.Background()
	jobUuid, err := uuid.Parse(payload.JobUuid)
	if err != nil {
		log.Printf("error while parsing job uuid: %s, %s", payload.JobUuid, err)
		return err
	}
	_, err = j.DB.CreateResult(ctx, db.CreateResultParams{
		Uuid:      uuid.New(),
		JobUuid:   jobUuid,
		Output:    payload.Output,
		Type:      db.Output(payload.Type),
		CreatedAt: payload.CreatedAt,
		UpdatedAt: payload.UpdatedAt,
	})
	if err != nil {
		log.Printf("error while inserting job result: %s, %s", payload.Output, err)
		return err
	}
	return nil
}
