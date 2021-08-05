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
	outputQueue, err := ch.QueueDeclare(
		common.ResultQueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   //args
	)
	if err != nil {
		return nil, err
	}

	_ = ch.QueueBind(outputQueue.Name, common.ResultRoutingKey, common.JobExchangeName, false, nil)

	return &outputQueue, nil
}
