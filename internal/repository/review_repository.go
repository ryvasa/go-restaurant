package repository

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type ReviewRepository interface {
	GetAllByMenuId(ctx context.Context, id string) ([]domain.Review, error)
	Create(ctx context.Context, review domain.Review) error
	GetOneById(ctx context.Context, id string) (domain.Review, error)
	Update(ctx context.Context, review domain.Review) error
}
