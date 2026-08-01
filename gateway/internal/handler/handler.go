// Package handler for the gateway service
package handler

import (
	shared "ecommerce-shared"
	"log"
	"net/http"

	pb "ecommerce-api/gen/order"
)

type handler struct {
	orderClient pb.OrderServiceClient
}

func NewHTTPHandler(orderClient pb.OrderServiceClient) *handler {
	return &handler{
		orderClient: orderClient,
	}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/ping", h.ping)
	mux.HandleFunc("POST /api/v1/orders", h.createOrder)
}

func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	_ = shared.WriteJSON(w, http.StatusOK, "pong")
}

type CreateOrderPayload struct {
	CustomerID string          `json:"customer_id"`
	Items      []*pb.OrderItem `json:"items"`
}

func (h *handler) createOrder(w http.ResponseWriter, r *http.Request) {
	var body CreateOrderPayload
	if err := shared.ReadJSON(r, body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	res, err := h.orderClient.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId: body.CustomerID,
		Items:      body.Items,
	})
	if err != nil {
		log.Printf("grpc create order failed: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := shared.WriteJSON(w, http.StatusCreated, res); err != nil {
		log.Printf("failed to write json: %v", err)
	}
}
