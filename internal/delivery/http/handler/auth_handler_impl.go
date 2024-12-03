package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type AuthHandlerImpl struct {
	authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) AuthHandler {
	return &AuthHandlerImpl{
		authUsecase: authUsecase,
	}
}

func (a *AuthHandlerImpl) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.LoginDto
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid request body"))
		return
	}

	createdUser, err := a.authUsecase.Login(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create user")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, createdUser, nil)
}
