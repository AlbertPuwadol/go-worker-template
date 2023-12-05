package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURI           string `env:"RABBITMQ_URI"`
	ConsumeQueueName      string `env:"CONSUME_QUEUE_NAME"`
	ConsumeErrorQueueName string `env:"CONSUME_ERROR_QUEUE_NAME"`
	PublishQueueName      string `env:"PUBLISH_QUEUE_NAME"`
	PublishExchangeName   string `env:"PUBLISH_EXCHANGE_NAME"`
	PublishExchangeType   string `env:"PUBLISH_EXCHANGE_TYPE"`
	GRPCUri               string `env:"GRPC_URI"`
	GRPCToken             string `env:"GRPC_TOKEN"`
}

func NewConfig() Config {
	godotenv.Load()
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}
	return config
}
