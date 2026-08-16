package main

import (
	"ecommerce-gateway/internal/handler"
	"log"
	"net/http"

	productspb "ecommerce-api/gen/products"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initProductService(mux *http.ServeMux) *grpc.ClientConn {
	productsServiceConn, err := grpc.NewClient(productsServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect with products grpc service: %s", err)
	}
	log.Printf("dialing products service at %s", productsServiceURL)
	productsClient := productspb.NewProductServiceClient(productsServiceConn)
	productsHandler := handler.NewProductsHTTPHandler(productsClient)
	productsHandler.RegisterRoutes(mux)
	return productsServiceConn
}
