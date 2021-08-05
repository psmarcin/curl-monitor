package main

import (
	"common"
	"github.com/streadway/amqp"
)

func SetupRabbitMQ(ch *amqp.Channel) (*amqp.Queue, error) {
	_ = ch.ExchangeDeclare(
		common.JobExchangeName,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)
	commandRunQueue, err := ch.QueueDeclare(
		common.CommandRunQueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   //args
	)
	if err != nil {
		return nil, err
	}

	_ = ch.QueueBind(commandRunQueue.Name, common.CommandRunRoutingKey, common.JobExchangeName, false, nil)

	return &commandRunQueue, nil
}
