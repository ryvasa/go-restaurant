package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type ReviewUsecase interface {
	GetAllByMenuId(ctx context.Context, id uuid.UUID) ([]domain.Review, error)
	Create(ctx context.Context, req dto.CreateReviewRequest, userId uuid.UUID) (domain.Review, error)
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Review, error)
	Update(ctx context.Context, id, userId uuid.UUID, req dto.UpdateReviewRequest) (domain.Review, error)
}
