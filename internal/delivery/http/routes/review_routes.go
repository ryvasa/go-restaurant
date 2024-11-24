package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func ReviewRoutes(r *mux.Router, reviewHandler handler.ReviewHandler) {

	r.HandleFunc("/review/menu/{id}", reviewHandler.GetAllByMenuId).Methods("GET")
	r.HandleFunc("/review", reviewHandler.Create).Methods("POST")
	r.HandleFunc("/review/{id}", reviewHandler.GetOneById).Methods("GET")
	r.HandleFunc("/review/{id}", reviewHandler.Update).Methods("PATCH")
	// r.HandleFunc("/review/{id}", reviewHandler.Delete).Methods("DELETE")
}
