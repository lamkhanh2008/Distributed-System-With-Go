package consistent

import (
	"hash/crc32"
	"sync"
)

type Ring struct {
	Nodes Nodes
	sync.Mutex
}

type Nodes []*Node

type Node struct {
	Id     string
	HashId uint32
}

func NewRing() *Ring {
	return &Ring{Nodes: Nodes{}}
}

func hashId(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func NewNode(id string) *Node {
	return &Node{
		Id:     id,
		HashId: hashId(id),
	}
}
