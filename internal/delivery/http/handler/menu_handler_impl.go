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
		logger.Log.WithError(err).Error("Error failed to get all menu")
		utils.Response(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.Response(w, http.StatusOK, menus, nil)
}

func (h *MenuHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateMenuRequest

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

	id, err := uuid.Parse(req.Restaurant)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		utils.Response(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	menu := domain.Menu{
		ID:          uuid.New(),
		Name:        req.Name,
		Price:       req.Price,
		Restaurant:  id,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
	}

	createdMenu, err := h.menuUsecase.Create(ctx, menu)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create menu")
		utils.Response(w, http.StatusInternalServerError, nil, err.Error())
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
		utils.Response(w, http.StatusNotFound, nil, err.Error())
		return
	}

	utils.Response(w, http.StatusOK, menu, nil)
}

func (h *MenuHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	stringID := mux.Vars(r)["id"]

	id, err := uuid.Parse(stringID)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid id format")
		utils.Response(w, http.StatusBadRequest, nil, err.Error())
		return
	}
	var req dto.UpdateMenuRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.Response(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); len(err) > 0 {
		logger.Log.WithField("validation_errors", err).Error("Error invalid request body")
		utils.Response(w, http.StatusBadRequest, nil, err)
		return
	}

	// Convert DTO to domain
	menu := domain.Menu{
		ID:          id,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
	}

	if req.Restaurant != "" {
		restaurantID, err := uuid.Parse(req.Restaurant)
		if err != nil {
			logger.Log.WithError(err).Error("Error invalid restaurant id format")
			utils.Response(w, http.StatusBadRequest, nil, err.Error())
			return
		}
		menu.Restaurant = restaurantID
	}

	// Update menu
	updatedtedMenu, err := h.menuUsecase.Update(ctx, menu)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update menu")
		utils.Response(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.Response(w, http.StatusOK, updatedtedMenu, nil)
}

func (h *MenuHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	if _, err := uuid.Parse(id); err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		utils.Response(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	err := h.menuUsecase.Delete(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Log.WithError(err).Error("Error menu not found")
			utils.Response(w, http.StatusNotFound, nil, err.Error())
			return
		}
		logger.Log.WithError(err).Error("Error failed to delete menu")
		utils.Response(w, http.StatusInternalServerError, nil, err.Error())
		return
	}

	utils.Response(w, http.StatusNoContent, "Menu deleted", nil)
}
