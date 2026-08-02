// Package handler
package handler

import (
	"context"
	"ecommerce-auth/internal/db"

	authpb "ecommerce-api/gen/auth"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthGRPCHanlder struct {
	authpb.UnimplementedAuthServiceServer
	pool *pgxpool.Pool
}

func NewAuthGRPCHandler() *AuthGRPCHanlder {
	return &AuthGRPCHanlder{
		pool: db.InitPool(),
	}
}

func (h *AuthGRPCHanlder) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
}
