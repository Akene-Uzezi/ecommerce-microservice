package handler

import (
	shared "ecommerce-shared"
	"log"
	"net/http"
	"testing"

	authpb "ecommerce-api/gen/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// gatewayPort     = shared.GetEnvString("GATEWAY_PORT", "3000")
// orderServiceURL = shared.GetEnvString("ORDER_SERVICE_URL", "localhost:4444")
var (
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
	mux = http.NewServeMux()
	authClient = authpb.NewAuthServiceClient(authConn)
	authHandler = NewAuthHTTPHandler(authClient)
	authHandler.RegisterRoutes(mux)
}
