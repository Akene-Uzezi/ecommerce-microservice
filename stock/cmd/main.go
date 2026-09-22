package main

import (
	"fmt"
	"log"
	"net"

	shared "ecommerce-shared"
	"ecommerce-stock/internal/db"
	"ecommerce-stock/internal/handler"

	stockpb "ecommerce-api/gen/stock"

	_ "github.com/joho/godotenv/autoload"

	"google.golang.org/grpc"
)

var stockPort = shared.GetEnvString("STOCK_PORT", "8888")

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", stockPort))
	if err != nil {
		log.Fatalf("error creating stock service listener: %s on port: %s", err, stockPort)
	}
	grpcServer := grpc.NewServer()
	stockDBConnStr := shared.GetEnvString("STOCK_DB_CONN_STR", "postgres://stock:stock@localhost:8433/stock_db")
	pool, err := shared.InitPool(stockDBConnStr)
	if err != nil {
		log.Fatalf("failed to init stock db pool: %s", err)
	}
	models := db.NewModels(pool)
	stockHandler := handler.NewStockGRPCHandler(models)
	stockpb.RegisterStockServiceServer(grpcServer, stockHandler)

	log.Printf("stock service running on: %s", stockPort)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve stock grpc: %v", err)
	}
}
