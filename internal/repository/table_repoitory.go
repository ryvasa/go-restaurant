package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type TableRepository interface {
	GetAll(ctx context.Context) ([]domain.Table, error)
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Table, error)
	Create(ctx context.Context, table domain.Table) error
	Update(ctx context.Context, id uuid.UUID, table domain.Table) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetDeleted(ctx context.Context, id uuid.UUID) (domain.Table, error)
	Restore(ctx context.Context, id uuid.UUID) error
}
