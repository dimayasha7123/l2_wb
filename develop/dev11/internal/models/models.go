package models

import "time"

// User struct
type User struct {
	ID       int64
	Nickname string
	Events   map[int64]Event
}

// Event struct
type Event struct {
	ID          int64
	Title       string
	Description string
	Date        time.Time
}
