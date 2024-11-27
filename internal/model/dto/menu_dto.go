package dto

import "mime/multipart"

type CreateMenuRequest struct {
	Name        string                `form:"name" validate:"required,min=3,max=100"`
	Price       float64               `form:"price" validate:"required,gt=0"`
	Description string                `form:"description" validate:"required,min=3,max=1000"`
	Category    string                `form:"category" validate:"required,oneof=main appetizer dessert drink snack vegetarian kids local special combo breakfast healthy international seafood spicy"`
	Image       *multipart.FileHeader `form:"image" validate:"required"`
}
type UpdateMenuRequest struct {
	Name        string                `form:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Price       float64               `form:"price,omitempty" validate:"omitempty,gt=0"`
	Description string                `form:"description,omitempty" validate:"omitempty,min=3,max=1000"`
	Category    string                `form:"category, omitempty" validate:"omitempty,oneof=main appetizer dessert drink snack vegetarian kids local special combo breakfast healthy international seafood spicy"`
	Image       *multipart.FileHeader `form:"image,omitempty" validate:"omitempty,required"`
}
