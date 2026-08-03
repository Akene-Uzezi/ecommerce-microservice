package handler

import (
	orderpb "ecommerce-api/gen/order"

	"github.com/golang-jwt/jwt/v5"
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

type Claims struct {
	Email string
	Name  string
	jwt.RegisteredClaims
}
