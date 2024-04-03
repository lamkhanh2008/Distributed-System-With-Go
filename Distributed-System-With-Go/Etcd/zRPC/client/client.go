package main

import (
	"context"
	"hello/model"
	"log"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

func main() {
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:2379"},
			Key:   "hello.rpc",
		},
	})

	conn := client.Conn()
	hello := model.NewGreeterClient(conn)
	reply, err := hello.SayHello(context.Background(), &model.HelloRequest{Name: "go-zero"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.Message)
}
