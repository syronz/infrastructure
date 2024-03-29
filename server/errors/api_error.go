// TODO: This part will be deleted. We will use error model and dictionery for replace
// this package
package errors

type APIError struct {
	Status int `json:"-"`
	ErrorCode string `json:"error_code"`
	Message string `json:"message"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func (e APIError) Error() string {
	return e.Message
}

func (e APIError) StatusCode() int {
	return e.Status
}
