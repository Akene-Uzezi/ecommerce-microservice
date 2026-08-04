package handler

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Email string
	Name  string
	jwt.RegisteredClaims
}
