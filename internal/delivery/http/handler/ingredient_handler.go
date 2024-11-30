package handler

import "net/http"

type IngredientHandler interface {
	GetOneById(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	Restore(w http.ResponseWriter, r *http.Request)
}
