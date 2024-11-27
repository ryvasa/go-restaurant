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
	menuRepository := repository.NewMenuRepository()
	menuUsecase := usecase.NewMenuUsecase(db, menuRepository)
	menuHandlerImpl := handler.NewMenuHandler(menuUsecase)
	userRepository := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(db, userRepository)
	userHandlerImpl := handler.NewUserHandler(userUsecase)
	reviewRepository := repository.NewReviewRepository()
	orderRepository := repository.NewOrderRepository()
	reviewUsecase := usecase.NewReviewUsecase(db, reviewRepository, userRepository, menuRepository, orderRepository)
	reviewHandlerImpl := handler.NewReviewHandler(reviewUsecase)
	tokenUtil := utils.NewTokenUtil(configConfig)
	authUsecase := usecase.NewAuthUsecase(db, userRepository, tokenUtil)
	authHandlerImpl := handler.NewAuthHandler(authUsecase)
	orderMenuRepository := repository.NewOrderMenuRepository()
	orderUsecase := usecase.NewOrderUsecase(db, orderRepository, menuRepository, userRepository, orderMenuRepository)
	orderHandlerImpl := handler.NewOrderHandler(orderUsecase)
	tableRepository := repository.NewTableRepository()
	tableUsecase := usecase.NewTableUsecase(db, tableRepository)
	tableHandlerImpl := handler.NewTableHandler(tableUsecase)
	reservationRepository := repository.NewReservationRepository()
	reservationUsecase := usecase.NewReservationUsecase(db, reservationRepository, tableRepository)
	reservationHandlerImpl := handler.NewReservationHandler(reservationUsecase)
	handlers := handler.NewHandlers(menuHandlerImpl, userHandlerImpl, reviewHandlerImpl, authHandlerImpl, orderHandlerImpl, tableHandlerImpl, reservationHandlerImpl)
	return handlers, nil
}

// wire.go:

var orderMenuSet = wire.NewSet(repository.NewOrderMenuRepository)

var orderSet = wire.NewSet(repository.NewOrderRepository, usecase.NewOrderUsecase, handler.NewOrderHandler)

var menuSet = wire.NewSet(repository.NewMenuRepository, usecase.NewMenuUsecase, handler.NewMenuHandler)

var reviewSet = wire.NewSet(repository.NewReviewRepository, usecase.NewReviewUsecase, handler.NewReviewHandler)

var authSet = wire.NewSet(usecase.NewAuthUsecase, handler.NewAuthHandler)

var userSet = wire.NewSet(repository.NewUserRepository, usecase.NewUserUsecase, handler.NewUserHandler)

var tableSet = wire.NewSet(repository.NewTableRepository, usecase.NewTableUsecase, handler.NewTableHandler)

var reservationSet = wire.NewSet(repository.NewReservationRepository, usecase.NewReservationUsecase, handler.NewReservationHandler)

var utilSet = wire.NewSet(utils.NewTokenUtil)
