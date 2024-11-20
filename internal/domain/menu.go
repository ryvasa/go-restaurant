package domain

import (
	"time"

	"github.com/google/uuid"
)

type Menu struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MenuRepository interface {
	GetAll() ([]Menu, error)
	Create(menu Menu) (Menu, error)
	// GetByID(id int) (Menu, error)
	// Update(menu Menu) error
	// Delete(id int) error
}

type MenuUsecase interface {
	GetAll() ([]Menu, error)
	Create(menu Menu) (Menu, error)
	// GetByID(id int) (Menu, error)
	// Update(menu Menu) error
	// Delete(id int) error
}
