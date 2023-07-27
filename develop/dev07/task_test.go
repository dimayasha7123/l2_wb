package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_or(t *testing.T) {
	tests := []struct {
		name     string
		channels []<-chan interface{}
		waitTime time.Duration
	}{
		{
			name:     "nil",
			channels: nil,
			waitTime: 0 * time.Second,
		},
		{
			name:     "zero channels",
			channels: []<-chan interface{}{},
			waitTime: 0 * time.Second,
		},
		{
			name: "one channel",
			channels: []<-chan interface{}{
				sig(1 * time.Second),
			},
			waitTime: 1 * time.Second,
		},
		{
			name: "few channel",
			channels: []<-chan interface{}{
				sig(5 * time.Second),
				sig(4 * time.Second),
				sig(3 * time.Second),
				sig(2 * time.Second),
				sig(1 * time.Second),
			},
			waitTime: 1 * time.Second,
		},
		{
			name: "few channel with one time",
			channels: []<-chan interface{}{
				sig(5 * time.Second),
				sig(4 * time.Second),
				sig(3 * time.Second),
				sig(2 * time.Second),
				sig(2 * time.Second),
			},
			waitTime: 2 * time.Second,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			init := time.Now()

			<-or(tt.channels...)

			assert.InDelta(t, tt.waitTime, time.Since(init), float64(1*time.Millisecond))
		})
	}
}
