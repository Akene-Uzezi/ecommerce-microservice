package db

import "github.com/jackc/pgx/v5/pgxpool"

type ProductModel struct {
	DB *pgxpool.Pool
}
