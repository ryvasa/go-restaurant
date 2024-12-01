package domain

import (
	"time"

	"github.com/google/uuid"
)

type Inventory struct {
	Id           uuid.UUID `json:"id" validate:"required"`
	IngredientId uuid.UUID `json:"ingredient_id" validate:"required"`
	Quantity     float64   `json:"quantity" validate:"required"`
	CreatedAt    time.Time `json:"created_at" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at" validate:"required"`
}

type InventoryMenu struct {
	TotalPortions float64                  `json:"total_portions"`
	Menu          Menu                     `json:"menu"`
	Recipe        Recipe                   `json:"recipe"`
	Ingredients   []SimpleRecipeIngredient `json:"ingredients"`
}
