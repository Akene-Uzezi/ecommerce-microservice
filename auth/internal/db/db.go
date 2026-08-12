package db

import (
	"context"
	shared "ecommerce-shared"
	"fmt"
	"time"

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

func InitPool() (*pgxpool.Pool, error) {
	authDBConnStr := shared.GetEnvString("AUTH_DB_CONN_STR", "postgres://auth:auth@localhost:6433/auth_db")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	config, err := pgxpool.ParseConfig(authDBConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conn str: %s", err)
	}
	config.MaxConns = 10
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool %s", err)
	}

	var pingErr error
	count := 0
	for i := 0; i < 10; i++ {
		pingErr = pool.Ping(ctx)
		if pingErr == nil {
			fmt.Println("DB connection established")
			break
		}
		count += 1
		fmt.Println("failed to connect to database retrying")
		time.Sleep(700 * time.Millisecond)
		if count == 10 {
			return nil, pingErr
		}
	}

	return pool, nil
}
