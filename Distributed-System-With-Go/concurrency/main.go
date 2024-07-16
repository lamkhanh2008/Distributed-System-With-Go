package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	"net"
	"sync"
	"testing"
	"time"
)

// type Button struct {
// 	Clicked *sync.Cond
// }
// func test_broadcast_cond(){

// }
func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}
func init() {
	daemonStarted := startNetworkDaemon()
	daemonStarted.Wait()
}

func startNetworkDaemon() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer server.Close()
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			connectToService()
			fmt.Fprintln(conn, "")
			conn.Close()
		}
	}()
	return &wg
}
func BenchmarkNetworkRequest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			b.Fatalf("cannot dial host: %v", err)
		}
		if _, err := ioutil.ReadAll(conn); err != nil {
			b.Fatalf("cannot read: %v", err)
		}
		conn.Close()
	}
}
func TestOnceSync() {
	var count int
	increment := func() {
		count++
	}

	var once sync.Once
	var increments sync.WaitGroup
	increments.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer increments.Done()
			once.Do(increment)
		}()
	}
	increments.Wait()
	fmt.Printf("Count is %d \n", count)
}
func UseCondSignalAndWait() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]interface{}, 0, 10)
	removeFromQueue := func(delay time.Duration, i int) {
		time.Sleep(delay)
		c.L.Lock()
		fmt.Println("before remove:", queue)
		queue = queue[1:]
		fmt.Println("after remove", queue)
		c.L.Unlock()
		c.Signal()
	}

	for i := 0; i < 10; i++ {
		fmt.Println("Start loop", i)
		c.L.Lock()
		for len(queue) == 2 {
			fmt.Println("len equal 2, waiting", i)
			c.Wait()
		}

		fmt.Println("Adding to queue", i)
		queue = append(queue, i)
		go removeFromQueue(1*time.Second, i)
		c.L.Unlock()
	}
}

func test_range_channel() {
	intStream := make(chan int)
	go func() {
		// defer close(intStream)
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()

	for inter := range intStream {
		fmt.Printf("%v", inter)
	}
}

func testbufferchannel() {
	test_buffer := func() <-chan int {
		result := make(chan int, 5)
		go func() {
			defer close(result)
			for i := 0; i < 7; i++ {
				result <- i
			}
		}()
		return result
	}
	rs_stream := test_buffer()
	for rs := range rs_stream {
		fmt.Printf("Received: %d\n", rs)
	}

}

func test_adhocconfinement() {
	data := make([]int, 4)

	loopdata := func(handledata chan<- int) {
		defer close(handledata)
		for i := range data {
			handledata <- data[i]
		}
	}

	handledata := make(chan int)
	go loopdata(handledata)
	for num := range handledata {
		fmt.Println(num)
	}

}

func test_nil_channel() {
	dowork := func(strings <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("dowork exited")
			defer close(completed)
			for s := range strings {
				fmt.Println(s)
			}
		}()
		return completed
	}

	dowork(nil)
	fmt.Println("Done")
}

func newRandStream() {
	newRandStream := func(done chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Intn(100):
				case <-done:
					return
				}
			}
		}()

		return randStream
	}
	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 0; i <= 4; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)
	time.Sleep(1 * time.Second)

}

func main() {
	// TestOnceSync()
	newRandStream()
	// stringStream := make(chan string)
	// go func() {
	// 	stringStream <- "Hello channels!"
	// }()
	// salutation, ok := <-stringStream
	// fmt.Printf("(%v): %v", ok, salutation)
	// begin := make(chan interface{})
	// var wg sync.WaitGroup
	// for i := 0; i < 5; i++ {
	// 	wg.Add(1)
	// 	go func(i int) {
	// 		defer wg.Done()
	// 		<-begin
	// 		fmt.Printf("%v has begun\n", i)
	// 	}(i)
	// }
	// testbufferchannel()
	// fmt.Println("Unblocking goroutines...")
	// close(begin)
	// wg.Wait()
}
