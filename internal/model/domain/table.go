package domain

import (
	"time"

	"github.com/google/uuid"
)

type Table struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	Number    string    `json:"number" validate:"required"`
	Capacity  int       `json:"capacity" validate:"required"`
	Location  string    `json:"location" validate:"required,oneof=indoor outdoor"`
	Status    string    `json:"status" validate:"required,oneof=available, reserved, out of service"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
