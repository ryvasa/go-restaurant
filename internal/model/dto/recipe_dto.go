package dto

import "github.com/google/uuid"

type CreateRecipeRequest struct {
	MenuId      uuid.UUID                 `json:"menu_id" validate:"required"`
	Name        string                    `json:"name" validate:"required"`
	Ingredients []CreateIngredientRequest `json:"ingredients" validate:"required"`
	Description string                    `json:"description" validate:"required"`
}

type UpdateRecipeRequest struct {
	Name        string                    `json:"name,omitempty" validate:"omitempty,required"`
	Ingredients []CreateIngredientRequest `json:"ingredients,omitempty" validate:"omitempty,required"`
	Description string                    `json:"description,omitempty" validate:"omitempty,required"`
}
