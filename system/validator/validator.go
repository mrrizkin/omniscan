package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	GlobalErrorResponse struct {
		Title   string `json:"title"`
		Message string `json:"message"`
		Status  string `json:"status"`
		Detail  string `json:"detail"`
	}

	Validator struct {
		validator *validator.Validate
	}
)

func New() *Validator {
	var validate = validator.New()
	return &Validator{validator: validate}
}

func (v *Validator) MustValidate(data interface{}) *fiber.Error {
	errs := v.Validate(data)
	if len(errs) == 0 {
		return nil
	}

	return &fiber.Error{
		Code:    fiber.ErrBadRequest.Code,
		Message: strings.Join(v.Format(errs), " and "),
	}
}

func (v *Validator) Validate(data interface{}) []ErrorResponse {
	errorResponse := make([]ErrorResponse, 0)

	errs := v.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var elem ErrorResponse

			elem.FailedField = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Value()
			elem.Error = true

			errorResponse = append(errorResponse, elem)
		}
	}

	return errorResponse
}

func (v *Validator) Format(errs []ErrorResponse) []string {
	errorMessages := make([]string, 0)
	for _, err := range errs {
		if !err.Error {
			continue
		}
		errorMessages = append(errorMessages, fmt.Sprintf(
			"[%s]: '%v' | Needs to implement '%s'",
			err.FailedField,
			err.Value,
			err.Tag,
		))
	}
	return errorMessages
}
