package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/labstack/echo/v4"
)

type Redis[T any] struct {
	client *redis.Client
}

func NewRedis[T any](rdb *redis.Client) *Redis[T] {
	return &Redis[T]{client: rdb}
}

func (r *Redis[T]) Start(ctx context.Context, idempotenKey string) (T, bool, error) {
	var t T
	tr := r.client.HSetNX(ctx, "idempotency:"+idempotenKey, "status", "started")
	if tr.Err() != nil {
		return t, false, tr.Err()
	}
	if tr.Val() {
		return t, false, nil
	}
	b, err := r.client.HGet(ctx, "idempotency:"+idempotenKey, "value").Bytes()
	if err != nil {
		return t, false, err
	}
	if err := json.Unmarshal(b, &t); err != nil {
		return t, false, err
	}
	return t, true, nil
}

func (r *Redis[T]) Store(ctx context.Context, idempotencyKey string, value T) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.HSet(ctx, "idempotency:"+idempotencyKey, "value", b).Err()
}

type shippingIdempotency struct {
	redis *Redis[PlaceShippingOrderRequest]
}

func (s *shippingIdempotency) Start(ctx context.Context, orderID string) (stored PlaceShippingOrderRequest, has bool, err error) {
	return
}

func TestIdempotency() {
	e := echo.New()
	e.POST("/shipping/order", func(c echo.Context) error { // this API is used to place a shipping order
		var request PlaceShippingOrderRequest
		if err := c.Bind(&request); err != nil {
			return err
		}
		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		Shipidempotency := &shippingIdempotency{
			redis: NewRedis[PlaceShippingOrderRequest](rdb),
		}
		stored, has, err := Shipidempotency.redis.Start(context.Background(), request.OrderID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"ok":    false,
				"error": "Internal server error",
			})
		}
		if has {
			return c.JSON(200, stored)
		}
		<-time.After(time.Second * 2)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		shippingRepository := Shipping{}
		createdOrder, err := shippingRepository.Save(ctx, ShippingOrder{
			OrderID: request.OrderID,
			Vendor:  request.Vendor,
			Address: request.Address,
		})
		//////////// saving the final value for future requests
		placeShipping := PlaceShippingOrderRequest{OrderID: request.OrderID,
			Vendor:  request.Vendor,
			Address: request.Address}
		if err := Shipidempotency.redis.Store(context.Background(), createdOrder.OrderID, placeShipping); err != nil {
			return err
		}
		////////////
		if err != nil {
			return err
		}
		return c.JSON(201, map[string]any{
			"ok":          true,
			"shipping_id": createdOrder.ID,
		})

	})
}
