package main

import (
	"common/config"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	"log"
	"net/url"
	"trigger/job"
)

type Cfg struct {
	RabbitMQConnectionString string `env:"RABBITMQ_CONNECTION_STRING,required"`
	JobConnectionString      string `env:"JOB_CONNECTION_STRING"`
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

	baseUrl, _ := url.Parse(cfg.JobConnectionString)
	jobClient := job.Job{URL: baseUrl}

	service := triggerService{
		Channel: ch,
		Job:     jobClient,
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

	ed := makeHandlerEndpoint(service)

	handler := amqptransport.NewSubscriber(
		ed,
		decodetriggerAMQPHandler,
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
