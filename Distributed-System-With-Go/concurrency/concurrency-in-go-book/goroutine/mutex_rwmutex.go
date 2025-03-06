package main

import (
	"fmt"
	"sync"
	"time"
)

func TestMutex() {
	var count int
	var lock sync.Mutex

	increment := func() {
		lock.Lock()
		defer lock.Unlock()
		count += 1
		fmt.Println("Increment: %d\n", count)
	}

	decrement := func() {
		lock.Lock()
		defer lock.Unlock()
		count -= 1
		fmt.Println("Decrement: %d\n", count)
	}

	go increment()
	go decrement()
	time.Sleep(2 * time.Second)
}

// func TestRWMutext() {
// 	producer := func(wg *sync.WaitGroup, l sync.Locker) {
// 		defer wg.Done()
// 		for i := 5; i > 0; i-- {
// 			l.Lock()
// 			l.Unlock()
// 			time.Sleep(1)
// 		}
// 	}

// 	var wg sync.WaitGroup
// 	var l s
// 	producer()
// }

func TestCond() {
	var read = false
	var cond = sync.NewCond(&sync.Mutex{})

	waitForEvent := func() {
		fmt.Println("START WAIT")
		cond.L.Lock()
		for !read {
			cond.Wait()
		}
		cond.L.Unlock()
		fmt.Println("done wait")
	}

	sendSignals := func() {
		time.Sleep(2 * time.Second)

		read = true
		cond.Signal()
		fmt.Println("set true")
	}

	go waitForEvent()
	go sendSignals()

	time.Sleep(5 * time.Second)
}

func TestOnce() {
	var (
		once  sync.Once
		wg    sync.WaitGroup
		count int
	)

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			once.Do(func() {
				count += 1
			})
		}()
	}

	fmt.Println(count)
}
