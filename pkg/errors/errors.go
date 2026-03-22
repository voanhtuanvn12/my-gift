package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

func NewNotFound(resource string) *AppError {
	return New(http.StatusNotFound, fmt.Sprintf("%s not found", resource))
}

func NewBadRequest(message string) *AppError {
	return New(http.StatusBadRequest, message)
}

func NewInternal(message string) *AppError {
	return New(http.StatusInternalServerError, message)
}

func NewUnauthorized(message string) *AppError {
	return New(http.StatusUnauthorized, message)
}
