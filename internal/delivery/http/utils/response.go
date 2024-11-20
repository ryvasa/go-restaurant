package utils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
    Status  int         `json:"status"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Errors  interface{} `json:"errors,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, payload interface{}) {
    response := Response{
        Status:  status,
        Message: http.StatusText(status),
        Data:    payload,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}

func WriteErrorJSON(w http.ResponseWriter, status int, message string, errors interface{}) {
    response := Response{
        Status:  status,
        Message: message,
        Errors:  errors,
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(response)
}
