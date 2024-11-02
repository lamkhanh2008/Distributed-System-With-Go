package main

import (
	"context"
	"hello/model"
	"hello/session_client"
	"log"

	"github.com/zeromicro/go-zero/core/discov"
	"github.com/zeromicro/go-zero/zrpc"
)

func main() {

	sess := session_client.NewSession(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:2379"},
			Key:   "hello.rpc",
		}})
	// client := zrpc.MustNewClient(zrpc.RpcClientConf{
	// 	Etcd: discov.EtcdConf{
	// 		Hosts: []string{"localhost:2379"},
	// 		Key:   "hello.rpc",
	// 	},
	// })

	// conn := client.Conn()
	// hello := model.NewGreeterClient(conn)
	service, err := sess.GetSessionClient("125", zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: []string{"localhost:2379"},
			Key:   "hello.rpc",
		}})
	if err != nil {
		log.Fatal(err)
	}
	reply, err := service.SayHello(context.Background(), &model.HelloRequest{Name: "go-zero"})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply.Message)
}
