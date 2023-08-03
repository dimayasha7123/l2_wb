package app

import (
	"fmt"
	"time"
)

// CreateEventReq request struct
type CreateEventReq struct {
	UserID           int64
	EventTitle       string
	EventDescription string
	EventDate        time.Time
}

// CreateEventResp response struct
type CreateEventResp struct {
	EventID int64
}

// CreateEvent method
func (a *App) CreateEvent(req CreateEventReq) (CreateEventResp, error) {
	id, err := a.repository.CreateEvent(req.UserID, req.EventTitle, req.EventDescription, req.EventDate)
	if err != nil {
		return CreateEventResp{}, fmt.Errorf("can't create event in repo: %v", err)
	}

	return CreateEventResp{EventID: id}, nil
}
