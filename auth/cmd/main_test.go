package main

import (
	"ecommerce-auth/internal/db"
	shared "ecommerce-shared"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testPool  *pgxpool.Pool
	UserModel *db.UserModel
)

func TestMain(m *testing.M) {
	pool, cleanup, err := shared.SetupTestDBSuite("/scripts/auth_init.sql")
	if err != nil {
		log.Fatalf("failed to initialize test container: %v", err)
	}
	testPool = pool
	UserModel = db.NewUserModel(testPool)
	exitCode := m.Run()

	cleanup()

	os.Exit(exitCode)
}
