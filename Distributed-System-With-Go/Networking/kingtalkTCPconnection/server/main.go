package main

import (
	"fmt"
	"kingtalk_tcp/tcp-server"
	"net"
	"os"
)

func main() {
	args := os.Args
	input := args[1]
	fmt.Println(input)
	lsn, err := net.Listen("tcp", input)
	if err != nil {
		fmt.Println("listen error: %v", err)
		panic(err)
	}
	server := tcp.NewTCPServer(tcp.TCPServerArgs{
		Listener:     lsn,
		ServerName:   "lamdeptrai",
		ProtoName:    "mtproto",
		SendChanSize: 1024,
		// ConnectionCallback: server,
	})
	server.Serve()
}
