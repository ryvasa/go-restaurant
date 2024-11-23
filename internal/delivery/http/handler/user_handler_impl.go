package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/dto"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/utils"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/logger"
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
		utils.Response(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.Response(w, http.StatusOK, users, nil)
}

func (h *UserHandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := mux.Vars(r)["id"]
	user, err := h.userUsecase.Get(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to find user")
		utils.Response(w, http.StatusNotFound, nil, err.Error())
		return
	}

	utils.Response(w, http.StatusOK, user, nil)
}

func (h *UserHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.Response(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid request body"))
		return
	}

	createdUser, err := h.userUsecase.Create(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create user")
		status := utils.GetErrorStatus(err)
		utils.Response(w, status, nil, err)
		return
	}

	utils.Response(w, http.StatusCreated, createdUser, nil)
}

func (h *UserHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.Response(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	updatedUser, err := h.userUsecase.Update(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update user")
		utils.Response(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.Response(w, http.StatusOK, updatedUser, nil)
}
