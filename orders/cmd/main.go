package main

import (
	"ecommerce-orders/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net"

	pb "ecommerce-api/gen/order"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/grpc"
)

var ordersPort = shared.GetEnvString("ORDERS_PORT", "4444")

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", ordersPort))
	if err != nil {
		log.Fatalf("Failed to listen on port 4444: %v", err)
	}
	grpcServer := grpc.NewServer()
	orderHandler := handler.NewOrderGRPCHandler()

	pb.RegisterOrderServiceServer(grpcServer, orderHandler)

	log.Println("Orders grpc server running on port: 4444")
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
