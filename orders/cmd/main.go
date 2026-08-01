package main

import (
	"ecommerce-orders/internal/handler"
	"log"
	"net"

	pb "ecommerce-api/gen/order"

	"google.golang.org/grpc"
)

func main() {
	l, err := net.Listen("tcp", ":4444")
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
