package memrepo

import (
	"time"

	"l2_wb/develop/dev11/internal/models"
)

// UpdateEvent method
func (r *Repo) UpdateEvent(userID int64, eventID int64, eventTitle string, eventDescription string, eventDate time.Time) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return NewUserNotFoundErr(userID)
	}

	event, ok := user.Events[eventID]
	if !ok {
		return NewEventNotFoundErr(userID, eventID)
	}

	event = models.Event{
		ID:          eventID,
		Title:       eventTitle,
		Description: eventDescription,
		Date:        eventDate,
	}

	user.Events[eventID] = event
	//r.users[userID] = user

	return nil
}
