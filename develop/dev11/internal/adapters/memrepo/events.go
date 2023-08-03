package memrepo

import (
	"time"

	"l2_wb/develop/dev11/internal/models"
)

// ReadEventsFromTo method
func (r *Repo) ReadEventsFromTo(userID int64, from, to *time.Time) (events []models.Event, err error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, ok := r.users[userID]
	if !ok {
		return nil, NewUserNotFoundErr(userID)
	}

	ret := make([]models.Event, 0)
	for _, event := range user.Events {
		if from != nil && event.Date.Before(*from) ||
			to != nil && !event.Date.Before(*to) {
			continue
		}
		ret = append(ret, event)
	}

	return ret, nil
}
