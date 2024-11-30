package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type IngredientRepository interface {
	Create(ctx context.Context, ingredient domain.Ingredient) error
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Ingredient, error)
	GetOneByName(ctx context.Context, name string) (domain.Ingredient, error)
	Update(ctx context.Context, id uuid.UUID, ingredient domain.Ingredient) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	GetDeletedById(ctx context.Context, id uuid.UUID) (domain.Ingredient, error)
}
