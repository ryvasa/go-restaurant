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

var userSet = wire.NewSet(
	repository.NewUserRepository,
	usecase.NewUserUsecase,
	handler.NewUserHandler,
)

// InitializeHandlers initializes Handlers with dependencies
func InitializeHandlers() (*handler.Handlers, error) {
	wire.Build(
		config.LoadConfig,
		database.ProvideDSN,
		database.NewMySQLConnection,
		menuSet,
		userSet,
		reviewSet,
		handler.NewHandlers,
	)
	return &handler.Handlers{}, nil
}
