package handler

import (
	"context"
	"ecommerce-auth/internal/db"
	"ecommerce-auth/internal/util"
	"errors"
	"fmt"
	"time"

	authpb "ecommerce-api/gen/auth"

	"github.com/golang-jwt/jwt/v5"
)

func (h *AuthGRPCHanlder) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
	hashPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %s", err)
	}
	user := &db.User{
		Email:    req.Email,
		Password: hashPassword,
		Name:     req.Name,
	}
	user, err = h.models.UserModel.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	response := &authpb.CreateUserResponse{
		Email: user.Email,
		Name:  user.Name,
	}
	return response, nil
}

func (h *AuthGRPCHanlder) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	user := &db.User{
		Email:    req.Email,
		Password: req.Password,
	}
	founduser, err := h.models.UserModel.GetUserByEmail(ctx, user)
	if err != nil {
		return nil, err
	}
	matchPassword := util.ComparePassword(founduser.Password, req.Password)
	if !matchPassword {
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
		return nil, fmt.Errorf("failed to generate jwt token: %s", err)
	}

	return &authpb.LoginResponse{
		Token: tokenString,
	}, nil
}
