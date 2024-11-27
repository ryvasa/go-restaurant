package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type TableUsecase interface {
	GetAll(ctx context.Context) ([]domain.Table, error)
	GetOneById(ctx context.Context, id string) (domain.Table, error)
	Create(ctx context.Context, req dto.CreateTableRequest) (domain.Table, error)
	Update(ctx context.Context, id string, req dto.UpdateTableRequest) (domain.Table, error)
	Delete(ctx context.Context, id string) error
	Restore(ctx context.Context, id string) (domain.Table, error)
}
