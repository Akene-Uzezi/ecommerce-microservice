// Package handler
package handler

import (
	"ecommerce-auth/internal/db"
	shared "ecommerce-shared"
	"log"

	_ "github.com/joho/godotenv/autoload"

	authpb "ecommerce-api/gen/auth"
)

var jwtSecret = shared.GetEnvString("jwt_secret", "secret")

type AuthGRPCHanlder struct {
	authpb.UnimplementedAuthServiceServer
	models *db.Models
}

func NewAuthGRPCHandler() *AuthGRPCHanlder {
	pool, err := db.InitPool()
	if err != nil {
		log.Fatalf("failed to init db pool: %s", err)
	}
	return &AuthGRPCHanlder{
		models: db.NewModels(pool),
	}
}
