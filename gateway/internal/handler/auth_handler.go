package handler

import (
	shared "ecommerce-shared"
	"log"
	"net/http"

	authpb "ecommerce-api/gen/auth"
)

type AuthHTTPHandler struct {
	authClient authpb.AuthServiceClient
}

func NewAuthHTTPHandler(client authpb.AuthServiceClient) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authClient: client,
	}
}

func (h *AuthHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/create/ping", h.helperPing)
}

func (h *AuthHTTPHandler) helperPing(w http.ResponseWriter, r *http.Request) {
	log.Println("received request from gateway")
	_ = shared.WriteJSON(w, http.StatusOK, "test")
}
