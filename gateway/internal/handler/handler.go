// Package handler for the gateway service
package handler

import (
	shared "ecommerce-shared"
	"log"
	"net/http"
	"time"

	orderpb "ecommerce-api/gen/order"
)

type handler struct {
	orderClient orderpb.OrderServiceClient
}

func NewOrderHTTPHandler(orderClient orderpb.OrderServiceClient) *handler {
	return &handler{
		orderClient: orderClient,
	}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/ping", h.ping)
	mux.HandleFunc("POST /api/v1/orders", h.createOrder)
}

func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	_ = shared.WriteJSON(w, http.StatusOK, "pong")
	shared.LogOK(r.Method, r.RequestURI, time.Since(start))
}

func (h *handler) createOrder(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var body CreateOrderPayload
	if err := shared.ReadJSON(r, &body); err != nil {
		shared.WriteErrorBadRequest(w, "Invalid request body", err)
		shared.LogBadRequest(r.Method, r.RequestURI, time.Since(start))
		return
	}

	res, err := h.orderClient.CreateOrder(r.Context(), &orderpb.CreateOrderRequest{
		CustomerId: body.CustomerID,
		Items:      body.Items,
	})
	if err != nil {
		log.Printf("grpc create order failed: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		shared.LogInternalServerError(r.Method, r.RequestURI, time.Since(start))
		return
	}

	if err := shared.WriteJSON(w, http.StatusCreated, res); err != nil {
		log.Printf("failed to write json: %v", err)
		return
	}
	shared.LogOK(r.Method, r.RequestURI, time.Since(start))
}
