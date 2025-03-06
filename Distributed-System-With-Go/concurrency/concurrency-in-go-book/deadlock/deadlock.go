package main

import (
	"fmt"
	"sync"
	"time"
)

type value struct {
	mu    sync.Mutex
	value int
}

var wg sync.WaitGroup

func main() {
	prinSum := func(v1, v2 *value) {
		// defer wg.Done()

		v1.mu.Lock()
		fmt.Println("----", v1)
		defer v1.mu.Unlock()

		time.Sleep(100 * time.Millisecond) //them time sleep vao de cho 1 goroutine giành được lock của v2

		v2.mu.Lock()
		defer v2.mu.Unlock()
		fmt.Printf("sum=%v +  v = %+v \n", v1.value+v2.value, v1)
	}

	a := value{value: 1}
	b := value{value: 2}
	// wg.Add(2)

	go prinSum(&a, &b)
	go prinSum(&b, &a)
	// wg.Wait()
	time.Sleep(1 * time.Second)
}
