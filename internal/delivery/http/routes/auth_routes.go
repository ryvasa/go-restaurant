package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func AuthRoutes(public *mux.Router, authHandler handler.AuthHandler) {
	public.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
}
