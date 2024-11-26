package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type OrderMenuRepository interface {
	Create(tx *sql.Tx, review domain.OrderMenu) (domain.OrderMenu, error)
	GetOneByOrderIdAndMenuId(tx *sql.Tx, orderId, menuId string) (domain.OrderMenu, error)
	GetAllByOrderId(tx *sql.Tx, orderId string) (domain.OrderMenu, error)
}
