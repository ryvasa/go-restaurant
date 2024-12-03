package handler

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type ReviewHandlerImpl struct {
	reviewUsecase usecase.ReviewUsecase
}

func NewReviewHandler(reviewUsecase usecase.ReviewUsecase) ReviewHandler {
	return &ReviewHandlerImpl{
		reviewUsecase: reviewUsecase,
	}
}

func (h *ReviewHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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

	var req dto.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid request body"))
	}

	id, err := uuid.Parse(userId)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid user id format")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid user id format"))
		return
	}

	createdReview, err := h.reviewUsecase.Create(ctx, req, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create review")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, createdReview, nil)
}

func (h *ReviewHandlerImpl) GetAllByMenuId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	reviews, err := h.reviewUsecase.GetAllByMenuId(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all reviews")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, reviews, nil)
}

func (h *ReviewHandlerImpl) GetOneById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	review, err := h.reviewUsecase.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get review")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, review, nil)
}

func (h *ReviewHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	var req dto.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid request body"))
	}
	claims, ok := ctx.Value("user").(jwt.MapClaims)
	if !ok {
		logger.Log.Error("Error getting user claims from context")
		utils.HttpResponse(w, http.StatusUnauthorized, nil,
			utils.NewUnauthorizedError("Invalid user context"))
		return
	}

	userIdStr, ok := claims["sub"].(string)
	if !ok {
		logger.Log.Error("Error getting user ID from claims")
		utils.HttpResponse(w, http.StatusUnauthorized, nil,
			utils.NewUnauthorizedError("Invalid user ID"))
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid id format"))
		return
	}

	review, err := h.reviewUsecase.Update(ctx, id, userId, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update review")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, review, nil)
}
