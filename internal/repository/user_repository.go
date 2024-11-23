package repository

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/domain"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
	Get(ctx context.Context, id string) (domain.User, error)
	Update(ctx context.Context, user domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}
