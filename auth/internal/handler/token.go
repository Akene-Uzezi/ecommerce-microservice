package handler

import (
	"context"
	"errors"
	"fmt"

	authpb "ecommerce-api/gen/auth"

	"github.com/golang-jwt/jwt/v5"
)

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
