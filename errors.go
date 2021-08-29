package clickup

import (
	"errors"
	"fmt"
)

var ErrValidation = errors.New("invalid input provided")

type HTTPError struct {
	Status     string
	StatusCode int
	URL        string
}

func (h *HTTPError) Error() string {
	return fmt.Sprintf("clickup response [%s] status: %s code: %d", h.URL, h.Status, h.StatusCode)
}

type ErrClickupResponse struct {
	ECode      string `json:"ECODE"`
	Err        string `json:"err"`
	StatusCode int
	Status     string
}

func (e *ErrClickupResponse) Error() string {
	return fmt.Sprintf("clickup response ECODE=%s Err=%s StatusCode=%d Status=%s", e.ECode, e.Err, e.StatusCode, e.Status)
}
