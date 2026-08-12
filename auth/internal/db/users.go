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

func newUserModel(db *pgxpool.Pool) *UserModel {
	return &UserModel{
		DB: db,
	}
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

func (m *UserModel) GetUserByEmail(ctx context.Context, user *User) (*User, error) {
	var founduser User
	err := m.DB.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", user.Email).Scan(
		&founduser.ID,
		&founduser.Email,
		&founduser.Password,
		&founduser.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("database query error: %s", err)
	}
	return &founduser, nil
}
