package main

import (
	"ecommerce-gateway/internal/handler"
	"log"
	"net/http"
)

func main() {
	addr := ":3000"
	mux := http.NewServeMux()
	handler := handler.NewHTTPHandler()
	handler.RegisterRoutes(mux)

	log.Printf("Server running on port%v", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
