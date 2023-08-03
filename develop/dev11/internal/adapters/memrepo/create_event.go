package memrepo

import (
	"time"

	"l2_wb/develop/dev11/internal/models"
)

// CreateEvent method
func (r *Repo) CreateEvent(userID int64, eventTitle string, eventDescription string, eventDate time.Time) (id int64, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return 0, NewUserNotFoundErr(userID)
	}

	eventID := r.eventsCounter.Inc()
	event := models.Event{
		ID:          eventID,
		Title:       eventTitle,
		Description: eventDescription,
		Date:        eventDate,
	}

	user.Events[eventID] = event
	//r.users[user.ID] = user

	return eventID, nil
}
