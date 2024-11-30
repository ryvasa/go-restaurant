package domain

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	Id            uuid.UUID `json:"id" validate:"required"`
	UserId        uuid.UUID `json:"user_id" validate:"required"`
	Status        string    `json:"status" validate:"required,oneof=pending processing success failed"`
	PaymentMethod *string   `json:"payment_method,omitempty"`
	PaymentStatus string    `json:"payment_status" validate:"required,oneof=paid unpaid"`
	Amount        float64   `json:"amount" validate:"required"`
	CreatedAt     time.Time `json:"created_at" validate:"required"`
	UpdatedAt     time.Time `json:"updated_at" validate:"required"`
}
