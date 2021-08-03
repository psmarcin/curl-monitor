package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"trigger/job"
)

type TriggerService interface {
	Handler(string) error
}

type triggerService struct {
	Channel *amqp.Channel
	Job     job.JobService
}

func (t triggerService) Handler(payload handlerInput) string {
	jobs, err := t.Job.ListJob()

	if err != nil {
		log.Printf("err: %s\n", err)
		return ""
	}

	for _, job := range jobs {
		payload, _ := json.Marshal(job)
		t.Channel.Publish(
			"CM.Job",
			"commandrun",
			true,
			false,
			amqp.Publishing{
				Body: payload,
			},
		)
	}

	return "some string"
}
