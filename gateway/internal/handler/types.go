package handler

import (
	orderpb "ecommerce-api/gen/order"
)

type CreateUserPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type CreateOrderPayload struct {
	CustomerID string               `json:"customer_id"`
	Items      []*orderpb.OrderItem `json:"items"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
