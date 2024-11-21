package domain

import (
	"time"

	"github.com/google/uuid"
)

type Menu struct {
	ID          uuid.UUID `json:"id"`
	Restaurant  uuid.UUID `json:"restaurant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
