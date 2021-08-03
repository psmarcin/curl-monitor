package job

import (
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"job/db"
	"net/http"
	"net/url"
	"time"
)

type Job struct {
	URL *url.URL
}

type JobService interface {
	ListJob() ([]db.Job, error)
}

type JobFromAPI struct {
	Uuid      string    `json:"uuid"`
	Name      string    `json:"name"`
	Command   string    `json:"command"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (j Job) ListJob() ([]db.Job, error) {
	var jobsFromAPI []JobFromAPI
	response, err := http.Get(j.URL.String())
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &jobsFromAPI)
	if err != nil {
		return nil, err
	}

	var jobs []db.Job
	for _, job := range jobsFromAPI {
		uid, err := uuid.Parse(job.Uuid)
		if err != nil {
			return jobs, err
		}
		jobs = append(jobs, db.Job{
			Uuid:      uid,
			Name:      job.Name,
			Command:   job.Command,
			CreatedAt: job.CreatedAt,
			UpdatedAt: job.UpdatedAt,
		})
	}

	return jobs, nil
}
