package domain

import (
	"time"

	"github.com/google/uuid"
)

type Menu struct {
	Id          uuid.UUID `json:"id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	Category    string    `json:"category" validate:"required"`
	ImageURL    string    `json:"image_url" validate:"required"`
	CreatedAt   time.Time `json:"created_at" validate:"required"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required"`
	Rating      int       `json:"rating" validate:"required"`
}
