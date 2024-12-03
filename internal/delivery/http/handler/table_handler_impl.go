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

type TableHandlerImpl struct {
	tableUsecase usecase.TableUsecase
}

func NewTableHandler(tableUsecase usecase.TableUsecase) TableHandler {
	return &TableHandlerImpl{
		tableUsecase: tableUsecase,
	}
}

func (h *TableHandlerImpl) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tables, err := h.tableUsecase.GetAll(ctx)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get all tables")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, tables, nil)

}

func (h *TableHandlerImpl) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req dto.CreateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	createdTable, err := h.tableUsecase.Create(ctx, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to create table")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusCreated, createdTable, nil)
}

func (h *TableHandlerImpl) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	table, err := h.tableUsecase.GetOneById(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to get table")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	utils.HttpResponse(w, http.StatusOK, table, nil)
}

func (h *TableHandlerImpl) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	var req dto.UpdateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Log.WithError(err).Error("Error invalid request body")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	updatedTable, err := h.tableUsecase.Update(ctx, id, req)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to update table")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, updatedTable, nil)
}

func (h *TableHandlerImpl) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	err := h.tableUsecase.Delete(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to delete table")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}
	res := map[string]string{"message": "Table deleted successfully"}

	utils.HttpResponse(w, http.StatusOK, res, nil)
}

func (h *TableHandlerImpl) Restore(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := mux.Vars(r)["id"]

	id := utils.ValidateIdParam(w, r, idStr)

	table, err := h.tableUsecase.Restore(ctx, id)
	if err != nil {
		logger.Log.WithError(err).Error("Error failed to restore table")
		utils.HttpResponse(w, utils.GetErrorStatus(err), nil, err)
		return
	}

	utils.HttpResponse(w, http.StatusOK, table, nil)
}
