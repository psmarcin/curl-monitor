package main

import (
	"database/sql"
	"flag"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"log"
	"result/db"
)

func main() {
	amqpURL := flag.String(
		"url",
		"amqp://localhost:5672",
		"URL to AMQP server",
	)

	// connect to AMQP
	conn, err := amqp.Dial(*amqpURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	urlExample := "postgres://user:pass@localhost:5432/job?sslmode=disable"
	connection, err := sql.Open("postgres", urlExample)
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
