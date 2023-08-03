package app

import (
	"time"

	"l2_wb/develop/dev11/internal/models"
)

// Repository interface
type Repository interface {
	CreateUser(nickname string) (id int64, err error)
	CreateEvent(userID int64, eventTitle string, eventDescription string, eventDate time.Time) (id int64, err error)
	UpdateEvent(userID int64, eventID int64, eventTitle string, eventDescription string, eventDate time.Time) error
	DeleteEvent(userID int64, eventID int64) error
	ReadEventsFromTo(userID int64, from, to *time.Time) (events []models.Event, err error)
}
