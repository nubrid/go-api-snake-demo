package utils

import (
	"github.com/go-playground/validator/v10"
)

// Abs returns the absolute value
func Abs(value int) int {
	if value < 0 {
		return -value
	}

	return value
}

type errorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ValidateStruct(s interface{}) []*errorResponse {
	var errors []*errorResponse

	validate := validator.New()

	err := validate.Struct(s)

	if err != nil {
		// _ = index
		for _, currentErr := range err.(validator.ValidationErrors) {
			var element errorResponse

			// e.g.
			// {
			// 	"FailedField": "MoveSet.Size",
			// 	"Tag": "min",
			// 	"Value": "2"
			// }
			element.FailedField = currentErr.StructNamespace()
			element.Tag = currentErr.Tag()
			element.Value = currentErr.Param()

			errors = append(errors, &element)
		}
	}

	return errors
}
