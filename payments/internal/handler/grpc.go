package handler

import (
	"context"
	"log"

	paymentpb "ecommerce-api/gen/payment"
)

type PaymentGRPCHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
}

func NewPaymentGRPCHandler() *PaymentGRPCHandler {
	return &PaymentGRPCHandler{}
}

func (h *PaymentGRPCHandler) ProcessPayment(ctx context.Context, req *paymentpb.ProcessPaymentRequest) (*paymentpb.ProcessPaymentResponse, error) {
	log.Printf("processing mock payment for order %s, customer %s, amount %.2f", req.OrderId, req.CustomerId, req.Amount)

	if req.Amount <= 0 {
		return &paymentpb.ProcessPaymentResponse{
			PaymentId: "mock-payment-" + req.OrderId,
			Status:    "failed",
		}, nil
	}

	if req.Amount > 10000 {
		return &paymentpb.ProcessPaymentResponse{
			PaymentId: "mock-payment-" + req.OrderId,
			Status:    "failed",
		}, nil
	}

	return &paymentpb.ProcessPaymentResponse{
		PaymentId: "mock-payment-" + req.OrderId,
		Status:    "success",
	}, nil
}
