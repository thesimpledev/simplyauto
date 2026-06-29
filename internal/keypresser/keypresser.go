// Package keypresser provides automated keyboard key pressing functionality.
package keypresser

import (
	"math/rand"
	"sync"
	"time"

	"simplyauto/internal/input"
)

type State int

const (
	StateStopped State = iota
	StateRunning
)

type KeyPresser struct {
	state      State
	config     *Config
	pressCount int
	stopChan   chan struct{}
	mu         sync.RWMutex
}

func New() *KeyPresser {
	return &KeyPresser{
		state:  StateStopped,
		config: DefaultConfig(),
	}
}

func (k *KeyPresser) SetConfig(cfg *Config) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	if err := cfg.Validate(); err != nil {
		return err
	}
	k.config = cfg
	return nil
}

func (k *KeyPresser) GetConfig() Config {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return *k.config
}

func (k *KeyPresser) GetState() State {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.state
}

func (k *KeyPresser) GetPressCount() int {
	k.mu.RLock()
	defer k.mu.RUnlock()
	return k.pressCount
}

func (k *KeyPresser) IsRunning() bool {
	return k.GetState() == StateRunning
}

func (k *KeyPresser) Start() error {
	k.mu.Lock()
	if k.state == StateRunning {
		k.mu.Unlock()
		return nil
	}

	k.state = StateRunning
	k.pressCount = 0
	k.stopChan = make(chan struct{})
	stopChan := k.stopChan
	cfg := *k.config
	k.mu.Unlock()

	go k.pressLoop(cfg, stopChan)
	return nil
}

func (k *KeyPresser) Stop() {
	k.mu.Lock()
	defer k.mu.Unlock()

	if k.state != StateRunning {
		return
	}

	close(k.stopChan)
	k.state = StateStopped
}

func (k *KeyPresser) Toggle() {
	if k.IsRunning() {
		k.Stop()
	} else {
		k.Start()
	}
}

func (k *KeyPresser) pressLoop(cfg Config, stopChan <-chan struct{}) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-stopChan:
			return
		default:
		}

		// Press each character in the configured text once, in order.
		for _, r := range cfg.Keys {
			select {
			case <-stopChan:
				return
			default:
			}
			input.KeyPressChar(r)
		}

		k.mu.Lock()
		k.pressCount++
		count := k.pressCount
		k.mu.Unlock()

		if cfg.RepeatMode == RepeatCount && count >= cfg.RepeatCount {
			k.Stop()
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
