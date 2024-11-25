package repository

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type OrderMenuRepository interface {
	Create(ctx context.Context, review domain.OrderMenu) (domain.OrderMenu, error)
	GetOneByOrderIdAndMenuId(ctx context.Context, orderId, menuId string) (domain.OrderMenu, error)
	GetAllByOrderId(ctx context.Context, orderId string) (domain.OrderMenu, error)
}
