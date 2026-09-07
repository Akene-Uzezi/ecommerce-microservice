package db

import (
	"context"
	"fmt"

	productspb "ecommerce-api/gen/products"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductModel struct {
	DB *pgxpool.Pool
}

type Product struct {
	ID       int
	Name     string
	Price    float32
	Quantity int
}

func (m *ProductModel) AddProduct(ctx context.Context, req *productspb.AddProductRequest) (*productspb.AddProductResponse, error) {
	var createdProduct productspb.AddProductResponse
	name := req.Product.Name
	price := req.Product.Price
	quantity := req.Product.Quantity
	query := `
		INSERT INTO products
		(name, price, quantity)
		VALUES ($1, $2, $3)
		RETURNING id, name, price, quantity
	`
	err := m.DB.QueryRow(ctx, query, name, price, quantity).Scan(
		&createdProduct.Product.Name, &createdProduct.Product.Price, &createdProduct.Product.Quantity,
	)
	if err != nil {
		return nil, fmt.Errorf("database query error %v", err)
	}
	return &createdProduct, nil
}

func (m *ProductModel) GetProducts(ctx context.Context, req *productspb.GetProductsRequest) (*productspb.GetProductsResponse, error) {
	return nil, nil
}

func (m *ProductModel) GetProduct(ctx context.Context, req *productspb.GetProductRequest) (*productspb.GetProductResponse, error) {
	return nil, nil
}
