package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/internal/model/domain"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/repository"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type OrderUsecaseImpl struct {
	orderRepo     repository.OrderRepository
	menuRepo      repository.MenuRepository
	userRepo      repository.UserRepository
	orderMenuRepo repository.OrderMenuRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository, menuRepo repository.MenuRepository, userRepo repository.UserRepository, orderMenuRepo repository.OrderMenuRepository) OrderUsecase {
	return &OrderUsecaseImpl{
		orderRepo,
		menuRepo,
		userRepo,
		orderMenuRepo,
	}
}

func (u *OrderUsecaseImpl) Create(ctx context.Context, req dto.CreateOrderDto, userId string) (domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		return domain.Order{}, utils.NewValidationError(err)
	}

	user, err := u.userRepo.Get(ctx, userId)
	if err != nil {
		logger.Log.WithError(err).Error("Error user not found")
		return domain.Order{}, utils.NewNotFoundError("User not found, order rejected")
	}
	var amount int
	for _, menu := range req.Menu {
		menuId := menu.MenuId

		existingMenu, err := u.menuRepo.Get(ctx, menuId)
		if err != nil {
			logger.Log.WithError(err).Error("Error user not found")
			return domain.Order{}, utils.NewNotFoundError("Menu not found, order rejected")
		}
		amount += (existingMenu.Price * menu.Quantity)
	}

	order := domain.Order{
		Id:     uuid.New(),
		UserId: user.Id,
		Amount: amount,
	}

	createdOrder, err := u.orderRepo.Create(ctx, order)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create order")
		return domain.Order{}, utils.NewInternalError("Failed to create order")
	}
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
		_, err = u.orderMenuRepo.Create(ctx, orderMenu)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to create order menu")
			return domain.Order{}, utils.NewInternalError("Failed to create order menu")
		}
	}
	return createdOrder, nil
}

func (u *OrderUsecaseImpl) GetOneById(ctx context.Context, id string) (domain.Order, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return u.orderRepo.GetOneById(ctx, id)
}
