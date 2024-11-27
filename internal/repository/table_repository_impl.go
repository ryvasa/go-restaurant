package repository

import (
	"database/sql"
	"time"

	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type TableRepositoryImpl struct {
}

func NewTableRepository() TableRepository {
	return &TableRepositoryImpl{}
}

func (r *TableRepositoryImpl) GetAll(tx *sql.Tx) ([]domain.Table, error) {
	tables := []domain.Table{}
	rows, err := tx.Query("SELECT id,number,capacity,location,status,created_at,updated_at FROM table WHERE deleted = false AND deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var table domain.Table
		err := rows.Scan(&table.Id, &table.Number, &table.Capacity, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func (r *TableRepositoryImpl) GetOneById(tx *sql.Tx, id string) (domain.Table, error) {
	table := domain.Table{}
	err := tx.QueryRow("SELECT id,number,capacity,location,status,created_at,updated_at FROM table WHERE id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&table.Id, &table.Number, &table.Capacity, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt)
	if err != nil {
		return domain.Table{}, err
	}
	return table, nil
}

func (r *TableRepositoryImpl) Create(tx *sql.Tx, table domain.Table) error {
	_, err := tx.Exec("INSERT INTO table (id,number,capacity,location,status) VALUES (?, ?, ?, ?, ?)",
		table.Id, table.Number, table.Capacity, table.Location, table.Status)
	if err != nil {
		return err
	}
	return nil
}

func (r *TableRepositoryImpl) Update(tx *sql.Tx, table domain.Table) error {
	_, err := tx.Exec("UPDATE table SET number = ?, capacity = ?, location = ?, status = ? WHERE id = ?",
		table.Number, table.Capacity, table.Location, table.Status, table.Id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TableRepositoryImpl) Delete(tx *sql.Tx, id string) error {
	_, err := tx.Exec("UPDATE table SET deleted = ?, deleted_at = ? WHERE id = ?", true, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TableRepositoryImpl) GetDeleted(tx *sql.Tx, id string) (domain.Table, error) {
	table := domain.Table{}
	err := tx.QueryRow("SELECT id,number,capacity,location,status,created_at,updated_at FROM table WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?", id).Scan(&table.Id, &table.Number, &table.Capacity, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt)
	if err != nil {
		return domain.Table{}, err
	}
	return table, nil
}

func (r *TableRepositoryImpl) Restore(tx *sql.Tx, id string) error {
	_, err := tx.Exec("UPDATE table SET deleted = ?, deleted_at = ? WHERE id = ?", false, nil, id)
	if err != nil {
		return err
	}
	return nil
}
