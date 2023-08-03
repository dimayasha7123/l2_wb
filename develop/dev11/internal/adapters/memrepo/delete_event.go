package memrepo

// DeleteEvent method
func (r *Repo) DeleteEvent(userID int64, eventID int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	user, ok := r.users[userID]
	if !ok {
		return NewUserNotFoundErr(userID)
	}

	_, ok = user.Events[eventID]
	if !ok {
		return NewEventNotFoundErr(userID, eventID)
	}

	delete(user.Events, eventID)
	//r.users[userID] = user

	return nil
}
