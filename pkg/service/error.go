package service

import (
	"context"

	"github.com/AlbertPuwadol/go-worker-template/pkg/adapter"
	"github.com/rabbitmq/amqp091-go"
)

func PublishError(ctx context.Context, queueAdapter adapter.RabbitMQ, queueName string, msg amqp091.Delivery, errInfo error) {
	queueAdapter.PublishError(ctx, queueName, msg.Body, errInfo)
	msg.Nack(false, false)
}
