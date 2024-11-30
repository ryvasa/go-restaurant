package domain

import (
	"time"

	"github.com/google/uuid"
)

type Ingredient struct {
	Id          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
}
