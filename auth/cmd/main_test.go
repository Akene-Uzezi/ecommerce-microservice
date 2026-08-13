package main

import (
	"ecommerce-auth/internal/db"
	"ecommerce-auth/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	authpb "ecommerce-api/gen/auth"

	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	pool, cleanup, err := shared.SetupTestDBSuite("/scripts/auth_init.sql")
	if err != nil {
		log.Fatalf("failed to init db: %s", err)
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", authPort))
	if err != nil {
		log.Fatalf("failed to listen on auth port: %s", err)
	}

	grpcServer := grpc.NewServer()
	authHandler := handler.NewAuthGRPCHandler(db.NewModels(pool))
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)

	go func() {
		log.Printf("auth service started %s", authPort)
		if err := grpcServer.Serve(l); err != nil {
			log.Printf("grpc server stopped: %v", err)
		}
	}()

	exitCode := m.Run()

	grpcServer.GracefulStop()
	cleanup()
	os.Exit(exitCode)
}
