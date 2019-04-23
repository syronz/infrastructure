package errors

import (
	"net/http"
	"sort"
	"fmt"
	"github.com/syronz/infrastructure/server/models"
	"github.com/syronz/infrastructure/server/dict"
	"github.com/syronz/infrastructure/server/app"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/syronz/ozzo-routing"
)

type validationError struct {
	Field string	`json:"field"`
	Error string	`json:"error"`
}

func InternalServerError(err error) *APIError {
	return NewAPIError(http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", Params{"error": err.Error()})
}

func NotFound(resource string) *APIError {
	return NewAPIError(http.StatusNotFound, "NOT_FOUND", Params{"resource": resource})
}

func Unauthorized(err string) *APIError {
	return NewAPIError(http.StatusUnauthorized, "UNAUTHORIZED", Params{"error": err})
}

func InvalidData(errs validation.Errors) *APIError {
	result := []validationError{}
	fields := []string{}
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		err := errs[field]
		result = append(result, validationError{
			Field: field,
			Error: err.Error(),
		})
	}

	err := NewAPIError(http.StatusBadRequest, "INVALID_DATA", nil)
	err.Details = result

	return err
}

func CustomError(c *routing.Context, code int, message string, fields []string, err error) *models.Result {
	var extra string
	if app.Config.Debug && err != nil {
		extra = fmt.Sprint(err.Error())
	}

	c.Response.WriteHeader(code)
	errCustome := &models.Error{
		Code: code,
		Message: dict.T(c, message),
		Fields: fields,
		Extra: extra,
	}

	result := &models.Result {
		Status: false,
		Message: dict.T(c, message),
		Error: errCustome,
	}
	return result
}
