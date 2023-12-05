package adapter

import (
	"context"
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ interface {
	Publish(ctx context.Context, routing string, payload []byte) error
	PublishError(ctx context.Context, routing string, rawPayload []byte, errInfo error) error
	Consume(queue string, prefetch int) (<-chan amqp.Delivery, error)
}

type QueueConfig struct {
	ConsumeQueueNames []string
	PublishQueueNames []string
	ErrorQueueNames   []string
	ExchangeName      string
	ExchangeType      string
}

type rabbitmq struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	config QueueConfig
}

func NewRabbitMQ(url string, config QueueConfig) (*rabbitmq, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q := &rabbitmq{conn, ch, config}

	err = q.createAndBindQueue(config)
	if err != nil {
		return nil, err
	}

	for _, v := range config.ErrorQueueNames {
		err := q.declareQueue(v)
		if err != nil {
			return nil, err
		}
	}

	return q, nil
}

func (q rabbitmq) createAndBindQueue(queueConfig QueueConfig) error {
	err := q.declareExchange(queueConfig.ExchangeName, queueConfig.ExchangeType)
	if err != nil {
		return err
	}
	for _, queue := range queueConfig.ConsumeQueueNames {
		err = q.declareQueue(queue)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	for _, queue := range queueConfig.PublishQueueNames {
		err = q.declareQueue(queue)
		if err != nil {
			return err
		}
		routingKey := ""
		if queueConfig.ExchangeType == "direct" {
			routingKey = queue
		}
		err = q.bindQueue(queue, routingKey, queueConfig.ExchangeName)
	}
	if err != nil {
		return err
	}
	return nil
}

func (q rabbitmq) declareExchange(exchangeName string, exchangeType string) error {
	err := q.ch.ExchangeDeclare(
		exchangeName, // name
		exchangeType, // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	return err
}

func (q rabbitmq) declareQueue(queueName string) error {
	_, err := q.ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return err
}

func (q rabbitmq) bindQueue(queueName string, routingKey string, exchangeName string) error {
	err := q.ch.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	return err
}

func (r *rabbitmq) CleanUp() {
	r.ch.Close()
	r.conn.Close()
}

func (r rabbitmq) Consume(queue string, prefetch int) (<-chan amqp.Delivery, error) {
	err := r.ch.Qos(prefetch, 0, false)
	if err != nil {
		return nil, err
	}

	msgs, err := r.ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (r rabbitmq) Publish(ctx context.Context, routing string, payload []byte) error {
	err := r.ch.PublishWithContext(ctx,
		r.config.ExchangeName,
		routing,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        payload,
		})

	return err
}

func (r rabbitmq) PublishError(ctx context.Context, routing string, rawPayload []byte, errInfo error) error {
	var payload map[string]interface{}
	err := json.Unmarshal(rawPayload, &payload)
	if err != nil {
		payload = make(map[string]interface{})
		payload["body"] = string(rawPayload)
	}

	payload["error"] = errInfo.Error()

	result, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = r.ch.PublishWithContext(ctx,
		"",
		routing,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        result,
		})

	return err
}
