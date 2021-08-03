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
	commandRunQueue, err := ch.QueueDeclare(
		"CM.Job.CommandRun",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   //args
	)
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		"CM.Job.Trigger",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   //args
	)

	_ = ch.QueueBind(commandRunQueue.Name, "commandrun", "CM.Job", false, nil)
	_ = ch.QueueBind(queue.Name, "trigger", "CM.Job", false, nil)

	return &queue, nil
}