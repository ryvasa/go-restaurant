package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type RecipeIngredientRepository interface {
	Create(ctx context.Context, recipeIngredient domain.RecipeIngredient) error
	Update(ctx context.Context, id uuid.UUID, recipeIngredient domain.RecipeIngredient) error

	GetIngredientsByRecipeId(ctx context.Context, recipeId uuid.UUID) ([]domain.SimpleRecipeIngredient, error)
	GetIngredientsByRecipeIdAndIngredientId(ctx context.Context, recipeId uuid.UUID, ingredientId uuid.UUID) (domain.RecipeIngredient, error)
}
