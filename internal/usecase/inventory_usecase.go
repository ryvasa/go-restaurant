package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type InventoryUsecase interface {
	Create(ctx context.Context, req dto.CreateInventoryRequest) (domain.Inventory, error)
	GetOneByIngredientId(ctx context.Context, ingredientId uuid.UUID) (domain.Inventory, error)
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Inventory, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateInventoryRequest) (domain.Inventory, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) (domain.Inventory, error)
}
