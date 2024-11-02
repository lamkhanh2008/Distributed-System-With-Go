package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var rdb *redis.Client

type Task struct {
	Id   string
	Data string
}

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

}
func scheduleTask(task Task, ttl time.Duration) error {
	err := rdb.Set(ctx, task.Id, task.Data, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to schedule task: %v", err)
	}

	fmt.Printf("Task %s scheduled successfully with TTL %v\n", task.Id, ttl)
	return nil
}

func ListenForExpirations(wg *sync.WaitGroup) {
	defer wg.Done()

	expiredSub := rdb.PSubscribe(ctx, "__keyevent@0__:expired")
	fmt.Println("Listening for task expiration events...")

	for msg := range expiredSub.Channel() {
		taskID := msg.Payload
		fmt.Printf("Task expired: %s\n", taskID)
		err := rdb.Publish(ctx, "task_notification", fmt.Sprintf("Task %s expired", taskID)).Err()
		if err != nil {
			log.Printf("Failed to publish expiration notification: %v", err)
		}
	}
}

func listenForNotifications(wg *sync.WaitGroup) {
	defer wg.Done()

	sub := rdb.Subscribe(ctx, "task_notification")

	fmt.Println("Listening for task notifications...")

	for msg := range sub.Channel() {
		fmt.Printf("Received notification: %s\n", msg.Payload)
	}
}

func main() {
	initRedis()
	defer rdb.Close()

	// err := rdb.ConfigSet(ctx, "notify-keyspace-events", "Ex").Err()
	// if err != nil {
	// 	log.Fatalf("Failed to set Redis config: %v", err)
	// }
	err := rdb.ConfigSet(ctx, "notify-keyspace-events", "Ex").Err()
	if err != nil {
		log.Fatalf("Không thể bật notify-keyspace-events: %v", err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go ListenForExpirations(&wg)

	wg.Add(1)
	go listenForNotifications(&wg)

	scheduleTask(Task{Id: "task1", Data: "Task 1 Data"}, 5*time.Second)
	scheduleTask(Task{Id: "task2", Data: "Task 2 Data"}, 10*time.Second)
	scheduleTask(Task{Id: "task3", Data: "Task 3 Data"}, 15*time.Second)
	wg.Wait()
	fmt.Println("Shutting down application...")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
