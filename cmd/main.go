package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/routes"
	"github.com/ryvasa/go-restaurant/pkg/di"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

func main() {
	// Initialize dependencies using Wire
	handlers, err := di.InitializeHandlers()
	if err != nil {
		logger.Log.Fatal("Failed to initialize dependencies:", err)
	}

	r := mux.NewRouter()
	routes.NewRoutes(r, handlers)

	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/",
		http.FileServer(http.Dir("uploads"))))

	// Start server
	logger.Log.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Log.Fatal("Server failed to start:", err)
	}
}
