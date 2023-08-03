package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"l2_wb/develop/dev11/internal/inputs/httpserver/middlewares"
	"l2_wb/develop/dev11/internal/models"
	"l2_wb/develop/dev11/internal/utils/converters"
)

func validateMethod(need, actual string) *middlewares.ServeHTTPError {
	if actual != need {
		message := fmt.Sprintf("using wrong method, need %s, but actual method is %s", need, actual)
		return &middlewares.ServeHTTPError{
			InternalError: errors.New(message),
			Message:       message,
			Code:          http.StatusBadRequest,
		}
	}
	return nil
}

// NewMissingFieldsErr constructor
func NewMissingFieldsErr(fields []string) error {
	return fmt.Errorf("missing fields: %s", strings.Join(fields, ", "))
}

func convertEvents(events []models.Event) []Event {
	ret := make([]Event, 0, len(events))

	for _, event := range events {
		ret = append(ret, Event{
			ID:          event.ID,
			Title:       event.Title,
			Description: event.Description,
			Date:        converters.DateToStr(event.Date),
		})
	}

	return ret
}
