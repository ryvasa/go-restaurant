package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type RecipeUsecase interface {
	Create(ctx context.Context, req dto.CreateRecipeRequest) (domain.RecipeAndIngredients, error)
	GetAll(ctx context.Context) ([]domain.Recipe, error)
	GetOneById(ctx context.Context, id uuid.UUID) (domain.RecipeAndIngredients, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateRecipeRequest) (domain.RecipeAndIngredients, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) (domain.Recipe, error)
}
