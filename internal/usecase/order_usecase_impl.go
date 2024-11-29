package usecase

import (
	"context"
	"fmt"

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
	txRepo        repository.TransactionRepository
}

func NewOrderUsecase(orderRepo repository.OrderRepository, menuRepo repository.MenuRepository, userRepo repository.UserRepository, orderMenuRepo repository.OrderMenuRepository, txRepo repository.TransactionRepository) OrderUsecase {
	return &OrderUsecaseImpl{
		orderRepo,
		menuRepo,
		userRepo,
		orderMenuRepo,
		txRepo,
	}
}

func (u *OrderUsecaseImpl) Create(ctx context.Context, req dto.CreateOrderDto, userId uuid.UUID) (domain.Order, error) {
	result := domain.Order{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {

		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		user, err := adapters.UserRepository.Get(ctx, userId)
		if err != nil {
			logger.Log.WithError(err).Error("Error user not found")
			return utils.NewNotFoundError("User not found, order rejected")
		}
		var amount float64
		for _, menu := range req.Menu {
			menuId, err := uuid.Parse(menu.MenuId)
			if err != nil {
				logger.Log.WithError(err).Error("Error invalid id format")
				return utils.NewValidationError("Invalid id format")
			}

			existingMenu, err := adapters.MenuRepository.Get(ctx, menuId)
			if err != nil {
				logger.Log.WithError(err).Error("Error menu not found")
				return utils.NewNotFoundError("Menu not found")
			}

			amount += (existingMenu.Price * float64(menu.Quantity))
		}

		order := domain.Order{
			Id:     uuid.New(),
			UserId: user.Id,
			Amount: amount,
		}

		err = adapters.OrderRepository.Create(ctx, order)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to create order")
			return utils.NewInternalError("Failed to create order")
		}

		createdOrder, err := adapters.OrderRepository.GetOneById(ctx, order.Id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get order")
			return utils.NewInternalError("Failed to get order")
		}

		for _, menu := range req.Menu {
			menuId, err := uuid.Parse(menu.MenuId)
			if err != nil {
				logger.Log.WithError(err).Error("Error invalid id format")
				return utils.NewValidationError("Invalid id format")
			}
			orderMenu := domain.OrderMenu{
				OrderId:  createdOrder.Id,
				MenuId:   menuId,
				Quantity: menu.Quantity,
			}
			err = adapters.OrderMenuRepository.Create(ctx, orderMenu)
			if err != nil {
				logger.Log.WithError(err).Error("Error failed to create order menu")
				return utils.NewInternalError("Failed to create order menu")
			}
		}
		result = createdOrder

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *OrderUsecaseImpl) GetOneById(ctx context.Context, id uuid.UUID) (domain.Order, error) {
	order, err := u.orderRepo.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error order not found")
		return domain.Order{}, utils.NewNotFoundError("Order not found")
	}
	return order, nil
}

func (u *OrderUsecaseImpl) UpdateOrderStatus(ctx context.Context, id uuid.UUID, req dto.UpdateOrderStatusDto) (domain.Order, error) {
	result := domain.Order{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		existingOrder, err := adapters.OrderRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error order not found")
			return utils.NewNotFoundError("Order not found")
		}

		if req.Status == "" {
			if existingOrder.Status == "pending" {
				req.Status = "processing"
			} else if existingOrder.Status == "processing" {
				logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
				return utils.NewBadRequestError("Invalid status, must include success or failed")
			} else {
				logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
				return utils.NewBadRequestError(fmt.Sprintf("Invalid update status, status already '%s'", existingOrder.Status))
			}
		}
		order := domain.Order{
			Status: req.Status,
		}
		err = adapters.OrderRepository.UpdateOrderStatus(ctx, id, order)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update order status")
			return utils.NewInternalError("Failed to update order status")
		}
		updatedOrder, err := adapters.OrderRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get order")
			return utils.NewInternalError("Failed to get order")
		}
		result = updatedOrder

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (u *OrderUsecaseImpl) UpdatePayment(ctx context.Context, id uuid.UUID, req dto.UpdatePaymentDto) (domain.Order, error) {
	result := domain.Order{}
	err := u.txRepo.Transact(func(adapters repository.Adapters) error {
		if err := utils.ValidateStruct(req); len(err) > 0 {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return utils.NewValidationError(err)
		}

		_, err := adapters.OrderRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error order not found")
			return utils.NewNotFoundError("Order not found")
		}
		var paymentStatus string
		if req.PaymentMethod == nil {
			paymentStatus = "unpaid"
		} else {
			paymentStatus = "paid"
		}

		order := domain.Order{
			PaymentMethod: req.PaymentMethod,
			PaymentStatus: paymentStatus,
		}

		err = adapters.OrderRepository.UpdatePayment(ctx, id, order)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update order payment")
			return utils.NewInternalError("Failed to update order payment")
		}

		if order.PaymentStatus == "paid" {

			status := domain.Order{
				Status: "success",
			}
			err = adapters.OrderRepository.UpdateOrderStatus(ctx, id, status)
			if err != nil {
				logger.Log.WithError(err).Error("Error failed to update order status")
				return utils.NewInternalError("Failed to update order status")
			}
		}
		updatedOrder, err := adapters.OrderRepository.GetOneById(ctx, id)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to get order")
			return utils.NewInternalError("Failed to get order")
		}
		result = updatedOrder

		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}
