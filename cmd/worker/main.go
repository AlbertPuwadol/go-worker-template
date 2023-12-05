package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AlbertPuwadol/go-worker-template/config"
	"github.com/AlbertPuwadol/go-worker-template/pkg/adapter"
	"github.com/AlbertPuwadol/go-worker-template/pkg/entity"
	"github.com/AlbertPuwadol/go-worker-template/pkg/repository"
	"github.com/AlbertPuwadol/go-worker-template/pkg/service"
	"github.com/AlbertPuwadol/go-worker-template/pkg/usecase"
	pb "github.com/AlbertPuwadol/grpc-clean/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumeQueues := []string{cfg.ConsumeQueueName}
	publishQueues := []string{cfg.PublishQueueName}
	errorQueues := []string{cfg.ConsumeErrorQueueName}

	rabbitmqAdapter, err := adapter.NewRabbitMQ(cfg.RabbitMQURI, adapter.QueueConfig{
		ConsumeQueueNames: consumeQueues,
		PublishQueueNames: publishQueues,
		ErrorQueueNames:   errorQueues,
		ExchangeName:      cfg.PublishExchangeName,
		ExchangeType:      cfg.PublishExchangeType,
	})
	if err != nil {
		panic(err)
	}
	defer rabbitmqAdapter.CleanUp()

	interceptor := service.NewAuthInterceptor(cfg.GRPCToken)
	conn, err := grpc.Dial(cfg.GRPCUri, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(interceptor.Unary()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	grpcClient := pb.NewGRPCCleanServiceClient(conn)

	healthcheck, err := grpcClient.Hello(context.Background(), &pb.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println("gRPC Healthcheck Status: ", healthcheck.Status)

	GRPCAdapter := adapter.NewGRPC(grpcClient)

	taskRepository := repository.NewTask(GRPCAdapter)
	taskUsecase := usecase.NewTask(taskRepository)

	msgs, err := rabbitmqAdapter.Consume(cfg.ConsumeQueueName, 1)
	fmt.Printf("Consuming from %s...\n", cfg.ConsumeQueueName)
	for msg := range msgs {
		var post entity.Post
		err := json.Unmarshal(msg.Body, &post)
		if err != nil {
			service.PublishError(ctx, rabbitmqAdapter, cfg.ConsumeErrorQueueName, msg, err)
			fmt.Printf("Sent %s to %s\n", post.ID, cfg.ConsumeErrorQueueName)
			continue
		}
		fmt.Printf("Job recieve: %+v\n", post)

		err = taskUsecase.GetPostTasks(&post)
		if err != nil {
			service.PublishError(ctx, rabbitmqAdapter, cfg.ConsumeErrorQueueName, msg, err)
			fmt.Printf("Sent %s to %s\n", post.ID, cfg.ConsumeErrorQueueName)
			continue
		}

		body, err := json.Marshal(post)
		if err != nil {
			service.PublishError(ctx, rabbitmqAdapter, cfg.ConsumeErrorQueueName, msg, err)
			fmt.Printf("Sent %s to %s\n", post.ID, cfg.ConsumeErrorQueueName)
			continue
		}

		err = rabbitmqAdapter.Publish(ctx, cfg.PublishQueueName, body)
		if err != nil {
			service.PublishError(ctx, rabbitmqAdapter, cfg.ConsumeErrorQueueName, msg, err)
			fmt.Printf("Sent %s to %s\n", post.ID, cfg.ConsumeErrorQueueName)
			continue
		}
		fmt.Printf("Job sent successfully, %s\n", post.ID)
		msg.Ack(false)
	}
}
