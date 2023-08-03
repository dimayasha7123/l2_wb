package memrepo

import (
	"testing"
	"time"

	"l2_wb/develop/dev11/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestRepo_ReadEventsFromTo(t *testing.T) {
	userID := int64(1)

	r := New()
	r.users = map[int64]models.User{
		userID: {
			ID:       userID,
			Nickname: "dimya",
			Events:   map[int64]models.Event{},
		},
	}

	t.Run("user not found", func(t *testing.T) {
		badUserID := int64(42)
		events, err := r.ReadEventsFromTo(badUserID, nil, nil)
		assert.Nil(t, events)
		assert.ErrorIs(t, err, NewUserNotFoundErr(badUserID))
	})

	type args struct {
		from *time.Time
		to   *time.Time
	}
	tests := []struct {
		name       string
		events     map[int64]models.Event
		args       args
		wantEvents []models.Event
	}{
		{
			name:   "nil borders, no events",
			events: map[int64]models.Event{},
			args: args{
				from: nil,
				to:   nil,
			},
			wantEvents: []models.Event{},
		},
		{
			name: "nil borders, some events",
			events: map[int64]models.Event{
				1: {
					ID:          1,
					Title:       "title 1",
					Description: "desc 1",
					Date:        time.Date(2003, 3, 2, 0, 0, 0, 0, time.UTC),
				},
				2: {
					ID:          2,
					Title:       "title 2",
					Description: "desc 2",
					Date:        time.Date(2003, 3, 3, 0, 0, 0, 0, time.UTC),
				},
			},
			args: args{
				from: nil,
				to:   nil,
			},
			wantEvents: []models.Event{
				{
					ID:          1,
					Title:       "title 1",
					Description: "desc 1",
					Date:        time.Date(2003, 3, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          2,
					Title:       "title 2",
					Description: "desc 2",
					Date:        time.Date(2003, 3, 3, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "events and borders",
			events: map[int64]models.Event{
				1: {
					ID:          1,
					Title:       "title 1",
					Description: "desc 1",
					Date:        time.Date(2003, 3, 2, 0, 0, 0, 0, time.UTC),
				},
				2: {
					ID:          2,
					Title:       "title 2",
					Description: "desc 2",
					Date:        time.Date(2003, 3, 3, 0, 0, 0, 0, time.UTC),
				},
				3: {
					ID:          3,
					Title:       "title 3",
					Description: "desc 3",
					Date:        time.Date(2003, 3, 4, 0, 0, 0, 0, time.UTC),
				},
				4: {
					ID:          4,
					Title:       "title 4",
					Description: "desc 4",
					Date:        time.Date(2003, 3, 5, 0, 0, 0, 0, time.UTC),
				},
				5: {
					ID:          5,
					Title:       "title 5",
					Description: "desc 5",
					Date:        time.Date(2003, 3, 6, 0, 0, 0, 0, time.UTC),
				},
			},
			args: args{
				from: timePointer(2003, 3, 3),
				to:   timePointer(2003, 3, 5),
			},
			wantEvents: []models.Event{
				{
					ID:          2,
					Title:       "title 2",
					Description: "desc 2",
					Date:        time.Date(2003, 3, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          3,
					Title:       "title 3",
					Description: "desc 3",
					Date:        time.Date(2003, 3, 4, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "events and right border",
			events: map[int64]models.Event{
				1: {
					ID:          1,
					Title:       "title 1",
					Description: "desc 1",
					Date:        time.Date(2003, 3, 2, 0, 0, 0, 0, time.UTC),
				},
				2: {
					ID:          2,
					Title:       "title 2",
					Description: "desc 2",
					Date:        time.Date(2003, 3, 3, 0, 0, 0, 0, time.UTC),
				},
				3: {
					ID:          3,
					Title:       "title 3",
					Description: "desc 3",
					Date:        time.Date(2003, 3, 4, 0, 0, 0, 0, time.UTC),
				},
				4: {
					ID:          4,
					Title:       "title 4",
					Description: "desc 4",
					Date:        time.Date(2003, 3, 5, 0, 0, 0, 0, time.UTC),
				},
				5: {
					ID:          5,
					Title:       "title 5",
					Description: "desc 5",
					Date:        time.Date(2003, 3, 6, 0, 0, 0, 0, time.UTC),
				},
			},
			args: args{
				from: nil,
				to:   timePointer(2003, 3, 5),
			},
			wantEvents: []models.Event{
				{
					ID:          1,
					Title:       "title 1",
					Description: "desc 1",
					Date:        time.Date(2003, 3, 2, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          2,
					Title:       "title 2",
					Description: "desc 2",
					Date:        time.Date(2003, 3, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          3,
					Title:       "title 3",
					Description: "desc 3",
					Date:        time.Date(2003, 3, 4, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "events and left border",
			events: map[int64]models.Event{
				1: {
					ID:          1,
					Title:       "title 1",
					Description: "desc 1",
					Date:        time.Date(2003, 3, 2, 0, 0, 0, 0, time.UTC),
				},
				2: {
					ID:          2,
					Title:       "title 2",
					Description: "desc 2",
					Date:        time.Date(2003, 3, 3, 0, 0, 0, 0, time.UTC),
				},
				3: {
					ID:          3,
					Title:       "title 3",
					Description: "desc 3",
					Date:        time.Date(2003, 3, 4, 0, 0, 0, 0, time.UTC),
				},
				4: {
					ID:          4,
					Title:       "title 4",
					Description: "desc 4",
					Date:        time.Date(2003, 3, 5, 0, 0, 0, 0, time.UTC),
				},
				5: {
					ID:          5,
					Title:       "title 5",
					Description: "desc 5",
					Date:        time.Date(2003, 3, 6, 0, 0, 0, 0, time.UTC),
				},
			},
			args: args{
				from: timePointer(2003, 3, 3),
				to:   nil,
			},
			wantEvents: []models.Event{
				{
					ID:          2,
					Title:       "title 2",
					Description: "desc 2",
					Date:        time.Date(2003, 3, 3, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          3,
					Title:       "title 3",
					Description: "desc 3",
					Date:        time.Date(2003, 3, 4, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          4,
					Title:       "title 4",
					Description: "desc 4",
					Date:        time.Date(2003, 3, 5, 0, 0, 0, 0, time.UTC),
				},
				{
					ID:          5,
					Title:       "title 5",
					Description: "desc 5",
					Date:        time.Date(2003, 3, 6, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := r.users[userID]
			user.Events = tt.events
			r.users[userID] = user

			events, err := r.ReadEventsFromTo(userID, tt.args.from, tt.args.to)

			assert.NoError(t, err)
			assert.ElementsMatch(t, tt.wantEvents, events)
		})
	}

}

func timePointer(year, month, day int) *time.Time {
	ret := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return &ret
}
