package main

import (
	"common"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os/exec"
	"time"
)

const (
	command = "curl"
)

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
	output, err := exec.Command(command, payload.Command).Output()
	if err != nil {
		result := Output{
			JobUuid:   payload.Uuid,
			Output:    string(output),
			Type:      "error",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		resultRaw, marshalErr := json.Marshal(result)
		if marshalErr != nil {
			log.Printf("error while marshaling: %s", marshalErr)
		}
		// TODO: handle error properly
		t.Channel.Publish(
			common.JobExchangeName,
			common.ResultRoutingKey,
			true,
			false,
			amqp.Publishing{Body: resultRaw},
		)
		log.Printf("error while running command curl %s: %s - error: %s", payload.Command, output, err)
		return ""
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
		common.JobExchangeName,
		common.ResultRoutingKey,
		true,
		false,
		amqp.Publishing{Body: resultRaw},
	)
	log.Printf("command curl %s: %s", payload.Command, output)
	return ""
}
