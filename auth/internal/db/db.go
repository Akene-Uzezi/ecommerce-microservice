package db

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
	UserModel UserModel
}

func NewModels(pool *pgxpool.Pool) *Models {
	return &Models{
		UserModel: UserModel{DB: pool},
	}
}
