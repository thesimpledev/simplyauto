// Package app coordinates all application components.
package app

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"simplyauto/internal/autoclicker"
	"simplyauto/internal/hooks"
	"simplyauto/internal/recorder"
	"simplyauto/internal/settings"
	"simplyauto/internal/storage"
)

var ErrNotIdle = errors.New("stop current operation first")

type HotkeyAction string

const (
	HotkeyAutoClicker HotkeyAction = "autoclicker"
	HotkeyRecord      HotkeyAction = "record"
	HotkeyPlayback    HotkeyAction = "playback"
	HotkeyStop        HotkeyAction = "stop"
)

type HotkeyBinding struct {
	Action   HotkeyAction
	Key      hooks.Key
	ID       int
	Bound    bool
	ErrorMsg string
}

type StateEvent struct {
	Type     string
	Running  bool
	Count    int
	Progress int
	Total    int
	Loop     int
}

type App struct {
	Log         *log.Logger
	AutoClicker *autoclicker.AutoClicker
	Recorder    *recorder.Recorder
	Player      *recorder.Player
	Storage     *storage.JSONStorage
	Hotkeys     *hooks.HotkeyManager
	Settings    settings.Settings

	hotkeyBindings map[HotkeyAction]*HotkeyBinding

	currentRecording *storage.Recording
	currentFilePath  string

	PlaybackSpeed float64
	PlaybackLoop  recorder.LoopMode
	PlaybackCount int

	// OnAutoClickerStart is called before the autoclicker starts to apply UI config
	OnAutoClickerStart func() error
	// OnPlaybackStart is called before playback starts to apply UI config
	OnPlaybackStart func()

	mu sync.Mutex

	EventChan chan StateEvent

	// Lazy logger - file only created on first error
	logMu      sync.Mutex
	logFile    *os.File
	logCreated bool
}

func New() *App {
	// Load settings from registry
	s := settings.Load()

	a := &App{
		Log:         log.New(&lazyLogWriter{}, "", log.Ldate|log.Ltime|log.Lshortfile),
		AutoClicker: autoclicker.New(),
		Recorder:    recorder.NewRecorder(recorder.DefaultRecorderOptions()),
		Player:      recorder.NewPlayer(),
		Storage:     storage.NewJSONStorage(),
		Hotkeys:     hooks.NewHotkeyManager(),
		Settings:    s,
		hotkeyBindings: map[HotkeyAction]*HotkeyBinding{
			HotkeyAutoClicker: {Action: HotkeyAutoClicker, Key: hooks.Key(s.HotkeyAutoClicker)},
			HotkeyRecord:      {Action: HotkeyRecord, Key: hooks.Key(s.HotkeyRecord)},
			HotkeyPlayback:    {Action: HotkeyPlayback, Key: hooks.Key(s.HotkeyPlayback)},
			HotkeyStop:        {Action: HotkeyStop, Key: hooks.Key(s.HotkeyStop)},
		},
		PlaybackSpeed: 1.0,
		PlaybackLoop:  recorder.LoopOnce,
		PlaybackCount: 1,
		EventChan:     make(chan StateEvent, 32),
	}

	// Set the app reference for lazy logging
	a.Log.SetOutput(&lazyLogWriter{app: a})

	a.Player.OnComplete = func() {
		a.sendEvent(StateEvent{Type: "player", Running: false})
	}

	return a
}

func (a *App) sendEvent(e StateEvent) {
	go func() {
		a.EventChan <- e
	}()
}

func (a *App) isIdle() bool {
	return !a.AutoClicker.IsRunning() && !a.Recorder.IsRecording() && !a.Player.IsPlaying()
}

func (a *App) isMacroActive() bool {
	return a.Recorder.IsRecording() || a.Player.IsPlaying()
}

func (a *App) stopRecorder() error {
	rec, err := a.Recorder.Stop()
	if err != nil {
		a.Log.Printf("failed to stop recording: %v", err)
		return err
	}
	a.currentRecording = rec
	a.currentFilePath = ""
	return nil
}

// lazyLogWriter creates the log file only when the first error is written.
type lazyLogWriter struct {
	app *App
}

func (w *lazyLogWriter) Write(p []byte) (n int, err error) {
	if w.app == nil {
		return len(p), nil
	}

	w.app.logMu.Lock()
	defer w.app.logMu.Unlock()

	// Create log file on first write
	if !w.app.logCreated {
		exe, err := os.Executable()
		if err != nil {
			return len(p), nil // Silently discard if we can't get executable path
		}

		logPath := filepath.Join(filepath.Dir(exe), "simplyauto_errors.log")
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return len(p), nil // Silently discard if we can't create log file
		}

		w.app.logFile = f
		w.app.logCreated = true
	}

	if w.app.logFile != nil {
		return w.app.logFile.Write(p)
	}

	return len(p), nil
}

func (a *App) RegisterDefaultHotkeys() error {
	var firstErr error
	for action, binding := range a.hotkeyBindings {
		if err := a.registerHotkey(action, binding.Key); err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (a *App) getCallbackForAction(action HotkeyAction) hooks.HotkeyCallback {
	switch action {
	case HotkeyAutoClicker:
		return a.ToggleAutoClicker
	case HotkeyRecord:
		return a.ToggleRecording
	case HotkeyPlayback:
		return a.TogglePlayback
	case HotkeyStop:
		return a.Stop
	}
	return nil
}

func (a *App) registerHotkey(action HotkeyAction, key hooks.Key) error {
	binding := a.hotkeyBindings[action]
	if binding == nil {
		return nil
	}

	callback := a.getCallbackForAction(action)
	if callback == nil {
		return nil
	}

	id, err := a.Hotkeys.Register(nil, key, callback)
	if err != nil {
		binding.Bound = false
		binding.ErrorMsg = err.Error()
		a.Log.Printf("failed to register %s hotkey (%s): %v", action, hooks.KeyName(key), err)
		return err
	}

	binding.ID = id
	binding.Key = key
	binding.Bound = true
	binding.ErrorMsg = ""
	return nil
}

func (a *App) RebindHotkey(action HotkeyAction, newKey hooks.Key) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	binding := a.hotkeyBindings[action]
	if binding == nil {
		return errors.New("unknown action")
	}

	// Check if another action is already using this key
	for otherAction, otherBinding := range a.hotkeyBindings {
		if otherAction != action && otherBinding.Key == newKey && otherBinding.Bound {
			return errors.New("key already in use by " + string(otherAction))
		}
	}

	// Unregister old hotkey if bound
	if binding.Bound {
		if err := a.Hotkeys.Unregister(binding.ID); err != nil {
			return fmt.Errorf("failed to unregister old hotkey: %w", err)
		}
		binding.Bound = false
	}

	// Register new hotkey
	if err := a.registerHotkey(action, newKey); err != nil {
		return err
	}

	// Save to registry
	registryName := "Hotkey" + string(action)
	registryName = map[HotkeyAction]string{
		HotkeyAutoClicker: "HotkeyAutoClicker",
		HotkeyRecord:      "HotkeyRecord",
		HotkeyPlayback:    "HotkeyPlayback",
		HotkeyStop:        "HotkeyStop",
	}[action]
	if registryName != "" {
		settings.SaveHotkey(registryName, uint16(newKey))
	}

	return nil
}

func (a *App) GetHotkeyBinding(action HotkeyAction) HotkeyBinding {
	a.mu.Lock()
	defer a.mu.Unlock()
	if b := a.hotkeyBindings[action]; b != nil {
		return *b
	}
	return HotkeyBinding{}
}

func (a *App) GetAllHotkeyBindings() map[HotkeyAction]HotkeyBinding {
	a.mu.Lock()
	defer a.mu.Unlock()
	result := make(map[HotkeyAction]HotkeyBinding)
	for k, v := range a.hotkeyBindings {
		result[k] = *v
	}
	return result
}

func (a *App) ToggleAutoClicker() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.isMacroActive() {
		return
	}

	// If starting (not currently running), apply config from UI first
	if !a.AutoClicker.IsRunning() && a.OnAutoClickerStart != nil {
		if err := a.OnAutoClickerStart(); err != nil {
			a.Log.Printf("failed to apply autoclicker config: %v", err)
			return
		}
	}

	a.AutoClicker.Toggle()
	a.sendEvent(StateEvent{
		Type:    "autoclicker",
		Running: a.AutoClicker.IsRunning(),
		Count:   a.AutoClicker.GetClickCount(),
	})
}

func (a *App) ToggleRecording() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.AutoClicker.IsRunning() || a.Player.IsPlaying() {
		return
	}

	if a.Recorder.IsRecording() {
		if err := a.stopRecorder(); err != nil {
			return
		}
	} else {
		if err := a.Recorder.Start(); err != nil {
			a.Log.Printf("failed to start recording: %v", err)
			return
		}
	}

	a.sendEvent(StateEvent{
		Type:    "recorder",
		Running: a.Recorder.IsRecording(),
		Count:   a.Recorder.GetEventCount(),
	})
}

func (a *App) TogglePlayback() {
	// Apply UI config before acquiring lock to avoid deadlock
	// (callback calls SetPlaybackSpeed/SetPlaybackLoop which also lock)
	if a.OnPlaybackStart != nil {
		a.OnPlaybackStart()
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.AutoClicker.IsRunning() || a.Recorder.IsRecording() {
		return
	}

	if a.Player.IsPlaying() {
		a.Player.Stop()
		a.sendEvent(StateEvent{Type: "player", Running: false})
		return
	}

	if a.currentRecording == nil || len(a.currentRecording.Events) == 0 {
		return
	}

	a.sendEvent(StateEvent{
		Type:     "player",
		Running:  true,
		Progress: 0,
		Total:    len(a.currentRecording.Events),
		Loop:     1,
	})

	config := recorder.PlaybackConfig{
		Speed:     a.PlaybackSpeed,
		LoopMode:  a.PlaybackLoop,
		LoopCount: a.PlaybackCount,
	}
	if err := a.Player.Play(a.currentRecording, config); err != nil {
		a.Log.Printf("failed to start playback: %v", err)
		a.sendEvent(StateEvent{Type: "player", Running: false})
	}
}

func (a *App) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.AutoClicker.IsRunning() {
		a.AutoClicker.Stop()
		a.sendEvent(StateEvent{Type: "autoclicker", Running: false, Count: a.AutoClicker.GetClickCount()})
		return
	}

	if a.Recorder.IsRecording() {
		if err := a.stopRecorder(); err != nil {
			return
		}
		a.sendEvent(StateEvent{Type: "recorder", Running: false, Count: a.Recorder.GetEventCount()})
		return
	}

	if a.Player.IsPlaying() {
		a.Player.Stop()
		a.sendEvent(StateEvent{Type: "player", Running: false})
	}
}

func (a *App) SaveRecording(path string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.currentRecording == nil {
		return nil
	}

	if err := a.Storage.Save(a.currentRecording, path); err != nil {
		a.Log.Printf("failed to save recording to %s: %v", path, err)
		return err
	}

	a.currentFilePath = path
	return nil
}

func (a *App) LoadRecording(path string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.isIdle() {
		return ErrNotIdle
	}

	rec, err := a.Storage.Load(path)
	if err != nil {
		a.Log.Printf("failed to load recording from %s: %v", path, err)
		return err
	}

	a.currentRecording = rec
	a.currentFilePath = path
	return nil
}

func (a *App) Cleanup() {
	a.Stop()
	a.Hotkeys.UnregisterAll()
	if a.logFile != nil {
		a.logFile.Close()
	}
}

func (a *App) GetCurrentRecording() *storage.Recording {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.currentRecording
}

func (a *App) GetCurrentFilePath() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.currentFilePath
}

func (a *App) HasRecording() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.currentRecording != nil
}

func (a *App) IsIdle() bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.isIdle()
}

func (a *App) SetPlaybackSpeed(speed float64) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if speed > 0 {
		a.PlaybackSpeed = speed
	}
}

func (a *App) SetPlaybackLoop(mode recorder.LoopMode, count int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.PlaybackLoop = mode
	if count > 0 {
		a.PlaybackCount = count
	}
}

func (a *App) GetPlaybackConfig() (speed float64, mode recorder.LoopMode, count int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.PlaybackSpeed, a.PlaybackLoop, a.PlaybackCount
}
