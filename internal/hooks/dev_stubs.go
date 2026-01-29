//go:build !windows

// This file provides stub implementations for non-Windows platforms.
// It exists only to allow Go tooling (IDE support, go list, go vet) to work
// during development on Linux. This file can be deleted once development is
// complete since the application only runs on Windows.

package hooks

import (
	"errors"

	"simplyauto/pkg/events"
)

var errNotWindows = errors.New("hooks only supported on Windows")

type Modifier uint8

const (
	ModCtrl  Modifier = 1
	ModShift Modifier = 2
	ModAlt   Modifier = 3
	ModWin   Modifier = 4
)

type Key uint16

const (
	KeyF1  Key = 0x70
	KeyF2  Key = 0x71
	KeyF3  Key = 0x72
	KeyF4  Key = 0x73
	KeyF5  Key = 0x74
	KeyF6  Key = 0x75
	KeyF7  Key = 0x76
	KeyF8  Key = 0x77
	KeyF9  Key = 0x78
	KeyF10 Key = 0x79
	KeyF11 Key = 0x7A
	KeyF12 Key = 0x7B
)

var keyNames = map[Key]string{
	KeyF1:  "F1",
	KeyF2:  "F2",
	KeyF3:  "F3",
	KeyF4:  "F4",
	KeyF5:  "F5",
	KeyF6:  "F6",
	KeyF7:  "F7",
	KeyF8:  "F8",
	KeyF9:  "F9",
	KeyF10: "F10",
	KeyF11: "F11",
	KeyF12: "F12",
}

func KeyName(k Key) string {
	if name, ok := keyNames[k]; ok {
		return name
	}
	return "Unknown"
}

func KeyFromName(name string) (Key, bool) {
	for k, n := range keyNames {
		if n == name {
			return k, true
		}
	}
	return 0, false
}

func AvailableKeys() []Key {
	return []Key{KeyF1, KeyF2, KeyF3, KeyF4, KeyF5, KeyF6, KeyF7, KeyF8, KeyF9, KeyF10, KeyF11, KeyF12}
}

type HotkeyCallback func()

type HotkeyManager struct{}

func NewHotkeyManager() *HotkeyManager {
	return &HotkeyManager{}
}

func (m *HotkeyManager) Register(mods []Modifier, key Key, callback HotkeyCallback) (int, error) {
	return 0, errNotWindows
}

func (m *HotkeyManager) Unregister(id int) error {
	return errNotWindows
}

func (m *HotkeyManager) UnregisterAll() {}

type EventCallback func(event events.InputEvent)

type InputCapture struct{}

type CaptureOptions struct {
	CaptureMouse    bool
	CaptureKeyboard bool
	FilterKeys      []Key
}

func NewInputCapture(opts CaptureOptions) *InputCapture {
	return &InputCapture{}
}

func (c *InputCapture) Start(callback EventCallback) error { return errNotWindows }
func (c *InputCapture) Stop() error                        { return errNotWindows }
func (c *InputCapture) IsRunning() bool                    { return false }
