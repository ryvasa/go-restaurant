package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

type MenuRepositoryImpl struct {
	db DB
}

func NewMenuRepository(db DB) MenuRepository {
	return &MenuRepositoryImpl{db}
}

func (r *MenuRepositoryImpl) GetAll(ctx context.Context) ([]domain.Menu, error) {
	menus := []domain.Menu{}
	rows, err := r.db.QueryContext(ctx, "SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE deleted = false AND deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var menu domain.Menu
		err := rows.Scan(&menu.Id, &menu.Name, &menu.Description, &menu.Price, &menu.Category, &menu.ImageURL, &menu.Rating, &menu.CreatedAt, &menu.UpdatedAt)
		if err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

func (r *MenuRepositoryImpl) Create(ctx context.Context, menu domain.Menu) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO menu (id,name,description,price,category,image_url) VALUES (?,  ?, ?, ?, ?, ?)",
		menu.Id, menu.Name, menu.Description, menu.Price, menu.Category, menu.ImageURL)

	if err != nil {
		return err
	}

	return nil
}

func (r *MenuRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (domain.Menu, error) {
	menu := domain.Menu{}
	err := r.db.QueryRowContext(ctx, "SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&menu.Id, &menu.Name, &menu.Description, &menu.Price, &menu.Category, &menu.ImageURL, &menu.Rating, &menu.CreatedAt, &menu.UpdatedAt)
	if err != nil {
		return domain.Menu{}, err

	}
	return menu, nil
}

func (r *MenuRepositoryImpl) Update(ctx context.Context, id uuid.UUID, menu domain.Menu) error {
	logger.Log.Info(id, menu)
	res, err := r.db.ExecContext(ctx,
		"UPDATE menu SET name = ?, price = ?,  description = ?, category = ?, image_url = ? WHERE id = ?",
		menu.Name,
		menu.Price,
		menu.Description,
		menu.Category,
		menu.ImageURL,
		id,
	)
	if err != nil {
		logger.Log.Error(err)

		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}

	return nil
}

func (r *MenuRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, "UPDATE menu SET deleted = ?, deleted_at = ? WHERE id = ?", true, time.Now(), id)
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}

	return nil
}

func (r *MenuRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx, "UPDATE menu SET deleted = ?, deleted_at = ? WHERE id = ?", false, nil, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

func (r *MenuRepositoryImpl) GetDeletedMenuById(ctx context.Context, id uuid.UUID) (domain.Menu, error) {
	menu := domain.Menu{}
	err := r.db.QueryRowContext(ctx, "SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?", id).Scan(&menu.Id, &menu.Name, &menu.Description, &menu.Price, &menu.Category, &menu.ImageURL, &menu.Rating, &menu.CreatedAt, &menu.UpdatedAt)
	if err != nil {
		return domain.Menu{}, err
	}
	return menu, nil
}

func (r *MenuRepositoryImpl) UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE menu SET rating = ? WHERE id = ?", rating, id)
	if err != nil {
		return err
	}
	return nil
}
