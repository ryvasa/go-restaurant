package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type IngredientUsecase interface {
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Ingredient, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateIngredientRequest) (domain.Ingredient, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) (domain.Ingredient, error)
}
