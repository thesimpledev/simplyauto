package recorder

import (
	"sync"
	"time"

	"simplyauto/internal/input"
	"simplyauto/internal/storage"
	"simplyauto/pkg/events"
)

type PlayerState int

const (
	PlayerIdle PlayerState = iota
	PlayerPlaying
	PlayerPaused
)

type LoopMode int

const (
	LoopOnce LoopMode = iota
	LoopCount
	LoopContinuous
)

type PlaybackConfig struct {
	Speed     float64
	LoopMode  LoopMode
	LoopCount int
}

func DefaultPlaybackConfig() PlaybackConfig {
	return PlaybackConfig{
		Speed:     1.0,
		LoopMode:  LoopOnce,
		LoopCount: 1,
	}
}

type Player struct {
	state       PlayerState
	recording   *storage.Recording
	config      PlaybackConfig
	currentIdx  int
	currentLoop int
	stopChan    chan struct{}
	pauseChan   chan struct{}
	resumeChan  chan struct{}
	mu          sync.RWMutex
	OnComplete  func()
}

func NewPlayer() *Player {
	return &Player{
		state:  PlayerIdle,
		config: DefaultPlaybackConfig(),
	}
}

func (p *Player) GetState() PlayerState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

func (p *Player) IsPlaying() bool {
	state := p.GetState()
	return state == PlayerPlaying || state == PlayerPaused
}

func (p *Player) GetProgress() (currentEvent, totalEvents, currentLoop int) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.recording == nil {
		return 0, 0, 0
	}
	return p.currentIdx, len(p.recording.Events), p.currentLoop
}

func (p *Player) Play(recording *storage.Recording, config PlaybackConfig) error {
	p.mu.Lock()
	if p.state == PlayerPlaying {
		p.mu.Unlock()
		return nil
	}

	if config.Speed <= 0 {
		config.Speed = 1.0
	}

	p.recording = recording
	p.config = config
	p.currentIdx = 0
	p.currentLoop = 1
	p.state = PlayerPlaying
	p.stopChan = make(chan struct{})
	p.pauseChan = make(chan struct{})
	p.resumeChan = make(chan struct{})

	// Capture session state for the goroutine
	rec := p.recording
	cfg := p.config
	stopChan := p.stopChan
	pauseChan := p.pauseChan
	resumeChan := p.resumeChan
	onComplete := p.OnComplete
	p.mu.Unlock()

	go p.playbackLoop(rec, cfg, stopChan, pauseChan, resumeChan, onComplete)
	return nil
}

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == PlayerIdle {
		return
	}

	close(p.stopChan)
	p.state = PlayerIdle
}

func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != PlayerPlaying {
		return
	}

	p.pauseChan <- struct{}{}
	p.state = PlayerPaused
}

func (p *Player) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != PlayerPaused {
		return
	}

	p.resumeChan <- struct{}{}
	p.state = PlayerPlaying
}

func (p *Player) playbackLoop(rec *storage.Recording, cfg PlaybackConfig, stopChan, pauseChan, resumeChan chan struct{}, onComplete func()) {
	currentLoop := 1

	for {
		shouldLoop := cfg.LoopMode == LoopContinuous ||
			(cfg.LoopMode == LoopCount && currentLoop <= cfg.LoopCount) ||
			(cfg.LoopMode == LoopOnce && currentLoop == 1)

		if !shouldLoop {
			p.Stop()
			if onComplete != nil {
				onComplete()
			}
			return
		}

		p.playEvents(rec, cfg.Speed, stopChan, pauseChan, resumeChan)

		select {
		case <-stopChan:
			return
		default:
		}

		currentLoop++
		p.mu.Lock()
		p.currentLoop = currentLoop
		p.currentIdx = 0
		p.mu.Unlock()
	}
}

func (p *Player) playEvents(rec *storage.Recording, speed float64, stopChan, pauseChan, resumeChan chan struct{}) {
	eventCount := len(rec.Events)
	var lastTimestamp time.Duration

	for i := 0; i < eventCount; i++ {
		select {
		case <-stopChan:
			return
		default:
		}

		select {
		case <-pauseChan:
			select {
			case <-stopChan:
				return
			case <-resumeChan:
			}
		default:
		}

		event := rec.Events[i]

		if i > 0 {
			delay := event.Timestamp - lastTimestamp
			adjustedDelay := time.Duration(float64(delay) / speed)
			if adjustedDelay > 0 {
				select {
				case <-stopChan:
					return
				case <-time.After(adjustedDelay):
				}
			}
		}
		lastTimestamp = event.Timestamp

		p.executeEvent(event)

		p.mu.Lock()
		p.currentIdx = i + 1
		p.mu.Unlock()
	}
}

func (p *Player) executeEvent(event events.InputEvent) {
	switch event.Type {
	case events.EventMouseMove:
		input.Move(event.X, event.Y)

	case events.EventMouseLeftDown:
		input.Move(event.X, event.Y)
		input.Toggle("left", "down")

	case events.EventMouseLeftUp:
		input.Move(event.X, event.Y)
		input.Toggle("left", "up")

	case events.EventMouseRightDown:
		input.Move(event.X, event.Y)
		input.Toggle("right", "down")

	case events.EventMouseRightUp:
		input.Move(event.X, event.Y)
		input.Toggle("right", "up")

	case events.EventMouseMiddleDown:
		input.Move(event.X, event.Y)
		input.Toggle("center", "down")

	case events.EventMouseMiddleUp:
		input.Move(event.X, event.Y)
		input.Toggle("center", "up")

	case events.EventMouseWheel:
		if event.Delta > 0 {
			input.ScrollDir(event.Delta, "up")
		} else {
			input.ScrollDir(-event.Delta, "down")
		}

	case events.EventKeyDown:
		input.KeyDownVK(event.KeyCode)

	case events.EventKeyUp:
		input.KeyUpVK(event.KeyCode)
	}
}
