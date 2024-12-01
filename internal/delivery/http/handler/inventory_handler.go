package handler

import "net/http"

type InventoryHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetOneByIngredientId(w http.ResponseWriter, r *http.Request)
	GetOneById(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Restore(w http.ResponseWriter, r *http.Request)
	CalculateMenuPortions(w http.ResponseWriter, r *http.Request)
}
