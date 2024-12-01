package usecase

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type InventoryUsecaseImpl struct {
	inventoryRepo repository.InventoryRepository
	txRepo        repository.TransactionRepository
}

func NewInventoryUsecase(inventoryRepo repository.InventoryRepository, txRepo repository.TransactionRepository) InventoryUsecase {
	return &InventoryUsecaseImpl{
		inventoryRepo: inventoryRepo,
		txRepo:        txRepo,
	}
}

func (u *InventoryUsecaseImpl) Create(ctx context.Context, req dto.CreateInventoryRequest) (domain.Inventory, error) {
	result := domain.Inventory{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}
		// TODO: apakah jika ingedient tidak ditemukan maka membuat ingredient baru?
		_, err := adapters.IngredientRepository.GetOneById(ctx, req.IngredientId)
		if err != nil {
			logger.Log.WithError(err).Error("Ingredient not found")
			return utils.NewNotFoundError("Ingredient not found")

		}
		inventory := domain.Inventory{
			Id:           uuid.New(),
			IngredientId: req.IngredientId,
			Quantity:     req.Quantity,
		}
		err = adapters.InventoryRepository.Create(ctx, inventory)
		if err != nil {
			logger.Log.WithError(err).Error("Error creating inventory")
			return utils.NewInternalError("Failed to create inventory")

		}

		createdInventory, err := adapters.InventoryRepository.GetOneById(ctx, inventory.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error getting inventory")
			return utils.NewInternalError("Failed to get inventory inventory")

		}

		result = createdInventory

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *InventoryUsecaseImpl) GetOneByIngredientId(ctx context.Context, ingredientId uuid.UUID) (domain.Inventory, error) {
	inventory, err := u.inventoryRepo.GetOneByIngredientId(ctx, ingredientId)
	if err != nil {
		logger.Log.WithError(err).Error("Error getting inventory")
		return inventory, utils.NewInternalError("Failed to get restored inventory")

	}
	return inventory, nil
}

func (u *InventoryUsecaseImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Inventory, error) {
	inventory, err := u.inventoryRepo.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error getting inventory")
		return inventory, utils.NewInternalError("Failed to get inventory")
	}
	return inventory, nil
}

func (u *InventoryUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateInventoryRequest) (domain.Inventory, error) {
	result := domain.Inventory{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		_, err := adapters.InventoryRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error getting inventory")
			return utils.NewNotFoundError("Inventory not found")
		}

		inventory := domain.Inventory{
			Quantity: req.Quantity,
		}

		err = adapters.InventoryRepository.Update(ctx, id, inventory)
		if err != nil {
			logger.Log.WithError(err).Error("Error updating inventory")
			return utils.NewInternalError("Failed to get restored inventory")
		}

		updatedInventory, err := adapters.InventoryRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error getting inventory")
			return utils.NewInternalError("Failed to get restored inventory")
		}

		result = updatedInventory

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *InventoryUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		_, err := adapters.InventoryRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error getting inventory")
			return utils.NewNotFoundError("Inventory not found")
		}
		err = adapters.InventoryRepository.Delete(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error deleting inventory")
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *InventoryUsecaseImpl) Restore(ctx context.Context, id uuid.UUID) (domain.Inventory, error) {
	result := domain.Inventory{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		_, err := adapters.InventoryRepository.GetDeletedById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error getting inventory")
			return utils.NewNotFoundError("Inventory not found")
		}
		err = adapters.InventoryRepository.Restore(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error restoring inventory")
			return err
		}

		restoredInventory, err := adapters.InventoryRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error getting inventory")
			return utils.NewInternalError("Failed to get restored inventory")
		}

		result = restoredInventory
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *InventoryUsecaseImpl) CalculateMenuPortions(ctx context.Context, menuId uuid.UUID) (domain.InventoryMenu, error) {
	result := domain.InventoryMenu{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		menu, err := adapters.MenuRepository.Get(ctx, menuId)
		if err != nil {
			logger.Log.WithError(err).Error("Error getting Menu")
			return utils.NewNotFoundError("Menu not found")
		}
		result.Menu = menu

		recipe, err := adapters.RecipeRepository.GetOneByMenuId(ctx, menuId)
		if err != nil {
			logger.Log.WithError(err).Error("Error getting Recipe")
			return utils.NewNotFoundError("Recipe not found")
		}
		result.Recipe = recipe

		recipeIngredients, err := adapters.RecipeIngredientRepository.GetIngredientsByRecipeId(ctx, recipe.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error getting Ingredients")
			return utils.NewInternalError("Failed to get ingredients")
		}
		result.Ingredients = recipeIngredients

		for _, recipeIngredient := range recipeIngredients {
			inventory, err := adapters.InventoryRepository.GetOneByIngredientId(ctx, recipeIngredient.IngredientId)
			if err != nil {
				logger.Log.WithError(err).Error("Error getting Inventory")
				return utils.NewInternalError("Failed to get inventory")
			}

			if recipeIngredient.Quantity > 0 {
				var minPortions float64 = math.MaxFloat64
				maxPortions := inventory.Quantity / recipeIngredient.Quantity

				result.TotalPortions = math.Floor(math.Min(minPortions, maxPortions))
			} else {
				return utils.NewInternalError("Invalid ingredient quantity in recipe")
			}
		}

		return nil
	})

	if err != nil {
		logger.Log.WithError(err).Error("Error calculating menu portions")
		return result, err
	}

	return result, nil
}
