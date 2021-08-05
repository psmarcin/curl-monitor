package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/streadway/amqp"
	"time"
)

type createResultRequest struct {
	JobUuid   string    `json:"jobUuid"`
	Output    string    `json:"output"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func makeCreateResultEndpoint(svc ResultService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createResultRequest)
		err := svc.CreateResult(req)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func decorateHandler(ctx context.Context, delivery *amqp.Delivery) (interface{}, error) {
	var request createResultRequest
	err := json.Unmarshal(delivery.Body, &request)
	if err != nil {
		// TODO: handle error properly
		delivery.Reject(true)
		return nil, err
	}
	return request, nil
}
