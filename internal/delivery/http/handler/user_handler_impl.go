package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type UserHandlerImpl struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandlerImpl {
	return &UserHandlerImpl{
		userUsecase: userUsecase,
	}
}

func (h *UserHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	users, err := h.userUsecase.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all users")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, users, nil)
}

func (h *UserHandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := mux.Vars(r)["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		utils.HttpResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid ID format: %w", err))
		return
	}
	user, err := h.userUsecase.Get(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to find user")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, user, nil)
}

func (h *UserHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// logger.Log.WithError(err).Error("Error invalid request body")
		// utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid request body"))
		// return
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	createdUser, err := h.userUsecase.Create(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create user")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, createdUser, nil)
}

func (h *UserHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		utils.HttpResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid ID format: %w", err))
		return
	}

	claims, ok := ctx.Value("user").(jwt.MapClaims)
	if !ok {
		logger.Log.Error("Error getting user claims from context")
		utils.HttpResponse(w, http.StatusUnauthorized, nil,
			utils.NewUnauthorizedError("Invalid user context"))
		return
	}

	authIdStr, ok := claims["sub"].(string)
	if !ok {
		logger.Log.Error("Error getting user ID from claims")
		utils.HttpResponse(w, http.StatusUnauthorized, nil,
			utils.NewUnauthorizedError("Invalid user ID"))
		return
	}

	authId, err := uuid.Parse(authIdStr)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		utils.HttpResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid ID format: %w", err))
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	updatedUser, err := h.userUsecase.Update(ctx, id, authId, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update user")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, updatedUser, nil)
}

func (h *UserHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		utils.HttpResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid ID format: %w", err))
		return
	}

	err = h.userUsecase.Delete(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete user")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	res := map[string]string{"message": "User deleted successfully"}

	utils.HttpResponse(w, http.StatusOK, res, nil)
}

func (h *UserHandlerImpl) Restore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		utils.HttpResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid ID format: %w", err))
		return
	}

	user, err := h.userUsecase.Restore(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore user")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, user, nil)
}
