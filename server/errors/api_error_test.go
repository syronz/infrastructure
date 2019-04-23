package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPIError_Error(t *testing.T) {
	e := APIError{
		Message: "abc",
	}
	t.Log("this is log>>>>>>>>> ", e.Error(), e.StatusCode())
	assert.Equal(t, "abc", e.Error())
}

func TestAPIError_StatusCode(t *testing.T) {
	e := APIError{
		Status: 201,
	}

	if e.StatusCode() != 200 {
		t.Errorf("StatusCode(200) == %v, want %v", e.StatusCode(), 200)
	}
}

func TestAPIError_Table(t *testing.T) {
	var tests = []struct {
		s, want APIError
	}{
		{s: APIError{Message: "abc"}, want: APIError{Message: "abc"}, },
		{s: APIError{Message: "OK", Status: 200}, want: APIError{Message: "OK", Status: 201}, },
	}

	for _, c := range tests {
		msg := c.s.Message
		if msg != c.want.Error() || c.s.Status != c.want.StatusCode() {
			t.Errorf("ERROR, Message(%v) != (%v) And StatusCode(%v) != (%v)" , msg, c.want.Error(), c.s.Status, c.want.StatusCode())
		}
	}
}

