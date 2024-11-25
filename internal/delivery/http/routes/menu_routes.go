package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func MenuRoutes(public, protected *mux.Router, handler handler.MenuHandler) {
	// all role tanpa auth
	public.HandleFunc("/menu", handler.GetAll).Methods("GET")
	public.HandleFunc("/menu/{id}", handler.Get).Methods("GET")

	// admin only
	protected.HandleFunc("/menu", handler.Create).Methods("POST")
	protected.HandleFunc("/menu/{id}", handler.Update).Methods("PATCH")
	protected.HandleFunc("/menu/{id}", handler.Delete).Methods("DELETE")
	protected.HandleFunc("/menu/{id}/restore", handler.Restore).Methods("PATCH")
}
