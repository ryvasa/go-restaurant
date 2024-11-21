package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/middleware"
)

func NewRoutes(r *mux.Router, handlers *handler.Handlers) {
	r.Use(middleware.LoggingMiddleware)
	MenuRoutes(r, handlers.MenuHandler)
}
