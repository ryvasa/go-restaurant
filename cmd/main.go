package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/routes"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/database"
)

func main() {
	// Setup database connection
	db, err := database.NewMySQLConnection("root:123@tcp(localhost:3306)/go_restaurant?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Init repositories
	menuRepo := repository.NewMenuRepository(db)

	// Init usecases
	menuUsecase := usecase.NewMenuUsecase(menuRepo)

	// // Init handlers
	// menuHandler := handler.NewMenuHandler(menuUsecase)

	// // Setup routes
	// router := route.SetupRoutes(menuHandler)

	r := mux.NewRouter()
	routes.SetupRoutes(r, menuUsecase)

	// Start server
	log.Fatal(http.ListenAndServe(":8080", r))
}
