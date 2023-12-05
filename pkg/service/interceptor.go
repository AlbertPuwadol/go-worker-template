package service

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct {
	accessToken string
}

func NewAuthInterceptor(accessToken string) *AuthInterceptor {
	return &AuthInterceptor{accessToken: accessToken}
}

func (a AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		return invoker(metadata.AppendToOutgoingContext(ctx, "authorization", a.accessToken), method, req, reply, cc, opts...)
	}
}
