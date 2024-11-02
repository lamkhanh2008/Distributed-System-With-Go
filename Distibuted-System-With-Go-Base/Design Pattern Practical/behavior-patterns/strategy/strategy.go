package main

import "fmt"

type EvictionAlgo interface {
	evict(c *Cache)
}

type Cache struct {
	storage      map[string]string
	evictionAlgo EvictionAlgo
	capacity     int
	maxCapacity  int
}

func initCache(e EvictionAlgo) *Cache {
	storage := make(map[string]string)
	return &Cache{
		storage:      storage,
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  2,
	}
}

func (c *Cache) setEvictionAlgo(e EvictionAlgo) {
	c.evictionAlgo = e
}

func (c *Cache) evict() {
	c.evictionAlgo.evict(c)
	c.capacity--
}

// cac loai chien luoc pkhac nhau ap dung strategy
type lru struct{}

func (l *lru) evict(c *Cache) {
	fmt.Println("LRU")
}

type lfu struct{}

func (l *lfu) evict(c *Cache) {
	fmt.Println("LFU")
}

func main() {
	lru := lru{}
	cache := initCache(&lru)
	cache.evict()
	lfu := lfu{}
	cache.setEvictionAlgo(&lfu)
	cache.evict()

}
