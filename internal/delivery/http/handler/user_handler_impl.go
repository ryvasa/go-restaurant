package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/dto"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/utils"
	"github.com/ryvasa/go-restaurant/internal/domain"
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

func (h *UserHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.Response(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		utils.Response(w, http.StatusBadRequest, nil, err)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Log.WithError(err).Error("Error hashing password")
		utils.Response(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	user := domain.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Password: hashedPassword,
		Email:    req.Email,
		Role:     "customer",
	}

	createdUser, err := h.userUsecase.Create(ctx, user)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create user")
		utils.Response(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.Response(w, http.StatusCreated, createdUser, nil)
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
func (h *UserHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stringID := mux.Vars(r)["id"]

	id, err := uuid.Parse(stringID)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		utils.Response(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.Response(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		utils.Response(w, http.StatusBadRequest, nil, err)
		return
	}

	user := domain.User{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
		Phone: req.Phone,
	}

	if req.Password != "" {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			logger.Log.WithError(err).Error("Error hashing password")
			utils.Response(w, http.StatusInternalServerError, nil, err.Error())
			return
		}
		user.Password = hashedPassword
	}

	updatedUser, err := h.userUsecase.Update(ctx, user)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update user")
		utils.Response(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.Response(w, http.StatusOK, updatedUser, nil)
}
