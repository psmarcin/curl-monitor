package main

import (
	"bytes"
	"common/config"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	"io/ioutil"
	"job/db"
	"log"
	"net/http"
)

type Cfg struct {
	RabbitMQConnectionString string `env:"RABBITMQ_CONNECTION_STRING,required"`
}

func main() {
	var cfg Cfg
	err := config.Load(&cfg)

	amqpURL := flag.String(
		"url",
		cfg.RabbitMQConnectionString,
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

	service := triggerService{
		Channel: ch,
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

// encodeHTTPGenericRequest is a transport/http.EncodeRequestFunc that
// JSON-encodes any request to the request body. Primarily useful in a client.
func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// encodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		//errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHTTPConcatResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded concat response from the HTTP response body. If the response
// has a non-200 status code, we will interpret that as an error and attempt to
// decode the specific error message from the response body. Primarily useful in
// a client.
func decodeHTTPConcatResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp db.Job
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}
