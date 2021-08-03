package main

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/streadway/amqp"
	"net/http"
	"time"
)

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type createResultRequest struct {
	JobUuid   string    `json:"jobUuid"`
	Output    string    `json:"output"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type createResultResponse struct{}

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

func decodeCreateResultRequest(ctx context.Context, delivery *amqp.Delivery) (interface{}, error) {
	var request createResultRequest
	err := json.Unmarshal(delivery.Body, &request)
	if err != nil {
		delivery.Reject(true)
		return nil, err
	}
	return request, nil
}
