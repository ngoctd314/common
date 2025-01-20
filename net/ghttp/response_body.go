package ghttp

import (
	"fmt"
	"net/http"
)

// ResponseBody represent common format response to client
type ResponseBody struct {
	Success    bool   `json:"success"`
	Data       any    `json:"data,omitempty"`
	StatusCode int    `json:"-"`
	Message    string `json:"message,omitempty"`
	Error      any    `json:"error,omitempty"`
	Paging     any    `json:"paging,omitempty"`
}

func ResponseBodyOK(data any, opts ...func(*ResponseBody)) *ResponseBody {
	res := &ResponseBody{
		Data:       data,
		Success:    true,
		StatusCode: http.StatusOK,
	}
	for _, opt := range opts {
		if opt != nil {
			opt(res)
		}
	}

	return res
}

func ResponseBodyCreated(data any, entity string) *ResponseBody {
	return &ResponseBody{
		Data:       data,
		Success:    true,
		StatusCode: http.StatusCreated,
		Message:    fmt.Sprintf("%s is created", entity),
	}
}

func ResponseBodyWithMessage(message string) func(*ResponseBody) {
	return func(res *ResponseBody) {
		res.Message = message
	}
}

func ResponseBodyWithStatusCode(statusCode int) func(*ResponseBody) {
	return func(res *ResponseBody) {
		res.StatusCode = statusCode
	}
}
