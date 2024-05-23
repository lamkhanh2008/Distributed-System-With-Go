package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

type Worker struct{}

type Result struct{}

type Request struct{}

func Process(w *Worker) Result {
	return Result{}
}

// func Remind_WorkerPool() {
// 	workCh := make(chan Worker)
// 	resultCh := make(chan Result)
// 	done := make(chan bool)

// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			for {
// 				work := <-workCh
// 				result := Process(&work)
// 				resultCh <- result
// 			}
// 		}()
// 	}

// 	results := make([]Result, 0)
// }

// func Test_Ideal() {
// 	stopCh := make(chan struct{})
// 	requestCh := make(chan Request)
// 	resultCh := make(chan Result)

//		go func(){
//			for {
//				var req Request
//				select {
//				case req = <-requestCh:
//				case <-stopCh:
//					return
//				}
//			}
//		}()
//	}
type Data struct {
}
type Cache struct {
	mu sync.Mutex
	m  map[string]*Data
}

func (c *Cache) Get(id string) (Data, bool) {
	c.mu.Lock()
	data, exists := c.m[id]
	c.mu.Unlock()
	if exists {
		if data == nil {
			return Data{}, false
		}
	}
	return *data, true
}

type Queue struct {
	elements    []int
	front, rear int
	len         int
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		elements: make([]int, capacity),
		front:    0,  // Read from elements[front]
		rear:     -1, // Write to elements[rear]
		len:      0,
	}
}

func (q *Queue) Enqueue(value int) bool {
	if q.len == len(q.elements) {
		return false
	}
	q.rear = (q.rear + 1) % len(q.elements)
	q.elements[q.rear] = value
	q.len++
	return true
}
func (q *Queue) Dequeue() (int, bool) {
	if q.len == 0 {
		return 0, false
	}
	// Read the value at the read pointer
	data := q.elements[q.front]
	// Advance the read pointer, go around in a circle
	q.front = (q.front + 1) % len(q.elements)
	q.len--
	return data, true
}

func Producer(lock *sync.Mutex, fullCond, emptyCond *sync.Cond, queue Queue) {
	for {
		value := rand.Int()
		lock.Lock()
		for !queue.Enqueue(value) {
			fmt.Println("Queue s full")
			fullCond.Wait()
		}
		lock.Unlock()
		emptyCond.Signal()
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
	}
}

func Consumer(lock *sync.Mutex, fullCond, emptyCond *sync.Cond, queue Queue) {
	for {
		lock.Lock()
		var v int
		for {
			var ok bool
			if v, ok = queue.Dequeue(); !ok {
				fmt.Println("Queue is empty")
				emptyCond.Wait()
				continue
			}
			break
		}

		lock.Unlock()
		fullCond.Signal()
		time.Sleep(time.Millisecond *
			time.Duration(rand.Intn(1000)))
		fmt.Println(v)
	}
}

func main() {

	TestAtomic()
	// lock := sync.Mutex{}
	// fullCond := sync.NewCond(&lock)
	// emptyCond := sync.NewCond(&lock)
	// queue := NewQueue(10)
	// for i := 0; i < 10; i++ {
	// 	go Producer(&lock, fullCond, emptyCond, *queue)
	// }
	// for i := 0; i < 10; i++ {
	// 	go Consumer(&lock, fullCond, emptyCond, *queue)
	// }
	// var wg sync.WaitGroup
	// resource := make(chan int, 10)
	// for i := 0; i < 10; i++ {
	// 	resource <- i
	// }
	// close(resource)

	// wg.Add(2)
	// go func() {
	// 	defer wg.Done()
	// 	for i := range resource {
	// 		fmt.Println(i)
	// 	}
	// }()
	// // go func() {
	// // 	for i := range resource {
	// // 		fmt.Println(i)
	// // 	}
	// // }()

	// go func() {
	// 	defer wg.Done()
	// 	for i := range resource {
	// 		fmt.Println("2  ", i)
	// 	}
	// }()
	// wg.Wait()
}

//test atomic package

type CacheSyncMap struct {
	values sync.Map
}

type CachedValue struct {
	sync.Once
	value *Data
}

func loadData(id string) *Data {
	return &Data{}
}
func (c *CacheSyncMap) Get(id string) *Data {

	v, _ := c.values.LoadOrStore(id, &CachedValue{})
	cv := v.(*CachedValue)
	cv.Do(func() {
		cv.value = loadData(id)
	})
	return cv.value
}

func TestAtomic() {
	var i int
	var v atomic.Value
	go func() {
		i = 2
		v.Store(2)
	}()

	go func() {
		for {
			if val, _ := v.Load().(int); val == 1 {
				fmt.Println(i)
				return
			}
		}
	}()
	select {}
}
