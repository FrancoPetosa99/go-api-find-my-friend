package utils

import (
	"net/http"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(statusCode int, message string, data interface{}) *Response {
	return &Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}
}

func NewCreatedResponse(message string, data interface{}) *Response {
	return NewSuccessResponse(http.StatusCreated, message, data)
}

func NewOKResponse(message string, data interface{}) *Response {
	return NewSuccessResponse(http.StatusOK, message, data)
}

func NewNoContentResponse() *Response {
	return NewSuccessResponse(http.StatusNoContent, "", nil)
}
