package repository

import (
	"context"
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type OrderRepositoryImpl struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &OrderRepositoryImpl{db}
}
func (r *OrderRepositoryImpl) Create(ctx context.Context, order domain.Order) (domain.Order, error) {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO orders (id, user_id, amount) VALUES (?, ?, ?)",
		order.Id, order.UserId, order.Amount)
	if err != nil {
		return domain.Order{}, err
	}
	return r.GetOneById(ctx, order.Id.String())
}
func (r *OrderRepositoryImpl) GetOneById(ctx context.Context, id string) (domain.Order, error) {
	order := domain.Order{}
	var paymentMethod sql.NullString

	err := r.db.QueryRowContext(ctx,
		"SELECT id, amount, payment_method, payment_status, status, user_id, created_at, updated_at FROM orders WHERE id = ? AND deleted = false AND deleted_at IS NULL",
		id).Scan(
		&order.Id,
		&order.Amount,
		&paymentMethod,
		&order.PaymentStatus,
		&order.Status,
		&order.UserId,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return order, err
	}

	if paymentMethod.Valid {
		order.PaymentMethod = &paymentMethod.String
	} else {
		order.PaymentMethod = nil
	}

	return order, nil
}
