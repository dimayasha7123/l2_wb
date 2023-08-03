package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"l2_wb/develop/dev11/internal/app"
	"l2_wb/develop/dev11/internal/inputs/httpserver/middlewares"
)

// EventsAllHandler handler struct
type EventsAllHandler struct {
	service *app.App
}

// NewEventsAllHandler constructor
func NewEventsAllHandler(service *app.App) *EventsAllHandler {
	return &EventsAllHandler{service: service}
}

func (h *EventsAllHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (any, *middlewares.ServeHTTPError) {
	valErr := validateMethod(http.MethodGet, r.Method)
	if valErr != nil {
		return nil, valErr
	}

	req, err := NewEventsAllReqFromURLValues(r.URL.Query())
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       fmt.Sprintf("not enough params: %v", err),
			Code:          http.StatusBadRequest,
		}
	}

	userID, err := strconv.Atoi(req.UserID)
	if err != nil {
		if err != nil {
			return nil, &middlewares.ServeHTTPError{
				InternalError: err,
				Message:       "can't convert user_id from string to integer",
				Code:          http.StatusBadRequest,
			}
		}
	}

	appReq := app.EventsAllReq{
		UserID: int64(userID),
	}

	appResp, err := h.service.EventsAll(appReq)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "service error",
			Code:          http.StatusServiceUnavailable,
		}
	}

	return EventsResp{Events: convertEvents(appResp.Events)}, nil
}
