// Package handler
package handler

import (
	"context"
	"ecommerce-auth/internal/db"
	"log"

	authpb "ecommerce-api/gen/auth"
)

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

func (h *AuthGRPCHanlder) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
	return nil, nil
}
