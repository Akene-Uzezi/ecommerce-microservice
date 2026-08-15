// Package db for products service
package db

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	ProductModel ProductModel
}

func NewProductModels(db *pgxpool.Pool) *ProductModel {
	return &ProductModel{DB: db}
}
