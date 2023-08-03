package memrepo

import (
	"testing"

	"l2_wb/develop/dev11/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestRepo_CreateUser(t *testing.T) {
	r := New()

	t.Run("create user 1", func(t *testing.T) {
		id, err := r.CreateUser("dimya")
		assert.NoError(t, err)
		assert.EqualValues(t, 1, id)
	})

	t.Run("create user 2 (duplicate)", func(t *testing.T) {
		id, err := r.CreateUser("dimya")
		assert.NoError(t, err)
		assert.EqualValues(t, 2, id)
	})

	t.Run("create user 3", func(t *testing.T) {
		id, err := r.CreateUser("vasya")
		assert.NoError(t, err)
		assert.EqualValues(t, 3, id)
	})

	t.Run("check data in repo", func(t *testing.T) {
		expectedUsers := map[int64]models.User{
			1: {
				ID:       1,
				Nickname: "dimya",
				Events:   map[int64]models.Event{},
			},
			2: {
				ID:       2,
				Nickname: "dimya",
				Events:   map[int64]models.Event{},
			},
			3: {
				ID:       3,
				Nickname: "vasya",
				Events:   map[int64]models.Event{},
			},
		}
		assert.Equal(t, expectedUsers, r.users)
	})
}
