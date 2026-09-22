// Package handler for the products service
package handler

import (
	"context"
	"ecommerce-products/internal/db"

	productspb "ecommerce-api/gen/products"
)

type ProductGRPCHandler struct {
	productspb.UnimplementedProductServiceServer
	models *db.Models
}

func NewProductGRPCHandler(models *db.Models) *ProductGRPCHandler {
	return &ProductGRPCHandler{
		models: models,
	}
}

func (p *ProductGRPCHandler) AddProduct(ctx context.Context, req *productspb.AddProductRequest) (*productspb.AddProductResponse, error) {
	product, err := p.models.ProductModel.AddProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *ProductGRPCHandler) GetProducts(ctx context.Context, req *productspb.GetProductsRequest) (*productspb.GetProductsResponse, error) {
	products, err := p.models.ProductModel.GetProducts(ctx, req)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *ProductGRPCHandler) GetProduct(ctx context.Context, req *productspb.GetProductRequest) (*productspb.GetProductResponse, error) {
	product, err := p.models.ProductModel.GetProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return product, nil
}
