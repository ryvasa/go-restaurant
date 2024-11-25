package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func ReviewRoutes(public, protected *mux.Router, handler handler.ReviewHandler) {
	// no auth
	public.HandleFunc("/review/menu/{id}", handler.GetAllByMenuId).Methods("GET")
	public.HandleFunc("/review/{id}", handler.GetOneById).Methods("GET")
	// all role
	protected.HandleFunc("/review", handler.Create).Methods("POST")
	protected.HandleFunc("/review/{id}", handler.Update).Methods("PATCH")
}
