package app

import "fmt"

// EventsForDay method
func (a *App) EventsForDay(req EventsReq) (EventsResp, error) {
	to := req.FirstDay.AddDate(0, 0, 1)
	events, err := a.repository.ReadEventsFromTo(req.UserID, &req.FirstDay, &to)
	if err != nil {
		return EventsResp{}, fmt.Errorf("can't read events from repo: %v", err)
	}

	return EventsResp{Events: events}, nil
}
