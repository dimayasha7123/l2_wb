package app

import (
	"time"

	"l2_wb/develop/dev11/internal/models"
)

// EventsReq request struct
type EventsReq struct {
	UserID   int64
	FirstDay time.Time
}

// EventsResp response struct
type EventsResp struct {
	Events []models.Event
}
