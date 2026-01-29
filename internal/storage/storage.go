// Package storage provides recording persistence.
package storage

import (
	"time"

	"simplyauto/pkg/events"
)

type Recording struct {
	Version     string              `json:"version"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	Duration    time.Duration       `json:"duration"`
	Events      []events.InputEvent `json:"events"`
	Metadata    RecordingMetadata   `json:"metadata,omitempty"`
}

type RecordingMetadata struct {
	ScreenWidth  int    `json:"screenWidth,omitempty"`
	ScreenHeight int    `json:"screenHeight,omitempty"`
	AppVersion   string `json:"appVersion,omitempty"`
	EventCount   int    `json:"eventCount,omitempty"`
}

func NewRecording(name string) *Recording {
	return &Recording{
		Version:   "1.0",
		Name:      name,
		CreatedAt: time.Now(),
		Events:    make([]events.InputEvent, 0),
	}
}

func (r *Recording) AddEvent(event events.InputEvent) {
	r.Events = append(r.Events, event)
	if event.Timestamp > r.Duration {
		r.Duration = event.Timestamp
	}
}

func (r *Recording) Finalize() {
	r.Metadata.EventCount = len(r.Events)
}
