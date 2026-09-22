package handler

import (
	"ecommerce-gateway/internal/middleware"
	shared "ecommerce-shared"
	"log"
	"net/http"
	"time"

	authpb "ecommerce-api/gen/auth"
	paymentpb "ecommerce-api/gen/payment"
)

type PaymentsHTTPHandler struct {
	paymentsClient paymentpb.PaymentServiceClient
	authClient     authpb.AuthServiceClient
}

func NewPaymentsHTTPHandler(paymentsClient paymentpb.PaymentServiceClient, authClient authpb.AuthServiceClient) *PaymentsHTTPHandler {
	return &PaymentsHTTPHandler{
		paymentsClient: paymentsClient,
		authClient:     authClient,
	}
}

func (h *PaymentsHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/payments", middleware.RequireAuth(h.authClient)(h.processPayment))
}

func (h *PaymentsHTTPHandler) processPayment(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var requestBody paymentpb.ProcessPaymentRequest
	if err := shared.ReadJSON(r, &requestBody); err != nil {
		shared.WriteErrorBadRequest(w, "Invalid request body", err)
		shared.LogBadRequest(r.Method, r.RequestURI, time.Since(start))
		return
	}

	ctx := r.Context()
	res, err := h.paymentsClient.ProcessPayment(ctx, &requestBody)
	if err != nil {
		log.Printf("failed to process payment: %s", err)
		shared.WriteErrorServerError(w, "failed to process payment", err)
		shared.LogInternalServerError(r.Method, r.RequestURI, time.Since(start))
		return
	}

	_ = shared.WriteJSON(w, http.StatusOK, res)
	shared.LogOK(r.Method, r.RequestURI, time.Since(start))
}
