package main

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/middleware"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/routes"
	"github.com/ryvasa/go-restaurant/pkg/config"
	"github.com/ryvasa/go-restaurant/pkg/di"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

func main() {
	// Initialize config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Fatal("Failed to load config:", err)
	}

	// Initialize handlers
	handlers, err := di.InitializeHandlers()
	if err != nil {
		logger.Log.Fatal("Failed to initialize handlers:", err)
	}
	// Initialize Casbin
	enforcer, err := casbin.NewEnforcer("pkg/config/casbin/model.conf", "pkg/config/casbin/policy.csv")
	if err != nil {
		logger.Log.Fatal("Failed to initialize Casbin:", err)
	}

	// Initialize auth middleware
	tokenUtil := utils.NewTokenUtil(cfg)
	authenticationMiddleware := middleware.NewAuthenticationMiddleware(tokenUtil)
	authorizationMiddleware := middleware.NewAuthorizationMiddleware(enforcer)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	// Public routes (tidak perlu auth)
	publicRouter := r.PathPrefix("/api").Subrouter()

	// Protected routes (perlu auth & authorization)
	protectedRouter := r.PathPrefix("/api").Subrouter()
	protectedRouter.Use(authenticationMiddleware.Handle)
	protectedRouter.Use(authorizationMiddleware.Handle)

	// Setup routes
	routes.NewRoutes(publicRouter, protectedRouter, handlers)

	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/",
		http.FileServer(http.Dir("uploads"))))

	// Start server
	logger.Log.Info("Server starting on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		logger.Log.Fatal("Server failed to start:", err)
	}
}
