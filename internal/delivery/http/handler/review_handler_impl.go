package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
)

type ReviewHandlerImpl struct {
	reviewUsecase usecase.ReviewUsecase
}

func NewReviewHandler(reviewUsecase usecase.ReviewUsecase) *ReviewHandlerImpl {
	return &ReviewHandlerImpl{
		reviewUsecase: reviewUsecase,
	}
}

func (h *ReviewHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid request body"))
	}

	createdReview, err := h.reviewUsecase.Create(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create review")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, createdReview, nil)
}

func (h *ReviewHandlerImpl) GetAllByMenuId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

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
	id := mux.Vars(r)["id"]

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
	id := mux.Vars(r)["id"]

	var req dto.UpdateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid request body"))
	}

	review, err := h.reviewUsecase.Update(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update review")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, review, nil)
}
