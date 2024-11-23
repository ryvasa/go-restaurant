package dto

type CreateMenuRequest struct {
	Restaurant  string `json:"restaurant_id" validate:"required,min=3,max=100"`
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Price       int    `json:"price" validate:"required,gt=0"`
	Description string `json:"description" validate:"required,min=3,max=1000"`
	Category    string `json:"category" validate:"required,oneof=main appetizer dessert drink snack vegetarian kids local special combo breakfast healthy international seafood spicy"`
	ImageURL    string `json:"image_url" validate:"required,min=3,max=1000"`
}

type UpdateMenuRequest struct {
	Restaurant  string `json:"restaurant_id,omitempty" validate:"omitempty,min=3,max=100"`
	Name        string `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Price       int    `json:"price,omitempty" validate:"omitempty,gt=0"`
	Description string `json:"description,omitempty" validate:"omitempty,min=3,max=1000"`
	Category    string `json:"category,omitempty" validate:"omitempty,oneof=main appetizer dessert drink snack vegetarian kids local special combo breakfast healthy international seafood spicy"`
	ImageURL    string `json:"image_url,omitempty" validate:"omitempty,min=3,max=1000"`
}
