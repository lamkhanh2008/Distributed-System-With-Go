package main

// var ctx = context.Background()

// func main() {
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     "localhost:6379", // Địa chỉ và cổng của Redis Server
// 		Password: "",               // Mặc định Redis không có mật khẩu
// 		DB:       0,                // Sử dụng DB 0 mặc định
// 	})

// 	// Bước 2: Gửi lệnh PING đến Redis để kiểm tra kết nối
// 	pong, err := rdb.Ping(ctx).Result()
// 	if err != nil {
// 		log.Fatalf("Không thể kết nối đến Redis: %v", err)
// 	}

// 	fmt.Printf("Kết nối thành công đến Redis! Phản hồi: %s\n", pong)
// }
