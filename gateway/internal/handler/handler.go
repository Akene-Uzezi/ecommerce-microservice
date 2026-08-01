// Package httphandler for the gateway service
package httphandler

import (
	shared "ecommerce-shared"
	"log"
	"net/http"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/ping", h.ping)
}

func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %v", r.Method, r.URL)
	if err := shared.WriteJSON(w, http.StatusOK, "pong"); err != nil {
		log.Printf("failed to write JSON: %v", err.Error())
		return
	}
}
