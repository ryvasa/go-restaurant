package routes

import (
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/domain"
)

func SetupRoutes(r *mux.Router, menuUsecase domain.MenuUsecase) {
	// Setup menu routes
	MenuRoutes(r, menuUsecase)

	// Bisa tambah route lain
	// SetupUserRoutes(r, userUsecase)
	// SetupOrderRoutes(r, orderUsecase)
}
