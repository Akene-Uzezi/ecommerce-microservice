// Package handler for the stock service
package handler

import stockpb "ecommerce-api/gen/stock"

type StockGRPCHandler struct {
	stockpb.UnimplementedStockServiceServer
}

func NewStockGRPCHandler() *StockGRPCHandler {
	return &StockGRPCHandler{}
}
