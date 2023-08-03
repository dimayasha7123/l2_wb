package memrepo

import (
	"testing"
	"time"

	"l2_wb/develop/dev11/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestRepo_DeleteEvent(t *testing.T) {
	r := New()
	r.users = map[int64]models.User{
		1: {
			ID:       1,
			Nickname: "dimya",
			Events: map[int64]models.Event{
				1: {
					ID:          1,
					Title:       "event 1",
					Description: "desc 1",
					Date:        time.Date(2023, 3, 4, 0, 0, 0, 0, time.UTC),
				},
				2: {
					ID:          2,
					Title:       "event 3",
					Description: "desc 3",
					Date:        time.Date(2023, 3, 6, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	t.Run("bad user id", func(t *testing.T) {
		badUserID := int64(42)
		err := r.DeleteEvent(badUserID, 3)
		assert.ErrorIs(t, err, NewUserNotFoundErr(badUserID))
	})

	goodUserID := int64(1)
	t.Run("bad event id", func(t *testing.T) {
		badEventID := int64(42)
		err := r.DeleteEvent(goodUserID, badEventID)
		assert.ErrorIs(t, err, NewEventNotFoundErr(goodUserID, badEventID))
	})

	t.Run("good case", func(t *testing.T) {
		err := r.DeleteEvent(goodUserID, 1)
		assert.NoError(t, err)
	})

	t.Run("check data in repo", func(t *testing.T) {
		expectedUsers := map[int64]models.User{
			1: {
				ID:       1,
				Nickname: "dimya",
				Events: map[int64]models.Event{
					2: {
						ID:          2,
						Title:       "event 3",
						Description: "desc 3",
						Date:        time.Date(2023, 3, 6, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		}
		assert.Equal(t, expectedUsers, r.users)
	})
}
