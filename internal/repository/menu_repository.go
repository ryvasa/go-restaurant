package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/domain"
)

type menuRepository struct {
	db *sql.DB
}

func NewMenuRepository(db *sql.DB) domain.MenuRepository {
	return &menuRepository{db}
}

func (r *menuRepository) GetAll() ([]domain.Menu, error) {
	menus := []domain.Menu{}
	rows, err := r.db.Query("SELECT id, name, price, created_at, updated_at FROM menu")
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

func (r *menuRepository) Create(menu domain.Menu) (domain.Menu, error) {
	menu.ID = uuid.New()

	_, err := r.db.Exec("INSERT INTO menu (id, name, price) VALUES (?, ?, ?)",
		menu.ID, menu.Name, menu.Price)
	if err != nil {
		return domain.Menu{}, err
	}

	// Ambil data yang baru dibuat
	var createdMenu domain.Menu
	err = r.db.QueryRow("SELECT id, name, price, created_at, updated_at FROM menu WHERE id = ?",
		menu.ID).Scan(&createdMenu.ID, &createdMenu.Name, &createdMenu.Price,
		&createdMenu.CreatedAt, &createdMenu.UpdatedAt)
	if err != nil {
		return domain.Menu{}, err
	}

	return createdMenu, nil
}

// func (r *menuRepository) GetByID(id int) (domain.Menu, error) {
// 	menu := domain.Menu{}
// 	err := r.db.QueryRow("SELECT id, name, price, created_at, updated_at FROM menus WHERE id = ?", id).Scan(&menu.ID, &menu.Name, &menu.Price, &menu.CreatedAt, &menu.UpdatedAt)
// 	if err != nil {
// 		return menu, err
// 	}
// 	return menu, nil
// }

// func (r *menuRepository) Update(menu domain.Menu) error {
// 	_, err := r.db.Exec("UPDATE menus SET name = ?, price = ?, updated_at = ? WHERE id = ?", menu.Name, menu.Price, menu.UpdatedAt, menu.ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (r *menuRepository) Delete(id int) error {
// 	_, err := r.db.Exec("DELETE FROM menus WHERE id = ?", id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
