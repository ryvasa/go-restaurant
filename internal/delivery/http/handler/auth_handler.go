package handler

import "net/http"

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}
