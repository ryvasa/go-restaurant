package usecase

import (
	"context"
	"database/sql"
	"fmt"
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
	var amount float64
	for _, menu := range req.Menu {
		menuId := menu.MenuId

		existingMenu, err := u.menuRepo.Get(tx, menuId)
		if err != nil {
			logger.Log.WithError(err).Error("Error menu not found")
			return domain.Order{}, utils.NewNotFoundError("Menu not found")
		}

		amount += (existingMenu.Price * float64(menu.Quantity))
	}

	order := domain.Order{
		Id:     uuid.New(),
		UserId: user.Id,
		Amount: amount,
	}

	createdOrder, err := u.orderRepo.Create(tx, order)
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
		_, err = u.orderMenuRepo.Create(tx, orderMenu)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to create order menu")
			return domain.Order{}, utils.NewInternalError("Failed to create order menu")
		}
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

	order, err := u.orderRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error order not found")
		return domain.Order{}, utils.NewNotFoundError("Order not found")
	}

	if err := tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Order{}, utils.NewInternalError("Failed to commit transaction")
	}
	return order, nil
}

func (u *OrderUsecaseImpl) UpdateOrderStatus(ctx context.Context, id string, req dto.UpdateOrderStatusDto) (domain.Order, error) {
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

	existingOrder, err := u.orderRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error order not found")
		return domain.Order{}, utils.NewNotFoundError("Order not found")
	}

	if req.Status == "" {
		if existingOrder.Status == "pending" {
			req.Status = "processing"
		} else if existingOrder.Status == "processing" {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return domain.Order{}, utils.NewBadRequestError("Invalid status, must include success or cancel")
		} else {
			logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
			return domain.Order{}, utils.NewBadRequestError(fmt.Sprintf("Invalid update status, status already '%s'", existingOrder.Status))
		}
	}
	order := domain.Order{
		Status: req.Status,
	}
	updatedOrder, err := u.orderRepo.UpdateOrderStatus(tx, id, order)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update order status")
		return domain.Order{}, utils.NewInternalError("Failed to update order status")
	}

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Order{}, utils.NewInternalError("Failed to commit transaction")
	}
	return updatedOrder, nil
}

func (u *OrderUsecaseImpl) UpdatePayment(ctx context.Context, id string, req dto.UpdatePaymentDto) (domain.Order, error) {
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

	_, err = u.orderRepo.GetOneById(tx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error order not found")
		return domain.Order{}, utils.NewNotFoundError("Order not found")
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

	updatedOrder, err := u.orderRepo.UpdatePayment(tx, id, order)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update order payment")
		return domain.Order{}, utils.NewInternalError("Failed to update order payment")
	}

	if updatedOrder.PaymentStatus == "paid" {

		status := domain.Order{
			Status: "success",
		}
		updatedOrder, err = u.orderRepo.UpdateOrderStatus(tx, id, status)
		if err != nil {
			logger.Log.WithError(err).Error("Error failed to update order status")
			return domain.Order{}, utils.NewInternalError("Failed to update order status")
		}
	}
	updatedOrder.PaymentStatus = paymentStatus

	if err = tx.Commit(); err != nil {
		logger.Log.WithError(err).Error("Error failed to commit transaction")
		return domain.Order{}, utils.NewInternalError("Failed to commit transaction")
	}
	return updatedOrder, nil
}
