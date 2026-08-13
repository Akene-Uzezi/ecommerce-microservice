package handler

import (
	shared "ecommerce-shared"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	authpb "ecommerce-api/gen/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// orderServiceURL = shared.GetEnvString("ORDER_SERVICE_URL", "localhost:4444")

var (
	gatewayPort    = shared.GetEnvString("GATEWAY_PORT", "3000")
	authServiceURL = shared.GetEnvString("AUTH_SERVICE_URL", "localhost:5555")
	mux            *http.ServeMux
	authClient     authpb.AuthServiceClient
	authHandler    *AuthHTTPHandler
)

func TestMain(m *testing.M) {
	authConn, err := grpc.NewClient(
		authServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Error creating auth client: %s", err)
	}
	addr := fmt.Sprintf(":%s", gatewayPort)
	mux = http.NewServeMux()
	authClient = authpb.NewAuthServiceClient(authConn)
	authHandler = NewAuthHTTPHandler(authClient)
	authHandler.RegisterRoutes(mux)
	go func() {
		log.Printf("server running on port %s", addr)
		if err := http.ListenAndServe(addr, mux); err != nil {
			log.Fatalf("unable to start server: %s", err)
		}
	}()
	exitCode := m.Run()
	os.Exit(exitCode)
}
