package main

import (
	shared "ecommerce-shared"
	"ecommerce-stock/internal/handler"
	"fmt"
	"log"
	"net"

	stockpb "ecommerce-api/gen/stock"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

var stockPort = shared.GetEnvString("STOCK_PORT", "8888")

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", stockPort))
	if err != nil {
		log.Fatalf("error creating stock service listner: %s", err)
	}

	grpcServer := grpc.NewServer()
	stockHandler := handler.NewStockGRPCHandler()
	stockpb.RegisterStockServiceServer(grpcServer, stockHandler)

	log.Printf("stock service running on port %s", stockPort)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to server stock service: %v", err)
	}
}
