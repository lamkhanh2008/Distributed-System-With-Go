package main

import (
	"context"
	"fmt"
	"log"
	"raft_grpc/proto/raft_grpc/proto"
	"time"

	"google.golang.org/grpc"
)

// RaftClient đại diện cho client gRPC
type RaftClient struct {
	client proto.RaftServiceClient
}

// NewRaftClient khởi tạo RaftClient
func NewRaftClient(addr string) (*RaftClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := proto.NewRaftServiceClient(conn)
	return &RaftClient{client: client}, nil
}

func (c *RaftClient) GenerateID(lastID int64) {
	req := &proto.GenerateIDRequest{LastID: lastID}
	resp, err := c.client.GenerateID(context.Background(), req)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("New ID: %d\n", resp.NewID)
}

func main() {
	// Địa chỉ của node Raft ban đầu
	nodeAddr := "localhost:8002"

	client, err := NewRaftClient(nodeAddr)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Gửi yêu cầu tạo ID
	for i := 0; i < 5; i++ {
		client.GenerateID(int64(i))
		time.Sleep(1 * time.Second)
	}
}
