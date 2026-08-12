package db

import (
	shared "ecommerce-shared"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testPool  *pgxpool.Pool
	userModel *UserModel
)

func TestMain(m *testing.M) {
	pool, cleanup, err := shared.SetupTestDBSuite("/scripts/auth_init.sql")
	if err != nil {
		log.Fatalf("failed to initialize test container: %v", err)
	}
	testPool = pool
	userModel = newUserModel(testPool)
	exitCode := m.Run()

	cleanup()

	os.Exit(exitCode)
}
