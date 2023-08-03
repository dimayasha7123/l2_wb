package memrepo

import "l2_wb/develop/dev11/internal/models"

// CreateUser method
func (r *Repo) CreateUser(nickname string) (id int64, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	newID := r.usersCounter.Inc()
	r.users[newID] = models.User{
		ID:       newID,
		Nickname: nickname,
		Events:   make(map[int64]models.Event),
	}

	return newID, nil
}
