package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
)

func NewRoutes(public, protected *mux.Router, handlers *handler.Handlers) {

	MenuRoutes(public, protected, handlers.MenuHandler)
	UserRoutes(public, protected, handlers.UserHandler)
	ReviewRoutes(public, protected, handlers.ReviewHandler)
	AuthRoutes(public, handlers.AuthHandler)
}
