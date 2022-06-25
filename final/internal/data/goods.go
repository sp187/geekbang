package data

import (
	"context"

	"github.com/sp187/geekbang/final/internal/biz"
	fw "github.com/sp187/geekbang/final/internal/framework"
)

type GoodsRPC struct {
	serviceURI string
}

type rpcGoods biz.Goods

type OrderGoodsResp struct {
	Code int        `json:"code"`
	Data []rpcGoods `json:"data"`
}

func (g *GoodsRPC) GetByIDs(ctx context.Context, id []uint) ([]biz.Goods, error) {
	req := struct {
		ID []uint `json:"id"`
	}{id}
	var resp OrderGoodsResp
	_, err := fw.GetHttpClient().RequestCtx(ctx, "POST", g.serviceURI, req, &resp)
	if err != nil {
		return nil, err
	}
	results := make([]biz.Goods, 0, len(resp.Data))
	for _, d := range resp.Data {
		results = append(results, biz.Goods(d))
	}
	return results, nil
}
