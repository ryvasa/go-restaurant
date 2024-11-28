package dto

type CreateReservationRequest struct {
	UserId          string `json:"user_id" validate:"required"`
	TableId         string `json:"table_id" validate:"required"`
	ReservationDate string `json:"reservation_date" validate:"required,datetime=2006-01-02"`
	ReservationTime string `json:"reservation_time" validate:"required,datetime=15:04:05"`
	NumberOfGuests  int    `json:"number_of_guests" validate:"required"`
}

type UpdateReservationRequest struct {
	ReservationDate string `json:"reservation_date,omitempty" validate:"omitempty,required,datetime=2006-01-02"`
	ReservationTime string `json:"reservation_time,omitempty" validate:"omitempty,required,datetime=15:04:05"`
	NumberOfGuests  int    `json:"number_of_guests,omitempty" validate:"omitempty,required"`
	Status          string `json:"status,omitempty" validate:"omitempty,oneof=pending confirmed canceled"`
}
