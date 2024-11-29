package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type UserUsecase interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, req dto.CreateUserRequest) (domain.User, error)
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id, authId uuid.UUID, req dto.UpdateUserRequest) (domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) (domain.User, error)
}
