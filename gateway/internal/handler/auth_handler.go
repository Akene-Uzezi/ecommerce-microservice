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
	mux.HandleFunc("POST /api/v1/create_user", h.createUser)
	mux.HandleFunc("POST /api/v1/login", h.login)
}

func (h *AuthHTTPHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var requestbody CreateUserPayload
	if err := shared.ReadJSON(r, &requestbody); err != nil {
		shared.WriteErrorBadRequest(w, "Invalid requestbody", err)
		shared.LogBadRequest(r.Method, r.RequestURI)
		return
	}
	hashPassword, err := shared.HashPassword(requestbody.Password)
	if err != nil {
		shared.WriteErrorServerError(w, "Failed to hashpassword", err)
		shared.LogInternalServerError(r.Method, r.RequestURI)
		return
	}

	ctx := r.Context()
	res, err := h.authClient.CreateUser(ctx, &authpb.CreateUserRequest{
		Email:    requestbody.Email,
		Password: hashPassword,
		Name:     requestbody.Name,
	})
	if err != nil {
		log.Printf("grpc create user failed: %v", err)
		shared.WriteErrorServerError(w, "internal server error", err)
		shared.LogInternalServerError(r.Method, r.RequestURI)
		return
	}

	_ = shared.WriteJSON(w, http.StatusCreated, res)
	shared.LogOK(r.Method, r.RequestURI)
}

func (h *AuthHTTPHandler) login(w http.ResponseWriter, r *http.Request) {
	var requestbody LoginRequest
	if err := shared.ReadJSON(r, &requestbody); err != nil {
		shared.WriteErrorBadRequest(w, "Invalid requestbody", err)
		shared.LogBadRequest(r.Method, r.RequestURI)
		return
	}

	ctx := r.Context()
	res, err := h.authClient.GetUserByEmail(ctx, &authpb.GetUserRequest{
		Email: requestbody.Email,
	})
	if err != nil {
		log.Printf("grpc failed to get user by email: %v", err)
		shared.WriteErrorServerError(w, "internal server error", err)
		shared.LogInternalServerError(r.Method, r.RequestURI)
		return
	}

	_ = shared.WriteJSON(w, http.StatusOK, res)
	shared.LogOK(r.Method, r.RequestURI)
}
