package adapter

import (
	"context"

	pb "github.com/AlbertPuwadol/grpc-clean/proto"
)

type GRPC interface {
	GetTask1(text string) (*pb.Task1Response, error)
	GetTask2(text string) (*pb.Task2Response, error)
	GetTask3(text string) (*pb.Task3Response, error)
}

type grpc struct {
	grpcClient pb.GRPCCleanServiceClient
	mdctx      context.Context
}

func NewGRPC(grpcClient pb.GRPCCleanServiceClient, mdctx context.Context) *grpc {
	return &grpc{grpcClient: grpcClient, mdctx: mdctx}
}

func (g grpc) GetTask1(text string) (*pb.Task1Response, error) {
	request := &pb.TaskRequest{Text: text}
	res, err := g.grpcClient.Task1(g.mdctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (g grpc) GetTask2(text string) (*pb.Task2Response, error) {
	request := &pb.TaskRequest{Text: text}
	res, err := g.grpcClient.Task2(g.mdctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (g grpc) GetTask3(text string) (*pb.Task3Response, error) {
	request := &pb.TaskRequest{Text: text}
	res, err := g.grpcClient.Task3(g.mdctx, request)
	if err != nil {
		return nil, err
	}
	return res, nil
}
