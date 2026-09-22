// Package handler for the stock service
package handler

import (
	"context"

	stockpb "ecommerce-api/gen/stock"
)

type StockGRPCHandler struct {
	stockpb.UnimplementedStockServiceServer
}

func NewStockGRPCHandler() *StockGRPCHandler {
	return &StockGRPCHandler{}
}

func (p *StockGRPCHandler) CheckProductInStore(ctx context.Context, req *stockpb.CheckProductInStoreRequest) (*stockpb.CheckProductInStoreResponse, error) {
	return nil, nil
}
