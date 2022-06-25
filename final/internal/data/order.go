package data

import (
	"context"

	"github.com/sp187/geekbang/final/internal/biz"
	fw "github.com/sp187/geekbang/final/internal/framework"
)

type OrderRPC struct {
	serviceURI string
}

type rpcOrder biz.Order

type UserOrderResp struct {
	Code int        `json:"code"`
	Data []rpcOrder `json:"data"`
}

func (o *OrderRPC) GetUserOrders(ctx context.Context, u uint) ([]biz.Order, error) {
	var resp UserOrderResp
	_, err := fw.GetHttpClient().RequestCtx(ctx, "GET", o.serviceURI, nil, &resp)
	if err != nil {
		return nil, err
	}
	results := make([]biz.Order, 0, len(resp.Data))
	for _, d := range resp.Data {
		results = append(results, biz.Order(d))
	}
	return results, nil
}

func (o *OrderRPC) GetOrders(ctx context.Context, id []uint) ([]biz.Order, error) {
	req := struct {
		ID []uint `json:"id"`
	}{id}
	var resp UserOrderResp
	_, err := fw.GetHttpClient().RequestCtx(ctx, "POST", o.serviceURI, req, &resp)
	if err != nil {
		return nil, err
	}
	results := make([]biz.Order, 0, len(resp.Data))
	for _, d := range resp.Data {
		results = append(results, biz.Order(d))
	}
	return results, nil
}
