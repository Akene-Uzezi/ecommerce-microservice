package handler

import (
	"context"

	"ecommerce-stock/internal/db"

	stockpb "ecommerce-api/gen/stock"
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
	return h.models.StockModel.CheckStock(ctx, req)
}

func (h *StockGRPCHandler) ReserveStock(ctx context.Context, req *stockpb.ReserveStockRequest) (*stockpb.ReserveStockResponse, error) {
	return h.models.StockModel.ReserveStock(ctx, req)
}
