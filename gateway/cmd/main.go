package main

import (
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net/http"

	authpb "ecommerce-api/gen/auth"

	_ "github.com/joho/godotenv/autoload"
)

var (
	gatewayPort        = shared.GetEnvString("GATEWAY_PORT", "3000")
	orderServiceURL    = shared.GetEnvString("ORDER_SERVICE_URL", "localhost:4444")
	authServiceURL     = shared.GetEnvString("AUTH_SERVICE_URL", "localhost:5555")
	productsServiceURL = shared.GetEnvString("PRODUCTS_SERVICE_URL", "localhost:7777")
	paymentsServiceURL = shared.GetEnvString("PAYMENTS_SERVICE_URL", "localhost:9000")
	stockServiceURL    = shared.GetEnvString("STOCK_SERVICE_URL", "localhost:8888")
	authClient         authpb.AuthServiceClient
)

func main() {
	addr := fmt.Sprintf(":%s", gatewayPort)
	mux := http.NewServeMux()
	authClient, authServiceConn := initAuthService(mux)
	defer authServiceConn.Close()
	_, ordersServiceConn := initOrdersService(mux, authClient)
	defer ordersServiceConn.Close()
	productsServiceConn := initProductService(mux, authClient)
	defer productsServiceConn.Close()
	paymentsServiceConn := initPaymentService(mux, authClient)
	defer paymentsServiceConn.Close()
	stockServiceConn := initStockService(mux, authClient)
	defer stockServiceConn.Close()
	log.Printf("Server running on port %s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
