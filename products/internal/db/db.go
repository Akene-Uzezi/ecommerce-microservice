// Package db for products service
package db

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	ProductModel ProductModel
}

func NewModels(db *pgxpool.Pool) *Models {
	return &ProductModel{DB: db}
}
