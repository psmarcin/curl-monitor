package main

import (
	"common"
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
		// TODO: handle error properly
		t.Channel.Publish(
			common.JobExchangeName,
			common.CommandRunRoutingKey,
			true,
			false,
			amqp.Publishing{
				Body: payload,
			},
		)
	}

	// TODO: get rid of unnecessary return value
	return "some string"
}
