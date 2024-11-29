package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
)

type OrderUsecase interface {
	Create(ctx context.Context, req dto.CreateOrderDto, userId uuid.UUID) (domain.Order, error)
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, req dto.UpdateOrderStatusDto) (domain.Order, error)
	UpdatePayment(ctx context.Context, id uuid.UUID, req dto.UpdatePaymentDto) (domain.Order, error)
}
