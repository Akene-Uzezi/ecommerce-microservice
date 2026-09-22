package main

import (
	"ecommerce-gateway/internal/handler"
	"log"
	"net/http"

	authpb "ecommerce-api/gen/auth"
	paymentpb "ecommerce-api/gen/payment"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func initPaymentService(mux *http.ServeMux, authClient authpb.AuthServiceClient) *grpc.ClientConn {
	paymentsServiceConn, err := grpc.NewClient(paymentsServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to payments grpc service: %s", err)
	}
	log.Printf("dialing payments service at %s", paymentsServiceURL)
	paymentsClient := paymentpb.NewPaymentServiceClient(paymentsServiceConn)
	paymentsHandler := handler.NewPaymentsHTTPHandler(paymentsClient, authClient)
	paymentsHandler.RegisterRoutes(mux)
	return paymentsServiceConn
}
