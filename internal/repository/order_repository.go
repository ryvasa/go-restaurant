package repository

import (
	"context"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type OrderRepository interface {
	Create(ctx context.Context, review domain.Order) (domain.Order, error)
	GetOneById(ctx context.Context, id string) (domain.Order, error)
}
