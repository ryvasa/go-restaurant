package domain

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	Id              uuid.UUID `json:"id" validate:"required"`
	TableId         uuid.UUID `json:"table_id" validate:"required"`
	UserId          uuid.UUID `json:"user_id" validate:"required"`
	ReservationDate time.Time `json:"reservation_date" validate:"required"`
	ReservationTime time.Time `json:"reservation_time" validate:"required"`
	NumberOfGuests  int       `json:"number_of_guests" validate:"required"`
	Status          string    `json:"status" validate:"required,oneof=pending confirmed canceled"`
	CreatedAt       time.Time `json:"created_at" validate:"required"`
	UpdatedAt       time.Time `json:"updated_at" validate:"required"`
}
