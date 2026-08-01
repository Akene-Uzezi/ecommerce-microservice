// Package handler
package handler

import (
	"context"
	"log"

	orderpb "ecommerce-api/gen/order"
)

type OrderGRPCHandler struct {
	orderpb.UnimplementedOrderServiceServer
}

func NewOrderGRPCHandler() *OrderGRPCHandler {
	return &OrderGRPCHandler{}
}

func (h *OrderGRPCHandler) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.OrderResponse, error) {
	log.Println("received create order request")

	return &orderpb.OrderResponse{
		Id:          "1",
		CustomerId:  "1",
		Status:      "Done",
		TotalAmount: 10.62,
		Items:       []*orderpb.OrderItem{{ProductId: "3", Quantity: 3, Price: 10.62}},
	}, nil
}
