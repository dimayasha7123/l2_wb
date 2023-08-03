package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"l2_wb/develop/dev11/internal/app"
	"l2_wb/develop/dev11/internal/inputs/httpserver/middlewares"
	"l2_wb/develop/dev11/internal/utils/converters"
)

// UpdateEventReq request struct
type UpdateEventReq struct {
	UserID      *int64 `json:"user_id"`
	EventID     *int64 `json:"event_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

// Validate method
func (r UpdateEventReq) Validate() error {
	fields := make([]string, 0, 5)
	if r.UserID == nil {
		fields = append(fields, "user_id")
	}
	if r.EventID == nil {
		fields = append(fields, "event_id")
	}
	if r.Title == "" {
		fields = append(fields, "title")
	}
	if r.Description == "" {
		fields = append(fields, "description")
	}
	if r.Date == "" {
		fields = append(fields, "date")
	}
	if len(fields) == 0 {
		return nil
	}
	return NewMissingFieldsErr(fields)
}

// UpdateEventHandler handler struct
type UpdateEventHandler struct {
	service *app.App
}

// NewUpdateEventHandler constructor
func NewUpdateEventHandler(service *app.App) *UpdateEventHandler {
	return &UpdateEventHandler{service: service}
}

func (h *UpdateEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (any, *middlewares.ServeHTTPError) {
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

	var req UpdateEventReq
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

	date, err := converters.StrToDate(req.Date)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "can't convert date from string",
			Code:          http.StatusBadRequest,
		}
	}

	appReq := app.UpdateEventReq{
		UserID:      *req.UserID,
		EventID:     *req.EventID,
		Title:       req.Title,
		Description: req.Description,
		Date:        date,
	}

	err = h.service.UpdateEvent(appReq)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "service error",
			Code:          http.StatusServiceUnavailable,
		}
	}

	return nil, nil
}
