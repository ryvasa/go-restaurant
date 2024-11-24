package domain

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	Id        uuid.UUID `json:"id" validate:"required"`
	Rating    int       `json:"rating" validate:"required"`
	Comment   string    `json:"comment" validate:"required"`
	UserId    uuid.UUID `json:"user_id" validate:"required"`
	MenuId    uuid.UUID `json:"menu_id" validate:"required"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}
