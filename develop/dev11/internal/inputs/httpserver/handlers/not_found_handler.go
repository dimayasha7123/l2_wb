package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"l2_wb/develop/dev11/internal/inputs/httpserver/middlewares"
)

// NotFoundHandler handler struct
type NotFoundHandler struct{}

// NewNotFoundHandler constructor
func NewNotFoundHandler() *NotFoundHandler {
	return &NotFoundHandler{}
}

func (h *NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (any, *middlewares.ServeHTTPError) {
	message := fmt.Sprintf("no registered handlers with path %s", r.URL.Path)
	return nil, &middlewares.ServeHTTPError{
		InternalError: errors.New(message),
		Message:       message,
		Code:          http.StatusNotFound,
	}
}
