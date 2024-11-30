package repository

import (
	"context"
	"fmt"

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
	query := `INSERT INTO order_menu (order_id, menu_id, quantity) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, orderMenu.OrderId, orderMenu.MenuId, orderMenu.Quantity)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

func (r *OrderMenuRepositoryImpl) GetOneByOrderIdAndMenuId(ctx context.Context, orderId, menuId uuid.UUID) (domain.OrderMenu, error) {
	orderMenu := domain.OrderMenu{}
	query := `SELECT  order_id, menu_id, quantity FROM order_menu WHERE order_id = ? AND menu_id = ?`
	err := r.db.QueryRowContext(ctx, query, orderId, menuId).Scan(&orderMenu.OrderId, &orderMenu.MenuId, &orderMenu.Quantity)

	if err != nil {
		return domain.OrderMenu{}, err
	}
	return orderMenu, nil
}

func (r *OrderMenuRepositoryImpl) GetAllByOrderId(ctx context.Context, id uuid.UUID) (domain.OrderMenu, error) {
	orderMenu := domain.OrderMenu{}
	query := `SELECT  order_id, menu_id, quantity FROM order_menu WHERE order_id = ? AND deleted = false AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&orderMenu.OrderId, &orderMenu.MenuId, &orderMenu.Quantity)

	if err != nil {
		return orderMenu, err
	}

	return orderMenu, nil
}
