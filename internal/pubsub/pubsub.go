package pubsub

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

var TransientQueue = SimpleQueueType("transient")
var DurableQueue = SimpleQueueType("durable")

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	rawData, err := json.Marshal(val)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        rawData,
	})

	return err
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}
	queue, err := channel.QueueDeclare(queueName, queueType == DurableQueue, queueType == TransientQueue, queueType == TransientQueue, false, nil)
	if err != nil {
		return nil, queue, err
	}

	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return channel, queue, err
	}

	return channel, queue, nil
}
