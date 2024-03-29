package errors

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

type (
	// Params is used to replace placeholders in an error template with the corresponding values.
	Params map[string]interface{}

	errorTemplate struct {
		Message				string `toml:"message"`
		DeveloperMessage	string `toml:"developer_message"`
	}
)

var templates map[string]errorTemplate

// LoadMessages reads a TOML file containing error templates.
func LoadMessages(file string) error {
	templates = map[string]errorTemplate{}
	_, err := toml.DecodeFile(file, &templates)
	return err
}


// NewAPIError creates a new APIError with the given HTTP status code, error code, and parameters for replacing placeholders in the error template.
// The param can be nil, indicating there is no need for placeholder replacement.
func NewAPIError(status int, code string, params Params) *APIError {
	err := &APIError{
		Status:    status,
		ErrorCode: code,
		Message:   code,
	}

	if template, ok := templates[code]; ok {
		err.Message = template.getMessage(params)
		err.DeveloperMessage = template.getDeveloperMessage(params)
	}

	return err
}

// getMessage returns the error message by replacing placeholders in the error template with the actual parameters.
func (e errorTemplate) getMessage(params Params) string {
	return replacePlaceholders(e.Message, params)
}

// getDeveloperMessage returns the developer message by replacing placeholders in the error template with the actual parameters.
func (e errorTemplate) getDeveloperMessage(params Params) string {
	return replacePlaceholders(e.DeveloperMessage, params)
}

func replacePlaceholders(message string, params Params) string {
	if len(message) == 0 {
		return ""
	}
	for key, value := range params {
		message = strings.Replace(message, "{"+key+"}", fmt.Sprint(value), -1)
	}
	return message
}

