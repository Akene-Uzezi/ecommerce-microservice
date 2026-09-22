package main

import (
	"ecommerce-gateway/internal/handler"
	"log"
	"net/http"

	authpb "ecommerce-api/gen/auth"
	stockpb "ecommerce-api/gen/stock"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initStockService(mux *http.ServeMux, authClient authpb.AuthServiceClient) *grpc.ClientConn {
	stockServiceConn, err := grpc.NewClient(stockServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to stock grpc service: %s", err)
	}
	log.Printf("dialing stock service at %s", stockServiceURL)
	stockClient := stockpb.NewStockServiceClient(stockServiceConn)
	stockHandler := handler.NewStockHTTPHandler(stockClient, authClient)
	stockHandler.RegisterRoutes(mux)
	return stockServiceConn
}
