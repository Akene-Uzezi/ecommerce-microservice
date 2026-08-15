package shared

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func SetupTestDBSuite(sql string) (*pgxpool.Pool, func(), error) {
	ctx := context.Background()
	pgContainer, err := postgres.Run(
		ctx, "postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		pgContainer.Terminate(ctx)
	}

	connStr, _ := pgContainer.ConnectionString(ctx, "sslmode=disable")
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	var PingErr error
	for range 10 {
		PingErr = pool.Ping(ctx)
		if PingErr == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if PingErr != nil {
		cleanup()
		return nil, nil, PingErr
	}

	dir, err := os.Getwd()
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			cleanup()
			return nil, nil, os.ErrNotExist
		}
		dir = parent
	}

	schemaPath := filepath.Join(dir, fmt.Sprint(sql))
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	_, err = pool.Exec(ctx, string(schema))
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	return pool, cleanup, nil
}

func InitPool(connStr string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	config, err := pgxpool.ParseConfig(connStr)
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
