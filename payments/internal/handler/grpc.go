package handler

import (
	"context"
	"time"

	paymentpb "ecommerce-api/gen/payment"
	shared "ecommerce-shared"
)

type PaymentGRPCHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
}

func NewPaymentGRPCHandler() *PaymentGRPCHandler {
	return &PaymentGRPCHandler{}
}

func (h *PaymentGRPCHandler) ProcessPayment(ctx context.Context, req *paymentpb.ProcessPaymentRequest) (*paymentpb.ProcessPaymentResponse, error) {
	start := time.Now()
	if req.Amount <= 0 {
		res := &paymentpb.ProcessPaymentResponse{
			PaymentId: "mock-payment-" + req.OrderId,
			Status:    "failed",
		}
		shared.LogRequest("ProcessPayment", "/payment.ProcessPayment", 400, time.Since(start))
		return res, nil
	}

	if req.Amount > 10000 {
		res := &paymentpb.ProcessPaymentResponse{
			PaymentId: "mock-payment-" + req.OrderId,
			Status:    "failed",
		}
		shared.LogRequest("ProcessPayment", "/payment.ProcessPayment", 400, time.Since(start))
		return res, nil
	}

	res := &paymentpb.ProcessPaymentResponse{
		PaymentId: "mock-payment-" + req.OrderId,
		Status:    "success",
	}
	shared.LogRequest("ProcessPayment", "/payment.ProcessPayment", 0, time.Since(start))
	return res, nil
}
