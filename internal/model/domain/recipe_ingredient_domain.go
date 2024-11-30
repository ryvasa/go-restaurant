package domain

import "github.com/google/uuid"

type RecipeIngredient struct {
	Id           uuid.UUID `json:"id" validate:"required"`
	RecipeId     uuid.UUID `json:"recipe_id" validate:"required"`
	IngredientId uuid.UUID `json:"ingredient_id" validate:"required"`
	Quantity     float64   `json:"quantity" validate:"required"`
}

type SimpleRecipeIngredient struct {
	IngredientId uuid.UUID `json:"ingredient_id"`
	Quantity     float64   `json:"quantity"`
	Name         string    `json:"name"`
}
