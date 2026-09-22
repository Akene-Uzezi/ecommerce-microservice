package handler

import (
	"ecommerce-gateway/internal/middleware"
	shared "ecommerce-shared"
	"log"
	"net/http"
	"time"

	authpb "ecommerce-api/gen/auth"
	stockpb "ecommerce-api/gen/stock"
)

type StockHTTPHandler struct {
	stockClient stockpb.StockServiceClient
	authClient  authpb.AuthServiceClient
}

func NewStockHTTPHandler(stockClient stockpb.StockServiceClient, authClient authpb.AuthServiceClient) *StockHTTPHandler {
	return &StockHTTPHandler{
		stockClient: stockClient,
		authClient:  authClient,
	}
}

func (h *StockHTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/stock", middleware.RequireAuth(h.authClient)(h.checkStock))
}

func (h *StockHTTPHandler) checkStock(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	productName := r.URL.Query().Get("product_name")
	if productName == "" {
		shared.WriteErrorBadRequest(w, "Missing product_name query parameter", nil)
		shared.LogBadRequest(r.Method, r.RequestURI, time.Since(start))
		return
	}

	ctx := r.Context()
	res, err := h.stockClient.CheckStock(ctx, &stockpb.CheckStockRequest{ProductName: productName})
	if err != nil {
		log.Printf("failed to check stock: %s", err)
		shared.WriteErrorServerError(w, "failed to check stock", err)
		shared.LogInternalServerError(r.Method, r.RequestURI, time.Since(start))
		return
	}

	_ = shared.WriteJSON(w, http.StatusOK, res)
	shared.LogOK(r.Method, r.RequestURI, time.Since(start))
}
