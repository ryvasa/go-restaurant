package usecase

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type MenuUsecaseImpl struct {
	menuRepo repository.MenuRepository
}

func NewMenuUsecase(menuRepo repository.MenuRepository) MenuUsecase {
	return &MenuUsecaseImpl{menuRepo}
}
func (u *MenuUsecaseImpl) GetAll(ctx context.Context) ([]domain.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	users, err := u.menuRepo.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all menus")
		return nil, utils.NewInternalError("Failed to get all menus")
	}
	return users, nil
}
func (u *MenuUsecaseImpl) Create(ctx context.Context, req dto.CreateMenuRequest, file multipart.File) (domain.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate request
	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Menu{}, utils.NewValidationError(err)
	}

	// Parse restaurant ID
	restaurantID, err := uuid.Parse(req.Restaurant)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid restaurant id format")
		return domain.Menu{}, utils.NewValidationError("Invalid restaurant id format")
	}

	// Upload file after validation
	imagePath, err := utils.UploadFile(file, req.Image, "menu")
	if err != nil {
		logger.Log.WithError(err).Error("Error uploading file")
		return domain.Menu{}, utils.NewInternalError("Failed to upload image")
	}

	menu := domain.Menu{
		ID:          uuid.New(),
		Name:        req.Name,
		Price:       req.Price,
		Restaurant:  restaurantID,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    imagePath,
	}

	return u.menuRepo.Create(ctx, menu)
}

func (u *MenuUsecaseImpl) Get(ctx context.Context, id string) (domain.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	user, err := u.menuRepo.Get(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		return domain.Menu{}, utils.NewNotFoundError("Menu not found")
	}
	return user, nil
}

func (u *MenuUsecaseImpl) Update(ctx context.Context, id string, req dto.UpdateMenuRequest, file multipart.File) (domain.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Menu{}, utils.NewValidationError(err)
	}

	menuId, err := uuid.Parse(id)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		return domain.Menu{}, utils.NewValidationError("Invalid id format")
	}

	existingMenu, err := u.menuRepo.Get(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		return domain.Menu{}, utils.NewNotFoundError("Menu not found")
	}

	menu := domain.Menu{
		ID:          menuId,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    existingMenu.ImageURL, // Keep existing image if no new image
	}

	// Upload new image if provided
	if file != nil && req.Image != nil {
		imagePath, err := utils.UploadFile(file, req.Image, "menu")
		if err != nil {
			logger.Log.WithError(err).Error("Error uploading file")
			return domain.Menu{}, utils.NewInternalError("Failed to upload image")
		}
		menu.ImageURL = imagePath
	}

	if req.Restaurant != "" {
		restaurantID, err := uuid.Parse(req.Restaurant)
		if err != nil {
			logger.Log.WithError(err).Error("Error invalid restaurant id format")
			return domain.Menu{}, utils.NewValidationError("Invalid id format")

		}
		menu.Restaurant = restaurantID
	}

	return u.menuRepo.Update(ctx, menu)
}

func (u *MenuUsecaseImpl) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return utils.NewValidationError("Invalid ID format")
	}

	if _, err := u.menuRepo.Get(ctx, id); err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		return utils.NewNotFoundError("Menu not found")
	}
	return u.menuRepo.Delete(ctx, id)
}

func (u *MenuUsecaseImpl) Restore(ctx context.Context, id string) (domain.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return domain.Menu{}, utils.NewValidationError("Invalid ID format")
	}

	if _, err := u.menuRepo.GetDeletedMenuById(ctx, id); err != nil {
		logger.Log.WithError(err).Error("Error menu not found to restore")
		return domain.Menu{}, utils.NewNotFoundError("Menu not found to restore")
	}
	return u.menuRepo.Restore(ctx, id)
}
