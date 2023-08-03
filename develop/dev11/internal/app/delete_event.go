package app

import "fmt"

// DeleteEventReq request struct
type DeleteEventReq struct {
	UserID  int64
	EventID int64
}

// DeleteEvent method
func (a *App) DeleteEvent(req DeleteEventReq) error {
	err := a.repository.DeleteEvent(req.UserID, req.EventID)
	if err != nil {
		return fmt.Errorf("can't delete event in repo: %v", err)
	}

	return nil
}
