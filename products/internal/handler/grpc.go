// Package handler for the products service
package handler

import (
	"context"

	prodcutspb "ecommerce-api/gen/products"
)

type ProductGRPCHandler struct {
	prodcutspb.UnimplementedProductServiceServer
}

func NewProductGRPCHandler() *ProductGRPCHandler {
	return &ProductGRPCHandler{}
}

func (p *ProductGRPCHandler) AddProduct(ctx context.Context, req *prodcutspb.AddProductRequest) (*prodcutspb.AddProductResponse, error) {
	return nil, nil
}

func (p *ProductGRPCHandler) GetProduct(ctx context.Context, req *prodcutspb.GetProductsRequest) (*prodcutspb.GetProductsResponse, error) {
	return nil, nil
}
