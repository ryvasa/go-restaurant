package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, order domain.Order) error
	GetOneById(ctx context.Context, id uuid.UUID) (domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, order domain.Order) error
	UpdatePayment(ctx context.Context, id uuid.UUID, order domain.Order) error
}
