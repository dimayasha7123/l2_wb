package app

import "fmt"

// EventsAllReq request struct
type EventsAllReq struct {
	UserID int64
}

// EventsAll method
func (a *App) EventsAll(req EventsAllReq) (EventsResp, error) {
	events, err := a.repository.ReadEventsFromTo(req.UserID, nil, nil)
	if err != nil {
		return EventsResp{}, fmt.Errorf("can't read events from repo: %v", err)
	}

	return EventsResp{Events: events}, nil
}
