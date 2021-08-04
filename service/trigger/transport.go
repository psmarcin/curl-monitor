package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/streadway/amqp"
	"time"
)

type handlerInput struct {
	Name string `json:"name"`
}

type handlerOutput struct {
	Timestamp time.Time `json:"timestamp"`
}

func makeHandlerEndpoint(service triggerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(handlerInput)
		res := service.Handler(req)

		return res, nil
	}
}

func decodetriggerAMQPHandler(ctx context.Context, delivery *amqp.Delivery) (interface{}, error) {
	var request handlerInput
	err := json.Unmarshal(delivery.Body, &request)
	if err != nil {
		delivery.Reject(true)
		return nil, err
	}
	return request, nil
}
