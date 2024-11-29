package usecase

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type MenuUsecaseImpl struct {
	menuRepo repository.MenuRepository
	txRepo   repository.TransactionRepository
}

func NewMenuUsecase(menuRepo repository.MenuRepository, txRepo repository.TransactionRepository) MenuUsecase {
	return &MenuUsecaseImpl{menuRepo, txRepo}
}
func (u *MenuUsecaseImpl) GetAll(ctx context.Context) ([]domain.Menu, error) {
	menu, err := u.menuRepo.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all menu")
		return []domain.Menu{}, utils.NewInternalError("Failed to get all menu")
	}

	return menu, nil
}
func (u *MenuUsecaseImpl) Create(ctx context.Context, req dto.CreateMenuRequest, file multipart.File) (domain.Menu, error) {
	result := domain.Menu{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		imagePath, err := utils.UploadFile(file, req.Image, "menu")
		if err != nil {
			logger.Log.WithError(err).Error("Error uploading file")
			return utils.NewInternalError("Failed to upload image")
		}

		menu := domain.Menu{
			Id:          uuid.New(),
			Name:        req.Name,
			Price:       req.Price,
			Description: req.Description,
			Category:    req.Category,
			ImageURL:    imagePath,
		}
		err = adapters.MenuRepository.Create(ctx, menu)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to create menu")
			return utils.NewInternalError("Failed to create menu")
		}
		createdMenu, err := adapters.MenuRepository.Get(ctx, menu.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get created menu")
			return utils.NewInternalError("Failed to get created menu")
		}
		result = createdMenu
		return nil
	})
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create menu")
		return domain.Menu{}, err
	}

	return result, nil
}

func (u *MenuUsecaseImpl) Get(ctx context.Context, id uuid.UUID) (domain.Menu, error) {
	menu, err := u.menuRepo.Get(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		return domain.Menu{}, utils.NewNotFoundError("Menu not found")
	}
	return menu, nil
}

func (u *MenuUsecaseImpl) Update(ctx context.Context, id uuid.UUID, req dto.UpdateMenuRequest, file multipart.File) (domain.Menu, error) {
	result := domain.Menu{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		existingMenu, err := adapters.MenuRepository.Get(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error menu not found")
			return utils.NewNotFoundError("Menu not found")
		}

		menu := domain.Menu{
			Name:        req.Name,
			Price:       req.Price,
			Description: req.Description,
			Category:    req.Category,
			ImageURL:    existingMenu.ImageURL,
		}

		if file != nil && req.Image != nil {
			imagePath, err := utils.UploadFile(file, req.Image, "menu")
			if err != nil {
				logger.Log.WithError(err).Error("Error uploading file")
				return utils.NewInternalError("Failed to upload image")
			}
			menu.ImageURL = imagePath
		}

		if req.Name == "" {
			menu.Name = existingMenu.Name
		}
		if req.Description == "" {
			menu.Description = existingMenu.Description
		}
		if req.Price == 0 {
			menu.Price = existingMenu.Price
		}
		if req.Category == "" {
			menu.Category = existingMenu.Category
		}

		if menu.ImageURL == "" {
			menu.ImageURL = existingMenu.ImageURL
		}

		err = adapters.MenuRepository.Update(ctx, id, menu)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update menu")
			return utils.NewInternalError("Failed to update menu")
		}
		updatedMenu, err := adapters.MenuRepository.Get(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get updated menu")
			return utils.NewInternalError("Failed to get updated menu")
		}
		result = updatedMenu
		return nil
	})
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update menu")
		return domain.Menu{}, err
	}
	return result, nil
}

func (u *MenuUsecaseImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		_, err := adapters.MenuRepository.Get(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error menu not found")
			return utils.NewNotFoundError("Menu not found")
		}

		err = adapters.MenuRepository.Delete(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to delete menu")
			return utils.NewInternalError("Failed to delete menu")
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *MenuUsecaseImpl) Restore(ctx context.Context, id uuid.UUID) (domain.Menu, error) {
	result := domain.Menu{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		_, err := adapters.MenuRepository.GetDeletedMenuById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error menu not found to restore")
			return utils.NewNotFoundError("Menu not found to restore")
		}

		err = adapters.MenuRepository.Restore(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to restore menu")
			return utils.NewInternalError("Failed to restore menu")
		}

		restoredMenu, err := adapters.MenuRepository.Get(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get restored menu")
			return utils.NewInternalError("Failed to get restored menu")
		}
		result = restoredMenu

		return nil
	})
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore menu")
		return domain.Menu{}, err
	}

	return result, nil
}
