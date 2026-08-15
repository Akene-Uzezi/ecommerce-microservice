package main

import (
	"ecommerce-gateway/internal/handler"
	"log"
	"net/http"

	authpb "ecommerce-api/gen/auth"
	orderpb "ecommerce-api/gen/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitOrdersService(mux *http.ServeMux, authClient authpb.AuthServiceClient) (orderpb.OrderServiceClient, *grpc.ClientConn) {
	orderServiceConn, err := grpc.NewClient(
		orderServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to connect to orders grpc service: %v", err)
	}
	log.Printf("dialing orderservice at %s", orderServiceURL)
	orderClient := orderpb.NewOrderServiceClient(orderServiceConn)
	orderHandler := handler.NewOrderHTTPHandler(orderClient, authClient)
	orderHandler.RegisterRoutes(mux)
	return orderClient, orderServiceConn
}
