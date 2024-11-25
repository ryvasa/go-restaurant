package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type OrderUsecase interface {
	Create(ctx context.Context, req dto.CreateOrderDto, userId string) (domain.Order, error)
	GetOneById(ctx context.Context, id string) (domain.Order, error)
}
