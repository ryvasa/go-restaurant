//go:build wireinject
// +build wireinject

package di

import (
	"database/sql"

	"github.com/google/wire"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/handler"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/internal/usecase"
)

// InitializeUserHandler initializes UserHandler with dependencies
func InitializeMenuHandler(db *sql.DB) *handler.MenuHandler {
	wire.Build(
		repository.NewMenuRepository,
		usecase.NewMenuUsecase,
		handler.NewMenuHandler,
	)
	return &handler.MenuHandler{}
}
