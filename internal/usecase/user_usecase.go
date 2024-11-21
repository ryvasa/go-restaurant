package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/domain"
)

type UserUsecase interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
	Get(ctx context.Context, id string) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
}
