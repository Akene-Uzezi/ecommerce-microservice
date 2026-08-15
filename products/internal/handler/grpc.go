// Package handler for the products service
package handler

import (
	"context"

	prodcutspb "ecommerce-api/gen/products"
)

type ProductGRPCHandler struct {
	prodcutspb.UnimplementedProductServiceServer
}

func (p *ProductGRPCHandler) NewProductGRPCHanler() *ProductGRPCHandler {
	return &ProductGRPCHandler{}
}

func (p *ProductGRPCHandler) CreateProduct(ctx context.Context, req *prodcutspb.CreateProductRequest) (*prodcutspb.CreateProductResponse, error) {
	return nil, nil
}

func (p *ProductGRPCHandler) GetProduct(ctx context.Context, req *prodcutspb.GetProductsRequest) (*prodcutspb.GetProductsResponse, error) {
	return nil, nil
}
