package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/middleware"
)

func NewRoutes(r *mux.Router, handlers *handler.Handlers) {
	// config := &middleware.RecoveryConfig{
	// 	EnableStackTrace: true,
	// 	LogError:         true,
	// }

	// r.Use(middleware.RecoveryMiddleware(config))
	r.Use(middleware.LoggingMiddleware)
	MenuRoutes(r, handlers.MenuHandler)
	UserRoutes(r, handlers.UserHandler)
}
