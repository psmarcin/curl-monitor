package main

import (
	"common"
	"github.com/streadway/amqp"
)

func SetupRabbitMQ(ch *amqp.Channel) (*amqp.Queue, error) {
	err := ch.ExchangeDeclare(
		common.JobExchangeName,
		amqp.ExchangeDirect,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

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

	queue, err := ch.QueueDeclare(
		common.TriggerQueueName,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   //args
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(commandRunQueue.Name, common.CommandRunRoutingKey, common.JobExchangeName, false, nil)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(queue.Name, common.TriggerRoutingKey, common.JobExchangeName, false, nil)
	if err != nil {
		return nil, err
	}

	return &queue, nil
}
