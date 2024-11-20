package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
	"github.com/ryvasa/go-restaurant/internal/domain"
)

func MenuRoutes(r *mux.Router, menuUsecase domain.MenuUsecase) {
    menuHandler := handler.NewMenuHandler(menuUsecase)

    r.HandleFunc("/menu", menuHandler.GetAll).Methods("GET")
    r.HandleFunc("/menu", menuHandler.Create).Methods("POST")
    // Tambahkan route lain untuk menu
    // r.HandleFunc("/menu/{id}", menuHandler.GetByID).Methods("GET")
    // r.HandleFunc("/menu/{id}", menuHandler.Update).Methods("PUT")
    // r.HandleFunc("/menu/{id}", menuHandler.Delete).Methods("DELETE")
}
