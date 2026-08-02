// Package db
package db

import "github.com/jackc/pgx/v5/pgxpool"

type UserModel struct {
	DB *pgxpool.Pool
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
