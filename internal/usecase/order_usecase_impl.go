package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type OrderUsecaseImpl struct {
	db            *sql.DB
	orderRepo     repository.OrderRepository
	menuRepo      repository.MenuRepository
	userRepo      repository.UserRepository
	orderMenuRepo repository.OrderMenuRepository
}

func NewOrderUsecase(db *sql.DB, orderRepo repository.OrderRepository, menuRepo repository.MenuRepository, userRepo repository.UserRepository, orderMenuRepo repository.OrderMenuRepository) OrderUsecase {
	return &OrderUsecaseImpl{
		db,
		orderRepo,
		menuRepo,
		userRepo,
		orderMenuRepo,
	}
}

func (u *OrderUsecaseImpl) Create(ctx context.Context, req dto.CreateOrderDto, userId string) (domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Order{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Order{}, utils.NewValidationError(err)
	}

	user, err := u.userRepo.Get(tx, userId)
	if err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return domain.Order{}, utils.NewNotFoundError("User not found, order rejected")
	}
	var amount int
	for _, menu := range req.Menu {
		menuId := menu.MenuId

		existingMenu, _ := u.menuRepo.Get(tx, menuId)

		amount += (existingMenu.Price * menu.Quantity)
	}

	order := domain.Order{
		Id:     uuid.New(),
		UserId: user.Id,
		Amount: amount,
	}

	createdOrder, _ := u.orderRepo.Create(tx, order)

	for _, menu := range req.Menu {
		menuId, err := uuid.Parse(menu.MenuId)
		if err != nil {
			logger.Log.WithError(err).Error("Error invalid id format")
			return domain.Order{}, utils.NewValidationError("Invalid id format")
		}
		orderMenu := domain.OrderMenu{
			OrderId:  createdOrder.Id,
			MenuId:   menuId,
			Quantity: menu.Quantity,
		}
		_, _ = u.orderMenuRepo.Create(tx, orderMenu)
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Order{}, utils.NewInternalError("Failed to commit transaction")
	}
	return createdOrder, nil
}

func (u *OrderUsecaseImpl) GetOneById(ctx context.Context, id string) (domain.Order, error) {
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.WithError(err).Error("Error begin transaction")
		return domain.Order{}, utils.NewInternalError("Failed to begin transaction")
	}
	defer func() {
		if err != nil {
			logger.Log.WithError(err).Error("Error when executing a transaction, rollback")
			tx.Rollback()
		}
	}()

	order, _ := u.orderRepo.GetOneById(tx, id)

	if err := tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Order{}, utils.NewInternalError("Failed to commit transaction")
	}
	return order, nil
}
