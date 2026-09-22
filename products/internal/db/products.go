package db

import (
	"context"
	"fmt"

	productspb "ecommerce-api/gen/products"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductModel struct {
	DB *pgxpool.Pool
}

type Product struct {
	ID       int
	Name     string
	Price    float64
	Quantity int
}

func NewProductModel(db *pgxpool.Pool) *ProductModel {
	return &ProductModel{DB: db}
}

func (m *ProductModel) AddProduct(ctx context.Context, req *productspb.AddProductRequest) (*productspb.AddProductResponse, error) {
	if req.Product == nil {
		return nil, fmt.Errorf("product is required")
	}
	var createdName string
	var createdPrice float64
	var createdQuantity uint32
	name := req.Product.Name
	price := req.Product.Price
	quantity := req.Product.Quantity
	query := `
		INSERT INTO products
		(name, price, quantity)
		VALUES ($1, $2, $3)
		RETURNING name, price, quantity
	`
	err := m.DB.QueryRow(ctx, query, name, price, quantity).Scan(
		&createdName, &createdPrice, &createdQuantity,
	)
	if err != nil {
		return nil, fmt.Errorf("database query error %v", err)
	}
	return &productspb.AddProductResponse{
		Product: &productspb.Product{
			Name:     createdName,
			Price:    createdPrice,
			Quantity: createdQuantity,
		},
	}, nil
}

func (m *ProductModel) GetProducts(ctx context.Context, req *productspb.GetProductsRequest) (*productspb.GetProductsResponse, error) {
	query := `
		SELECT * FROM products
	`
	rows, err := m.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("database query error %v", err)
	}
	products, err := pgx.CollectRows(rows, pgx.RowToStructByName[Product])
	if err != nil {
		return nil, fmt.Errorf("database query error %v", err)
	}
	pbproducts := make([]*productspb.Product, 0, len(products))
	for _, p := range products {
		pbproducts = append(pbproducts, &productspb.Product{
			Name:     p.Name,
			Price:    float64(p.Price),
			Quantity: uint32(p.Quantity),
		})
	}
	return &productspb.GetProductsResponse{
		Products: pbproducts,
	}, nil
}

func (m *ProductModel) GetProduct(ctx context.Context, req *productspb.GetProductRequest) (*productspb.GetProductResponse, error) {
	query := `
		SELECT * FROM products
		WHERE name = $1
	`
	rows, err := m.DB.Query(ctx, query, req.Name)
	if err != nil {
		return nil, fmt.Errorf("database query error %v", err)
	}
	product, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Product])
	if err != nil {
		return nil, fmt.Errorf("database query error %v", err)
	}
	return &productspb.GetProductResponse{
		Product: &productspb.Product{
			Name:     product.Name,
			Price:    float64(product.Price),
			Quantity: uint32(product.Quantity),
		},
	}, nil
}
