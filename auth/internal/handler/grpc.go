// Package handler
package handler

import (
	"context"
	"ecommerce-auth/internal/db"
	"log"

	authpb "ecommerce-api/gen/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthGRPCHanlder struct {
	authpb.UnimplementedAuthServiceServer
	pool *pgxpool.Pool
}

func NewAuthGRPCHandler() *AuthGRPCHanlder {
	pool, err := db.InitPool()
	if err != nil {
		log.Fatalf("failed to init db pool: %s", err)
	}
	return &AuthGRPCHanlder{
		pool: pool,
	}
}

func (h *AuthGRPCHanlder) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
	return nil, nil
}
