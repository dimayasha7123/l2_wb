package memrepo

import (
	"testing"
	"time"

	"l2_wb/develop/dev11/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepo_CreateEvent(t *testing.T) {
	r := New()

	userID, err := r.CreateUser("dimya")
	require.NoError(t, err, "can't continue without creating user")

	date1 := time.Date(2023, 3, 4, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2023, 3, 5, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2023, 3, 6, 0, 0, 0, 0, time.UTC)

	t.Run("create good event 1", func(t *testing.T) {
		eventID, err := r.CreateEvent(userID, "event 1", "desc 1", date1)
		assert.NoError(t, err)
		assert.EqualValues(t, 1, eventID)
	})

	t.Run("create bad event", func(t *testing.T) {
		badUserID := int64(42)
		eventID, err := r.CreateEvent(badUserID, "event 2", "desc 2", date2)
		assert.ErrorIs(t, err, NewUserNotFoundErr(badUserID))
		assert.EqualValues(t, 0, eventID)
	})

	t.Run("create good event 2", func(t *testing.T) {
		eventID, err := r.CreateEvent(userID, "event 3", "desc 3", date3)
		assert.NoError(t, err)
		assert.EqualValues(t, 2, eventID)
	})

	t.Run("check data in repo", func(t *testing.T) {
		expectedUsers := map[int64]models.User{
			1: {
				ID:       1,
				Nickname: "dimya",
				Events: map[int64]models.Event{
					1: {
						ID:          1,
						Title:       "event 1",
						Description: "desc 1",
						Date:        date1,
					},
					2: {
						ID:          2,
						Title:       "event 3",
						Description: "desc 3",
						Date:        date3,
					},
				},
			},
		}
		assert.Equal(t, expectedUsers, r.users)
	})
}
