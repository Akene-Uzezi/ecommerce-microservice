package db

import "github.com/jackc/pgx/v5/pgxpool"

type Models struct {
	OrderModel OrderModel
}

func NewModels(db *pgxpool.Pool) *Models {
	return &Models{OrderModel: OrderModel{DB: db}}
}
