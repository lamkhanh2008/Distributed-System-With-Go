package structure

import (
	"hash_sharding/src/handler"
	"sync"
	"time"
)

type KeyValue struct {
	Value       string
	Expriration time.Time
}

type Shard struct {
	data map[string]KeyValue
	mu   sync.RWMutex
}

type KeyValueStore struct {
	shards   []*Shard
	replicas int
}

func NewKeyValueStore(numShards, numReplicas int) *KeyValueStore {
	store := &KeyValueStore{
		shards:   make([]*Shard, numShards),
		replicas: numReplicas,
	}
	for i := 0; i < numShards; i++ {
		store.shards[i] = &Shard{data: make(map[string]KeyValue)}
	}
	return store
}
func (kv *KeyValueStore) GetShardIndex(key string) int {
	hash := handler.FnvHash(key)
	return int(hash) % len(kv.shards)
}

func (kv *KeyValueStore) Set(key, value string, ttl time.Duration) {
	shardIndex := kv.GetShardIndex(key)
	shard := kv.shards[shardIndex]
	shard.mu.Lock()
	defer shard.mu.Unlock()
	expiration := time.Now().Add(ttl)
	shard.data[key] = KeyValue{Value: value, Expriration: expiration}
}

func (kv *KeyValueStore) Get(key string) (string, bool) {
	shardIndex := kv.GetShardIndex(key)
	shard := kv.shards[shardIndex]
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	item, ok := shard.data[key]

	if !ok || time.Now().After(item.Expriration) {
		return "", false
	}
	return item.Value, true
}
