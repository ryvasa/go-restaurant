package domain

import "github.com/google/uuid"

type OrderMenu struct {
	OrderId uuid.UUID `json:"order_id" validate:"required"`
	MenuId  uuid.UUID `json:"menu_id" validate:"required"`
	Quantity int       `json:"quantity" validate:"required"`
}
