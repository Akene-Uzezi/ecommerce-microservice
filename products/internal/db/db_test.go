package db

import (
	shared "ecommerce-shared"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testPool     *pgxpool.Pool
	productModel *ProductModel
)

func TestMain(m *testing.M) {
	pool, cleanup, err := shared.SetupTestDBSuite("/scripts/products_init.sql")
	if err != nil {
		log.Fatal("failed to init test db %v", err)
	}
	testPool = pool
	productModel = NewProductModel(testPool)

	exitcode := m.Run()
	cleanup()
	os.Exit(exitcode)
}
