package handler

import (
	shared "ecommerce-shared"
	"log"
	"net/http"
	"time"

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
	mux.HandleFunc("GET /api/v1/ping_products", h.pingProductsHandler)
	mux.HandleFunc("POST /api/v1/add_product", h.addProduct)
	mux.HandleFunc("GET /api/v1/products", h.getProducts)
}

func (h *ProductsHTTPHandler) pingProductsHandler(w http.ResponseWriter, r *http.Request) {
	_ = shared.WriteJSON(w, http.StatusOK, "Products PONG")
}

func (h *ProductsHTTPHandler) addProduct(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var requestBody productspb.AddProductRequest
	if err := shared.ReadJSON(r, &requestBody); err != nil {
		shared.WriteErrorBadRequest(w, "Invalid requestBody", err)
		shared.LogBadRequest(r.Method, r.RequestURI, time.Since(start))
		return
	}

	ctx := r.Context()
	res, err := h.productsClient.AddProduct(ctx, &requestBody)
	if err != nil {
		log.Printf("failed to add a product: %s", err)
		shared.WriteErrorServerError(w, "failed to add a product", err)
		shared.LogInternalServerError(r.Method, r.RequestURI, time.Since(start))
		return
	}

	_ = shared.WriteJSON(w, http.StatusCreated, res)
	shared.LogOK(r.Method, r.RequestURI, time.Since(start))
}

func (h *ProductsHTTPHandler) getProducts(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	var requestBody productspb.GetProductsRequest
	if err := shared.ReadJSON(r, &requestBody); err != nil {
		shared.WriteErrorBadRequest(w, "Invalid requestBody", err)
		shared.LogBadRequest(r.Method, r.RequestURI, time.Since(start))
		return
	}

	ctx := r.Context()
	res, err := h.productsClient.GetProducts(ctx, &requestBody)
	if err != nil {
		log.Printf("failed to init request: %s", err)
		shared.WriteErrorServerError(w, "failed to init request", err)
		shared.LogInternalServerError(r.Method, r.RequestURI, time.Since(start))
		return
	}

	_ = shared.WriteJSON(w, http.StatusOK, res)
	shared.LogOK(r.Method, r.RequestURI, time.Since(start))
}
