package main

import (
	"fmt"
	"time"
)

//work same queue

func main() {
	Run_worker_pool()
}

func SendReceiveValueToBufferChannel() {
	bufferChannel := make(chan int, 10)
	bufferChannel <- 1
	bufferChannel <- 2
	bufferChannel <- 3
	close(bufferChannel)
	fmt.Print(len(bufferChannel))
	fmt.Printf("value in channel: %v \n", <-bufferChannel)
	fmt.Printf("value in channel: %v \n", <-bufferChannel)
	fmt.Printf("value in channel: %v \n", <-bufferChannel)
	fmt.Print(len(bufferChannel))
}
func Run_worker_pool() {
	const numJobs = 10
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	go worker("ti", jobs, results)
	go worker("teo", jobs, results)
	go worker("tun", jobs, results)
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)
	for a := 1; a <= numJobs; a++ {
		fmt.Println(<-results)
	}
}
func worker(id string, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}
