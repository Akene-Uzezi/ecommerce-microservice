// Package handler
package handler

import (
	"context"
	"ecommerce-auth/internal/db"
	"ecommerce-auth/internal/util"
	shared "ecommerce-shared"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/joho/godotenv/autoload"

	authpb "ecommerce-api/gen/auth"

	"github.com/golang-jwt/jwt/v5"
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

func (h *AuthGRPCHanlder) VerifyToken(ctx context.Context, req *authpb.VerifyTokenRequest) (*authpb.VerifyTokenResponse, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(req.Token, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing error")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %s", err)
	}

	return &authpb.VerifyTokenResponse{
		Valid: token.Valid,
	}, nil
}
