// Package events defines input event types for recording and playback.
package events

import "time"

type EventType uint16

// Windows message codes for input events
const (
	EventMouseMove       EventType = 0x0200
	EventMouseLeftDown   EventType = 0x0201
	EventMouseLeftUp     EventType = 0x0202
	EventMouseRightDown  EventType = 0x0204
	EventMouseRightUp    EventType = 0x0205
	EventMouseMiddleDown EventType = 0x0207
	EventMouseMiddleUp   EventType = 0x0208
	EventMouseWheel      EventType = 0x020A
	EventKeyDown         EventType = 0x0100
	EventKeyUp           EventType = 0x0101
)

type InputEvent struct {
	Type      EventType     `json:"type"`
	Timestamp time.Duration `json:"timestamp"`
	X         int           `json:"x,omitempty"`
	Y         int           `json:"y,omitempty"`
	KeyCode   uint16        `json:"keyCode,omitempty"`
	ScanCode  uint16        `json:"scanCode,omitempty"`
	Delta     int           `json:"delta,omitempty"`
}

type MouseButton int

const (
	MouseLeft MouseButton = iota
	MouseRight
	MouseMiddle
)

func (b MouseButton) String() string {
	switch b {
	case MouseLeft:
		return "left"
	case MouseRight:
		return "right"
	case MouseMiddle:
		return "middle"
	default:
		return "unknown"
	}
}

type ClickType int

const (
	ClickSingle ClickType = iota
	ClickDouble
)

func (c ClickType) String() string {
	switch c {
	case ClickSingle:
		return "single"
	case ClickDouble:
		return "double"
	default:
		return "unknown"
	}
}
