package repository

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type MenuRepository interface {
	GetAll(ctx context.Context) ([]domain.Menu, error)
	Create(ctx context.Context, menu domain.Menu) (domain.Menu, error)
	Get(ctx context.Context, id string) (domain.Menu, error)
	Update(ctx context.Context, menu domain.Menu) (domain.Menu, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) (domain.Menu, error)
	GetDeletedMenuById(ctx context.Context, id string) ([]domain.Menu, error)
}
