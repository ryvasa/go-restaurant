package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type TableUsecase interface {
	GetAll(ctx context.Context) ([]domain.Table, error)
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Table, error)
	Create(ctx context.Context, req dto.CreateTableRequest) (domain.Table, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateTableRequest) (domain.Table, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Restore(ctx context.Context, id uuid.UUID) (domain.Table, error)
}
