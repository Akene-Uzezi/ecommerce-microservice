package handler

import (
	shared "ecommerce-shared"
	"log"
	"net/http"
	"time"

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
	start := time.Now()
	var requestbody CreateUserPayload
	if err := shared.ReadJSON(r, &requestbody); err != nil {
		shared.WriteErrorBadRequest(w, "Invalid requestbody", err)
		shared.LogBadRequest(r.Method, r.RequestURI, time.Since(start))
		return
	}

	ctx := r.Context()
	res, err := h.authClient.CreateUser(ctx, &authpb.CreateUserRequest{
		Email:    requestbody.Email,
		Password: requestbody.Password,
		Name:     requestbody.Name,
	})
	if err != nil {
		log.Printf("grpc create user failed: %v", err)
		shared.WriteErrorServerError(w, "internal server error", err)
		shared.LogInternalServerError(r.Method, r.RequestURI, time.Since(start))
		return
	}

	_ = shared.WriteJSON(w, http.StatusCreated, res)
	shared.LogOK(r.Method, r.RequestURI, time.Since(start))
}

func (h *AuthHTTPHandler) login(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var loginrequest authpb.LoginRequest
	if err := shared.ReadJSON(r, &loginrequest); err != nil {
		shared.WriteErrorBadRequest(w, "Invalid request", err)
		shared.LogBadRequest(r.Method, r.RequestURI, time.Since(start))
		return
	}

	ctx := r.Context()

	res, err := h.authClient.Login(ctx, &loginrequest)
	if err != nil {
		shared.WriteErrorServerError(w, "grpc failed to login", err)
		shared.LogInternalServerError(r.Method, r.RequestURI, time.Since(start))
		return
	}

	_ = shared.WriteJSON(w, http.StatusOK, res.Token)
	shared.LogOK(r.Method, r.RequestURI, time.Since(start))
}
