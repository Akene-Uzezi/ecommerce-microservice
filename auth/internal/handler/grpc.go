// Package handler
package handler

import (
	"ecommerce-auth/internal/db"
	shared "ecommerce-shared"

	_ "github.com/joho/godotenv/autoload"

	authpb "ecommerce-api/gen/auth"
)

var jwtSecret = shared.GetEnvString("jwt_secret", "secret")

type AuthGRPCHandler struct {
	authpb.UnimplementedAuthServiceServer
	models *db.Models
}

func NewAuthGRPCHandler(models *db.Models) *AuthGRPCHandler {
	return &AuthGRPCHandler{
		models: models,
	}
}
