package handler

import (
	"net/http"

	productspb "ecommerce-api/gen/products"
)

type ProductsHTTPHandler struct {
	productsClient productspb.ProductServiceClient
}

func NewProductsHTTPHandler(productsClient productspb.ProductServiceClient) *ProductsHTTPHandler {
	return &ProductsHTTPHandler{
		productsClient: productsClient,
	}
}

func (h *ProductsHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("ping /api/v1/ping_products", h.pingProductsHandler)
}

func (h *ProductsHTTPHandler) pingProductsHandler(w http.ResponseWriter, r *http.Request) {}
