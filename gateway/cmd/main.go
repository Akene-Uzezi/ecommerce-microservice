package main

import (
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
)

var (
	gatewayPort     = shared.GetEnvString("GATEWAY_PORT", "3000")
	orderServiceURL = shared.GetEnvString("ORDER_SERVICE_URL", "localhost:4444")
	authServiceURL  = shared.GetEnvString("AUTH_SERVICE_URL", "localhost:5555")
)

func main() {
	addr := fmt.Sprintf(":%s", gatewayPort)
	mux := http.NewServeMux()
	authClient, authServiceConn := initAuthService(mux)
	defer authServiceConn.Close()
	_, ordersServiceConn := initOrdersService(mux, authClient)
	defer ordersServiceConn.Close()
	log.Printf("Server running on port%v", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
