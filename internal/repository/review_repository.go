package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type ReviewRepository interface {
	GetAllByMenuId(ctx context.Context, id uuid.UUID) ([]domain.Review, error)
	Create(ctx context.Context, review domain.Review) error
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Review, error)
	Update(ctx context.Context, id uuid.UUID, review domain.Review) error
	CheckReviewedItem(ctx context.Context, userId, menuId, orderId uuid.UUID) bool
	CountReviewByMenuId(ctx context.Context, menuId uuid.UUID) (int, float64, error)
}
