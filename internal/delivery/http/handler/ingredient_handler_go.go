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

type IngredientHandlerImpl struct {
	ingredientUsecase usecase.IngredientUsecase
}

func NewIngredientHandler(ingredientUsecase usecase.IngredientUsecase) IngredientHandler {
	return &IngredientHandlerImpl{
		ingredientUsecase: ingredientUsecase,
	}
}

func (h *IngredientHandlerImpl) GetOneById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)
	ingredient, err := h.ingredientUsecase.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get ingredient")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	utils.HttpResponse(w, http.StatusOK, ingredient, nil)

}

func (h *IngredientHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	var req dto.UpdateIngredientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	updatedIngredient, err := h.ingredientUsecase.Update(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update ingredient")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	utils.HttpResponse(w, http.StatusOK, updatedIngredient, nil)
}

func (h *IngredientHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	err := h.ingredientUsecase.Delete(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete ingredient")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	res := map[string]string{"message": "Ingredient deleted successfully"}

	utils.HttpResponse(w, http.StatusOK, res, nil)
}

func (h *IngredientHandlerImpl) Restore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	ingredient, err := h.ingredientUsecase.Restore(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore ingredient")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, ingredient, nil)
}
