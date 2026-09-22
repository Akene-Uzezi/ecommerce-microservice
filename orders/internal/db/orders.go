package db

import (
	"context"
	"fmt"

	orderpb "ecommerce-api/gen/order"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderModel struct {
	DB *pgxpool.Pool
}

func NewOrderModel(db *pgxpool.Pool) *OrderModel {
	return &OrderModel{DB: db}
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID string
	Quantity  int
	Price     float64
}

func (m *OrderModel) CreateOrder(ctx context.Context, customerID string, items []*orderpb.OrderItem) (string, float64, error) {
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return "", 0, fmt.Errorf("failed to begin transaction: %s", err)
	}
	defer tx.Rollback(ctx)

	var orderID int
	err = tx.QueryRow(ctx, `INSERT INTO orders (customer_id, status) VALUES ($1, 'pending') RETURNING id`, customerID).Scan(&orderID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create order: %s", err)
	}

	var totalAmount float64
	for _, item := range items {
		_, err = tx.Exec(ctx, `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)`, orderID, item.ProductId, item.Quantity, item.Price)
		if err != nil {
			return "", 0, fmt.Errorf("failed to create order item: %s", err)
		}
		totalAmount += float64(item.Quantity) * item.Price
	}

	_, err = tx.Exec(ctx, `UPDATE orders SET total_amount = $1 WHERE id = $2`, totalAmount, orderID)
	if err != nil {
		return "", 0, fmt.Errorf("failed to update order total: %s", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return "", 0, fmt.Errorf("failed to commit transaction: %s", err)
	}

	return fmt.Sprintf("%d", orderID), totalAmount, nil
}

func (m *OrderModel) GetOrder(ctx context.Context, orderID string) (*orderpb.OrderResponse, error) {
	var customerID string
	var status string
	var totalAmount float64
	err := m.DB.QueryRow(ctx, `SELECT customer_id, status, total_amount FROM orders WHERE id = $1`, orderID).Scan(&customerID, &status, &totalAmount)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %s", err)
	}

	rows, err := m.DB.Query(ctx, `SELECT product_id, quantity, price FROM order_items WHERE order_id = $1`, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %s", err)
	}
	defer rows.Close()

	var items []*orderpb.OrderItem
	for rows.Next() {
		var productID string
		var quantity int32
		var price float64
		if err := rows.Scan(&productID, &quantity, &price); err != nil {
			return nil, fmt.Errorf("failed to scan order item: %s", err)
		}
		items = append(items, &orderpb.OrderItem{
			ProductId: productID,
			Quantity:  quantity,
			Price:     price,
		})
	}

	return &orderpb.OrderResponse{
		Id:          orderID,
		CustomerId:  customerID,
		Status:      status,
		TotalAmount: totalAmount,
		Items:       items,
	}, nil
}

func (m *OrderModel) CheckProductInStore(ctx context.Context, productName string) (uint32, error) {
	var quantity uint32
	err := m.DB.QueryRow(ctx, `SELECT quantity FROM products WHERE name = $1`, productName).Scan(&quantity)
	if err != nil {
		return 0, fmt.Errorf("failed to check product in store: %s", err)
	}
	return quantity, nil
}
