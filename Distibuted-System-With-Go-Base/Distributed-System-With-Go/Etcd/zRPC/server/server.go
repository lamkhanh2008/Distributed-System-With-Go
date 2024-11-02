package main

import (
	"context"
	"hello/model"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type Config struct {
	zrpc.RpcServerConf
}

// var cfgFile = flag.String("f", "./hello.yaml", "cfg file")
var args = os.Args

func main() {
	// args := os.Args
	var c zrpc.RpcServerConf
	conf.MustLoad(string(args[1]), &c)

	server := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		model.RegisterGreeterServer(grpcServer, &Hello{})
	})
	server.Start()

	// flag.Parse()
	// var cfg Config
	// conf.MustLoad(*cfgFile, &cfg)
	// srv, err := zrpc.MustNewServer(cfg.RpcServerConf, func(s *grpc.Server) {
	// 	model.RegisterGreeterServer(s, &Hello{})
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// srv.Start()
}

type Hello struct {
	model.UnimplementedGreeterServer
}

func (h *Hello) SayHello(ctx context.Context, in *model.HelloRequest) (*model.HelloReply, error) {
	return &model.HelloReply{Message: "hello " + in.Name + " with config:" + string(args[1])}, nil
}

// func (h *Hello) mustEmbedUnimplementedGreeterServer() {}
