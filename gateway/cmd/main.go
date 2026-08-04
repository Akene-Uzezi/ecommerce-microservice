package main

import (
	"ecommerce-gateway/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net/http"

	authpb "ecommerce-api/gen/auth"
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
	authServiceConn, err := grpc.NewClient(authServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to auth grpc service: %v", err)
	}
	log.Printf("dialing auth service at %s", authServiceURL)

	defer authServiceConn.Close()

	orderClient := orderpb.NewOrderServiceClient(orderServiceConn)
	authClient := authpb.NewAuthServiceClient(authServiceConn)
	addr := fmt.Sprintf(":%s", gatewayPort)
	mux := http.NewServeMux()
	orderHandler := handler.NewOrderHTTPHandler(orderClient, authClient)
	authHandler := handler.NewAuthHTTPHandler(authClient)
	orderHandler.RegisterRoutes(mux)
	authHandler.RegisterRoutes(mux)

	log.Printf("Server running on port%v", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
