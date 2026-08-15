package db

import (
	"context"

	productspb "ecommerce-api/gen/products"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductModel struct {
	DB *pgxpool.Pool
}

func (m *ProductModel) AddProduct(ctx context.Context, req *productspb.AddProductRequest) (*productspb.AddProductResponse, error) {
	return nil, nil
}

func (m *ProductModel) GetProduct(ctx context.Context, req *productspb.GetProductsRequest) (*productspb.GetProductsResponse, error) {
	return nil, nil
}
