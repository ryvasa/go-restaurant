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

type ReservationHandlerImpl struct {
	reservationUsecase usecase.ReservationUsecase
}

func NewReservationHandler(reservationUsecase usecase.ReservationUsecase) *ReservationHandlerImpl {
	return &ReservationHandlerImpl{
		reservationUsecase: reservationUsecase,
	}
}

func (h *ReservationHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reservations, err := h.reservationUsecase.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all reservations")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
	}

	utils.HttpResponse(w, http.StatusOK, reservations, nil)
}

func (h *ReservationHandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r,idStr)

	reservation, err := h.reservationUsecase.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get reservation")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, reservation, nil)
}

func (h *ReservationHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
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

	var req dto.CreateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	req.UserId = userId

	createdReservation, err := h.reservationUsecase.Create(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error Failed to create reservation")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, createdReservation, nil)
}

func (h *ReservationHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.UpdateReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r,idStr)

	updatedReservation, err := h.reservationUsecase.Update(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update reservation")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, updatedReservation, nil)
}

func (h *ReservationHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r,idStr)

	err := h.reservationUsecase.Delete(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete reservation")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	res := map[string]string{"message": "Reservation deleted successfully"}

	utils.HttpResponse(w, http.StatusOK, res, nil)
}

func (h *ReservationHandlerImpl) Restore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r,idStr)

	reservation, err := h.reservationUsecase.Restore(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore reservation")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, reservation, nil)
}
