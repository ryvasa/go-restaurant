package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func ReservationRoutes(public, protected *mux.Router, handler handler.ReservationHandler) {
	// no auth
	public.HandleFunc("/reservations", handler.GetAll).Methods("GET")
	public.HandleFunc("/reservations/{id}", handler.Get).Methods("GET")
	// all role
	protected.HandleFunc("/reservations", handler.Create).Methods("POST")
	protected.HandleFunc("/reservations/{id}", handler.Update).Methods("PATCH")

	// admin and staff
	protected.HandleFunc("/reservations/{id}", handler.Delete).Methods("DELETE")
	protected.HandleFunc("/reservations/{id}/restore", handler.Restore).Methods("PATCH")
}
