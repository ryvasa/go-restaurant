package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func TableRoutes(public, protected *mux.Router, handler handler.TableHandler) {
	// no auth
	public.HandleFunc("/tables", handler.GetAll).Methods("GET")
	public.HandleFunc("/tables/{id}", handler.Get).Methods("GET")
	// staff role
	protected.HandleFunc("/tables", handler.Create).Methods("POST")
	protected.HandleFunc("/tables/{id}", handler.Update).Methods("PATCH")
	protected.HandleFunc("/tables/{id}", handler.Delete).Methods("DELETE")
	protected.HandleFunc("/tables/{id}/restore", handler.Restore).Methods("PATCH")
}
