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

// CreateEventReq request struct
type CreateEventReq struct {
	UserID      *int64 `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

// Validate method
func (r CreateEventReq) Validate() error {
	fields := make([]string, 0, 4)
	if r.UserID == nil {
		fields = append(fields, "user_id")
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

// CreateEventResp response struct
type CreateEventResp struct {
	EventID int64 `json:"event_id"`
}

// CreateEventHandler handler struct
type CreateEventHandler struct {
	service *app.App
}

// NewCreateEventHandler constructor
func NewCreateEventHandler(service *app.App) *CreateEventHandler {
	return &CreateEventHandler{service: service}
}

func (h *CreateEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (any, *middlewares.ServeHTTPError) {
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

	var req CreateEventReq
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

	appReq := app.CreateEventReq{
		UserID:           *req.UserID,
		EventTitle:       req.Title,
		EventDescription: req.Description,
		EventDate:        date,
	}

	// TODO: здесь можно было бы сделать обработку ошибок, и отрабатывать не найденные user_id и event_id как 400
	appResp, err := h.service.CreateEvent(appReq)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "service error",
			Code:          http.StatusServiceUnavailable,
		}
	}

	return CreateEventResp{EventID: appResp.EventID}, nil
}
