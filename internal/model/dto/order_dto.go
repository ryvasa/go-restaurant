package dto

type OrderMenuDto struct {
	MenuId   string `json:"menu_id" validate:"required"`
	Quantity int    `json:"quantity" validate:"required"`
}

type CreateOrderDto struct {
	Menu []OrderMenuDto `json:"menu" validate:"required"`
}

type UpdateOrderStatusDto struct {
	Status string `json:"status,omitempty" validate:"omitempty,oneof=pending processing success failed"`
}

type UpdatePaymentDto struct {
	PaymentMethod *string `json:"payment_method" validate:"required"`
}
