package main

import (
	"ecommerce-auth/internal/db"
	"ecommerce-auth/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net"

	authpb "ecommerce-api/gen/auth"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/grpc"
)

var authPort = shared.GetEnvString("AUTH_PORT", "5555")

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", authPort))
	if err != nil {
		log.Fatalf("Failed to listen for auth service on port:%s, %v", authPort, err)
	}

	grpcServer := grpc.NewServer()
	authDBConnStr := shared.GetEnvString("AUTH_DB_CONN_STR", "postgres://auth:auth@localhost:6433/auth_db")

	pool, err := shared.InitPool(authDBConnStr)
	if err != nil {
		log.Fatalf("failed to init auth db pool: %s", err)
	}
	authHandler := handler.NewAuthGRPCHandler(db.NewModels(pool))
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Printf("auth service running on port:%s", authPort)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
