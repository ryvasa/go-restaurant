package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/domain"
)

func MenuRoutes(r *mux.Router, menuHandler domain.MenuHandler) {

	r.HandleFunc("/menu", menuHandler.GetAll).Methods("GET")
	r.HandleFunc("/menu", menuHandler.Create).Methods("POST")
	r.HandleFunc("/menu/{id}", menuHandler.Get).Methods("GET")
	r.HandleFunc("/menu/{id}", menuHandler.Update).Methods("PATCH")
	r.HandleFunc("/menu/{id}", menuHandler.Delete).Methods("DELETE")
}
