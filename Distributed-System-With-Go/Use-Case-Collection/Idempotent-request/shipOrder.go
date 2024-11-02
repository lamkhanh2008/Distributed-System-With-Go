package main

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ShippingOrder struct {
	gorm.Model
	OrderID string `json:"order_id" gorm:"index"`
	Vendor  string `json:"vendor"`
	Address string `json:"address"`
}

type Shipping struct {
	db *gorm.DB
}

func (s *Shipping) Save(ctx context.Context, order ShippingOrder) (ShippingOrder, error) {
	err := s.db.WithContext(ctx).Save(&order).Error
	return order, err
}

func (s *Shipping) ByID(ctx context.Context, id uint) (*ShippingOrder, error) {
	return s.by(ctx, "id", id)
}

func (s *Shipping) ByOrderID(ctx context.Context, id uint) (*ShippingOrder, error) {
	return s.by(ctx, "order_id", id)
}
func (s *Shipping) by(ctx context.Context, key string, val any) (*ShippingOrder, error) {
	var order ShippingOrder
	if tr := s.db.WithContext(ctx).Where(key+"=?", val).First(&order); tr.Error != nil {
		return nil, tr.Error
	}
	return &order, nil
}

type PlaceShippingOrderRequest struct {
	OrderID string `json:"order_id"`
	Vendor  string `json:"vendor"`
	Address string `json:"address"`
	// and etc ...
}

func main() {
	var e = echo.New()
	e.POST("/shipping/order", func(c echo.Context) error { // this API is used to place a shipping order
		<-time.After(time.Second * 2)
		var request PlaceShippingOrderRequest
		if err := c.Bind(&request); err != nil {
			return err
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		shippingRepository := Shipping{}
		createdOrder, err := shippingRepository.Save(ctx, ShippingOrder{
			OrderID: request.OrderID,
			Vendor:  request.Vendor,
			Address: request.Address,
		})
		if err != nil {
			return err
		}
		return c.JSON(201, map[string]any{
			"ok":          true,
			"shipping_id": createdOrder.ID,
		})

	})
}
