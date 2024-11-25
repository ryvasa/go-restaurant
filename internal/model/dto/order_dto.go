package dto

type OrderMenuDto struct {
	MenuId   string `json:"menu_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

type CreateOrderDto struct {
	Menu []OrderMenuDto `json:"menu" validate:"required"`
}
