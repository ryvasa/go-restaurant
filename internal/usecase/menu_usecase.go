package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/domain"
)

type MenuUsecase interface {
	GetAll(ctx context.Context) ([]domain.Menu, error)
	Create(ctx context.Context, menu domain.Menu) (domain.Menu, error)
	Get(ctx context.Context, id string) (domain.Menu, error)
	Update(ctx context.Context, menu domain.Menu) (domain.Menu, error)
	Delete(ctx context.Context, id string) error
}
