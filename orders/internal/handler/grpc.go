package handler

import (
	"context"
	"ecommerce-orders/internal/db"
	"time"

	orderpb "ecommerce-api/gen/order"
	shared "ecommerce-shared"
)

type OrderGRPCHandler struct {
	orderpb.UnimplementedOrderServiceServer
	models *db.Models
}

func NewOrderGRPCHandler(models *db.Models) *OrderGRPCHandler {
	return &OrderGRPCHandler{
		models: models,
	}
}

func (h *OrderGRPCHandler) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.OrderResponse, error) {
	start := time.Now()
	orderID, totalAmount, err := h.models.OrderModel.CreateOrder(ctx, req.CustomerId, req.Items)
	if err != nil {
		shared.LogRequest("CreateOrder", "/order.CreateOrder", 500, time.Since(start))
		return nil, err
	}
	shared.LogRequest("CreateOrder", "/order.CreateOrder", 0, time.Since(start))
	return &orderpb.OrderResponse{
		Id:          orderID,
		CustomerId:  req.CustomerId,
		Status:      "pending",
		TotalAmount: totalAmount,
		Items:       req.Items,
	}, nil
}

func (h *OrderGRPCHandler) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.OrderResponse, error) {
	start := time.Now()
	order, err := h.models.OrderModel.GetOrder(ctx, req.Id)
	if err != nil {
		shared.LogRequest("GetOrder", "/order.GetOrder", 500, time.Since(start))
		return nil, err
	}
	shared.LogRequest("GetOrder", "/order.GetOrder", 0, time.Since(start))
	return order, nil
}

func (h *OrderGRPCHandler) CheckProductInStore(ctx context.Context, req *orderpb.CheckProductInStoreRequest) (*orderpb.CheckProductInStoreResponse, error) {
	start := time.Now()
	quantity, err := h.models.OrderModel.CheckProductInStore(ctx, req.ProductName)
	if err != nil {
		shared.LogRequest("CheckProductInStore", "/order.CheckProductInStore", 500, time.Since(start))
		return nil, err
	}
	shared.LogRequest("CheckProductInStore", "/order.CheckProductInStore", 0, time.Since(start))
	return &orderpb.CheckProductInStoreResponse{Quantity: quantity}, nil
}
