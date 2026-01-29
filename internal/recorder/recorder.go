// Package recorder provides macro recording and playback.
package recorder

import (
	"sync"
	"time"

	"simplyauto/internal/hooks"
	"simplyauto/internal/storage"
	"simplyauto/pkg/events"
)

type RecorderState int

const (
	RecorderIdle RecorderState = iota
	RecorderRecording
)

type RecorderOptions struct {
	RecordMouse    bool
	RecordKeyboard bool
	FilterKeys     []hooks.Key
}

func DefaultRecorderOptions() RecorderOptions {
	return RecorderOptions{
		RecordMouse:    true,
		RecordKeyboard: true,
		FilterKeys:     []hooks.Key{hooks.KeyF9, hooks.KeyF10, hooks.KeyF11},
	}
}

type Recorder struct {
	state     RecorderState
	recording *storage.Recording
	capture   *hooks.InputCapture
	options   RecorderOptions
	startTime time.Time
	mu        sync.RWMutex
}

func NewRecorder(opts RecorderOptions) *Recorder {
	return &Recorder{
		state:   RecorderIdle,
		options: opts,
	}
}

func (r *Recorder) GetState() RecorderState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.state
}

func (r *Recorder) IsRecording() bool {
	return r.GetState() == RecorderRecording
}

func (r *Recorder) GetEventCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.recording == nil {
		return 0
	}
	return len(r.recording.Events)
}

func (r *Recorder) GetDuration() time.Duration {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.recording == nil || r.state != RecorderRecording {
		return 0
	}
	return time.Since(r.startTime)
}

func (r *Recorder) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state == RecorderRecording {
		return nil
	}

	r.recording = storage.NewRecording("Untitled Recording")
	r.startTime = time.Now()

	r.capture = hooks.NewInputCapture(hooks.CaptureOptions{
		CaptureMouse:    r.options.RecordMouse,
		CaptureKeyboard: r.options.RecordKeyboard,
		FilterKeys:      r.options.FilterKeys,
	})

	if err := r.capture.Start(r.onEvent); err != nil {
		return err
	}

	r.state = RecorderRecording
	return nil
}

func (r *Recorder) Stop() (*storage.Recording, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.state != RecorderRecording {
		return r.recording, nil
	}

	if err := r.capture.Stop(); err != nil {
		return nil, err
	}

	r.state = RecorderIdle
	r.recording.Duration = time.Since(r.startTime)
	r.recording.Finalize()

	return r.recording, nil
}

func (r *Recorder) onEvent(event events.InputEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.recording != nil {
		r.recording.AddEvent(event)
	}
}

func (r *Recorder) SetName(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.recording != nil {
		r.recording.Name = name
	}
}
