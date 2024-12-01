package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type RecipeRepository interface {
	Create(ctx context.Context, recipe domain.Recipe) error
	GetAll(ctx context.Context) ([]domain.Recipe, error)
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Recipe, error)
	Update(ctx context.Context, id uuid.UUID, recipe domain.Recipe) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	GetDeletedById(ctx context.Context, id uuid.UUID) (domain.Recipe, error)
	GetOneByMenuId(ctx context.Context, menuId uuid.UUID) (domain.Recipe, error)
}
