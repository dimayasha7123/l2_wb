package main

import (
	"testing"

	"github.com/beevik/ntp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPreciseTime(t *testing.T) {
	sampleTime, err := ntp.Time(ntpAddress)
	require.NoError(t, err, "can't get sample time")

	got, err := GetPreciseTime()
	require.NoError(t, err, "can't get test time")

	deltaMilli := 100.0
	assert.InDelta(t,
		sampleTime.UnixMilli(),
		got.UnixMilli(),
		deltaMilli,
	)
}
