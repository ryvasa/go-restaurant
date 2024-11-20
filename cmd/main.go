package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/routes"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/database"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

func main() {
	// Setup database connection
	db, err := database.NewMySQLConnection("root:123@tcp(localhost:3306)/go_restaurant?parseTime=true")
	if err != nil {
		logger.Log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Init repositories
	menuRepo := repository.NewMenuRepository(db)

	// Init usecases
	menuUsecase := usecase.NewMenuUsecase(menuRepo)

	r := mux.NewRouter()
	routes.SetupRoutes(r, menuUsecase)

	// Start server
	logger.Log.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Log.Fatal("Server failed to start:", err)
	}
}
