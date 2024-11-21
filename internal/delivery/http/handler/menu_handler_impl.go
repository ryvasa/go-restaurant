package handler

import (
	"database/sql"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, menus)
}

func (h *MenuHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateMenuRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	if err := utils.ValidateStruct(req); len(err) > 0 {
		utils.WriteErrorJSON(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	menu := domain.Menu{
		Name:  req.Name,
		Price: req.Price,
	}

	createdMenu, err := h.menuUsecase.Create(ctx, menu)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to create menu", nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdMenu)
}

func (h *MenuHandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := mux.Vars(r)["id"]
	menu, err := h.menuUsecase.Get(ctx, id)

	if err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		utils.WriteErrorJSON(w, http.StatusNotFound, "Menu not found", nil)
		return
	}

	utils.WriteJSON(w, http.StatusOK, menu)
}

func (h *MenuHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stringID := mux.Vars(r)["id"]

	id, err := uuid.Parse(stringID)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		utils.WriteErrorJSON(w, http.StatusBadRequest, "Invalid request id", nil)
		return
	}
	var req dto.UpdateMenuRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.WriteErrorJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		utils.WriteErrorJSON(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Convert DTO to domain
	menu := domain.Menu{
		ID:    id,
		Name:  req.Name,
		Price: req.Price,
	}

	// Update menu
	updatedtedMenu, err := h.menuUsecase.Update(ctx, menu)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update menu")
		utils.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to update menu", nil)
		return
	}

	utils.WriteJSON(w, http.StatusOK, updatedtedMenu)
}

func (h *MenuHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		utils.WriteErrorJSON(w, http.StatusBadRequest, "Invalid ID format", nil)
		return
	}

	err := h.menuUsecase.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.WithError(err).Error("Error menu not found")
			utils.WriteErrorJSON(w, http.StatusNotFound, "Menu not found", nil)
			return
		}
		logger.Log.WithError(err).Error("Error failed to delete menu")
		utils.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to delete menu", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
