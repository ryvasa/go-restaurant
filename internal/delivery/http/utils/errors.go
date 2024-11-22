package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

type ErrorResponse struct {
	Status  int         `json:"status"`
	Code    string           `json:"code"`
	Message string           `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// AppError adalah custom error type untuk aplikasi
type AppError struct {
	Err     error
	Status  int
	Code    string
	Message string
	Details interface{}
}

func (e *AppError) Error() string {
	return e.Message
}

// HandleError adalah fungsi utama untuk handling berbagai jenis error
func HandleError(err error) *ErrorResponse {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return &ErrorResponse{
			Status:  appErr.Status,
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		}
	}

	// Handle different types of errors
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return &ErrorResponse{
			Status:  http.StatusNotFound,
			Code:    "NOT_FOUND",
			Message: "Resource not found",
		}

	case strings.Contains(err.Error(), "Duplicate entry"):
		field, value := extractDuplicateError(err.Error())
		return &ErrorResponse{
			Status:  http.StatusConflict,
			Code:    "DUPLICATE_ENTRY",
			Message: fmt.Sprintf("%s %s already exists", field, value),
		}

	case strings.Contains(err.Error(), "foreign key constraint fails"):
		return &ErrorResponse{
			Status:  http.StatusBadRequest,
			Code:    "FOREIGN_KEY_VIOLATION",
			Message: "Related resource does not exist",
		}

	default:
		// Log unknown errors here
		return &ErrorResponse{
			Status:  http.StatusInternalServerError,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An unexpected error occurred",
		}
	}
}

// NewAppError creates a new AppError
func NewAppError(status int, code, message string, details interface{}) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Helper functions
func extractDuplicateError(errStr string) (field, value string) {
	re := regexp.MustCompile(`Duplicate entry '(.+?)' for key '(.+?)'`)
	matches := re.FindStringSubmatch(errStr)

	if len(matches) > 2 {
		value = matches[1]
		field = matches[2]

		fieldParts := strings.Split(field, ".")
		if len(fieldParts) > 1 {
			field = fieldParts[1]
		}
		field = strings.Title(field)
	}

	return field, value
}
