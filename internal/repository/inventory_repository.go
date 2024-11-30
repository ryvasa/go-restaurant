package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type InventoryRepository interface {
	Create(ctx context.Context, inventory domain.Inventory) error
	GetOneByIngredientId(ctx context.Context, ingredientId uuid.UUID) (domain.Inventory, error)
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Inventory, error)
	Update(ctx context.Context, id uuid.UUID, inventory domain.Inventory) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	GetDeletedById(ctx context.Context, id uuid.UUID) (domain.Inventory, error)
}
