package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type OrderHandlerImpl struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) *OrderHandlerImpl {
	return &OrderHandlerImpl{
		orderUsecase: orderUsecase,
	}
}

func (h *OrderHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// disini errornya
	claims, ok := ctx.Value("user").(jwt.MapClaims)
	if !ok {
		logger.Log.Error("Error getting user claims from context")
		utils.HttpResponse(w, http.StatusUnauthorized, nil,
			utils.NewUnauthorizedError("Invalid user context"))
		return
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		logger.Log.Error("Error getting user ID from claims")
		utils.HttpResponse(w, http.StatusUnauthorized, nil,
			utils.NewUnauthorizedError("Invalid user ID"))
		return
	}
	var req dto.CreateOrderDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid request body"))
	}

	createdOrder, err := h.orderUsecase.Create(ctx, req, userId)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create order")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, createdOrder, nil)
}

func (h *OrderHandlerImpl) GetOneById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	order, err := h.orderUsecase.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get order")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, order, nil)
}
