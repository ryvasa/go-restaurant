package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/dto"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/utils"
	"github.com/ryvasa/go-restaurant/internal/domain"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
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

func (u *MenuUsecaseImpl) Create(ctx context.Context, req dto.CreateMenuRequest) (domain.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := utils.ValidateStruct(req); len(err) > 0 {
		//logger
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Menu{}, utils.NewValidationError(err)
	}

	id, err := uuid.Parse(req.Restaurant)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid restaurant id format")
		return domain.Menu{}, utils.NewValidationError("Invalid restaurant id format")
	}

	menu := domain.Menu{
		ID:          uuid.New(),
		Name:        req.Name,
		Price:       req.Price,
		Restaurant:  id,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
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

func (u *MenuUsecaseImpl) Update(ctx context.Context, id string, req dto.UpdateMenuRequest) (domain.Menu, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := utils.ValidateStruct(req); len(err) > 0 {
		//logger
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Menu{}, utils.NewValidationError(err)
	}

	menuId, err := uuid.Parse(id)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		return domain.Menu{}, utils.NewValidationError("Invalid id format")
	}

	_, err = u.menuRepo.Get(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		return domain.Menu{}, utils.NewNotFoundError("Menu not found")
	}

	// Convert DTO to domain
	menu := domain.Menu{
		ID:          menuId,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
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
