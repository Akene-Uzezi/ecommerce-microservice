package handler

import (
	shared "ecommerce-shared"
	"errors"
	"net/http"
	"time"

	authpb "ecommerce-api/gen/auth"
)

func (h *AuthHTTPHandler) searchUsersByEmail(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	userQuery := r.URL.Query().Get("email")
	if userQuery == "" {
		shared.WriteErrorBadRequest(w, "Pass an email query parameter", errors.New("query parameter email was not passed"))
		shared.LogBadRequest(r.Method, r.RequestURI, time.Since(start))
		return
	}

	res, err := h.authClient.SearchUsersByEmail(r.Context(), &authpb.SearchUserByEmailRequest{
		Email: userQuery,
	})
	if err != nil {
		shared.WriteErrorServerError(w, "An error occured", err)
		shared.LogInternalServerError(r.Method, r.RequestURI, time.Since(start))
		return
	}

	response := map[string]any{
		"message": "User found",
		"user":    res,
	}

	_ = shared.WriteJSON(w, http.StatusOK, response)
	shared.LogOK(r.Method, r.RequestURI, time.Since(start))
}
