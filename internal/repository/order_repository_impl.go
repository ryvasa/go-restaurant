package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type OrderRepositoryImpl struct {
}

func NewOrderRepository() OrderRepository {
	return &OrderRepositoryImpl{}
}
func (r *OrderRepositoryImpl) Create(tx *sql.Tx, order domain.Order) (domain.Order, error) {
	_, err := tx.Exec(
		"INSERT INTO orders (id, user_id, amount) VALUES (?, ?, ?)",
		order.Id, order.UserId, order.Amount)
	if err != nil {
		logger.Log.WithError(err).Error("Error create order")
		return domain.Order{}, utils.NewInternalError("Failed to create order")
	}
	return r.GetOneById(tx, order.Id.String())
}

func (r *OrderRepositoryImpl) GetOneById(tx *sql.Tx, id string) (domain.Order, error) {
	order := domain.Order{}
	var paymentMethod sql.NullString

	err := tx.QueryRow(
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
		logger.Log.WithError(err).Error("Error order not found")
		return domain.Order{}, utils.NewNotFoundError("Order not found")
	}

	if paymentMethod.Valid {
		order.PaymentMethod = &paymentMethod.String
	} else {
		order.PaymentMethod = nil
	}

	return order, nil
}
