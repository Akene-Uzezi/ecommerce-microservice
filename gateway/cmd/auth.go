package main

import (
	"ecommerce-gateway/internal/handler"
	"log"
	"net/http"

	authpb "ecommerce-api/gen/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitAuthService(mux *http.ServeMux) {
	authServiceConn, err := grpc.NewClient(authServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to auth grpc service: %v", err)
	}
	defer authServiceConn.Close()
	log.Printf("dialing auth service at %s", authServiceURL)

	authClient := authpb.NewAuthServiceClient(authServiceConn)
	authHandler := handler.NewAuthHTTPHandler(authClient)
	authHandler.RegisterRoutes(mux)
}
