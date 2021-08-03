package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type ListJobRequest struct{}
type ListJobResponse []ListJobResponseJob
type ListJobResponseJob struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func makeListJobEndpoint(svc JobService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		jobs, _ := svc.ListJob()
		var response ListJobResponse
		for _, job := range jobs {
			response = append(response, ListJobResponseJob{
				Uuid:      job.Uuid,
				Name:      job.Name,
				Command:   job.Command,
				CreatedAt: job.CreatedAt,
				UpdatedAt: job.UpdatedAt,
			})
		}
		return response, nil
	}
}

func decodeListJobRequest(_ context.Context, r *http.Request) (interface{}, error) {
	request := ListJobRequest{}
	return request, nil
}

type getJobRequest struct {
	Uuid string `json:"id"`
}
type getJobResponse struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func makeGetJobEndpoint(svc JobService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(getJobRequest)
		response, _ := svc.GetJob(req.Uuid)
		return getJobResponse{
			Uuid:      response.Uuid,
			Name:      response.Name,
			Command:   response.Command,
			CreatedAt: response.CreatedAt,
			UpdatedAt: response.UpdatedAt,
		}, nil
	}
}

func decodeGetJobRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		return nil, errors.New("id is not provided")
	}

	req := getJobRequest{Uuid: id}
	return req, nil
}

type createJobRequest struct {
	Name    string `json:"name"`
	Command string `json:"command"`
}
type createJobResponse struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func makeCreateJobEndpoint(svc JobService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(createJobRequest)
		response, _ := svc.CreateJob(req)
		return createJobResponse{
			Uuid:      response.Uuid,
			Name:      response.Name,
			Command:   response.Command,
			CreatedAt: response.CreatedAt,
			UpdatedAt: response.UpdatedAt,
		}, nil
	}
}

func decodeCreateJobRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}
	return request, nil
}

type updateJobRequest struct {
	Uuid    string `json:"uuid"`
	Name    string `json:"name"`
	Command string `json:"command"`
}
type updateJobResponse struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func makeUpdateJobEndpoint(svc JobService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(updateJobRequest)
		response, _ := svc.UpdateJob(req)
		return updateJobResponse{
			Uuid:      response.Uuid,
			Name:      response.Name,
			Command:   response.Command,
			CreatedAt: response.CreatedAt,
			UpdatedAt: response.UpdatedAt,
		}, nil
	}
}

func decodeUpdateJobRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request updateJobRequest
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		return nil, errors.New("id is not provided")
	}
	request.Uuid = id
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return request, err
	}
	return request, nil
}
