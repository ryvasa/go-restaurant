package dto

type CreateReservationRequest struct {
	TableId         string `json:"table_id" validate:"required"`
	ReservationDate string `json:"reservation_date" validate:"required"`
	ReservationTime string `json:"reservation_time" validate:"required"`
	NumberOfGuests  int    `json:"number_of_guests" validate:"required"`
}

type UpdateReservationRequest struct {
	TableId         string `json:"table_id,omitempty" validate:"omitempty,required"`
	ReservationDate string `json:"reservation_date,omitempty" validate:"omitempty,required"`
	ReservationTime string `json:"reservation_time,omitempty" validate:"omitempty,required"`
	NumberOfGuests  int    `json:"number_of_guests,omitempty" validate:"omitempty,required"`
	Status          string `json:"status,omitempty" validate:"omitempty,oneof=pending confirmed canceled"`
}
