package db

import (
	"context"
	"fmt"

	stockpb "ecommerce-api/gen/stock"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StockModel struct {
	DB *pgxpool.Pool
}

func NewStockModel(db *pgxpool.Pool) *StockModel {
	return &StockModel{DB: db}
}

type Models struct {
	StockModel StockModel
}

func NewModels(db *pgxpool.Pool) *Models {
	return &Models{StockModel: StockModel{DB: db}}
}

func (m *StockModel) CheckStock(ctx context.Context, req *stockpb.CheckStockRequest) (*stockpb.CheckStockResponse, error) {
	var quantity uint32
	query := `SELECT quantity FROM products WHERE name = $1`
	err := m.DB.QueryRow(ctx, query, req.ProductName).Scan(&quantity)
	if err != nil {
		return nil, fmt.Errorf("database query error: %s", err)
	}
	return &stockpb.CheckStockResponse{Quantity: quantity}, nil
}

func (m *StockModel) ReserveStock(ctx context.Context, req *stockpb.ReserveStockRequest) (*stockpb.ReserveStockResponse, error) {
	tx, err := m.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %s", err)
	}
	defer tx.Rollback(ctx)

	for _, item := range req.Items {
		var currentQty uint32
		err := tx.QueryRow(ctx, `SELECT quantity FROM products WHERE name = $1 FOR UPDATE`, item.ProductName).Scan(&currentQty)
		if err != nil {
			return nil, fmt.Errorf("failed to check stock for %s: %s", item.ProductName, err)
		}
		if currentQty < uint32(item.Quantity) {
			return &stockpb.ReserveStockResponse{Success: false}, nil
		}
		_, err = tx.Exec(ctx, `UPDATE products SET quantity = quantity - $1 WHERE name = $2`, item.Quantity, item.ProductName)
		if err != nil {
			return nil, fmt.Errorf("failed to reserve stock for %s: %s", item.ProductName, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %s", err)
	}

	return &stockpb.ReserveStockResponse{Success: true}, nil
}
