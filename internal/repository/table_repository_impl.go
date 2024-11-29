package repository

import (
	"context"
	"log"
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
	rows, err := r.db.QueryContext(ctx, "SELECT id,number,capacity,location,status,created_at,updated_at FROM tables WHERE deleted = false AND deleted_at IS NULL")
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
	err := r.db.QueryRowContext(ctx, "SELECT id,number,capacity,location,status,created_at,updated_at FROM tables WHERE id = ? AND deleted = false AND deleted_at IS NULL", id).Scan(&table.Id, &table.Number, &table.Capacity, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt)
	if err != nil {
		return domain.Table{}, err
	}
	return table, nil
}

func (r *TableRepositoryImpl) Create(ctx context.Context, table domain.Table) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO tables (id,number,capacity,location) VALUES (?, ?, ?, ?)",
		table.Id, table.Number, table.Capacity, table.Location)
	if err != nil {
		return err
	}
	return nil
}

func (r *TableRepositoryImpl) Update(ctx context.Context, id uuid.UUID, table domain.Table) error {
	log.Println(table)
	result, err := r.db.ExecContext(ctx, "UPDATE tables SET number = ?, capacity = ?, location = ?, status = ? WHERE id = ?",
		table.Number, table.Capacity, table.Location, table.Status, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	log.Printf("Rows affected: %d", rowsAffected)

	return nil
}

func (r *TableRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tables SET deleted = ?, deleted_at = ? WHERE id = ?", true, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TableRepositoryImpl) GetDeleted(ctx context.Context, id uuid.UUID) (domain.Table, error) {
	table := domain.Table{}
	err := r.db.QueryRowContext(ctx, "SELECT id,number,capacity,location,status,created_at,updated_at FROM tables WHERE deleted = true AND deleted_at IS NOT NULL AND id = ?", id).Scan(&table.Id, &table.Number, &table.Capacity, &table.Location, &table.Status, &table.CreatedAt, &table.UpdatedAt)
	if err != nil {
		return domain.Table{}, err
	}
	return table, nil
}

func (r *TableRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "UPDATE tables SET deleted = ?, deleted_at = ? WHERE id = ?", false, nil, id)
	if err != nil {
		return err
	}
	return nil
}
