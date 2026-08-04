package handler

import (
	orderpb "ecommerce-api/gen/order"
)

type CreateOrderPayload struct {
	CustomerID string               `json:"customer_id"`
	Items      []*orderpb.OrderItem `json:"items"`
}
