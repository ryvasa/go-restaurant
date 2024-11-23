package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/delivery/http/dto"
	"github.com/ryvasa/go-restaurant/internal/domain"
)

type MenuUsecase interface {
	GetAll(ctx context.Context) ([]domain.Menu, error)
	Create(ctx context.Context, req dto.CreateMenuRequest) (domain.Menu, error)
	Get(ctx context.Context, id string) (domain.Menu, error)
	Update(ctx context.Context, id string, req dto.UpdateMenuRequest) (domain.Menu, error)
	Delete(ctx context.Context, id string) error
}
