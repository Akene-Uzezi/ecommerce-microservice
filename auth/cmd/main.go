package main

import (
	"ecommerce-auth/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net"

	authpb "ecommerce-api/gen/auth"

	"google.golang.org/grpc"
)

var authPort = shared.GetEnvString("AUTH_PORT", "5555")

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", authPort))
	if err != nil {
		log.Fatalf("Failed to listen for auth service on port:%s, %v", authPort, err)
	}

	grpcServer := grpc.NewServer()
	authHandler := handler.NewAuthGRPCHandler()
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Printf("auth service running on port:%s", authPort)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
