package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func RecipeRoutes(protected *mux.Router, handler handler.RecipeHandler) {
	protected.HandleFunc("/recipes", handler.Create).Methods("POST")
	protected.HandleFunc("/recipes", handler.GetAll).Methods("GET")
	protected.HandleFunc("/recipes/{id}", handler.GetOneById).Methods("GET")
	protected.HandleFunc("/recipes/{id}", handler.Update).Methods("PATCH")
	protected.HandleFunc("/recipes/{id}", handler.Delete).Methods("DELETE")
	protected.HandleFunc("/recipes/{id}/restore", handler.Restore).Methods("PATCH")
}
