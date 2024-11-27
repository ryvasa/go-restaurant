package repository

import (
	"database/sql"
	"time"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type MenuRepositoryImpl struct {
}

func NewMenuRepository() MenuRepository {
	return &MenuRepositoryImpl{}
}

func (r *MenuRepositoryImpl) GetAll(tx *sql.Tx) ([]domain.Menu, error) {
	menus := []domain.Menu{}
	rows, err := tx.Query("SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE deleted = false AND deleted_at IS NULL")
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

func (r *MenuRepositoryImpl) Create(tx *sql.Tx, menu domain.Menu) (domain.Menu, error) {
	_, err := tx.Exec("INSERT INTO menu (id,name,description,price,category,image_url) VALUES (?,  ?, ?, ?, ?, ?)",
		menu.Id, menu.Name, menu.Description, menu.Price, menu.Category, menu.ImageURL)

	if err != nil {
		return domain.Menu{}, err
	}

	createdMenu, _ := r.Get(tx, menu.Id.String())

	return createdMenu, nil
}

func (r *MenuRepositoryImpl) Get(tx *sql.Tx, id string) (domain.Menu, error) {
	menu := domain.Menu{}
	err := tx.QueryRow("SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&menu.Id, &menu.Name, &menu.Description, &menu.Price, &menu.Category, &menu.ImageURL, &menu.Rating, &menu.CreatedAt, &menu.UpdatedAt)
	if err != nil {
		return domain.Menu{}, err

	}
	return menu, nil
}

func (r *MenuRepositoryImpl) Update(tx *sql.Tx, menu domain.Menu) (domain.Menu, error) {
	// Ambil data menu yang ada
	existingMenu, err := r.Get(tx, menu.Id.String())
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
	_, err = tx.Exec(
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

	return r.Get(tx, existingMenu.Id.String())
}

func (r *MenuRepositoryImpl) Delete(tx *sql.Tx, id string) error {
	_, err := tx.Exec("UPDATE menu SET deleted = ?, deleted_at = ? WHERE id = ?", true, time.Now(), id)
	if err != nil {
		return err
		// logger.Log.WithError(err).Error("Error failed to delete menu")
		// return utils.NewInternalError("Failed to delete menu")
	}

	return nil
}

func (r *MenuRepositoryImpl) Restore(tx *sql.Tx, id string) (domain.Menu, error) {
	_, err := tx.Exec("UPDATE menu SET deleted = ?, deleted_at = ? WHERE id = ?", false, nil, id)
	if err != nil {
		return domain.Menu{}, err
		// logger.Log.WithError(err).Error("Error failed to restore menu")
		// return domain.Menu{}, utils.NewInternalError("Failed to restore menu")
	}
	menu, _ := r.Get(tx, id)
	return menu, nil
}

func (r *MenuRepositoryImpl) GetDeletedMenuById(tx *sql.Tx, id string) (domain.Menu, error) {
	menus := domain.Menu{}
	err := tx.QueryRow("SELECT id,name,description,price,category,image_url,rating,created_at,updated_at FROM menu WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?", id).Scan(&menus)
	if err != nil {
		return domain.Menu{}, err
	}
	return menus, nil
}

func (r *MenuRepositoryImpl) UpdateRating(tx *sql.Tx, id string, rating float64) error {
	_, err := tx.Exec("UPDATE menu SET rating = ? WHERE id = ?", rating, id)
	if err != nil {
		return err
	}
	return nil
}
