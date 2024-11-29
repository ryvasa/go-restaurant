package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func OrderRoutes(protected *mux.Router, handler handler.OrderHandler) {
	protected.HandleFunc("/orders", handler.Create).Methods("POST")
	protected.HandleFunc("/orders/{id}", handler.GetOneById).Methods("GET")

	// staff and admin only
	protected.HandleFunc("/orders/{id}/status", handler.UpdateOrderStatus).Methods("PATCH")
	protected.HandleFunc("/orders/{id}/payment", handler.UpdatePayment).Methods("PATCH")
}
