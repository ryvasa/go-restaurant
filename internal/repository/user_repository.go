package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	Create(ctx context.Context, user domain.User) error
	Get(ctx context.Context, id uuid.UUID) (domain.User, error)
	Update(ctx context.Context, id uuid.UUID, user domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) error
	GetDeletedUserById(ctx context.Context, id uuid.UUID) (domain.User, error)
}
