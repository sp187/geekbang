package biz

import "context"

type Goods struct {
	ID          uint
	Description string
	Price       float64
}

type GoodsRepo interface {
	GetByIDs(context.Context, []uint) ([]Goods, error)
}
