package domain

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	Id          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	MenuId      uuid.UUID `json:"menu_id" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
}

type RecipeAndIngredients struct {
	Id          uuid.UUID                `json:"id" validate:"required"`
	Name        string                   `json:"name" validate:"required"`
	Description string                   `json:"description" validate:"required"`
	MenuId      uuid.UUID                `json:"menu_id" validate:"required"`
	Ingredients []SimpleRecipeIngredient `json:"ingredients"`
	CreatedAt   time.Time                `json:"created_at" validate:"required"`
	UpdatedAt   time.Time                `json:"updated_at" validate:"required"`
}
