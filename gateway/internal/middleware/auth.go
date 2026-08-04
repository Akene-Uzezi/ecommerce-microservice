// Package middleware
// NOTE: this package is the middleware for verifying a token
package middleware

import (
	shared "ecommerce-shared"
	"errors"
	"net/http"
	"strings"
	"time"

	authpb "ecommerce-api/gen/auth"
)

func RequireAuth(authClient authpb.AuthServiceClient) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			authHeader := r.Header.Get("Authorization")
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				shared.WriteErrorUnauthorized(w, "missing or malformed token", errors.New("invalid token format"))
				shared.LogUnauthorized(r.Method, r.RequestURI, time.Since(start))
				return
			}

			res, err := authClient.VerifyToken(r.Context(), &authpb.VerifyTokenRequest{Token: parts[1]})
			if err != nil || !res.Valid {
				shared.WriteErrorUnauthorized(w, "invalid or expired token", err)
				shared.LogUnauthorized(r.Method, r.RequestURI, time.Since(start))
				return
			}

			next(w, r)
		}
	}
}
