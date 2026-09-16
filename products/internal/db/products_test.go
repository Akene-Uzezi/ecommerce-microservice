package db

import (
	"context"
	"testing"

	productspb "ecommerce-api/gen/products"

	"github.com/stretchr/testify/assert"
)

func TestFunction(t *testing.T) {
	t.Log("test func works")
}

func TestAddProduct(t *testing.T) {
	ctx := context.Background()
	req := &productspb.AddProductRequest{
		Product: &productspb.Product{
			Name:     "testproduct",
			Price:    10.32,
			Quantity: 4,
		},
	}
	_, err := productModel.AddProduct(ctx, req)
	assert.NoError(t, err)
}

func TestGetProduct(t *testing.T) {
	ctx := context.Background()
	req := &productspb.GetProductRequest{
		Name: "testproduct",
	}
	_, err := productModel.GetProduct(ctx, req)
	assert.NoError(t, err)
}

func TestGetProducts(t *testing.T) {
	ctx := context.Background()
	req := &productspb.GetProductsRequest{}
	_, err := productModel.GetProducts(ctx, req)
	assert.NoError(t, err)
}
