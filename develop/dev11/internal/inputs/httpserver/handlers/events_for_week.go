package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"l2_wb/develop/dev11/internal/app"
	"l2_wb/develop/dev11/internal/inputs/httpserver/middlewares"
	"l2_wb/develop/dev11/internal/utils/converters"
)

// EventsForWeekHandler handler struct
type EventsForWeekHandler struct {
	service *app.App
}

// NewEventsForWeekHandler constructor
func NewEventsForWeekHandler(service *app.App) *EventsForWeekHandler {
	return &EventsForWeekHandler{service: service}
}

func (h *EventsForWeekHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) (any, *middlewares.ServeHTTPError) {
	valErr := validateMethod(http.MethodGet, r.Method)
	if valErr != nil {
		return nil, valErr
	}

	req, err := NewEventsForPeriodReqFromURLValues(r.URL.Query())
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

	date, err := converters.StrToDate(req.Date)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "can't convert date from string",
			Code:          http.StatusBadRequest,
		}
	}

	appReq := app.EventsReq{
		UserID:   int64(userID),
		FirstDay: date,
	}

	appResp, err := h.service.EventsForWeek(appReq)
	if err != nil {
		return nil, &middlewares.ServeHTTPError{
			InternalError: err,
			Message:       "service error",
			Code:          http.StatusServiceUnavailable,
		}
	}

	return EventsResp{Events: convertEvents(appResp.Events)}, nil
}
