//go:build wireinject
// +build wireinject

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
	repository.NewAuthRepository,
	usecase.NewAuthUsecase,
	handler.NewAuthHandler,
)

var userSet = wire.NewSet(
	repository.NewUserRepository,
	usecase.NewUserUsecase,
	handler.NewUserHandler,
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
		utilSet,
		menuSet,
		userSet,
		reviewSet,
		authSet,
		handler.NewHandlers,
	)
	return &handler.Handlers{}, nil
}
