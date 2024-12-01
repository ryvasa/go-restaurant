package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func InventoryRoutes(protected *mux.Router, handler handler.InventoryHandler) {
	protected.HandleFunc("/inventory", handler.Create).Methods("POST")
	protected.HandleFunc("/inventory/{id}", handler.GetOneById).Methods("GET")
	protected.HandleFunc("/inventory/{id}/ingredient", handler.GetOneByIngredientId).Methods("GET")
	protected.HandleFunc("/inventory/{id}", handler.Update).Methods("PATCH")
	protected.HandleFunc("/inventory/{id}", handler.Delete).Methods("DELETE")
	protected.HandleFunc("/inventory/{id}/restore", handler.Restore).Methods("PATCH")
	protected.HandleFunc("/inventory/menu/{menu_id}", handler.CalculateMenuPortions).Methods("GET")
}
