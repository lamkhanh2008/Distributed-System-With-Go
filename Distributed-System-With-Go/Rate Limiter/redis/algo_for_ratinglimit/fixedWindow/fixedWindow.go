package algoRate

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func fixWindowAlgorithm(client *redis.Client, key string, limit int64, window time.Duration) bool {
	currentTime := time.Now()
	keyWindow := fmt.Sprintf("%s_%d", key, currentTime.Unix()/int64(window))

	client.Incr(key)
	count, err := client.Get(keyWindow).Int64()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if count >= limit {
		return false
	}
	pipe := client.TxPipeline()
	pipe.Incr(keyWindow)
	pipe.Expire(keyWindow, window)
	_, err = pipe.Exec()
	if err != nil {
		panic(err)
	}

	return true
}

func SlidingWindowAlgorithm(client *redis.Client, key string, limit int64, window time.Duration) bool {
	keyWindow := fmt.Sprintf("%s_window", key)
	currentTime := time.Now().UnixNano()

	pipe := client.TxPipeline()
	pipe.ZRemRangeByScore(keyWindow, "0", fmt.Sprintf("%d", currentTime-int64(window)))

	pipe.ZAdd(keyWindow, redis.Z{
		Score:  float64(currentTime),
		Member: currentTime,
	})

	pipe.ZCard(keyWindow)
	_, err := pipe.Exec()
	if err != nil {
		panic(err)
	}

	count, err := client.Get(keyWindow).Int64()
	if err != nil && err != redis.Nil {
		panic(err)
	}
	if count > limit {
		return false
	}
	return true
}
