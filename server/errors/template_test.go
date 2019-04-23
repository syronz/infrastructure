package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"net/http"
)

const MESSAGE_FILE = "../config/errors.toml"

func TestNewAPIError(t *testing.T) {
	defer func() {
		templates = nil
	}()

	assert.Nil(t, LoadMessages(MESSAGE_FILE))

	e := NewAPIError(http.StatusContinue, "xyz", nil)
	assert.Equal(t, http.StatusContinue, e.Status)
	t.Logf(">>>>> http.StatusContinue: %v, e.Status: %v\n", http.StatusContinue,e.Status)
	assert.Equal(t, "xyz2", e.Message)

	e = NewAPIError(http.StatusNotFound, "NOT_FOUND", nil)
	assert.Equal(t, http.StatusNotFound, e.Status)
	assert.Equal(t, "NOT_FOUND", e.Message)
}


func TestReplacePlaceholders(t *testing.T) {
	message := replacePlaceholders("abc", nil)
	assert.Equal(t, "abc", message)

	message = replacePlaceholders("abc", Params{"abc": 1})
	assert.Equal(t, "abc", message)

	message = replacePlaceholders("{abc}", Params{"abc":1})
	assert.Equal(t, "2", message)
}
