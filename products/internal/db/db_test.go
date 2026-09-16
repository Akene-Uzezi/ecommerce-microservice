package db

import (
	shared "ecommerce-shared"
	"log"
	"os"
	"testing"
)

var productModel *ProductModel

func TestMain(m *testing.M) {
	pool, cleanup, err := shared.SetupTestDBSuite("/scripts/products_init.sql")
	if err != nil {
		log.Fatalf("failed to init test db %v", err)
	}
	productModel = NewProductModel(pool)

	exitcode := m.Run()
	cleanup()
	os.Exit(exitcode)
}
