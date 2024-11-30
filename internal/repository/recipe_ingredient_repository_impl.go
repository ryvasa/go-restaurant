package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

type RecipeIngredientRepositoryImpl struct {
	db DB
}

func NewRecipeIngredientRepository(db DB) RecipeIngredientRepository {
	return &RecipeIngredientRepositoryImpl{
		db: db,
	}
}

func (r *RecipeIngredientRepositoryImpl) Create(ctx context.Context, recipeIngredient domain.RecipeIngredient) error {
	query := `INSERT INTO recipes_ingredients (id, recipe_id, ingredient_id, quantity) VALUES (?, ?, ?, ?)`
	res, err := r.db.ExecContext(ctx, query,
		recipeIngredient.Id, recipeIngredient.RecipeId, recipeIngredient.IngredientId, recipeIngredient.Quantity)
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

func (r *RecipeIngredientRepositoryImpl) Update(ctx context.Context, id uuid.UUID, recipeIngredient domain.RecipeIngredient) error {
	query := `UPDATE recipes_ingredients SET  quantity = ? WHERE id = ?`
	res, err := r.db.ExecContext(ctx, query, recipeIngredient.Quantity, id)
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

func (r *RecipeIngredientRepositoryImpl) GetIngredientsByRecipeId(ctx context.Context, recipeId uuid.UUID) ([]domain.SimpleRecipeIngredient, error) {
	query := `SELECT ingredient_id, ingredients.name, quantity FROM recipes_ingredients INNER JOIN ingredients ON recipes_ingredients.ingredient_id = ingredients.id WHERE recipe_id = ?`
	rows, err := r.db.QueryContext(ctx, query, recipeId)
	if err != nil {
		logger.Log.Error(err)
		return nil, err
	}
	defer rows.Close()

	var recipeIngredients []domain.SimpleRecipeIngredient
	for rows.Next() {
		var recipeIngredient domain.SimpleRecipeIngredient
		err = rows.Scan(&recipeIngredient.IngredientId, &recipeIngredient.Name, &recipeIngredient.Quantity)
		if err != nil {
			logger.Log.Error(err)
			return nil, err
		}
		recipeIngredients = append(recipeIngredients, recipeIngredient)
	}

	return recipeIngredients, nil
}

func (r *RecipeIngredientRepositoryImpl) GetIngredientsByRecipeIdAndIngredientId(ctx context.Context, recipeId, ingredientId uuid.UUID) (domain.RecipeIngredient, error) {
	recipeWithIngredients := domain.RecipeIngredient{}
	query := `SELECT id, recipe_id, ingredient_id, quantity FROM recipes_ingredients WHERE recipe_id = ? AND ingredient_id = ?`
	err := r.db.QueryRowContext(ctx, query, recipeId, ingredientId).Scan(&recipeWithIngredients.Id, &recipeWithIngredients.RecipeId, &recipeWithIngredients.IngredientId, &recipeWithIngredients.Quantity)
	if err != nil {
		logger.Log.Error(err)
		return domain.RecipeIngredient{}, err
	}
	return recipeWithIngredients, nil
}
