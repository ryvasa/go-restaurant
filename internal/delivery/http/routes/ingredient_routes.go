package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func IngredientRoutes(protected *mux.Router, handler handler.IngredientHandler) {
	protected.HandleFunc("/ingredients/{id}", handler.GetOneById).Methods("GET")
	protected.HandleFunc("/ingredients/{id}", handler.Update).Methods("PATCH")
	protected.HandleFunc("/ingredients/{id}", handler.Delete).Methods("DELETE")
	protected.HandleFunc("/ingredients/{id}/restore", handler.Restore).Methods("PATCH")
}
