// Package handler
package handler

import (
	"context"
	"log"

	pb "ecommerce-api/gen/order"
)

type OrderGRPCHandler struct {
	pb.UnimplementedOrderServiceServer
}

func NewOrderGRPCHandler() *OrderGRPCHandler {
	return &OrderGRPCHandler{}
}

func (h *OrderGRPCHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	log.Println("received create order request")

	return &pb.OrderResponse{
		Id:          "1",
		CustomerId:  "1",
		Status:      "Done",
		TotalAmount: 10.62,
		Items:       []*pb.OrderItem{{ProductId: "3", Quantity: 3, Price: 10.62}},
	}, nil
}
