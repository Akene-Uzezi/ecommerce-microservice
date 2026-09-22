package main

import (
	"ecommerce-orders/internal/db"
	"ecommerce-orders/internal/handler"
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net"

	orderpb "ecommerce-api/gen/order"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/grpc"
)

var ordersPort = shared.GetEnvString("ORDERS_PORT", "4444")

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", ordersPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", ordersPort, err)
	}
	grpcServer := grpc.NewServer()
	ordersDBConnStr := shared.GetEnvString("ORDERS_DB_CONN_STR", "postgres://orders:orders@orders-db:5432/orders_db")
	pool, err := shared.InitPool(ordersDBConnStr)
	if err != nil {
		log.Fatalf("failed to init orders db pool: %s", err)
	}
	models := db.NewModels(pool)
	orderHandler := handler.NewOrderGRPCHandler(models)

	orderpb.RegisterOrderServiceServer(grpcServer, orderHandler)

	log.Printf("Orders grpc server running on port: %s", ordersPort)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}
