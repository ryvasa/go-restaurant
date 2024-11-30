package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
)

type TableRepositoryImpl struct {
	db DB
}

func NewTableRepository(db DB) TableRepository {
	return &TableRepositoryImpl{db}
}

func (r *TableRepositoryImpl) GetAll(ctx context.Context) ([]domain.Table, error) {
	tables := []domain.Table{}
	query := `SELECT id,number,capacity,location,status,created_at,updated_at FROM tables WHERE deleted = false AND deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query)
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

func (r *TableRepositoryImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Table, error) {
	table := domain.Table{}
	query := `SELECT id,number,capacity,location,status,created_at,updated_at FROM tables WHERE id = ? AND deleted = false AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&table.Id, &table.Number, &table.Capacity, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt)
	if err != nil {
		return domain.Table{}, err
	}
	return table, nil
}

func (r *TableRepositoryImpl) Create(ctx context.Context, table domain.Table) error {
	query := `INSERT INTO tables (id,number,capacity,location) VALUES (?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, table.Id, table.Number, table.Capacity, table.Location)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

func (r *TableRepositoryImpl) Update(ctx context.Context, id uuid.UUID, table domain.Table) error {
	query := `UPDATE tables SET number = ?, capacity = ?, location = ?, status = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, table.Number, table.Capacity, table.Location, table.Status, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No rows affected")
	}
	return nil
}

func (r *TableRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tables SET deleted = ?, deleted_at = ? WHERE id = ?`
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

func (r *TableRepositoryImpl) GetDeleted(ctx context.Context, id uuid.UUID) (domain.Table, error) {
	table := domain.Table{}
	query := `SELECT id,number,capacity,location,status,created_at,updated_at FROM tables WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&table.Id, &table.Number, &table.Capacity, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt)
	if err != nil {
		return domain.Table{}, err
	}
	return table, nil
}

func (r *TableRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tables SET deleted = ?, deleted_at = ? WHERE id = ?`
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
