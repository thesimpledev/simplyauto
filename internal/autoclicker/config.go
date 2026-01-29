package autoclicker

import (
	"errors"
	"time"

	"simplyauto/pkg/events"
)

type RepeatMode int

const (
	RepeatUntilStopped RepeatMode = iota
	RepeatCount
)

type PositionMode int

const (
	PositionCurrent PositionMode = iota
	PositionFixed
)

type Config struct {
	Interval       time.Duration
	RandomOffsetMs int
	Button         events.MouseButton
	ClickType      events.ClickType
	RepeatMode     RepeatMode
	RepeatCount    int
	PositionMode   PositionMode
	FixedX         int
	FixedY         int
}

func DefaultConfig() *Config {
	return &Config{
		Interval:       100 * time.Millisecond,
		RandomOffsetMs: 0,
		Button:         events.MouseLeft,
		ClickType:      events.ClickSingle,
		RepeatMode:     RepeatUntilStopped,
		RepeatCount:    1,
		PositionMode:   PositionCurrent,
		FixedX:         0,
		FixedY:         0,
	}
}

func (c *Config) Validate() error {
	if c.Interval < time.Millisecond {
		return errors.New("interval must be at least 1 millisecond")
	}
	if c.RandomOffsetMs < 0 {
		return errors.New("random offset cannot be negative")
	}
	if c.RepeatCount < 1 {
		return errors.New("repeat count must be at least 1")
	}
	return nil
}
