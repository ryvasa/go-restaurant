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

type RecipeUsecaseImpl struct {
	recipeRepo repository.RecipeRepository
	menuRepo   repository.MenuRepository
	txRepo     repository.TransactionRepository
}

func NewRecipeUsecase(recipeRepo repository.RecipeRepository, menuRepo repository.MenuRepository, txRepo repository.TransactionRepository) RecipeUsecase {
	return &RecipeUsecaseImpl{
		recipeRepo: recipeRepo,
		menuRepo:   menuRepo,
		txRepo:     txRepo,
	}
}

func (u *RecipeUsecaseImpl) Create(ctx context.Context, req dto.CreateRecipeRequest) (domain.RecipeAndIngredients, error) {
	result := domain.RecipeAndIngredients{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		recipe := domain.Recipe{
			Id:          uuid.New(),
			Name:        req.Name,
			MenuId:      req.MenuId,
			Description: req.Description,
		}

		err := adapters.RecipeRepository.Create(ctx, recipe)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to create recipe")
			return utils.NewInternalError("Failed to create recipe")
		}

		// loop through ingredients and create if not found
		for _, ingredientReq := range req.Ingredients {
			ingredient, err := adapters.IngredientRepository.GetOneByName(ctx, ingredientReq.Name)

			// if ingredient not found create new ingredient
			if err != nil {
				ingredient = domain.Ingredient{
					Id:          uuid.New(),
					Name:        ingredientReq.Name,
					Description: ingredientReq.Name,
				}
				err = adapters.IngredientRepository.Create(ctx, ingredient)
				if err != nil {
					logger.Log.WithError(err).Error("Error failed to create ingredient")
					return utils.NewInternalError("Failed to create ingredient")
				}
			}

			// create recipe ingredient / table junction
			recipeIngredient := domain.RecipeIngredient{
				Id:           uuid.New(),
				RecipeId:     recipe.Id,
				IngredientId: ingredient.Id,
				Quantity:     ingredientReq.Quantity,
			}
			err = adapters.RecipeIngredientRepository.Create(ctx, recipeIngredient)
			if err != nil {
				logger.Log.WithError(err).Error("Error failed to create recipe ingredient")
				return utils.NewInternalError("Failed to create recipe ingredient")
			}
		}

		// get created recipe
		createdRecipe, err := adapters.RecipeRepository.GetOneById(ctx, recipe.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get created recipe")
			return utils.NewInternalError("Failed to get created recipe")
		}
		result.Id = createdRecipe.Id
		result.Name = createdRecipe.Name
		result.Description = createdRecipe.Description
		result.MenuId = createdRecipe.MenuId
		result.CreatedAt = createdRecipe.CreatedAt
		result.UpdatedAt = createdRecipe.UpdatedAt

		// get created ingredients
		createdIngredients, err := adapters.RecipeIngredientRepository.GetIngredientsByRecipeId(ctx, createdRecipe.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get created ingredients")
			return utils.NewInternalError("Failed to get created ingredients")
		}
		result.Ingredients = createdIngredients

		return nil
	})
	if err != nil {
		logger.Log.Error(err)
		return result, err
	}
	return result, nil
}

func (u *RecipeUsecaseImpl) GetAll(ctx context.Context) ([]domain.Recipe, error) {
	recipes, err := u.recipeRepo.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all recipes")
		return []domain.Recipe{}, utils.NewInternalError("Failed to get all recipes")
	}
	return recipes, nil
}

func (u *RecipeUsecaseImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.RecipeAndIngredients, error) {
	result := domain.RecipeAndIngredients{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		recipe, err := adapters.RecipeRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error recipe not found")
			return utils.NewNotFoundError("Recipe not found")
		}
		result.Id = recipe.Id
		result.Name = recipe.Name
		result.Description = recipe.Description
		result.MenuId = recipe.MenuId
		result.CreatedAt = recipe.CreatedAt
		result.UpdatedAt = recipe.UpdatedAt

		recipeIngredients, err := adapters.RecipeIngredientRepository.GetIngredientsByRecipeId(ctx, recipe.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get ingredients")
			return utils.NewInternalError("Failed to get ingredients")
		}
		result.Ingredients = recipeIngredients
		return nil
	})
	if err != nil {
		logger.Log.Error(err)
		return result, err
	}

	return result, nil
}

func (u *RecipeUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateRecipeRequest) (domain.RecipeAndIngredients, error) {
	result := domain.RecipeAndIngredients{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}
		// check if recipe exists
		existingRecipe, err := adapters.RecipeRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error recipe not found")
			return utils.NewNotFoundError("Recipe not found")
		}

		if req.Name != "" {
			existingRecipe.Name = req.Name
		}
		if req.Description != "" {
			existingRecipe.Description = req.Description
		}

		// update recipe
		err = adapters.RecipeRepository.Update(ctx, id, existingRecipe)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update recipe")
			return utils.NewInternalError("Failed to update recipe")
		}

		// get updated recipe
		updatedRecipe, err := adapters.RecipeRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get updated recipe")
			return utils.NewInternalError("Failed to get updated recipe")
		}
		result.Id = updatedRecipe.Id
		result.Name = updatedRecipe.Name
		result.Description = updatedRecipe.Description
		result.MenuId = updatedRecipe.MenuId
		result.CreatedAt = updatedRecipe.CreatedAt
		result.UpdatedAt = updatedRecipe.UpdatedAt

		// if ingredients are provided in request body
		if req.Ingredients != nil {

			// loop for each ingredient in request body
			for _, ingredientReq := range req.Ingredients {
				ingredient, err := adapters.IngredientRepository.GetOneByName(ctx, ingredientReq.Name)

				// if ingredient not found create new ingredient
				if err != nil {
					ingredient = domain.Ingredient{
						Id:          uuid.New(),
						Name:        ingredientReq.Name,
						Description: ingredientReq.Name,
					}
					err = adapters.IngredientRepository.Create(ctx, ingredient)
					if err != nil {
						logger.Log.WithError(err).Error("Error failed to create ingredient")
						return utils.NewInternalError("Failed to create ingredient")
					}
				}

				// if ingredient exists in database, find in table junction
				recepIngredients, err := adapters.RecipeIngredientRepository.GetIngredientsByRecipeIdAndIngredientId(ctx, existingRecipe.Id, ingredient.Id)

				// if in table junction does not exist, create new
				if err != nil {
					// create recipe ingredient / table junction
					recipeIngredient := domain.RecipeIngredient{
						Id:           uuid.New(),
						RecipeId:     existingRecipe.Id,
						IngredientId: ingredient.Id,
						Quantity:     ingredientReq.Quantity,
					}
					err = adapters.RecipeIngredientRepository.Create(ctx, recipeIngredient)
					if err != nil {
						logger.Log.WithError(err).Error("Error failed to create recipe ingredient")
						return utils.NewInternalError("Failed to create recipe ingredient")
					}
				}

				// if in table junction, ingredient quantity does not match with request body, update,
				if recepIngredients.Quantity != ingredientReq.Quantity && ingredientReq.Name == ingredient.Name {
					recepIngredients.Quantity = ingredientReq.Quantity
					err = adapters.RecipeIngredientRepository.Update(ctx, recepIngredients.Id, recepIngredients)
					if err != nil {
						logger.Log.WithError(err).Error("Error failed to update recipe ingredient")
						return utils.NewInternalError("Failed to update recipe ingredient")
					}
				}

			}
		}

		// get updated ingredients
		updatedIngredients, err := adapters.RecipeIngredientRepository.GetIngredientsByRecipeId(ctx, updatedRecipe.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get updated ingredients")
			return utils.NewInternalError("Failed to get updated ingredients")
		}
		result.Ingredients = updatedIngredients
		return nil
	})
	if err != nil {
		logger.Log.Error(err)
		return result, err
	}
	return result, nil
}

func (u *RecipeUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		_, err := adapters.RecipeRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error recipe not found")
			return utils.NewNotFoundError("Recipe not found")
		}

		err = adapters.RecipeRepository.Delete(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to delete recipe")
			return utils.NewInternalError("Failed to delete recipe")
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *RecipeUsecaseImpl) Restore(ctx context.Context, id uuid.UUID) (domain.Recipe, error) {
	recipe := domain.Recipe{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		_, err := adapters.RecipeRepository.GetDeletedById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error recipe not found to restore")
			return utils.NewNotFoundError("Recipe not found to restore")
		}

		err = adapters.RecipeRepository.Restore(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to restore recipe")
			return utils.NewInternalError("Failed to restore recipe")
		}

		restoredRecipe, err := adapters.RecipeRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get restored recipe")
			return utils.NewInternalError("Failed to get restored recipe")
		}
		recipe = restoredRecipe

		return nil
	})
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore recipe")
		return domain.Recipe{}, err
	}
	return recipe, nil
}
