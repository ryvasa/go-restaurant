package utils

import (
	"database/sql"
	"errors"
)
type AppError struct {
    Code    string      `json:"code"`
    Message string      `json:"message"`
    Details interface{} `json:"details,omitempty"`
}

func NewAppError(code string, message string, details interface{}) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Details: details,
    }
}

const (
    ErrInvalidRequest     = "INVALID_REQUEST"
    ErrValidationFailed   = "VALIDATION_FAILED"
    ErrResourceNotFound   = "RESOURCE_NOT_FOUND"
    ErrInternalServer     = "INTERNAL_SERVER_ERROR"
    ErrDuplicateResource  = "DUPLICATE_RESOURCE"
)
func HandleError(err error) *AppError {
    switch {
    case errors.Is(err, sql.ErrNoRows):
        return NewAppError(ErrResourceNotFound, "Resource not found", nil)
    // Handle more specific errors here
    default:
        return NewAppError(ErrInternalServer, "Internal server error", nil)
    }
}
