package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type OrderRepository interface {
	Create(tx *sql.Tx, order domain.Order) (domain.Order, error)
	GetOneById(tx *sql.Tx, id string) (domain.Order, error)
	UpdateOrderStatus(tx *sql.Tx, id string, order domain.Order) (domain.Order, error)
	UpdatePayment(tx *sql.Tx, id string, order domain.Order) (domain.Order, error)
}
