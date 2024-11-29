package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type OrderMenuRepository interface {
	Create(ctx context.Context, review domain.OrderMenu) error
	GetOneByOrderIdAndMenuId(ctx context.Context, orderId, menuId uuid.UUID) (domain.OrderMenu, error)
	GetAllByOrderId(ctx context.Context, orderId uuid.UUID) (domain.OrderMenu, error)
}
