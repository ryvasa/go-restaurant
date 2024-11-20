package dto

type CreateMenuRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Price int    `json:"price" validate:"required,gt=0"`
}

type UpdateMenuRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=100"`
	Price int    `json:"price" validate:"required,gt=0"`
}
