package utils

import "net/http"

type AppError struct {
	HttpStatus int         `json:"-"`                 // HTTP status code (tidak ditampilkan di response)
	Code       string      `json:"code"`              // Error code internal
	Message    string      `json:"message"`           // Pesan error untuk user
	Details    interface{} `json:"details,omitempty"` // Detail tambahan (optional)
}

func (e AppError) Error() string {
	return e.Message
}

// Helper functions untuk membuat error
func NewValidationError(details interface{}) error {
	return AppError{
		HttpStatus: http.StatusBadRequest,
		Code:       "VALIDATION_ERROR",
		Message:    "Validation failed",
		Details:    details,
	}
}

func NewNotFoundError(message string) error {
	return AppError{
		HttpStatus: http.StatusNotFound,
		Code:       "NOT_FOUND",
		Message:    message,
	}
}

func NewConflictError(message string) error {
	return AppError{
		HttpStatus: http.StatusConflict,
		Code:       "CONFLICT",
		Message:    message,
	}
}

func NewInternalError(message string) error {
	return AppError{
		HttpStatus: http.StatusInternalServerError,
		Code:       "INTERNAL_ERROR",
		Message:    message,
	}
}

func NewUnauthorizedError(message string) error {
	return AppError{
		HttpStatus: http.StatusUnauthorized,
		Code:       "UNAUTHORIZED",
		Message:    message,
	}
}

func GetErrorStatus(err error) int {
	if appErr, ok := err.(AppError); ok {
		return appErr.HttpStatus
	}
	return http.StatusInternalServerError
}
