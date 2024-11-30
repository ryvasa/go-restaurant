package dto

import "github.com/google/uuid"

type CreateInventoryRequest struct {
	IngredientId uuid.UUID `json:"ingredient_id" validate:"required"`
	Quantity     float64   `json:"quantity" validate:"required"`
}

type UpdateInventoryRequest struct {
	Quantity float64 `json:"quantity" validate:"required"`
}
