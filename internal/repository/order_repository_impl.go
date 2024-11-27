package repository

import (
	"database/sql"
	"log"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
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
		return domain.Order{}, err
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
		return domain.Order{}, err
	}

	if paymentMethod.Valid {
		order.PaymentMethod = &paymentMethod.String
	} else {
		order.PaymentMethod = nil
	}

	return order, nil
}

func (r *OrderRepositoryImpl) UpdateOrderStatus(tx *sql.Tx, id string, order domain.Order) (domain.Order, error) {
	_, err := tx.Exec(
		"UPDATE orders SET status = ? WHERE id = ? AND deleted = false AND deleted_at IS NULL",
		order.Status, id)
	if err != nil {
		return domain.Order{}, err
	}
	return r.GetOneById(tx, id)
}

func (r *OrderRepositoryImpl) UpdatePayment(tx *sql.Tx, id string, order domain.Order) (domain.Order, error) {
	log.Println(id)
	_, err := tx.Exec(
		"UPDATE orders SET payment_status = ?, payment_method = ? WHERE id = ? AND deleted = false AND deleted_at IS NULL",
		order.PaymentStatus, order.PaymentMethod, id)
	if err != nil {
		return domain.Order{}, err
	}
	return r.GetOneById(tx, id)
}
