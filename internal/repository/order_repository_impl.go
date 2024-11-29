package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type OrderRepositoryImpl struct {
	db DB
}

func NewOrderRepository(db DB) OrderRepository {
	return &OrderRepositoryImpl{db}
}
func (r *OrderRepositoryImpl) Create(ctx context.Context, order domain.Order) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO orders (id, user_id, amount) VALUES (?, ?, ?)",
		order.Id, order.UserId, order.Amount)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepositoryImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Order, error) {
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
		return domain.Order{}, err
	}

	if paymentMethod.Valid {
		order.PaymentMethod = &paymentMethod.String
	} else {
		order.PaymentMethod = nil
	}

	return order, nil
}

func (r *OrderRepositoryImpl) UpdateOrderStatus(ctx context.Context, id uuid.UUID, order domain.Order) error {
	res, err := r.db.ExecContext(ctx,
		"UPDATE orders SET status = ? WHERE id = ? AND deleted = false AND deleted_at IS NULL",
		order.Status, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

func (r *OrderRepositoryImpl) UpdatePayment(ctx context.Context, id uuid.UUID, order domain.Order) error {
	res, err := r.db.ExecContext(ctx,
		"UPDATE orders SET payment_status = ?, payment_method = ? WHERE id = ? AND deleted = false AND deleted_at IS NULL",
		order.PaymentStatus, order.PaymentMethod, id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}
