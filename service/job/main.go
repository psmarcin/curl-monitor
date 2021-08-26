//go:generate sqlc generate --file ./../../sqlc.yaml
package main

import (
	"database/sql"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	"job/db"
	"log"
	"net/http"
	"time"

	"common/config"
	_ "github.com/lib/pq"
)

type Job struct {
	Uuid      string
	Name      string
	Command   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Cfg struct {
	PostgresConnectionString string `env:"POSTGRES_CONNECTION_STRING,required"`
	PORT                     string `env:"PORT,required" envDefault:"8080"`
	HealthCheckPORT          string `env:"HEALTHCHECK_PORT,required" envDefault:"8081"`
}

func main() {
	var cfg Cfg
	err := config.Load(&cfg)

	connection, err := sql.Open("postgres", cfg.PostgresConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	err = connection.Ping()
	if err != nil {
		log.Fatal(err)
	}

	queries := db.New(connection)
	svc := jobService{
		DB: queries,
	}

	getJobHandler := httptransport.NewServer(
		makeGetJobEndpoint(svc),
		decodeGetJobRequest,
		encodeResponse,
	)

	listJobHandler := httptransport.NewServer(
		makeListJobEndpoint(svc),
		decodeListJobRequest,
		encodeResponse,
	)

	createJobHandler := httptransport.NewServer(
		makeCreateJobEndpoint(svc),
		decodeCreateJobRequest,
		encodeResponse,
	)
	updateJobHandler := httptransport.NewServer(
		makeUpdateJobEndpoint(svc),
		decodeUpdateJobRequest,
		encodeResponse,
	)

	// health check
	health := healthcheck.NewHandler()
	health.AddReadinessCheck("postgres", healthcheck.DatabasePingCheck(connection, time.Second))
	go func() {
		log.Printf("Healthcheck listening on port :%s", cfg.HealthCheckPORT)
		err := http.ListenAndServe("0.0.0.0:"+cfg.HealthCheckPORT, health)
		if err != nil {
			log.Fatalf("error on healtcheck %s", err)
		}
	}()

	r := mux.NewRouter()
	r.Methods(http.MethodGet).Path("/{id}").Handler(getJobHandler)
	r.Methods(http.MethodGet).Path("/").Handler(listJobHandler)
	r.Methods(http.MethodPut).Path("/{id}").Handler(updateJobHandler)
	r.Methods(http.MethodPost).Path("/").Handler(createJobHandler)
	http.Handle("/", r)
	log.Printf("App listening on port :%s", cfg.PORT)
	log.Fatal(http.ListenAndServe(":"+cfg.PORT, nil))
}
