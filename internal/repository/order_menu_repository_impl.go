package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type OrderMenuRepositoryImpl struct {
	db DB
}

func NewOrderMenuRepository(db DB) OrderMenuRepository {
	return &OrderMenuRepositoryImpl{db}
}

func (r *OrderMenuRepositoryImpl) Create(ctx context.Context, orderMenu domain.OrderMenu) error {
	_, err := r.db.ExecContext(ctx,
		"INSERT INTO order_menu (order_id, menu_id, quantity) VALUES (?, ?, ?)",
		orderMenu.OrderId, orderMenu.MenuId, orderMenu.Quantity)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderMenuRepositoryImpl) GetOneByOrderIdAndMenuId(ctx context.Context, orderId, menuId uuid.UUID) (domain.OrderMenu, error) {
	orderMenu := domain.OrderMenu{}

	err := r.db.QueryRowContext(ctx,
		"SELECT  order_id, menu_id, quantity FROM order_menu WHERE order_id = ? AND menu_id = ?",
		orderId, menuId).Scan(
		&orderMenu.OrderId,
		&orderMenu.MenuId,
		&orderMenu.Quantity,
	)

	if err != nil {
		return domain.OrderMenu{}, err
	}
	return orderMenu, nil
}

func (r *OrderMenuRepositoryImpl) GetAllByOrderId(ctx context.Context, id uuid.UUID) (domain.OrderMenu, error) {
	orderMenu := domain.OrderMenu{}

	err := r.db.QueryRowContext(ctx,
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
