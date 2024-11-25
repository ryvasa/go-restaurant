package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func UserRoutes(public, protected *mux.Router, handler handler.UserHandler) {
	// no auth
	public.HandleFunc("/users", handler.Create).Methods("POST")

	//all role
	protected.HandleFunc("/users/{id}", handler.Get).Methods("GET")
	protected.HandleFunc("/users/{id}", handler.Update).Methods("PATCH")

	//protected admin and staff
	protected.HandleFunc("/users", handler.GetAll).Methods("GET")
	protected.HandleFunc("/users/{id}", handler.Delete).Methods("DELETE")
	protected.HandleFunc("/users/{id}/restore", handler.Restore).Methods("PATCH")
}
