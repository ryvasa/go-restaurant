package utils

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ryvasa/go-restaurant/pkg/logger"
)

var validate = validator.New()

type errorResponse struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(data interface{}) []*errorResponse {
	var errors []*errorResponse
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errorResponse
			element.Field = err.Field()
			element.Message = getErrorMsg(err)
			errors = append(errors, &element)
		}
	}
	return errors
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Field ini wajib diisi"
	case "min":
		return fmt.Sprintf("Minimal panjang karakter adalah %s", fe.Param())
	case "max":
		return fmt.Sprintf("Maksimal panjang karakter adalah %s", fe.Param())
	case "gt":
		return fmt.Sprintf("Nilai harus lebih besar dari %s", fe.Param())
	}
	return "Unknown error"
}

func ValidateIdParam(w http.ResponseWriter, r *http.Request, idStr string) uuid.UUID {
	id, err := uuid.Parse(idStr)
	if err != nil {
		logger.Log.WithError(err).Error("Error invalid ID format")
		HttpResponse(w, http.StatusBadRequest, nil, fmt.Errorf("invalid ID format: %w", err))
		return uuid.UUID{}
	}
	return id
}
