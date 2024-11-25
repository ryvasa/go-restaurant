package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func OrderRoutes(protected *mux.Router, handler handler.OrderHandler) {
	// admin only
	protected.HandleFunc("/order", handler.Create).Methods("POST")
	protected.HandleFunc("/order/{id}", handler.GetOneById).Methods("GET")
}
