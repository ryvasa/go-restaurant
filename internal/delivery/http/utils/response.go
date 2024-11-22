package utils

import (
	"encoding/json"
	"net/http"

	"github.com/ryvasa/go-restaurant/internal/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

func Response(w http.ResponseWriter, status int, payload interface{}, error interface{}) {
	response := domain.Response{
		Status:  status,
		Success: status < 400,
		Message: http.StatusText(status),
		Data:    payload,
		Errors:  error,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.WithError(err).Error("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
