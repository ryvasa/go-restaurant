package usecase

import (
	"context"
	"database/sql"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type MenuUsecaseImpl struct {
	db       *sql.DB
	menuRepo repository.MenuRepository
}

func NewMenuUsecase(db *sql.DB, menuRepo repository.MenuRepository) MenuUsecase {
	return &MenuUsecaseImpl{db: db, menuRepo: menuRepo}
}
func (u *MenuUsecaseImpl) GetAll(ctx context.Context) ([]domain.Menu, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return []domain.Menu{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	menu, _ := u.menuRepo.GetAll(tx)

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return []domain.Menu{}, utils.NewInternalError("Failed to commit transaction")
	}
	return menu, nil
}
func (u *MenuUsecaseImpl) Create(ctx context.Context, req dto.CreateMenuRequest, file multipart.File) (domain.Menu, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Menu{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()
	// Validate request
	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Menu{}, utils.NewValidationError(err)
	}

	// Upload file after validation
	imagePath, err := utils.UploadFile(file, req.Image, "menu")
	if err != nil {
		logger.Log.WithError(err).Error("Error uploading file")
		return domain.Menu{}, utils.NewInternalError("Failed to upload image")
	}

	menu := domain.Menu{
		Id:          uuid.New(),
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    imagePath,
	}
	createdMenu, _ := u.menuRepo.Create(tx, menu)

	if err := tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Menu{}, utils.NewInternalError("Failed to commit transaction")
	}

	return createdMenu, nil
}

func (u *MenuUsecaseImpl) Get(ctx context.Context, id string) (domain.Menu, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Menu{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	menu, _ := u.menuRepo.Get(tx, id)

	if err := tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Menu{}, utils.NewInternalError("Failed to commit transaction")
	}
	return menu, nil
}

func (u *MenuUsecaseImpl) Update(ctx context.Context, id string, req dto.UpdateMenuRequest, file multipart.File) (domain.Menu, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Menu{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Menu{}, utils.NewValidationError(err)
	}

	menuId, err := uuid.Parse(id)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		return domain.Menu{}, utils.NewValidationError("Invalid id format")
	}

	existingMenu, _ := u.menuRepo.Get(tx, id)

	menu := domain.Menu{
		Id:          menuId,
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
	updatedMenu, _ := u.menuRepo.Update(tx, menu)

	if err := tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Menu{}, utils.NewInternalError("Failed to commit transaction")
	}
	return updatedMenu, nil
}

func (u *MenuUsecaseImpl) Delete(ctx context.Context, id string) error {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return utils.NewValidationError("Invalid ID format")
	}

	u.menuRepo.Get(tx, id)

	u.menuRepo.Delete(tx, id)

	if err := tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return utils.NewInternalError("Failed to commit transaction")
	}
	return nil
}

func (u *MenuUsecaseImpl) Restore(ctx context.Context, id string) (domain.Menu, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Menu{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		return domain.Menu{}, utils.NewValidationError("Invalid ID format")
	}

	u.menuRepo.GetDeletedMenuById(tx, id)

	restoredMenu, _ := u.menuRepo.Restore(tx, id)

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Menu{}, utils.NewInternalError("Failed to commit transaction")
	}
	return restoredMenu, nil
}
