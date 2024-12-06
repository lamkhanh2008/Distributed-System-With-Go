package main

import (
	"sync"
)

func testConcurrent(m sync.Map) {
	var wg sync.WaitGroup
	// Concurrent write operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			m.Store(n, n*10)
		}(i)
	}

	wg.Wait()
	// Concurrent read operations
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			if val, ok := m.Load(n); ok {
				println(val.(int))
			}
		}(i)
	}

	wg.Wait()
}

func main() {
	var m sync.Map

	m.Store("key1", "value1")
	m.Store("key2", "value2")

	// Iterating over sync.Map
	m.Range(func(key, value interface{}) bool {
		println(key.(string), value.(string))
		return true
	})

	testConcurrent(m)
}
