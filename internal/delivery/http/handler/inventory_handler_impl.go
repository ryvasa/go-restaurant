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

type InventoryHandlerImpl struct {
	invetoryUsecase usecase.InventoryUsecase
}

func NewInventoryHandler(invetoryUsecase usecase.InventoryUsecase) *InventoryHandlerImpl {
	return &InventoryHandlerImpl{
		invetoryUsecase: invetoryUsecase,
	}
}

func (h *InventoryHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.CreateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	createdInventory, err := h.invetoryUsecase.Create(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create inventory")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	utils.HttpResponse(w, http.StatusCreated, createdInventory, nil)

}

func (h *InventoryHandlerImpl) GetOneByIngredientId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]
	ingredientId := utils.ValidateIdParam(w, r, idStr)
	inventory, err := h.invetoryUsecase.GetOneByIngredientId(ctx, ingredientId)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get inventory")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	utils.HttpResponse(w, http.StatusOK, inventory, nil)

}

func (h *InventoryHandlerImpl) GetOneById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)
	inventory, err := h.invetoryUsecase.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get inventory")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	utils.HttpResponse(w, http.StatusOK, inventory, nil)

}

func (h *InventoryHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)
	var req dto.UpdateInventoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	updatedInventory, err := h.invetoryUsecase.Update(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update inventory")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	utils.HttpResponse(w, http.StatusOK, updatedInventory, nil)

}

func (h *InventoryHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	err := h.invetoryUsecase.Delete(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete inventory")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	res := map[string]string{"message": "Inventory deleted successfully"}

	utils.HttpResponse(w, http.StatusOK, res, nil)
}

func (h *InventoryHandlerImpl) Restore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	recipe, err := h.invetoryUsecase.Restore(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore recipe")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, recipe, nil)
}

func (h *InventoryHandlerImpl) CalculateMenuPortions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idStr := mux.Vars(r)["menu_id"]
	menuId := utils.ValidateIdParam(w, r, idStr)
	total, err := h.invetoryUsecase.CalculateMenuPortions(ctx, menuId)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get total menu portions")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	utils.HttpResponse(w, http.StatusOK, total, nil)
}
