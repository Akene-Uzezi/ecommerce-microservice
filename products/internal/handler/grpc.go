// Package handler for the products service
package handler

import (
	"context"
	"ecommerce-products/internal/db"

	prodcutspb "ecommerce-api/gen/products"
)

type ProductGRPCHandler struct {
	prodcutspb.UnimplementedProductServiceServer
	models *db.Models
}

func NewProductGRPCHandler(models *db.Models) *ProductGRPCHandler {
	return &ProductGRPCHandler{
		models: models,
	}
}

func (p *ProductGRPCHandler) AddProduct(ctx context.Context, req *prodcutspb.AddProductRequest) (*prodcutspb.AddProductResponse, error) {
	return nil, nil
}

func (p *ProductGRPCHandler) GetProduct(ctx context.Context, req *prodcutspb.GetProductsRequest) (*prodcutspb.GetProductsResponse, error) {
	return nil, nil
}
