package dto

type CreateIngredientRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Quantity    float64 `json:"quantity" validate:"required"`
}

type UpdateIngredientRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,required"`
	Description string `json:"description,omitempty" validate:"omitempty,required"`
}
