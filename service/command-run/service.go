package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os/exec"
	"time"
)

type TriggerService interface {
	Handler(string) error
}

type triggerService struct {
	Channel *amqp.Channel
}

type Output struct {
	Uuid      string
	JobUuid   string
	Output    string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (t triggerService) Handler(payload handlerInput) string {
	log.Printf("command run received message: %s", payload)
	output, err := exec.Command("curl", payload.Command).Output()
	if err != nil {
		result := Output{
			JobUuid:   payload.Uuid,
			Output:    string(output),
			Type:      "error",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		resultRaw, err := json.Marshal(result)
		if err != nil {
			log.Printf("error while marshaling: %s", err)
		}
		t.Channel.Publish(
			"CM.Job",
			"output",
			true,
			false,
			amqp.Publishing{Body: resultRaw},
		)
		log.Printf("error while running command curl %s: %s - error: %s", payload.Command, output, err)
	}
	result := Output{
		JobUuid:   payload.Uuid,
		Output:    string(output),
		Type:      "success",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	resultRaw, err := json.Marshal(result)
	if err != nil {
		log.Printf("error while marshaling: %s", err)
	}
	t.Channel.Publish(
		"CM.Job",
		"output",
		true,
		false,
		amqp.Publishing{Body: resultRaw},
	)
	log.Printf("command curl %s: %s", payload.Command, output)
	return ""
}
