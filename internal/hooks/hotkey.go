//go:build windows

// Package hooks provides global hotkeys and input capture.
package hooks

import (
	"sync"

	"golang.design/x/hotkey"
)

type Modifier = hotkey.Modifier

const (
	ModCtrl  = hotkey.ModCtrl
	ModShift = hotkey.ModShift
	ModAlt   = hotkey.ModAlt
	ModWin   = hotkey.ModWin
)

type Key = hotkey.Key

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

type HotkeyManager struct {
	hotkeys map[int]*hotkeyEntry
	nextID  int
	mu      sync.Mutex
}

type hotkeyEntry struct {
	hk       *hotkey.Hotkey
	callback HotkeyCallback
	stopChan chan struct{}
}

func NewHotkeyManager() *HotkeyManager {
	return &HotkeyManager{
		hotkeys: make(map[int]*hotkeyEntry),
		nextID:  1,
	}
}

func (m *HotkeyManager) Register(mods []Modifier, key Key, callback HotkeyCallback) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	hk := hotkey.New(mods, key)
	if err := hk.Register(); err != nil {
		return 0, err
	}

	id := m.nextID
	m.nextID++

	entry := &hotkeyEntry{
		hk:       hk,
		callback: callback,
		stopChan: make(chan struct{}),
	}
	m.hotkeys[id] = entry

	go func() {
		for {
			select {
			case <-entry.stopChan:
				return
			case <-hk.Keydown():
				callback()
			}
		}
	}()

	return id, nil
}

func (m *HotkeyManager) Unregister(id int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	entry, ok := m.hotkeys[id]
	if !ok {
		return nil
	}

	close(entry.stopChan)
	err := entry.hk.Unregister()
	delete(m.hotkeys, id)
	return err
}

func (m *HotkeyManager) UnregisterAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id, entry := range m.hotkeys {
		close(entry.stopChan)
		entry.hk.Unregister()
		delete(m.hotkeys, id)
	}
}
