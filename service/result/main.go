package main

import (
	"common/config"
	"database/sql"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"log"
	"result/db"
)

type Cfg struct {
	PostgresConnectionString string `env:"POSTGRES_CONNECTION_STRING,required"`
	RabbitMQConnectionString string `env:"RABBITMQ_CONNECTION_STRING,required"`
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

	queue, err := AMQPDeclaration(ch)
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

	ed := makeCreateResultEndpoint(service)

	handler := amqptransport.NewSubscriber(
		ed,
		decodeCreateResultRequest,
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
