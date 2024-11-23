package utils

import (
	"encoding/json"
	"net/http"

	"github.com/ryvasa/go-restaurant/internal/domain"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

func Response(w http.ResponseWriter, status int, payload interface{}, err interface{}) {
	var errorResponse interface{}

	if err != nil {
		switch v := err.(type) {
		case AppError:
			errorResponse = v
		default:
			if e, ok := err.(error); ok {
				errorResponse = AppError{
					Code:    "UNKNOWN_ERROR",
					Message: e.Error(),
				}
			} else {
				errorResponse = err
			}
		}
	}

	response := domain.Response{
		Status:  status,
		Success: status < 400,
		Message: http.StatusText(status),
		Data:    payload,
		Errors:  errorResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Log.WithError(err).Error("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
