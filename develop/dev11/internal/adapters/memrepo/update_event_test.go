package memrepo

import (
	"testing"
	"time"

	"l2_wb/develop/dev11/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestRepo_UpdateEvent(t *testing.T) {
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
		err := r.UpdateEvent(badUserID, 3, "some title", "some desc", time.Now())
		assert.ErrorIs(t, err, NewUserNotFoundErr(badUserID))
	})

	goodUserID := int64(1)
	t.Run("bad event id", func(t *testing.T) {
		badEventID := int64(42)
		err := r.UpdateEvent(goodUserID, badEventID, "some title", "some desc", time.Now())
		assert.ErrorIs(t, err, NewEventNotFoundErr(goodUserID, badEventID))
	})

	newDate := time.Date(2023, 3, 8, 0, 0, 0, 0, time.UTC)
	t.Run("good case", func(t *testing.T) {
		err := r.UpdateEvent(goodUserID, 1, "new title", "new desc", newDate)
		assert.NoError(t, err)
	})

	t.Run("check data in repo", func(t *testing.T) {
		expectedUsers := map[int64]models.User{
			1: {
				ID:       1,
				Nickname: "dimya",
				Events: map[int64]models.Event{
					1: {
						ID:          1,
						Title:       "new title",
						Description: "new desc",
						Date:        newDate,
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
		assert.Equal(t, expectedUsers, r.users)
	})
}
