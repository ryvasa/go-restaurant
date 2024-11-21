package domain

import (
	"context"
	"net/http"
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
	GetAll(ctx context.Context) ([]Menu, error)
	Create(ctx context.Context, menu Menu) (Menu, error)
	Get(ctx context.Context, id string) (Menu, error)
	Update(ctx context.Context, menu Menu) (Menu, error)
	Delete(ctx context.Context, id string) error
}

type MenuUsecase interface {
	GetAll(ctx context.Context) ([]Menu, error)
	Create(ctx context.Context, menu Menu) (Menu, error)
	Get(ctx context.Context, id string) (Menu, error)
	Update(ctx context.Context, menu Menu) (Menu, error)
	Delete(ctx context.Context, id string) error
}

type MenuHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}
