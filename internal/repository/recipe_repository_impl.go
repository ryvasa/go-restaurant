package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

type RecipeRepositoryImpl struct {
	db DB
}

func NewRecipeRepository(db DB) RecipeRepository {
	return &RecipeRepositoryImpl{db: db}
}

func (r *RecipeRepositoryImpl) Create(ctx context.Context, recipe domain.Recipe) error {
	query := `INSERT INTO recipes (id, menu_id, name, description) VALUES (?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query, recipe.Id, recipe.MenuId, recipe.Name, recipe.Description)
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

func (r *RecipeRepositoryImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Recipe, error) {
	recipe := domain.Recipe{}
	query := `SELECT id, menu_id, name, description, created_at, updated_at FROM recipes WHERE id = ? AND deleted = false AND deleted_at IS NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&recipe.Id, &recipe.MenuId, &recipe.Name, &recipe.Description, &recipe.CreatedAt, &recipe.UpdatedAt)
	if err != nil {
		logger.Log.Error(err)
		return domain.Recipe{}, err
	}
	return recipe, nil
}

func (r *RecipeRepositoryImpl) GetAll(ctx context.Context) ([]domain.Recipe, error) {
	recipes := []domain.Recipe{}
	query := `SELECT id, menu_id, name, description FROM recipes WHERE deleted = false AND deleted_at IS NULL`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var recipe domain.Recipe
		err := rows.Scan(&recipe.Id, &recipe.MenuId, &recipe.Name, &recipe.Description)
		if err != nil {
			logger.Log.Error(err)
			return nil, err
		}
		recipes = append(recipes, recipe)
	}
	return recipes, nil
}

func (r *RecipeRepositoryImpl) Update(ctx context.Context, id uuid.UUID, recipe domain.Recipe) error {
	query := `UPDATE recipes SET name = ?, description = ? WHERE id = ? AND deleted = false AND deleted_at IS NULL`
	res, err := r.db.ExecContext(ctx, query, recipe.Name, recipe.Description, id)
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

func (r *RecipeRepositoryImpl) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE recipes SET deleted = true, deleted_at = ? WHERE id = ?`
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

func (r *RecipeRepositoryImpl) Restore(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE recipes SET deleted = false, deleted_at = NULL WHERE id = ?`
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

func (r *RecipeRepositoryImpl) GetDeletedById(ctx context.Context, id uuid.UUID) (domain.Recipe, error) {
	recipe := domain.Recipe{}
	query := `SELECT id, menu_id, name, description FROM recipes WHERE id = ? AND deleted = true AND deleted_at IS NOT NULL`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&recipe.Id, &recipe.MenuId, &recipe.Name, &recipe.Description)
	if err != nil {
		logger.Log.Error(err)
		return domain.Recipe{}, err
	}
	return recipe, nil
}
