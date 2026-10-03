// Package handler for the products service
package handler

import (
	"context"
	"ecommerce-products/internal/db"
	"time"

	productspb "ecommerce-api/gen/products"
	shared "ecommerce-shared"
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
	start := time.Now()
	product, err := p.models.ProductModel.AddProduct(ctx, req)
	if err != nil {
		shared.LogRequest("AddProduct", "/product.AddProduct", 500, time.Since(start))
		return nil, err
	}
	shared.LogRequest("AddProduct", "/product.AddProduct", 0, time.Since(start))
	return product, nil
}

func (p *ProductGRPCHandler) GetProducts(ctx context.Context, req *productspb.GetProductsRequest) (*productspb.GetProductsResponse, error) {
	start := time.Now()
	products, err := p.models.ProductModel.GetProducts(ctx, req)
	if err != nil {
		shared.LogRequest("GetProducts", "/product.GetProducts", 500, time.Since(start))
		return nil, err
	}
	shared.LogRequest("GetProducts", "/product.GetProducts", 0, time.Since(start))
	return products, nil
}

func (p *ProductGRPCHandler) GetProduct(ctx context.Context, req *productspb.GetProductRequest) (*productspb.GetProductResponse, error) {
	start := time.Now()
	product, err := p.models.ProductModel.GetProduct(ctx, req)
	if err != nil {
		shared.LogRequest("GetProduct", "/product.GetProduct", 404, time.Since(start))
		return nil, err
	}
	shared.LogRequest("GetProduct", "/product.GetProduct", 0, time.Since(start))
	return product, nil
}
