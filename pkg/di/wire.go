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

// InitializeMenuHandler initializes MenuHandler with dependencies
func InitializeHandlers() (*handler.Handlers, error) {
	wire.Build(
		config.LoadConfig,
		database.ProvideDSN,
		database.NewMySQLConnection,
		menuSet,
		// orderSet,
		// userSet,
		handler.NewHandlers,
	)
	return &handler.Handlers{}, nil
}
