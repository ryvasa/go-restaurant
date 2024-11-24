package handler

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ryvasa/go-restaurant/internal/model/dto"
	"github.com/ryvasa/go-restaurant/internal/usecase"
	"github.com/ryvasa/go-restaurant/pkg/logger"
	"github.com/ryvasa/go-restaurant/utils"
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
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, menus, nil)
}
func (h *MenuHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logger.Log.WithError(err).Error("Error parsing multipart form")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid form data"))
		return
	}

	// Create request object
	req := dto.CreateMenuRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Category:    r.FormValue("category"),
	}

	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid price format"))
		return
	}
	req.Price = price

	// Get file
	file, handler, err := r.FormFile("image")
	if err != nil {
		logger.Log.WithError(err).Error("Error getting file")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Image is required"))
		return
	}
	defer file.Close()

	req.Image = handler

	// Pass file to usecase
	createdMenu, err := h.menuUsecase.Create(ctx, req, file)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create menu")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, createdMenu, nil)
}

func (h *MenuHandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	menu, err := h.menuUsecase.Get(ctx, id)

	if err != nil {
		logger.Log.WithError(err).Error("Error menu not found")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, menu, nil)
}

func (h *MenuHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		logger.Log.WithError(err).Error("Error parsing multipart form")
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid form data"))
		return
	}

	req := dto.UpdateMenuRequest{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Category:    r.FormValue("category"),
	}

	// Parse price
	price, err := strconv.Atoi(r.FormValue("price"))
	if err != nil {
		utils.HttpResponse(w, http.StatusBadRequest, nil, utils.NewValidationError("Invalid price format"))
		return
	}
	req.Price = price

	// Get file (optional)
	file, handler, err := r.FormFile("image")
	var multipartFile multipart.File
	if err == nil {
		multipartFile = file
		req.Image = handler
		defer file.Close()
	}

	updatedtedMenu, err := h.menuUsecase.Update(ctx, id, req, multipartFile)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update menu")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, updatedtedMenu, nil)
}

func (h *MenuHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	err := h.menuUsecase.Delete(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete menu")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	res := map[string]string{"message": "Menu deleted successfully"}

	utils.HttpResponse(w, http.StatusOK, res, nil)
}

func (h *MenuHandlerImpl) Restore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]

	menu, err := h.menuUsecase.Restore(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore menu")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, menu, nil)
}
