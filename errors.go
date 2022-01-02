// Copyright (c) 2021, John Moore
// All rights reserved.

// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package clickup

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrValidation = errors.New("invalid input provided")
var ErrCall = errors.New("failed to make request")

type RateLimitError struct {
	msg       string
	limit     string
	remaining string
	resetAt   string
	cause     error
}

func (r *RateLimitError) Error() string {
	return fmt.Sprintf("%s - limit: %s, remaining: %s, reset at: %s", r.msg, r.limit, r.remaining, r.resetAt)
}

func (r *RateLimitError) Unwrap() error {
	return r.cause
}

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

func errorFromResponse(res *http.Response, decoder *json.Decoder) error {
	if res.StatusCode == http.StatusTooManyRequests {
		return &RateLimitError{
			msg:       "rate limit exceeded",
			limit:     res.Header.Get("x-ratelimit-limit"),
			remaining: res.Header.Get("x-ratelimit-remaining"),
			resetAt:   res.Header.Get("x-ratelimit-reset"),
		}
	}

	var errResponse ErrClickupResponse
	if err := decoder.Decode(&errResponse); err != nil {
		return &HTTPError{
			Status:     res.Status,
			StatusCode: res.StatusCode,
			URL:        res.Request.URL.String(),
		}
	}
	errResponse.StatusCode = res.StatusCode
	errResponse.Status = res.Status
	return &errResponse
}
