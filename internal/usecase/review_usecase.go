package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type ReviewUsecase interface {
	GetAllByMenuId(ctx context.Context, id string) ([]domain.Review, error)
	Create(ctx context.Context, req dto.CreateReviewRequest, userId string) (domain.Review, error)
	GetOneById(ctx context.Context, id string) (domain.Review, error)
	Update(ctx context.Context, id, userId string, req dto.UpdateReviewRequest) (domain.Review, error)
}
