package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type ReviewUsecase interface {
	GetAllByMenuId(ctx context.Context, id string) ([]domain.Review, error)
	Create(ctx context.Context, req dto.CreateReviewRequest) (domain.Review, error)
	GetOneById(ctx context.Context, id string) (domain.Review, error)
	Update(ctx context.Context, id string, req dto.UpdateReviewRequest) (domain.Review, error)
}
