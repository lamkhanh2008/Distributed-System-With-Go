package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

type user struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Score int    `json:"score"`
}

func addUser(name string) (user, error) {
	id := uuid.NewString()
	fmt.Println("--id: ", id)
	newUser := user{
		Name:  name,
		ID:    id,
		Score: len(name),
	}

	serializeUser, err := json.Marshal(newUser)
	if err != nil {
		return user{}, err
	}

	err = rdb.Set(ctx, id, serializeUser, 0).Err()
	if err != nil {
		return user{}, err
	}
	val, err := rdb.Get(ctx, "e7b12fae-09ec-438f-9ba7-43b81faf7acb").Result()
	if err != nil {
		log.Fatalf("Không thể lấy key 'username': %v", err)
	}
	fmt.Println("val: ", val)
	err = rdb.ZAdd(ctx, "ranking", redis.Z{Score: float64(newUser.Score), Member: id}).Err()
	if err != nil {
		return user{}, err
	}

	rank, err := rdb.ZRevRank(ctx, "ranking", id).Result()
	if err != nil {
		return user{}, err
	}
	fmt.Println(rank)
	newUser.Rank = int(rank + 1)
	return newUser, nil
}
