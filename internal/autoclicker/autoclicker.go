// Package autoclicker provides automated mouse clicking functionality.
package autoclicker

import (
	"math/rand"
	"sync"
	"time"

	"simplyauto/internal/input"
	"simplyauto/pkg/events"
)

type State int

const (
	StateStopped State = iota
	StateRunning
)

type AutoClicker struct {
	state      State
	config     *Config
	clickCount int
	stopChan   chan struct{}
	mu         sync.RWMutex
}

func New() *AutoClicker {
	return &AutoClicker{
		state:  StateStopped,
		config: DefaultConfig(),
	}
}

func (a *AutoClicker) SetConfig(cfg *Config) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := cfg.Validate(); err != nil {
		return err
	}
	a.config = cfg
	return nil
}

func (a *AutoClicker) GetConfig() Config {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return *a.config
}

func (a *AutoClicker) GetState() State {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.state
}

func (a *AutoClicker) GetClickCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.clickCount
}

func (a *AutoClicker) IsRunning() bool {
	return a.GetState() == StateRunning
}

func (a *AutoClicker) Start() error {
	a.mu.Lock()
	if a.state == StateRunning {
		a.mu.Unlock()
		return nil
	}

	a.state = StateRunning
	a.clickCount = 0
	a.stopChan = make(chan struct{})
	stopChan := a.stopChan
	cfg := *a.config
	a.mu.Unlock()

	go a.clickLoop(cfg, stopChan)
	return nil
}

func (a *AutoClicker) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.state != StateRunning {
		return
	}

	close(a.stopChan)
	a.state = StateStopped
}

func (a *AutoClicker) Toggle() {
	if a.IsRunning() {
		a.Stop()
	} else {
		a.Start()
	}
}

func (a *AutoClicker) clickLoop(cfg Config, stopChan <-chan struct{}) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-stopChan:
			return
		default:
		}

		if cfg.PositionMode == PositionFixed {
			input.Move(cfg.FixedX, cfg.FixedY)
		}

		input.Click(cfg.Button.String(), cfg.ClickType == events.ClickDouble)

		a.mu.Lock()
		a.clickCount++
		count := a.clickCount
		a.mu.Unlock()

		if cfg.RepeatMode == RepeatCount && count >= cfg.RepeatCount {
			a.Stop()
			return
		}

		interval := cfg.Interval
		if cfg.RandomOffsetMs > 0 {
			offset := rng.Intn(cfg.RandomOffsetMs*2) - cfg.RandomOffsetMs
			interval += time.Duration(offset) * time.Millisecond
			if interval < time.Millisecond {
				interval = time.Millisecond
			}
		}

		select {
		case <-stopChan:
			return
		case <-time.After(interval):
		}
	}
}
