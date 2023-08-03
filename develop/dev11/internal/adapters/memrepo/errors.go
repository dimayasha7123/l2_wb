package memrepo

import "fmt"

// UserNotFoundErr error
type UserNotFoundErr struct {
	userID int64
}

// NewUserNotFoundErr constructor
func NewUserNotFoundErr(userID int64) UserNotFoundErr {
	return UserNotFoundErr{userID: userID}
}

func (e UserNotFoundErr) Error() string {
	return fmt.Sprintf("user (id = %d) not found", e.userID)
}

// EventNotFoundErr error
type EventNotFoundErr struct {
	userID  int64
	eventID int64
}

// NewEventNotFoundErr constructor
func NewEventNotFoundErr(userID, eventID int64) EventNotFoundErr {
	return EventNotFoundErr{userID: userID, eventID: eventID}
}

func (e EventNotFoundErr) Error() string {
	return fmt.Sprintf("event (id = %d) for user (id = %d) not found", e.eventID, e.userID)
}
