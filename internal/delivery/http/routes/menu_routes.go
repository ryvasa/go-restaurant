package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func MenuRoutes(r *mux.Router, menuHandler handler.MenuHandler) {

	r.HandleFunc("/menu", menuHandler.GetAll).Methods("GET")
	r.HandleFunc("/menu", menuHandler.Create).Methods("POST")
	r.HandleFunc("/menu/{id}", menuHandler.Get).Methods("GET")
	r.HandleFunc("/menu/{id}", menuHandler.Update).Methods("PATCH")
	r.HandleFunc("/menu/{id}", menuHandler.Delete).Methods("DELETE")
	r.HandleFunc("/menu/{id}/restore", menuHandler.Restore).Methods("PATCH")
}
