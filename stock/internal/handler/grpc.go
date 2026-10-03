package handler

import (
	"context"
	"time"

	"ecommerce-stock/internal/db"

	stockpb "ecommerce-api/gen/stock"
	shared "ecommerce-shared"
)

type StockGRPCHandler struct {
	stockpb.UnimplementedStockServiceServer
	models *db.Models
}

func NewStockGRPCHandler(models *db.Models) *StockGRPCHandler {
	return &StockGRPCHandler{
		models: models,
	}
}

func (h *StockGRPCHandler) CheckStock(ctx context.Context, req *stockpb.CheckStockRequest) (*stockpb.CheckStockResponse, error) {
	start := time.Now()
	result, err := h.models.StockModel.CheckStock(ctx, req)
	if err != nil {
		shared.LogRequest("CheckStock", "/stock.CheckStock", 500, time.Since(start))
		return nil, err
	}
	shared.LogRequest("CheckStock", "/stock.CheckStock", 0, time.Since(start))
	return result, nil
}

func (h *StockGRPCHandler) ReserveStock(ctx context.Context, req *stockpb.ReserveStockRequest) (*stockpb.ReserveStockResponse, error) {
	start := time.Now()
	result, err := h.models.StockModel.ReserveStock(ctx, req)
	if err != nil {
		shared.LogRequest("ReserveStock", "/stock.ReserveStock", 500, time.Since(start))
		return nil, err
	}
	shared.LogRequest("ReserveStock", "/stock.ReserveStock", 0, time.Since(start))
	return result, nil
}
