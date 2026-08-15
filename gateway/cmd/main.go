package main

import (
	"ecommerce-gateway/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net/http"

	orderpb "ecommerce-api/gen/order"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	gatewayPort     = shared.GetEnvString("GATEWAY_PORT", "3000")
	orderServiceURL = shared.GetEnvString("ORDER_SERVICE_URL", "localhost:4444")
	authServiceURL  = shared.GetEnvString("AUTH_SERVICE_URL", "localhost:5555")
)

func main() {
	orderServiceConn, err := grpc.NewClient(
		orderServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to orders grpc service: %v", err)
	}
	defer orderServiceConn.Close()
	log.Printf("dialing orderservice at %s", orderServiceURL)

	orderClient := orderpb.NewOrderServiceClient(orderServiceConn)
	addr := fmt.Sprintf(":%s", gatewayPort)
	mux := http.NewServeMux()
	authClient, authServiceConn := InitAuthService(mux)
	defer authServiceConn.Close()
	orderHandler := handler.NewOrderHTTPHandler(orderClient, authClient)
	orderHandler.RegisterRoutes(mux)

	log.Printf("Server running on port%v", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
