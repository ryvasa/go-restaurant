package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type MenuRepository interface {
	GetAll(ctx context.Context) ([]domain.Menu, error)
	Create(ctx context.Context, menu domain.Menu) error
	Get(ctx context.Context, id uuid.UUID) (domain.Menu, error)
	Update(ctx context.Context, id uuid.UUID, menu domain.Menu) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	GetDeletedMenuById(ctx context.Context, id uuid.UUID) (domain.Menu, error)
	UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error
}
