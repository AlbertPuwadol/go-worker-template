package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AlbertPuwadol/go-worker-template/config"
	"github.com/AlbertPuwadol/go-worker-template/pkg/adapter"
	"github.com/AlbertPuwadol/go-worker-template/pkg/repository"
	"github.com/AlbertPuwadol/go-worker-template/pkg/service"
	"github.com/AlbertPuwadol/go-worker-template/pkg/usecase"
	formatter "github.com/wisesight/spider-go-formatter"
	"github.com/wisesight/spider-go-utilities/queue"
)

func main() {
	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queues := []string{cfg.ImageProcessingSplitMessageQueueName}
	errorQueues := []string{cfg.EnrichmentErrorQueueName}

	rabbitmqAdapter, err := queue.NewRabbitMQ(cfg.RabbitMQURI, queue.QueueConfig{
		QueueNames:      queues,
		ErrorQueueNames: errorQueues,
		ExchangeName:    cfg.ImageProcessingSplitMessageExchangeName,
		ExchangeType:    cfg.ImageProcessingSplitMessageExchangeType,
	})
	if err != nil {
		panic(err)
	}
	defer rabbitmqAdapter.CleanUp()

	kiringRPCAdapter := adapter.NewGRPC()

	enrichmentRepository := repository.NewEnrichment(kiringRPCAdapter)
	enrichmentUsecase := usecase.NewEnrichment(enrichmentRepository)

	msgs, err := rabbitmqAdapter.Consume(cfg.EnrichmentQueueName, 1)
	fmt.Printf("Consuming from %s...\n", cfg.EnrichmentQueueName)
	for msg := range msgs {
		var post formatter.Spider
		err := json.Unmarshal(msg.Body, &post)
		if err != nil {
			service.PublishError(ctx, rabbitmqAdapter, cfg.EnrichmentErrorQueueName, msg, err)
			fmt.Printf("Sent %s to %s\n", post.ID, cfg.EnrichmentErrorQueueName)
			continue
		}
		fmt.Printf("Job recieve: %+v\n", post)

		err = enrichmentUsecase.GetPostEnrichment(&post)
		if err != nil {
			service.PublishError(ctx, rabbitmqAdapter, cfg.EnrichmentErrorQueueName, msg, err)
			fmt.Printf("Sent %s to %s\n", post.ID, cfg.EnrichmentErrorQueueName)
			continue
		}

		body, err := json.Marshal(post)
		if err != nil {
			service.PublishError(ctx, rabbitmqAdapter, cfg.EnrichmentErrorQueueName, msg, err)
			fmt.Printf("Sent %s to %s\n", post.ID, cfg.EnrichmentErrorQueueName)
			continue
		}

		err = rabbitmqAdapter.Publish(ctx, cfg.ImageProcessingSplitMessageQueueName, body)
		if err != nil {
			service.PublishError(ctx, rabbitmqAdapter, cfg.EnrichmentErrorQueueName, msg, err)
			fmt.Printf("Sent %s to %s\n", post.ID, cfg.EnrichmentErrorQueueName)
			continue
		}
		fmt.Printf("Job sent successfully, %s\n", post.ID)
		msg.Ack(false)
	}
}
