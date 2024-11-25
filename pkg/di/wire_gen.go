// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/google/wire"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/config"
	"github.com/ryvasa/go-restaurant/pkg/database"
	"github.com/ryvasa/go-restaurant/utils"
)

// Injectors from wire.go:

// InitializeHandlers initializes Handlers with dependencies
func InitializeHandlers() (*handler.Handlers, error) {
	configConfig, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	string2 := database.ProvideDSN(configConfig)
	db, err := database.NewMySQLConnection(string2)
	if err != nil {
		return nil, err
	}
	menuRepository := repository.NewMenuRepository(db)
	menuUsecase := usecase.NewMenuUsecase(menuRepository)
	menuHandlerImpl := handler.NewMenuHandler(menuUsecase)
	userRepository := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userHandlerImpl := handler.NewUserHandler(userUsecase)
	reviewRepository := repository.NewReviewRepository(db)
	reviewUsecase := usecase.NewReviewUsecase(reviewRepository, userRepository, menuRepository)
	reviewHandlerImpl := handler.NewReviewHandler(reviewUsecase)
	tokenUtil := utils.NewTokenUtil(configConfig)
	authUsecase := usecase.NewAuthUsecase(userRepository, tokenUtil)
	authHandlerImpl := handler.NewAuthHandler(authUsecase)
	orderRepository := repository.NewOrderRepository(db)
	orderMenuRepository := repository.NewOrderMenuRepository(db)
	orderUsecase := usecase.NewOrderUsecase(orderRepository, menuRepository, userRepository, orderMenuRepository)
	orderHandlerImpl := handler.NewOrderHandler(orderUsecase)
	handlers := handler.NewHandlers(menuHandlerImpl, userHandlerImpl, reviewHandlerImpl, authHandlerImpl, orderHandlerImpl)
	return handlers, nil
}

// wire.go:

var orderMenuSet = wire.NewSet(repository.NewOrderMenuRepository)

var orderSet = wire.NewSet(repository.NewOrderRepository, usecase.NewOrderUsecase, handler.NewOrderHandler)

var menuSet = wire.NewSet(repository.NewMenuRepository, usecase.NewMenuUsecase, handler.NewMenuHandler)

var reviewSet = wire.NewSet(repository.NewReviewRepository, usecase.NewReviewUsecase, handler.NewReviewHandler)

var authSet = wire.NewSet(usecase.NewAuthUsecase, handler.NewAuthHandler)

var userSet = wire.NewSet(repository.NewUserRepository, usecase.NewUserUsecase, handler.NewUserHandler)

var utilSet = wire.NewSet(utils.NewTokenUtil)
