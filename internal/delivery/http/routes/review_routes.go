package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func ReviewRoutes(public, protected *mux.Router, handler handler.ReviewHandler) {
	// no auth
	public.HandleFunc("/reviews/menu/{id}", handler.GetAllByMenuId).Methods("GET")
	public.HandleFunc("/reviews/{id}", handler.GetOneById).Methods("GET")
	// all role
	protected.HandleFunc("/reviews", handler.Create).Methods("POST")
	protected.HandleFunc("/reviews/{id}", handler.Update).Methods("PATCH")
}
