package session_client

import (
	"context"
	"hello/model"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type SessionClient interface {
	SayHello(ctx context.Context, in *model.HelloRequest, opts ...grpc.CallOption) (*model.HelloReply, error)
}

type defaultSessionClient struct {
	cli zrpc.Client
}

func NewSessionClient(cli zrpc.Client) SessionClient {
	return &defaultSessionClient{
		cli: cli}
}

func (c *defaultSessionClient) SayHello(ctx context.Context, in *model.HelloRequest, opts ...grpc.CallOption) (*model.HelloReply, error) {
	client := model.NewGreeterClient(c.cli.Conn())
	return client.SayHello(ctx, in)
}
