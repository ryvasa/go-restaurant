package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

type InventoryRepositoryImpl struct {
	db DB
}

func NewInventoryRepository(db DB) InventoryRepository {
	return &InventoryRepositoryImpl{
		db: db,
	}
}

func (r *InventoryRepositoryImpl) Create(ctx context.Context, inventory domain.Inventory) error {
	query := `INSERT INTO inventory (id, ingredient_id, quantity) VALUES (?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, inventory.Id, inventory.IngredientId, inventory.Quantity)
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

func (r *InventoryRepositoryImpl) GetOneByIngredientId(ctx context.Context, ingredientId uuid.UUID) (domain.Inventory, error) {
	inventory := domain.Inventory{}
	query := `SELECT id, ingredient_id, quantity, created_at, updated_at FROM inventory WHERE ingredient_id = ?`
	err := r.db.QueryRowContext(ctx, query, ingredientId).Scan(&inventory.Id, &inventory.IngredientId, &inventory.Quantity, &inventory.CreatedAt, &inventory.UpdatedAt)
	if err != nil {
		logger.Log.Error(err)
		return inventory, err
	}
	return inventory, nil
}

func (r *InventoryRepositoryImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Inventory, error) {
	inventory := domain.Inventory{}
	query := `SELECT id, ingredient_id, quantity, created_at, updated_at FROM inventory WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&inventory.Id, &inventory.IngredientId, &inventory.Quantity, &inventory.CreatedAt, &inventory.UpdatedAt)
	if err != nil {
		logger.Log.Error(err)
		return inventory, err
	}
	return inventory, nil
}

func (r *InventoryRepositoryImpl) Update(ctx context.Context, id uuid.UUID, inventory domain.Inventory) error {
	query := `UPDATE inventory SET quantity = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, inventory.Quantity, id)
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

func (r *InventoryRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE inventory SET deleted_at = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, time.Now(), id)
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

func (r *InventoryRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE inventory SET deleted = false, deleted_at = NULL WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, id)
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

func (r *InventoryRepositoryImpl) GetDeletedById(ctx context.Context, id uuid.UUID) (domain.Inventory, error) {
	inventory := domain.Inventory{}
	query := `SELECT id, ingredient_id, quantity FROM inventory WHERE id = ? AND deleted = true AND deleted_at IS NOT NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&inventory.Id, &inventory.IngredientId, &inventory.Quantity)
	if err != nil {
		logger.Log.Error(err)
		return inventory, err
	}
	return inventory, nil
}
