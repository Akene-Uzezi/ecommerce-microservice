// Package db
package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	DB *pgxpool.Pool
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Name     string `json:"name"`
}

func (m *UserModel) CreateUser(ctx context.Context, user *User) (*User, error) {
	var createduser User
	query := `
		INSERT INTO users 
		(email, password, name)
		VALUES ($1, $2, $3)
		RETURNING id, email, password, name
	`
	err := m.DB.QueryRow(ctx, query, user.Email, user.Password, user.Name).Scan(
		&createduser.ID, &createduser.Email, &createduser.Password, &createduser.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("database query error: %s", err)
	}
	return &createduser, nil
}
