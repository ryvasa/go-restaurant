package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ErrorResponse struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

func ValidateStruct(data interface{}) []*ErrorResponse {
    var errors []*ErrorResponse
    err := validate.Struct(data)
    if err != nil {
        for _, err := range err.(validator.ValidationErrors) {
            var element ErrorResponse
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
