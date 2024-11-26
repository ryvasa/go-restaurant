package repository

import (
	"database/sql"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type OrderMenuRepositoryImpl struct {
}

func NewOrderMenuRepository() OrderMenuRepository {
	return &OrderMenuRepositoryImpl{}
}

func (r *OrderMenuRepositoryImpl) Create(tx *sql.Tx, orderMenu domain.OrderMenu) (domain.OrderMenu, error) {
	_, err := tx.Exec(
		"INSERT INTO order_menu (order_id, menu_id, quantity) VALUES (?, ?, ?)",
		orderMenu.OrderId, orderMenu.MenuId, orderMenu.Quantity)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create order menu")
		return domain.OrderMenu{}, utils.NewInternalError("Failed to create order menu")
	}
	return r.GetOneByOrderIdAndMenuId(tx, orderMenu.OrderId.String(), orderMenu.MenuId.String())
}

func (r *OrderMenuRepositoryImpl) GetOneByOrderIdAndMenuId(tx *sql.Tx, orderId, menuId string) (domain.OrderMenu, error) {
	orderMenu := domain.OrderMenu{}

	err := tx.QueryRow(
		"SELECT  order_id, menu_id, quantity FROM order_menu WHERE order_id = ? AND menu_id = ?",
		orderId, menuId).Scan(
		&orderMenu.OrderId,
		&orderMenu.MenuId,
		&orderMenu.Quantity,
	)

	if err != nil {
		logger.Log.WithError(err).Error("Error order menu not found")
		return domain.OrderMenu{}, utils.NewNotFoundError("Order menu not found")
	}
	return orderMenu, nil
}

func (r *OrderMenuRepositoryImpl) GetAllByOrderId(tx *sql.Tx, id string) (domain.OrderMenu, error) {
	orderMenu := domain.OrderMenu{}

	err := tx.QueryRow(
		"SELECT  order_id, menu_id, quantity FROM order_menu WHERE order_id = ? AND deleted = false AND deleted_at IS NULL",
		id).Scan(
		&orderMenu.OrderId,
		&orderMenu.MenuId,
		&orderMenu.Quantity,
	)

	if err != nil {
		return orderMenu, err
	}

	return orderMenu, nil
}
