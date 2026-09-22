package main

import (
	"ecommerce-payments/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net"

	paymentpb "ecommerce-api/gen/payment"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/grpc"
)

var paymentsPort = shared.GetEnvString("PAYMENTS_PORT", "9000")

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", paymentsPort))
	if err != nil {
		log.Fatalf("error creating payments service listener: %s on port: %s", err, paymentsPort)
	}
	grpcServer := grpc.NewServer()
	paymentsHandler := handler.NewPaymentGRPCHandler()
	paymentpb.RegisterPaymentServiceServer(grpcServer, paymentsHandler)

	log.Printf("payments service running on: %s", paymentsPort)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve payments grpc: %v", err)
	}
}
