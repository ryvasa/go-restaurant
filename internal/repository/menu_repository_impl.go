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
	query := `SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE deleted = false AND deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query)
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
	query := `INSERT INTO menu (id,name,description,price,category,image_url) VALUES (?,  ?, ?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query,
		menu.Id, menu.Name, menu.Description, menu.Price, menu.Category, menu.ImageURL)

	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}

	return nil
}

func (r *MenuRepositoryImpl) Get(ctx context.Context, id uuid.UUID) (domain.Menu, error) {
	menu := domain.Menu{}
	query := `SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE id = ? AND deleted = false AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&menu.Id, &menu.Name, &menu.Description, &menu.Price, &menu.Category, &menu.ImageURL, &menu.Rating, &menu.CreatedAt, &menu.UpdatedAt)
	if err != nil {
		return domain.Menu{}, err

	}
	return menu, nil
}

func (r *MenuRepositoryImpl) Update(ctx context.Context, id uuid.UUID, menu domain.Menu) error {
	query := `UPDATE menu SET name = ?, price = ?,  description = ?, category = ?, image_url = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx,
		query,
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
	query := `UPDATE menu SET deleted = ?, deleted_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, true, time.Now(), id)
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
	query := `UPDATE menu SET deleted = ?, deleted_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, false, nil, id)
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
	query := `SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&menu.Id, &menu.Name, &menu.Description, &menu.Price, &menu.Category, &menu.ImageURL, &menu.Rating, &menu.CreatedAt, &menu.UpdatedAt)
	if err != nil {
		return domain.Menu{}, err
	}
	return menu, nil
}

func (r *MenuRepositoryImpl) UpdateRating(ctx context.Context, id uuid.UUID, rating float64) error {
	query := `UPDATE menu SET rating = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, rating, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}
