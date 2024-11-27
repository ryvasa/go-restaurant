package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func OrderRoutes(protected *mux.Router, handler handler.OrderHandler) {
	protected.HandleFunc("/order", handler.Create).Methods("POST")
	protected.HandleFunc("/order/{id}", handler.GetOneById).Methods("GET")

	// staff and admin only
	protected.HandleFunc("/order/{id}/status", handler.UpdateOrderStatus).Methods("PATCH")
	protected.HandleFunc("/order/{id}/payment", handler.UpdatePayment).Methods("PATCH")
}
