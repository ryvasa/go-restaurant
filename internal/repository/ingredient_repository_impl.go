package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

type IngredientRepositoryImpl struct {
	db DB
}

func NewIngredientRepository(db DB) IngredientRepository {
	return &IngredientRepositoryImpl{
		db: db,
	}
}

func (r *IngredientRepositoryImpl) Create(ctx context.Context, ingredient domain.Ingredient) error {
	query := `
		INSERT INTO ingredients (id, name, description)
		VALUES (?, ?, ?)
	`
	res, err := r.db.ExecContext(ctx, query, ingredient.Id, ingredient.Name, ingredient.Description)
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

func (r *IngredientRepositoryImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Ingredient, error) {
	recipe := domain.Ingredient{}
	query := `
		SELECT id, name, description, created_at, updated_at FROM ingredients WHERE id = ? AND deleted = false AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&recipe.Id, &recipe.Name, &recipe.Description, &recipe.CreatedAt, &recipe.UpdatedAt)
	if err != nil {
		logger.Log.Error(err)
		return domain.Ingredient{}, err
	}
	return recipe, nil
}

func (r *IngredientRepositoryImpl) GetOneByName(ctx context.Context, name string) (domain.Ingredient, error) {
	recipe := domain.Ingredient{}
	query := `
		SELECT id, name, description FROM ingredients WHERE name = ? AND deleted = false AND deleted_at IS NULL
	`
	err := r.db.QueryRowContext(ctx, query, name).Scan(&recipe.Id, &recipe.Name, &recipe.Description)
	if err != nil {
		logger.Log.Error(err)
		return domain.Ingredient{}, err
	}
	return recipe, nil
}

func (r *IngredientRepositoryImpl) Update(ctx context.Context, id uuid.UUID, ingredient domain.Ingredient) error {
	query := `
		UPDATE ingredients SET name = ?, description = ? WHERE id = ?
	`
	res, err := r.db.ExecContext(ctx, query, ingredient.Name, ingredient.Description, id)
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

func (r *IngredientRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE ingredients SET deleted = true, deleted_at = ? WHERE id = ?
	`
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

func (r *IngredientRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE ingredients SET deleted = false, deleted_at = NULL WHERE id = ?
	`
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

func (r *IngredientRepositoryImpl) GetDeletedById(ctx context.Context, id uuid.UUID) (domain.Ingredient, error) {
	ingredient := domain.Ingredient{}
	query := `
		SELECT id, name, description FROM ingredients WHERE id = ? AND deleted = true AND deleted_at IS NOT NULL
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&ingredient.Id, &ingredient.Name, &ingredient.Description)
	if err != nil {
		logger.Log.Error(err)
		return domain.Ingredient{}, err
	}
	return ingredient, nil
}
