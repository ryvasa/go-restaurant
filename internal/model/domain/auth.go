package domain

import (
	"time"

	"github.com/google/uuid"
)

type Auth struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password,omitempty" validate:"required,omitempty,min=6,max=100"`
	Phone     string    `json:"phone" validate:"required"`
	Role      string    `json:"role" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
	Token     string    `json:"token" validate:"required"`
}
