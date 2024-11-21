package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func UserRoutes(r *mux.Router, userHandler handler.UserHandler) {
	r.HandleFunc("/users", userHandler.GetAll).Methods("GET")
	r.HandleFunc("/users", userHandler.Create).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.Get).Methods("Get")
	r.HandleFunc("/users/{id}", userHandler.Update).Methods("PATCH")
}
