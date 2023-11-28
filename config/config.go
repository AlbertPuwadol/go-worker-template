package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURI                             string `env:"RABBITMQ_URI"`
	EnrichmentQueueName                     string `env:"ENRICHMENT_QUEUE_NAME"`
	EnrichmentErrorQueueName                string `env:"ENRICHMENT_ERROR_QUEUE_NAME"`
	ImageProcessingSplitMessageQueueName    string `env:"IMAGE_PROCESSING_SPLIT_MESSAGE_QUEUE_NAME"`
	ImageProcessingSplitMessageExchangeName string `env:"IMAGE_PROCESSING_SPLIT_MESSAGE_EXCHANGE_NAME"`
	ImageProcessingSplitMessageExchangeType string `env:"IMAGE_PROCESSING_SPLIT_MESSAGE_EXCHANGE_TYPE"`
}

func NewConfig() Config {
	godotenv.Load()
	config := Config{}
	if err := env.Parse(&config); err != nil {
		panic(err)
	}
	return config
}
