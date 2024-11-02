package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"raft_grpc/proto/raft_grpc/proto"
	"time"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"google.golang.org/grpc"
)

// RaftServiceServer triển khai gRPC server cho Raft
type RaftServiceServer struct {
	proto.UnimplementedRaftServiceServer
	raft *raft.Raft
}

// GenerateID xử lý yêu cầu tạo ID mới
func (s *RaftServiceServer) GenerateID(ctx context.Context, req *proto.GenerateIDRequest) (*proto.GenerateIDResponse, error) {
	// Kiểm tra xem node hiện tại có phải là leader không
	if s.raft.State() != raft.Leader {
		leader := s.raft.Leader()
		if leader == "" {
			return nil, fmt.Errorf("no leader found")
		}

		// Chuyển tiếp yêu cầu đến leader
		conn, err := grpc.Dial(string(leader), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := proto.NewRaftServiceClient(conn)
		return client.GenerateID(ctx, req)
	}

	// Nếu là leader, xử lý yêu cầu
	newID := req.LastID + 1
	fmt.Printf("Leader handling request: new ID = %d\n", newID)
	return &proto.GenerateIDResponse{NewID: newID}, nil
}

// RaftNode đại diện cho một node Raft
type RaftNode struct {
	raft *raft.Raft
	addr string
}

// NewRaftNode khởi tạo node Raft
func NewRaftNode(nodeID, raftAddr, raftDir string) (*RaftNode, *transport.Manager, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)
	config.HeartbeatTimeout = 2 * time.Second
	config.ElectionTimeout = 2 * time.Second
	config.LeaderLeaseTimeout = 1 * time.Second
	config.CommitTimeout = 500 * time.Millisecond
	// Thiết lập bộ lưu trữ bền vững
	err := os.MkdirAll(raftDir, 0755)

	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-log.bolt"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create log store: %v", err)
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-stable.bolt"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create stable store: %v", err)
	}

	snapshotStore := raft.NewDiscardSnapshotStore()

	// Tạo địa chỉ TCP cho Raft
	// transport, err := raft.NewTCPTransport(raftAddr, nil, 2, time.Second*10, os.Stderr)
	grpcTransport := transport.New(raft.ServerAddress(raftAddr), []grpc.DialOption{grpc.WithInsecure()})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transport: %v", err)
	}

	fsm := &FSM{}

	raftNode, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, grpcTransport.Transport())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create raft node: %v", err)
	}

	// Bootstrap cụm nếu đây là node đầu tiên
	if nodeID == "node1" {
		config := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raft.ServerID(nodeID),
					Address: raft.ServerAddress(raftAddr),
				},
			},
		}
		raftNode.BootstrapCluster(config)
	}

	return &RaftNode{
		raft: raftNode,
		addr: raftAddr,
	}, grpcTransport, nil
}

// FSM đại diện cho state machine của Raft
type FSM struct{}

func (f *FSM) Apply(log *raft.Log) interface{}     { return nil }
func (f *FSM) Snapshot() (raft.FSMSnapshot, error) { return nil, nil }
func (f *FSM) Restore(io.ReadCloser) error         { return nil }

// runGRPCServer chạy gRPC server cho mỗi node
func runGRPCServer(node *RaftNode, port string, tm *transport.Manager) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterRaftServiceServer(grpcServer, &RaftServiceServer{raft: node.raft})
	tm.Register(grpcServer)
	log.Printf("gRPC server running on %s", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func joinCluster(raftNode *raft.Raft, nodeID, raftAddr string) {
	for {
		leader := raftNode.Leader()
		if leader == "" {
			log.Println("No leader found, waiting to join the cluster...")
			time.Sleep(2 * time.Second)
			continue
		}

		// Gửi yêu cầu join cluster đến leader hiện tại
		err := raftNode.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(raftAddr), 0, 0).Error()
		if err != nil {
			log.Printf("Failed to join cluster: %v", err)
			time.Sleep(2 * time.Second)
		} else {
			log.Println("Successfully joined the cluster")
			break
		}
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run main.go <nodeID> <raftAddr> <grpcPort>")
		return
	}

	nodeID := os.Args[1]
	raftAddr := os.Args[2]
	grpcPort := os.Args[3]
	raftDir := fmt.Sprintf("./data/raft_data_%s", nodeID)

	// Khởi tạo node Raft
	node, tm, err := NewRaftNode(nodeID, raftAddr, raftDir)
	if err != nil {
		log.Fatalf("failed to create raft node: %v", err)
	}

	if nodeID != "node1" {
		joinCluster(node.raft, nodeID, raftAddr)

	}

	// Chạy gRPC server cho node
	go runGRPCServer(node, grpcPort, tm)

	// Giữ chương trình chạy mãi mãi
	select {}
}
