//go:build wireinject
// +build wireinject

package di

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/config"
	"github.com/ryvasa/go-restaurant/pkg/database"
	"github.com/ryvasa/go-restaurant/utils"
)

var orderMenuSet = wire.NewSet(
	repository.NewOrderMenuRepository,
)

var orderSet = wire.NewSet(
	repository.NewOrderRepository,
	usecase.NewOrderUsecase,
	handler.NewOrderHandler,
)

var menuSet = wire.NewSet(
	repository.NewMenuRepository,
	usecase.NewMenuUsecase,
	handler.NewMenuHandler,
)

var reviewSet = wire.NewSet(
	repository.NewReviewRepository,
	usecase.NewReviewUsecase,
	handler.NewReviewHandler,
)

var authSet = wire.NewSet(
	usecase.NewAuthUsecase,
	handler.NewAuthHandler,
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	usecase.NewUserUsecase,
	handler.NewUserHandler,
)

func ProvideDBConnection(db *sql.DB) repository.DB {
	return db
}

var tableSet = wire.NewSet(
	repository.NewTableRepository,
	usecase.NewTableUsecase,
	handler.NewTableHandler,
)

var reservationSet = wire.NewSet(
	repository.NewReservationRepository,
	usecase.NewReservationUsecase,
	handler.NewReservationHandler,
)

var txSet = wire.NewSet(
	repository.NewTransactionRepository,
)

var utilSet = wire.NewSet(
	utils.NewTokenUtil,
)

// InitializeHandlers initializes Handlers with dependencies
func InitializeHandlers() (*handler.Handlers, error) {
	wire.Build(
		config.LoadConfig,
		database.ProvideDSN,
		database.NewMySQLConnection,
		ProvideDBConnection,
		utilSet,
		menuSet,
		userSet,
		reviewSet,
		authSet,
		orderSet,
		orderMenuSet,
		tableSet,
		reservationSet,
		txSet,
		handler.NewHandlers,
	)
	return &handler.Handlers{}, nil
}
