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

type MenuHandlerImpl struct {
	menuUsecase usecase.MenuUsecase
}

func NewMenuHandler(menuUsecase usecase.MenuUsecase) *MenuHandlerImpl {
	return &MenuHandlerImpl{
		menuUsecase: menuUsecase,
	}
}

func (h *MenuHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	menus, err := h.menuUsecase.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all menu")
		utils.Response(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.Response(w, http.StatusOK, menus, nil)
}

func (h *MenuHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateMenuRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.Response(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	createdMenu, err := h.menuUsecase.Create(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create menu")
		utils.Response(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.Response(w, http.StatusCreated, createdMenu, nil)
}

func (h *MenuHandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	menu, err := h.menuUsecase.Get(ctx, id)

	if err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		utils.Response(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.Response(w, http.StatusOK, menu, nil)
}

func (h *MenuHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	var req dto.UpdateMenuRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.Response(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	// Update menu
	updatedtedMenu, err := h.menuUsecase.Update(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update menu")
		utils.Response(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.Response(w, http.StatusOK, updatedtedMenu, nil)
}

func (h *MenuHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	err := h.menuUsecase.Delete(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete menu")
		utils.Response(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.Response(w, http.StatusNoContent, "Menu deleted", nil)
}
