package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type UserUsecase interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, req dto.CreateUserRequest) (domain.User, error)
	Get(ctx context.Context, id string) (domain.User, error)
	Update(ctx context.Context, id string, req dto.UpdateUserRequest) (domain.User, error)
}
