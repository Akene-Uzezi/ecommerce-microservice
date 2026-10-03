package handler

import (
	"context"
	"ecommerce-auth/internal/db"
	"ecommerce-auth/internal/util"
	"errors"
	"fmt"
	"time"

	authpb "ecommerce-api/gen/auth"
	shared "ecommerce-shared"

	"github.com/golang-jwt/jwt/v5"
)

func (h *AuthGRPCHandler) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
	start := time.Now()
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		shared.LogRequest("CreateUser", "/auth.CreateUser", 500, time.Since(start))
		return nil, fmt.Errorf("error hashing password: %s", err)
	}
	user := &db.User{
		Email:    req.Email,
		Password: hashPassword,
		Name:     req.Name,
	}
	user, err = h.models.UserModel.CreateUser(ctx, user)
	if err != nil {
		shared.LogRequest("CreateUser", "/auth.CreateUser", 500, time.Since(start))
		return nil, err
	}
	shared.LogRequest("CreateUser", "/auth.CreateUser", 0, time.Since(start))
	response := &authpb.CreateUserResponse{
		Email: user.Email,
		Name:  user.Name,
	}
	return response, nil
}

func (h *AuthGRPCHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	start := time.Now()
	user := &db.User{
		Email:    req.Email,
		Password: req.Password,
	}
	founduser, err := h.models.UserModel.GetUserByEmail(ctx, user)
	if err != nil {
		shared.LogRequest("Login", "/auth.Login", 401, time.Since(start))
		return nil, err
	}
	matchPassword := util.ComparePassword(founduser.Password, req.Password)
	if !matchPassword {
		shared.LogRequest("Login", "/auth.Login", 401, time.Since(start))
		return nil, errors.New("invalid credentials")
	}

	claims := &Claims{
		Email: founduser.Email,
		Name:  founduser.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		shared.LogRequest("Login", "/auth.Login", 500, time.Since(start))
		return nil, fmt.Errorf("failed to generate jwt token: %s", err)
	}

	shared.LogRequest("Login", "/auth.Login", 0, time.Since(start))
	return &authpb.LoginResponse{
		Token: tokenString,
	}, nil
}

func (h *AuthGRPCHandler) SearchUsersByEmail(ctx context.Context, req *authpb.SearchUserByEmailRequest) (*authpb.SearchUserByEmailResponse, error) {
	start := time.Now()
	user := &db.User{
		Email: req.Email,
	}

	foundUser, err := h.models.UserModel.GetUserByEmail(ctx, user)
	if err != nil {
		shared.LogRequest("SearchUsersByEmail", "/auth.SearchUsersByEmail", 404, time.Since(start))
		return nil, fmt.Errorf("user not found: %s", err)
	}

	shared.LogRequest("SearchUsersByEmail", "/auth.SearchUsersByEmail", 0, time.Since(start))
	return &authpb.SearchUserByEmailResponse{
		Email: foundUser.Email,
		Name:  foundUser.Name,
	}, nil
}
