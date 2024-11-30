package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type IngredientUsecaseImpl struct {
	ingredientRepo repository.IngredientRepository
	recipeRepo     repository.RecipeRepository
	txRepo         repository.TransactionRepository
}

func NewIngredientUsecase(ingredientRepo repository.IngredientRepository, recipeRepo repository.RecipeRepository, txRepo repository.TransactionRepository) IngredientUsecase {
	return &IngredientUsecaseImpl{
		ingredientRepo: ingredientRepo,
		recipeRepo:     recipeRepo,
		txRepo:         txRepo,
	}
}
func (u *IngredientUsecaseImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Ingredient, error) {
	ingredient, err := u.ingredientRepo.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get ingredient")
		return domain.Ingredient{}, utils.NewNotFoundError("Ingredient not found")
	}
	return ingredient, nil
}

func (u *IngredientUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateIngredientRequest) (domain.Ingredient, error) {
	result := domain.Ingredient{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		ingredient, err := adapters.IngredientRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get ingredient")
			return utils.NewNotFoundError("Ingredient not found")
		}
		if req.Name != "" {
			ingredient.Name = req.Name
		}
		if req.Description != "" {
			ingredient.Description = req.Description
		}
		err = adapters.IngredientRepository.Update(ctx, id, ingredient)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update ingredient")
			return utils.NewInternalError("Failed to update ingredient")
		}
		result = ingredient
		return nil
	})
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update ingredient")
		return domain.Ingredient{}, err
	}

	return result, nil
}

func (u *IngredientUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		_, err := adapters.IngredientRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get ingredient")
			return utils.NewNotFoundError("Ingredient not found")
		}
		err = adapters.IngredientRepository.Delete(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to delete ingredient")
			return utils.NewInternalError("Failed to delete ingedient")
		}
		return nil
	})
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete ingredient")
		return err
	}
	return nil
}

func (u *IngredientUsecaseImpl) Restore(ctx context.Context, id uuid.UUID) (domain.Ingredient, error) {
	result := domain.Ingredient{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		_, err := adapters.IngredientRepository.GetDeletedById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get ingredient")
			return utils.NewNotFoundError("Ingredient not found")
		}
		err = adapters.IngredientRepository.Restore(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to restore ingredient")
			return utils.NewInternalError("Failed to restore ingredient")
		}
		ingredient, err := adapters.IngredientRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get ingredient")
			return utils.NewNotFoundError("Ingredient not found")
		}
		result = ingredient
		return nil
	})
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore ingredient")
		return domain.Ingredient{}, err
	}
	return result, nil
}
