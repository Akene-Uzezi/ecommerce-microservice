package main

import (
	"ecommerce-products/internal/handler"
	shared "ecommerce-shared"
	"log"
	"net"

	productspb "ecommerce-api/gen/products"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
)

var productsPort = shared.GetEnvString("PRODUCTS_PORT", "7777")

func main() {
	l, err := net.Listen("tcp", productsPort)
	if err != nil {
		log.Fatalf("error creating product service listner: %s on port: %s", err, productsPort)
	}
	grpcServer := grpc.NewServer()
	// productsDBConnStr := shared.GetEnvString("PRODUCTS_DB_CONN_STR", "postgres://product:product@localhost:7433/products_db")
	// // pool, err := shared.InitPool(productsDBConnStr)
	// // if err != nil {
	// // 	log.Fatalf("failed to init producs db pool: %s", err)
	// }
	productsHandler := handler.NewProductGRPCHandler()
	productspb.RegisterProductServiceServer(grpcServer, productsHandler)

	log.Printf("products service running on: %s", productsPort)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("failed to serve products grpc: %v", err)
	}
}
