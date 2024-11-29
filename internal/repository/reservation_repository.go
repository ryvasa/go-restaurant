package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type ReservationRepository interface {
	GetAll(ctx context.Context) ([]domain.Reservation, error)
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Reservation, error)
	GetOneByTableId(ctx context.Context, tableId uuid.UUID) (domain.Reservation, error)
	Create(ctx context.Context, reservation domain.Reservation) error
	Update(ctx context.Context, id uuid.UUID, reservation domain.Reservation) error
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	GetDeleted(ctx context.Context, id uuid.UUID) (domain.Reservation, error)
}
