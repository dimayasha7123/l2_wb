package memrepo

import (
	"sync"

	"l2_wb/develop/dev11/internal/models"
	"l2_wb/develop/dev11/internal/utils/counter"
)

// Repo struct
type Repo struct {
	mutex         *sync.RWMutex
	users         map[int64]models.User
	usersCounter  *counter.Counter
	eventsCounter *counter.Counter
}

// New constructor for Repo
func New() *Repo {
	return &Repo{
		mutex:         &sync.RWMutex{},
		users:         make(map[int64]models.User),
		usersCounter:  counter.New(0),
		eventsCounter: counter.New(0),
	}
}
