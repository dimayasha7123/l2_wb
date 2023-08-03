package app

import "fmt"

// EventsForMonth method
func (a *App) EventsForMonth(req EventsReq) (EventsResp, error) {
	to := req.FirstDay.AddDate(0, 1, 0)
	events, err := a.repository.ReadEventsFromTo(req.UserID, &req.FirstDay, &to)
	if err != nil {
		return EventsResp{}, fmt.Errorf("can't read events from repo: %v", err)
	}

	return EventsResp{Events: events}, nil
}
