package main

import (
	"ecommerce-products/internal/db"
	"ecommerce-products/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	productspb "ecommerce-api/gen/products"

	"google.golang.org/grpc"
)

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
	productsHandler := handler.NewProductGRPCHandler(db.NewModels(pool))
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
