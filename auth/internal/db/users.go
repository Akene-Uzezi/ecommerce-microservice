// Package db
package db

import "github.com/jackc/pgx/v5/pgxpool"

type UserModel struct {
	DB *pgxpool.Pool
}

type User struct {
	email    string `json:"email"`
	password string `json:"password"`
	name     string `json:"name"`
}
