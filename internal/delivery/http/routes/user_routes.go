package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func UserRoutes(r *mux.Router, userHandler handler.UserHandler) {
	r.HandleFunc("/users", userHandler.GetAll).Methods("GET")
	r.HandleFunc("/users", userHandler.Create).Methods("POST")
	r.HandleFunc("/users/{id}", userHandler.Get).Methods("GET")
	r.HandleFunc("/users/{id}", userHandler.Update).Methods("PATCH")
	r.HandleFunc("/users/{id}", userHandler.Delete).Methods("DELETE")
	r.HandleFunc("/users/{id}/restore", userHandler.Restore).Methods("PATCH")
}
