package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
)

// FSM đại diện cho state machine của Raft
type FSM struct{}

func (f *FSM) Apply(log *raft.Log) interface{} {
	return nil
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	return nil, nil
}

func (f *FSM) Restore(io.ReadCloser) error {
	return nil
}

// NewRaftNode khởi tạo node Raft với TCP transport
func NewRaftNode(nodeID, raftAddr, raftDir string, isBootstrap bool) (*raft.Raft, error) {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	// Cấu hình thời gian timeout cho Raft
	config.HeartbeatTimeout = 1 * time.Second
	config.ElectionTimeout = 2 * time.Second
	config.LeaderLeaseTimeout = 1 * time.Second
	config.CommitTimeout = 500 * time.Millisecond

	// Thiết lập bộ lưu trữ bền vững cho log và stable store
	err := os.MkdirAll(raftDir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Could not create data directory: %s", err)
	}
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-log.bolt"))
	if err != nil {
		return nil, fmt.Errorf("failed to create log store: %v", err)
	}

	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(raftDir, "raft-stable.bolt"))
	if err != nil {
		return nil, fmt.Errorf("failed to create stable store: %v", err)
	}

	snapshotStore := raft.NewDiscardSnapshotStore()

	// Thiết lập TCP transport cho Raft
	transport, err := raft.NewTCPTransport(raftAddr, nil, 2, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to create TCP transport: %v", err)
	}

	// Tạo node Raft
	fsm := &FSM{}
	raftNode, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, fmt.Errorf("failed to create raft node: %v", err)
	}

	// Bootstrap cụm nếu là node đầu tiên
	if isBootstrap {
		config := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raft.ServerID(nodeID),
					Address: raft.ServerAddress(raftAddr),
				},
			},
		}
		err = raftNode.BootstrapCluster(config).Error()
		if err != nil && err != raft.ErrCantBootstrap {
			return nil, fmt.Errorf("failed to bootstrap cluster: %v", err)
		}
		log.Println("Node bootstrapped as the initial leader")
	}

	return raftNode, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <nodeID> <raftAddr>")
		return
	}

	nodeID := os.Args[1]
	raftAddr := os.Args[2]
	raftDir := fmt.Sprintf("raft_data_%s", nodeID)

	// Kiểm tra nếu đây là node đầu tiên trong cụm
	isBootstrap := checkIfFirstNode(nodeID)

	// Khởi tạo node Raft với TCP transport
	raftNode, err := NewRaftNode(nodeID, raftAddr, raftDir, isBootstrap)
	if err != nil {
		log.Fatalf("failed to create raft node: %v", err)
	}

	// Nếu không phải node đầu tiên, cố gắng tham gia cụm Raft
	if !isBootstrap {
		print("------------------------------------------------------------------------")
		go joinCluster(raftNode, nodeID, raftAddr)
	}

	// Theo dõi trạng thái của Raft và hiển thị thông tin leader
	go monitorRaftState(raftNode)

	// Giữ chương trình chạy mãi mãi
	select {}
}

// checkIfFirstNode kiểm tra nếu node hiện tại là node đầu tiên trong cụm
func checkIfFirstNode(nodeID string) bool {
	return nodeID == "node1"
}

// joinCluster cố gắng tham gia vào cụm Raft
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

// monitorRaftState theo dõi trạng thái Raft và hiển thị thông tin leader
func monitorRaftState(raftNode *raft.Raft) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		leader := raftNode.Leader()
		if leader == "" {
			log.Println("No leader detected, waiting for election...")
		} else {
			log.Printf("Current leader is: %s", leader)
		}
	}
}
