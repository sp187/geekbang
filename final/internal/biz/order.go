package biz

import (
	"context"
	"time"
)

type Order struct {
	ID          uint      `json:"id"`
	GoodsID     []uint    `json:"goods_id"`
	GoodsAmount []uint    `json:"goods_amount"`
	Created     time.Time `json:"created"`
	Status      int       `json:"status"`
}

type OrderRepo interface {
	GetUserOrders(context.Context, uint) ([]Order, error)
	GetOrders(context.Context, []uint) ([]Order, error)
}
