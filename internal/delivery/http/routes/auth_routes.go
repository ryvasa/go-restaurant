package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func AuthRoutes(r *mux.Router, authHandler handler.AuthHandler) {
	r.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
}
