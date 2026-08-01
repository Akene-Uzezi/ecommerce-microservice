// Package handler
package handler

import (
	"context"
	"log"

	authpb "ecommerce-api/gen/auth"
)

type AuthGRPCHanlder struct {
	authpb.UnimplementedAuthServiceServer
}

func NewAuthGRPCHandler() *AuthGRPCHanlder {
	return &AuthGRPCHanlder{}
}

func (h *AuthGRPCHanlder) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
	log.Println("received create user request")

	return &authpb.CreateUserResponse{
		Email: "test@test.com",
		Name:  "test name",
	}, nil
}
