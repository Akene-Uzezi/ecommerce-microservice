package main

import (
	"ecommerce-gateway/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net/http"

	pb "ecommerce-api/gen/order"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	gatewayPort     = shared.GetEnvString("GATEWAY_PORT", "3000")
	orderServiceURL = shared.GetEnvString("ORDER_SERVICE_URL", "localhost:4444")
)

func main() {
	conn, err := grpc.NewClient(
		orderServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to orders grpc service: %v", err)
	}
	defer conn.Close()

	orderClient := pb.NewOrderServiceClient(conn)
	addr := fmt.Sprintf(":%s", gatewayPort)
	mux := http.NewServeMux()
	handler := handler.NewHTTPHandler(orderClient)
	handler.RegisterRoutes(mux)

	log.Printf("Server running on port%v", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
