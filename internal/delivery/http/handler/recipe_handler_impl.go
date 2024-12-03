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

type RecipeHandlerImpl struct {
	recipeUsecase usecase.RecipeUsecase
}

func NewRecipeHandler(recipeUsecase usecase.RecipeUsecase) RecipeHandler {
	return &RecipeHandlerImpl{
		recipeUsecase: recipeUsecase,
	}
}

func (h *RecipeHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.CreateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	createdRecipe, err := h.recipeUsecase.Create(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create recipe")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, createdRecipe, nil)

}

func (h *RecipeHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	recipes, err := h.recipeUsecase.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all recipes")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
	}

	utils.HttpResponse(w, http.StatusOK, recipes, nil)
}

func (h *RecipeHandlerImpl) GetOneById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)
	recipe, err := h.recipeUsecase.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get recipe")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, recipe, nil)
}

func (h *RecipeHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	var req dto.UpdateRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	updatedRecipe, err := h.recipeUsecase.Update(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update recipe")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, updatedRecipe, nil)
}

func (h *RecipeHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	err := h.recipeUsecase.Delete(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete recipe")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	res := map[string]string{"message": "Recipe deleted successfully"}

	utils.HttpResponse(w, http.StatusOK, res, nil)

}

func (h *RecipeHandlerImpl) Restore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	recipe, err := h.recipeUsecase.Restore(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore recipe")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, recipe, nil)

}
