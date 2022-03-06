package http_service

import (
	"context"
	"google.golang.org/grpc"
	"simple_grpc/cmd/global"
)

type Service struct {
	ctx    context.Context
	client *grpc.ClientConn
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.client = global.GRPCClient
	return svc
}
