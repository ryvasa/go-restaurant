package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ryvasa/go-restaurant/internal/delivery/http/dto"
	"github.com/ryvasa/go-restaurant/internal/delivery/http/utils"
	"github.com/ryvasa/go-restaurant/internal/domain"
)

type MenuHandler struct {
	menuUsecase domain.MenuUsecase
}

func NewMenuHandler(menuUsecase domain.MenuUsecase) *MenuHandler {
	return &MenuHandler{
		menuUsecase: menuUsecase,
	}
}

func (h *MenuHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	menus, err := h.menuUsecase.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, menus)

}

func (h *MenuHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateMenuRequest

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorJSON(w, http.StatusBadRequest, "Invalid request body", nil)
		return
	}

	// Validate request
	if errors := utils.ValidateStruct(req); len(errors) > 0 {
		utils.WriteErrorJSON(w, http.StatusBadRequest, "Validation failed", errors)
		return
	}

	// Convert DTO to domain
	menu := domain.Menu{
		Name:  req.Name,
		Price: req.Price,
	}

	// Create menu
	createdMenu, err := h.menuUsecase.Create(menu)
	if err != nil {
		utils.WriteErrorJSON(w, http.StatusInternalServerError, "Failed to create menu", nil)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, createdMenu)
}

// func (h *MenuHandler) GetByID(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
// 	menuID, _ := strconv.Atoi(id)

// 	menu, err := h.menuUsecase.GetByID(r.Context(), menuID)
// 	if err != nil {
// 		helpers.ResponseError(w, err)
// 		return
// 	}

// 	helpers.ResponseJSON(w, menu, http.StatusOK)
// }

// func (h *MenuHandler) Update(w http.ResponseWriter, r *http.Request) {
// 	id := chi.URLParam(r, "id")
// 	menuID, _ := strconv.Atoi(id)

// 	var menu domain.Menu
// 	err := json.NewDecoder(r.Body).Decode(&menu)
//
