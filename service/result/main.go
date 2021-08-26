package main

import (
	"common/config"
	"database/sql"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/heptiolabs/healthcheck"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"result/db"
	"time"
)

type Cfg struct {
	PostgresConnectionString string `env:"POSTGRES_CONNECTION_STRING,required"`
	RabbitMQConnectionString string `env:"RABBITMQ_CONNECTION_STRING,required"`
	HealthCheckPORT          string `env:"HEALTHCHECK_PORT,required" envDefault:"8081"`
}

func main() {
	var cfg Cfg
	err := config.Load(&cfg)

	// connect to AMQP
	conn, err := amqp.Dial(cfg.RabbitMQConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	connection, err := sql.Open("postgres", cfg.PostgresConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	queries := db.New(connection)

	service := resultService{
		DB: queries,
	}

	queue, err := SetupRabbitMQ(ch)
	if err != nil {
		log.Fatalf("err: %+v", err)
	}

	consumer, err := ch.Consume(
		queue.Name,
		"",    // consumer
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
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

	ed := makeCreateResultEndpoint(service)

	handler := amqptransport.NewSubscriber(
		ed,
		decorateHandler,
		amqptransport.EncodeJSONResponse,
	)

	listener := handler.ServeDelivery(ch)

	forever := make(chan bool)

	go func() {
		for true {
			select {
			case message := <-consumer:
				log.Printf("received trigger event\n")
				listener(&message)
				message.Ack(false)
			}
		}
	}()

	log.Printf("listening")
	<-forever
}
