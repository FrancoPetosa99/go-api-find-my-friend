package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Type    string `json:"error"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	err := e.Type
	msg := e.Message
	return fmt.Sprintf("%s:%s", err, msg)

}

func (e *AppError) GetMessage() string {
	return e.Message
}

func (e *AppError) GetStatusCode() int {
	return e.Code
}

func NewAppError(code int, message, errorType string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Type:    errorType,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: message,
		Type:    "Not Found",
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Type:    "Bad Request",
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
		Type:    "Unauthorized",
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
		Type:    "Forbidden",
	}
}

func NewConflictError(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
		Type:    "Conflict",
	}
}

func NewUnprocessableEntityError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
		Type:    "Unprocessable Entity",
	}
}

func NewInternalServerError(message string) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Type:    "Internal Server Error",
	}
}

func NewServiceUnavailableError(message string) *AppError {
	return &AppError{
		Code:    http.StatusServiceUnavailable,
		Message: message,
		Type:    "Service Unavailable",
	}
}
