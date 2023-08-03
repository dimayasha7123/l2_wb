package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"l2_wb/develop/dev11/internal/app"
	"l2_wb/develop/dev11/internal/inputs/httpserver/middlewares"
)

// DeleteEventReq request struct
type DeleteEventReq struct {
	UserID  *int64 `json:"user_id"`
	EventID *int64 `json:"event_id"`
}

// Validate method
func (r DeleteEventReq) Validate() error {
	fields := make([]string, 0, 2)
	if r.UserID == nil {
		fields = append(fields, "user_id")
	}
	if r.EventID == nil {
		fields = append(fields, "event_id")
	}
	if len(fields) == 0 {
		return nil
	}
	return NewMissingFieldsErr(fields)
}

// DeleteEventHandler handler struct
type DeleteEventHandler struct {
	service *app.App
}

// NewDeleteEventHandler constructor
func NewDeleteEventHandler(service *app.App) *UpdateEventHandler {
	return &UpdateEventHandler{service: service}
}

func (h *DeleteEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (any, *middlewares.ServeHTTPError) {
	valErr := validateMethod(http.MethodPost, r.Method)
	if valErr != nil {
		return nil, valErr
	}

	reader := r.Body
	if reader == nil {
		message := "no request body"
		return nil, &middlewares.ServeHTTPError{
			InternalError: errors.New(message),
			Message:       message,
			Code:          http.StatusBadRequest,
		}
	}

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "can't read request body",
			Code:          http.StatusInternalServerError,
		}
	}

	var req DeleteEventReq
	err = json.Unmarshal(bytes, &req)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "can't unmarshal request body",
			Code:          http.StatusInternalServerError,
		}
	}

	err = req.Validate()
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       err.Error(),
			Code:          http.StatusBadRequest,
		}
	}

	appReq := app.DeleteEventReq{
		UserID:  *req.UserID,
		EventID: *req.EventID,
	}

	err = h.service.DeleteEvent(appReq)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "service error",
			Code:          http.StatusServiceUnavailable,
		}
	}

	return nil, nil
}
