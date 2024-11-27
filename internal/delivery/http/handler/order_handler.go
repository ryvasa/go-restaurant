package handler

import "net/http"

type OrderHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetOneById(w http.ResponseWriter, r *http.Request)
	UpdateOrderStatus(w http.ResponseWriter, r *http.Request)
	UpdatePayment(w http.ResponseWriter, r *http.Request)
}
