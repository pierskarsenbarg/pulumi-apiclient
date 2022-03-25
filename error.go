package pulumiapi

import "fmt"

type ErrorResponse struct {
	StatusCode int `json:"code"`
	Message    string
}

func (err *ErrorResponse) Error() string {
	return fmt.Sprintf("%d API Error: %s", err.StatusCode, err.Message)
}
