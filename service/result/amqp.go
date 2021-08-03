package main

import (
	"github.com/streadway/amqp"
)

func AMQPDeclaration(ch *amqp.Channel) (*amqp.Queue, error) {
	_ = ch.ExchangeDeclare(
		"CM.Job",
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)
	outputQueue, err := ch.QueueDeclare(
		"CM.Job.Output",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   //args
	)
	if err != nil {
		return nil, err
	}

	_ = ch.QueueBind(outputQueue.Name, "output", "CM.Job", false, nil)

	return &outputQueue, nil
}
