package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
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
	rows, err := r.db.QueryContext(ctx, "SELECT id,restaurant_id,name,description,price,category,image_url,created_at,updated_at FROM menu WHERE deleted = false AND deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var menu domain.Menu
		err := rows.Scan(&menu.ID, &menu.Restaurant, &menu.Name, &menu.Description, &menu.Price, &menu.Category, &menu.ImageURL, &menu.CreatedAt, &menu.UpdatedAt)
		if err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

func (r *MenuRepositoryImpl) Create(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	_, err := r.db.ExecContext(ctx, "INSERT INTO menu (id,restaurant_id,name,description,price,category,image_url) VALUES (?, ?, ?, ?, ?, ?, ?)",
		menu.ID, menu.Restaurant, menu.Name, menu.Description, menu.Price, menu.Category, menu.ImageURL)

	if err != nil {
		return domain.Menu{}, err
	}

	createdMenu, err := r.Get(ctx, menu.ID.String())
	if err != nil {
		return domain.Menu{}, err
	}

	return createdMenu, nil
}

func (r *MenuRepositoryImpl) Get(ctx context.Context, id string) (domain.Menu, error) {
	menu := domain.Menu{}
	err := r.db.QueryRowContext(ctx, "SELECT id,restaurant_id,name,description,price,category,image_url,created_at,updated_at FROM menu WHERE id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&menu.ID, &menu.Restaurant, &menu.Name, &menu.Description, &menu.Price, &menu.Category, &menu.ImageURL, &menu.CreatedAt, &menu.UpdatedAt)
	if err != nil {
		return menu, err
	}
	return menu, nil
}

func (r *MenuRepositoryImpl) Update(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	// Ambil data menu yang ada
	existingMenu, err := r.Get(ctx, menu.ID.String())
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
	if menu.Restaurant != uuid.Nil {
		existingMenu.Restaurant = menu.Restaurant
	}

	existingMenu.UpdatedAt = time.Now()

	// Eksekusi query update
	_, err = r.db.ExecContext(ctx,
		"UPDATE menu SET name = ?, price = ?, updated_at = ?, restaurant_id = ?, description = ?, category = ?, image_url = ? WHERE id = ?",
		existingMenu.Name,
		existingMenu.Price,
		existingMenu.UpdatedAt,
		existingMenu.Restaurant,
		existingMenu.Description,
		existingMenu.Category,
		existingMenu.ImageURL,
		existingMenu.ID,
	)
	if err != nil {
		return domain.Menu{}, err
	}

	return r.Get(ctx, existingMenu.ID.String())
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
	rows, err := r.db.QueryContext(ctx, "SELECT id,restaurant_id,name,description,price,category,image_url,created_at,updated_at FROM menu WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var menu domain.Menu
		err := rows.Scan(&menu.ID, &menu.Restaurant, &menu.Name, &menu.Description, &menu.Price, &menu.Category, &menu.ImageURL, &menu.CreatedAt, &menu.UpdatedAt)
		if err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	return menus, nil
}
