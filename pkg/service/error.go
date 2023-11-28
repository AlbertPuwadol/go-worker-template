package service

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/wisesight/spider-go-utilities/queue"
)

func PublishError(ctx context.Context, queueAdapter queue.RabbitMQ, queueName string, msg amqp091.Delivery, errInfo error) {
	var payload map[string]interface{}
	payload = make(map[string]interface{})
	payload["body"] = string(msg.Body)
	body, err := json.Marshal(payload)
	if err != nil {
		log.Println(err.Error())
	}

	queueAdapter.PublishError(ctx, queueName, body, errInfo)
	msg.Nack(false, false)
}
