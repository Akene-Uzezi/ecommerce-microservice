package main

import (
	"ecommerce-auth/internal/handler"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	authpb "ecommerce-api/gen/auth"

	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", authPort))
	if err != nil {
		log.Fatalf("failed to listen on auth port: %s", err)
	}

	grpcServer := grpc.NewServer()
	authHandler := handler.NewAuthGRPCHandler()
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)

	go func() {
		log.Printf("auth service started %s", authPort)
		if err := grpcServer.Serve(l); err != nil {
			log.Printf("grpc server stopped: %v", err)
		}
	}()

	exitCode := m.Run()

	grpcServer.GracefulStop()
	os.Exit(exitCode)
}
