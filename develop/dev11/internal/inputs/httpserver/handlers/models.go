package handlers

import (
	"net/url"
)

// EventsForPeriodReq request struct
type EventsForPeriodReq struct {
	UserID string
	Date   string
}

const (
	userIDKey = "user_id"
	dateKey   = "date"
)

// NewEventsForPeriodReqFromURLValues constructor
func NewEventsForPeriodReqFromURLValues(values url.Values) (EventsForPeriodReq, error) {
	fields := make([]string, 0, 2)

	userID := values.Get(userIDKey)
	if userID == "" {
		fields = append(fields, userIDKey)
	}

	date := values.Get(dateKey)
	if date == "" {
		fields = append(fields, dateKey)
	}

	if len(fields) != 0 {
		return EventsForPeriodReq{}, NewMissingFieldsErr(fields)
	}

	return EventsForPeriodReq{
		UserID: userID,
		Date:   date,
	}, nil
}

// EventsAllReq request struct
type EventsAllReq struct {
	UserID string
}

// NewEventsAllReqFromURLValues constructor
func NewEventsAllReqFromURLValues(values url.Values) (EventsAllReq, error) {
	fields := make([]string, 0, 1)

	userID := values.Get(userIDKey)
	if userID == "" {
		fields = append(fields, userIDKey)
	}

	if len(fields) != 0 {
		return EventsAllReq{}, NewMissingFieldsErr(fields)
	}

	return EventsAllReq{UserID: userID}, nil
}

// EventsResp response struct
type EventsResp struct {
	Events []Event `json:"events"`
}

// Event struct
type Event struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"`
}
