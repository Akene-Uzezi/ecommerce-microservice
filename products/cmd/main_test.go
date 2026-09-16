package main

import (
	"context"
	"ecommerce-products/internal/db"
	"ecommerce-products/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	productspb "ecommerce-api/gen/products"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

var productsHandler *handler.ProductGRPCHandler

func TestMain(m *testing.M) {
	pool, cleanup, err := shared.SetupTestDBSuite("/scripts/products_init.sql")
	if err != nil {
		log.Fatalf("failed to init db: %s", err)
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", productsPort))
	if err != nil {
		log.Fatalf("failed to listen on products port: %s", err)
	}

	grpcServer := grpc.NewServer()
	productsHandler = handler.NewProductGRPCHandler(db.NewModels(pool))
	productspb.RegisterProductServiceServer(grpcServer, productsHandler)

	go func() {
		log.Printf("products service started on port %s", productsPort)
		if err := grpcServer.Serve(l); err != nil {
			log.Printf("grpc server stopped, %s", err)
		}
	}()

	exitcode := m.Run()
	grpcServer.GracefulStop()
	cleanup()
	os.Exit(exitcode)
}

func TestAddProductInt(t *testing.T) {
	ctx := context.Background()
	req := &productspb.AddProductRequest{
		Product: &productspb.Product{
			Name:     "testproduct",
			Price:    10.32,
			Quantity: 4,
		},
	}
	_, err := productsHandler.AddProduct(ctx, req)
	assert.NoError(t, err)
}

func TestGetProductInt(t *testing.T) {
	ctx := context.Background()
	req := &productspb.GetProductRequest{
		Name: "testproduct",
	}
	_, err := productsHandler.GetProduct(ctx, req)
	assert.NoError(t, err)
}

func TestGetProductsInt(t *testing.T) {
	ctx := context.Background()
	req := &productspb.GetProductsRequest{}
	_, err := productsHandler.GetProducts(ctx, req)
	assert.NoError(t, err)
}
