package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/domain"
)

type menuRepository struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) domain.MenuRepository {
	return &menuRepository{db}
}

func (r *menuRepository) GetAll(ctx context.Context) ([]domain.Menu, error) {
	menus := []domain.Menu{}
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, price, created_at, updated_at FROM menu")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var menu domain.Menu
		err := rows.Scan(&menu.ID, &menu.Name, &menu.Price, &menu.CreatedAt, &menu.UpdatedAt)
		if err != nil {
			return nil, err
		}
		menus = append(menus, menu)
	}
	return menus, nil
}

func (r *menuRepository) Create(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	menu.ID = uuid.New()

	_, err := r.db.ExecContext(ctx, "INSERT INTO menu (id, name, price) VALUES (?, ?, ?)",
		menu.ID, menu.Name, menu.Price)
	if err != nil {
		return domain.Menu{}, err
	}

	createdMenu, err := r.Get(ctx, menu.ID.String())
	if err != nil {
		return domain.Menu{}, err
	}

	return createdMenu, nil
}

func (r *menuRepository) Get(ctx context.Context, id string) (domain.Menu, error) {
	menu := domain.Menu{}
	err := r.db.QueryRowContext(ctx, "SELECT id, name, price, created_at, updated_at FROM menu WHERE id = ?", id).Scan(&menu.ID, &menu.Name, &menu.Price, &menu.CreatedAt, &menu.UpdatedAt)
	if err != nil {
		return menu, err
	}
	return menu, nil
}

// func (r *menuRepository) Create(menu domain.Menu) (domain.Menu, error) {
// 	menu.ID = uuid.New()

// 	_, err := r.db.Exec("INSERT INTO menu (id, name, price) VALUES (?, ?, ?)",
// 		menu.ID, menu.Name, menu.Price)
// 	if err != nil {
// 		return domain.Menu{}, err
// 	}

// 	createdMenu, err := r.Get(menu.ID.String())
// 	if err != nil {
// 		return domain.Menu{}, err
// 	}

// 	return createdMenu, nil
// }

func (r *menuRepository) Update(ctx context.Context, menu domain.Menu) (domain.Menu, error) {
	menu.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, "UPDATE menu SET name = ?, price = ?, updated_at = ? WHERE id = ?",
		menu.Name, menu.Price, menu.UpdatedAt, menu.ID)
	if err != nil {
		return domain.Menu{}, err
	}

	updatedMenu, err := r.Get(ctx, menu.ID.String())
	if err != nil {
		return domain.Menu{}, err
	}

	return updatedMenu, nil
}

func (r *menuRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM menu WHERE id = ?", id)
	if err != nil {
		return err
	}

	// Cek apakah ada baris yang terhapus
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// Jika tidak ada baris yang terhapus, return error
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
