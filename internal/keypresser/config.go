package keypresser

import (
	"errors"
	"fmt"
	"time"

	"simplyauto/internal/input"
)

type RepeatMode int

const (
	RepeatUntilStopped RepeatMode = iota
	RepeatCount
)

type Config struct {
	Interval       time.Duration
	RandomOffsetMs int
	Keys           string
	RepeatMode     RepeatMode
	RepeatCount    int
}

func DefaultConfig() *Config {
	return &Config{
		Interval:       100 * time.Millisecond,
		RandomOffsetMs: 0,
		Keys:           "",
		RepeatMode:     RepeatUntilStopped,
		RepeatCount:    1,
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
	if c.Keys == "" {
		return errors.New("enter at least one key to press")
	}
	for _, r := range c.Keys {
		if !input.CanType(r) {
			return fmt.Errorf("cannot type character %q", r)
		}
	}
	return nil
}
