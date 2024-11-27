package dto

type CreateTableRequest struct {
	Number   string `json:"number" validate:"required"`
	Capacity int    `json:"capacity" validate:"required"`
	Location string `json:"location" validate:"required,oneof=indoor outdoor"`
	Status   string `json:"status" validate:"required,oneof=avaible, reserved, out of service"`
}

type UpdateTableRequest struct {
	Number   string `json:"number,omitempty" validate:"omitempty,required"`
	Capacity int    `json:"capacity,omitempty" validate:"omitempty,required"`
	Location string `json:"location,omitempty" validate:"omitempty,oneof=indoor outdoor"`
	Status   string `json:"status,omitempty" validate:"omitempty,oneof=avaible, reserved, out of service"`
}
