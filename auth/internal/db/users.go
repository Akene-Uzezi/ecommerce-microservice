// Package db
package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	DB *pgxpool.Pool
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"-"`
	Name     string `json:"name"`
}

func (m *UserModel) CreateUser(ctx context.Context, user User) (*User, error) {
	return nil, nil
}
