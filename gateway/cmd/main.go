package main

import (
	"ecommerce-gateway/internal/handler"
	"log"
	"net/http"

	pb "ecommerce-api/gen/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:4444",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to orders grpc service: %v", err)
	}
	defer conn.Close()

	orderClient := pb.NewOrderServiceClient(conn)
	addr := ":3000"
	mux := http.NewServeMux()
	handler := handler.NewHTTPHandler(orderClient)
	handler.RegisterRoutes(mux)

	log.Printf("Server running on port%v", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
