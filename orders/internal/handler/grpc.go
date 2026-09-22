package handler

import (
	"context"
	"ecommerce-orders/internal/db"
	"log"

	orderpb "ecommerce-api/gen/order"
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
	log.Printf("received create order request for customer %s with %d items", req.CustomerId, len(req.Items))

	orderID, totalAmount, err := h.models.OrderModel.CreateOrder(ctx, req.CustomerId, req.Items)
	if err != nil {
		return nil, err
	}

	return &orderpb.OrderResponse{
		Id:          orderID,
		CustomerId:  req.CustomerId,
		Status:      "pending",
		TotalAmount: totalAmount,
		Items:       req.Items,
	}, nil
}

func (h *OrderGRPCHandler) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.OrderResponse, error) {
	log.Printf("received get order request for order %s", req.Id)
	return h.models.OrderModel.GetOrder(ctx, req.Id)
}

func (h *OrderGRPCHandler) CheckProductInStore(ctx context.Context, req *orderpb.CheckProductInStoreRequest) (*orderpb.CheckProductInStoreResponse, error) {
	log.Printf("received check product in store request for %s", req.ProductName)
	quantity, err := h.models.OrderModel.CheckProductInStore(ctx, req.ProductName)
	if err != nil {
		return nil, err
	}
	return &orderpb.CheckProductInStoreResponse{Quantity: quantity}, nil
}
