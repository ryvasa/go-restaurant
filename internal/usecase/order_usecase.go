package usecase

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type OrderUsecase interface {
	Create(ctx context.Context, req dto.CreateOrderDto, userId string) (domain.Order, error)
	GetOneById(ctx context.Context, id string) (domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id string, req dto.UpdateOrderStatusDto) (domain.Order, error)
	UpdatePayment(ctx context.Context, id string, req dto.UpdatePaymentDto) (domain.Order, error)
}
