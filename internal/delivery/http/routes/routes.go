package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/middleware"
	"github.com/ryvasa/go-restaurant/internal/domain"
)

func SetupRoutes(r *mux.Router, menuHandler domain.MenuHandler) {
	// Setup menu routes
	r.Use(middleware.LoggingMiddleware)
	MenuRoutes(r, menuHandler)

	// Bisa tambah route lain
	// SetupUserRoutes(r, userUsecase)
	// SetupOrderRoutes(r, orderUsecase)
}
