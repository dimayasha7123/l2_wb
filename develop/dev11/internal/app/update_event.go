package app

import (
	"fmt"
	"time"
)

// UpdateEventReq request struct
type UpdateEventReq struct {
	UserID      int64
	EventID     int64
	Title       string
	Description string
	Date        time.Time
}

// UpdateEvent method
func (a *App) UpdateEvent(req UpdateEventReq) error {
	err := a.repository.UpdateEvent(req.UserID, req.EventID, req.Title, req.Description, req.Date)
	if err != nil {
		return fmt.Errorf("can't update event in repo: %v", err)
	}
	return nil
}
