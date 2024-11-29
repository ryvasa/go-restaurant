package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type ReservationUsecase interface {
	GetAll(ctx context.Context) ([]domain.Reservation, error)
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Reservation, error)
	Create(ctx context.Context, req dto.CreateReservationRequest) (domain.Reservation, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateReservationRequest) (domain.Reservation, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) (domain.Reservation, error)
}
