package main

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Địa chỉ Redis server
		DB:   0,
	})
	ticketID := "ticket:123"
	success := tryLockTicket(rdb, ticketID, 10*time.Second)
	if success {
		fmt.Println("Vé đã được giữ thành công!")
		// Mô phỏng thời gian xử lý
		time.Sleep(5 * time.Second)

		// Người dùng hoàn tất giữ vé hoặc thanh toán
		releaseTicket(rdb, ticketID)
		fmt.Println("Vé đã được giải phóng!")
	} else {
		fmt.Println("Vé hiện không khả dụng, vui lòng thử lại sau!")
	}
}

func tryLockTicket(rdb *redis.Client, ticketID string, ttl time.Duration) bool {
	// Sử dụng SETNX để đặt khóa nếu key chưa tồn tại
	result, err := rdb.SetNX(ctx, ticketID, "locked", ttl).Result()
	if err != nil {
		fmt.Printf("Lỗi khi thực hiện SETNX: %v\n", err)
		return false
	}
	return result // Trả về true nếu SETNX thành công, false nếu thất bại
}

// releaseTicket: Giải phóng vé bằng cách xóa key trong Redis
func releaseTicket(rdb *redis.Client, ticketID string) {
	_, err := rdb.Del(ctx, ticketID).Result()
	if err != nil {
		fmt.Printf("Lỗi khi giải phóng vé: %v\n", err)
	}
}
