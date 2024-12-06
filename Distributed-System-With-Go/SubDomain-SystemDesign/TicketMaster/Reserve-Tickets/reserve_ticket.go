package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	rdb        *redis.Client
	ctx        = context.Background()
	expiration = time.Second
)

func acquireLock(ticketID string) bool {
	lockKey := "ticket:" + ticketID
	set, err := rdb.SetNX(ctx, lockKey, "locked", expiration).Result()
	if err != nil {
		log.Println("Error acquiring lock: ", err)
		return false
	}

	if set {
		fmt.Println("Locked ok")
		return true
	}
	fmt.Printf("Ticket %s is already reserved by another user.\n", ticketID)
	return false
}

func releaseLock(ticketID string) bool {
	lockKey := "ticket:" + ticketID
	err := rdb.Del(ctx, lockKey).Err()
	if err != nil {
		log.Println("Error releaseLock lock: ", err)
		return false
	}
	fmt.Println("Success release")
	return true
}
func updateTicketStatus(ticketID, status string) {
	// In a real-world scenario, you'd update the status in a SQL/NoSQL database
	fmt.Printf("Ticket %s status updated to '%s'.\n", ticketID, status)
}

func processCheckout(ticketID string, completePurchase bool) {
	if acquireLock(ticketID) {
		if completePurchase {
			updateTicketStatus(ticketID, "Booked")
			releaseLock(ticketID)
		} else {
			fmt.Printf("User abandoned the checkout for ticket %s.\n", ticketID)
		}
	} else {
		fmt.Printf("Unable to book ticket %s, it is already locked.\n", ticketID)
	}
}
func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})
}
func main() {
	initRedis()
	defer rdb.Close()
	ticketID := "12345" // Example ticket ID

	userCompletePurchase := true // Change to false to simulate abandonment

	processCheckout(ticketID, userCompletePurchase)
}
