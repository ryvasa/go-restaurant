package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type MenuRepositoryImpl struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) MenuRepository {
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

func (r *MenuRepositoryImpl) Create(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	_, err := r.db.ExecContext(ctx, "INSERT INTO menu (id,name,description,price,category,image_url) VALUES (?,  ?, ?, ?, ?, ?)",
		menu.Id, menu.Name, menu.Description, menu.Price, menu.Category, menu.ImageURL)

	if err != nil {
		return domain.Menu{}, err
	}

	createdMenu, err := r.Get(ctx, menu.Id.String())
	if err != nil {
		return domain.Menu{}, err
	}

	return createdMenu, nil
}

func (r *MenuRepositoryImpl) Get(ctx context.Context, id string) (domain.Menu, error) {
	menu := domain.Menu{}
	err := r.db.QueryRowContext(ctx, "SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&menu.Id, &menu.Name, &menu.Description, &menu.Price, &menu.Category, &menu.ImageURL, &menu.Rating, &menu.CreatedAt, &menu.UpdatedAt)
	if err != nil {
		return menu, err
	}
	return menu, nil
}

func (r *MenuRepositoryImpl) Update(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	// Ambil data menu yang ada
	existingMenu, err := r.Get(ctx, menu.Id.String())
	if err != nil {
		return domain.Menu{}, err
	}

	// Update hanya field yang tidak kosong
	if menu.Name != "" {
		existingMenu.Name = menu.Name
	}
	if menu.Description != "" {
		existingMenu.Description = menu.Description
	}
	if menu.Price > 0 {
		existingMenu.Price = menu.Price
	}
	if menu.Category != "" {
		existingMenu.Category = menu.Category
	}
	if menu.ImageURL != "" {
		existingMenu.ImageURL = menu.ImageURL
	}

	existingMenu.UpdatedAt = time.Now()

	// Eksekusi query update
	_, err = r.db.ExecContext(ctx,
		"UPDATE menu SET name = ?, price = ?, updated_at = ?, description = ?, category = ?, image_url = ? WHERE id = ?",
		existingMenu.Name,
		existingMenu.Price,
		existingMenu.UpdatedAt,
		existingMenu.Description,
		existingMenu.Category,
		existingMenu.ImageURL,
		existingMenu.Id,
	)
	if err != nil {
		return domain.Menu{}, err
	}

	return r.Get(ctx, existingMenu.Id.String())
}

func (r *MenuRepositoryImpl) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE menu SET deleted = ?, deleted_at = ? WHERE id = ?", true, time.Now(), id)
	if err != nil {
		return err
	}

	return nil
}

func (r *MenuRepositoryImpl) Restore(ctx context.Context, id string) (domain.Menu, error) {
	_, err := r.db.ExecContext(ctx, "UPDATE menu SET deleted = ?, deleted_at = ? WHERE id = ?", false, nil, id)
	if err != nil {
		return domain.Menu{}, err
	}
	menu, _ := r.Get(ctx, id)
	return menu, nil
}

func (r *MenuRepositoryImpl) GetDeletedMenuById(ctx context.Context, id string) ([]domain.Menu, error) {
	menus := []domain.Menu{}
	rows, err := r.db.QueryContext(ctx, "SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?", id)
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
