//go:build windows

package settings

import (
	"golang.org/x/sys/windows/registry"
)

const registryPath = `Software\SimplyAuto`

// Load reads settings from the Windows Registry.
// Returns default settings if the registry key doesn't exist.
func Load() Settings {
	s := Default()

	key, err := registry.OpenKey(registry.CURRENT_USER, registryPath, registry.QUERY_VALUE)
	if err != nil {
		// Key doesn't exist, return defaults
		return s
	}
	defer key.Close()

	// Hotkeys
	if v, _, err := key.GetIntegerValue("HotkeyAutoClicker"); err == nil {
		s.HotkeyAutoClicker = uint16(v)
	}
	if v, _, err := key.GetIntegerValue("HotkeyRecord"); err == nil {
		s.HotkeyRecord = uint16(v)
	}
	if v, _, err := key.GetIntegerValue("HotkeyPlayback"); err == nil {
		s.HotkeyPlayback = uint16(v)
	}
	if v, _, err := key.GetIntegerValue("HotkeyStop"); err == nil {
		s.HotkeyStop = uint16(v)
	}

	// Window
	if v, _, err := key.GetIntegerValue("AlwaysOnTop"); err == nil {
		s.AlwaysOnTop = v != 0
	}

	// Auto-clicker interval
	if v, _, err := key.GetIntegerValue("ClickIntervalHours"); err == nil {
		s.ClickIntervalHours = int(v)
	}
	if v, _, err := key.GetIntegerValue("ClickIntervalMins"); err == nil {
		s.ClickIntervalMins = int(v)
	}
	if v, _, err := key.GetIntegerValue("ClickIntervalSecs"); err == nil {
		s.ClickIntervalSecs = int(v)
	}
	if v, _, err := key.GetIntegerValue("ClickIntervalMs"); err == nil {
		s.ClickIntervalMs = int(v)
	}

	// Auto-clicker random
	if v, _, err := key.GetIntegerValue("ClickRandomEnabled"); err == nil {
		s.ClickRandomEnabled = v != 0
	}
	if v, _, err := key.GetIntegerValue("ClickRandomOffsetMs"); err == nil {
		s.ClickRandomOffsetMs = int(v)
	}

	// Auto-clicker options
	if v, _, err := key.GetStringValue("ClickButton"); err == nil {
		s.ClickButton = v
	}
	if v, _, err := key.GetStringValue("ClickType"); err == nil {
		s.ClickType = v
	}
	if v, _, err := key.GetStringValue("ClickRepeatMode"); err == nil {
		s.ClickRepeatMode = v
	}
	if v, _, err := key.GetIntegerValue("ClickRepeatCount"); err == nil {
		s.ClickRepeatCount = int(v)
	}

	// Playback
	if v, _, err := key.GetStringValue("PlaybackSpeed"); err == nil {
		s.PlaybackSpeed = v
	}
	if v, _, err := key.GetStringValue("PlaybackLoopMode"); err == nil {
		s.PlaybackLoopMode = v
	}
	if v, _, err := key.GetIntegerValue("PlaybackLoopCount"); err == nil {
		s.PlaybackLoopCount = int(v)
	}

	return s
}

// Save writes settings to the Windows Registry.
func Save(s Settings) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, registryPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	// Hotkeys
	key.SetDWordValue("HotkeyAutoClicker", uint32(s.HotkeyAutoClicker))
	key.SetDWordValue("HotkeyRecord", uint32(s.HotkeyRecord))
	key.SetDWordValue("HotkeyPlayback", uint32(s.HotkeyPlayback))
	key.SetDWordValue("HotkeyStop", uint32(s.HotkeyStop))

	// Window
	if s.AlwaysOnTop {
		key.SetDWordValue("AlwaysOnTop", 1)
	} else {
		key.SetDWordValue("AlwaysOnTop", 0)
	}

	// Auto-clicker
	key.SetDWordValue("ClickIntervalHours", uint32(s.ClickIntervalHours))
	key.SetDWordValue("ClickIntervalMins", uint32(s.ClickIntervalMins))
	key.SetDWordValue("ClickIntervalSecs", uint32(s.ClickIntervalSecs))
	key.SetDWordValue("ClickIntervalMs", uint32(s.ClickIntervalMs))
	if s.ClickRandomEnabled {
		key.SetDWordValue("ClickRandomEnabled", 1)
	} else {
		key.SetDWordValue("ClickRandomEnabled", 0)
	}
	key.SetDWordValue("ClickRandomOffsetMs", uint32(s.ClickRandomOffsetMs))
	key.SetStringValue("ClickButton", s.ClickButton)
	key.SetStringValue("ClickType", s.ClickType)
	key.SetStringValue("ClickRepeatMode", s.ClickRepeatMode)
	key.SetDWordValue("ClickRepeatCount", uint32(s.ClickRepeatCount))

	// Playback
	key.SetStringValue("PlaybackSpeed", s.PlaybackSpeed)
	key.SetStringValue("PlaybackLoopMode", s.PlaybackLoopMode)
	key.SetDWordValue("PlaybackLoopCount", uint32(s.PlaybackLoopCount))

	return nil
}

// SaveHotkey saves a single hotkey setting to the registry.
func SaveHotkey(name string, vk uint16) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, registryPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	return key.SetDWordValue(name, uint32(vk))
}

// SaveAlwaysOnTop saves the always-on-top setting to the registry.
func SaveAlwaysOnTop(enabled bool) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, registryPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	v := uint32(0)
	if enabled {
		v = 1
	}
	return key.SetDWordValue("AlwaysOnTop", v)
}

// SaveAutoClicker saves auto-clicker settings to the registry.
func SaveAutoClicker(hours, mins, secs, ms int, randomEnabled bool, randomOffset int,
	button, clickType, repeatMode string, repeatCount int) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, registryPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	key.SetDWordValue("ClickIntervalHours", uint32(hours))
	key.SetDWordValue("ClickIntervalMins", uint32(mins))
	key.SetDWordValue("ClickIntervalSecs", uint32(secs))
	key.SetDWordValue("ClickIntervalMs", uint32(ms))
	if randomEnabled {
		key.SetDWordValue("ClickRandomEnabled", 1)
	} else {
		key.SetDWordValue("ClickRandomEnabled", 0)
	}
	key.SetDWordValue("ClickRandomOffsetMs", uint32(randomOffset))
	key.SetStringValue("ClickButton", button)
	key.SetStringValue("ClickType", clickType)
	key.SetStringValue("ClickRepeatMode", repeatMode)
	key.SetDWordValue("ClickRepeatCount", uint32(repeatCount))

	return nil
}

// SavePlayback saves playback settings to the registry.
func SavePlayback(speed, loopMode string, loopCount int) error {
	key, _, err := registry.CreateKey(registry.CURRENT_USER, registryPath, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	key.SetStringValue("PlaybackSpeed", speed)
	key.SetStringValue("PlaybackLoopMode", loopMode)
	key.SetDWordValue("PlaybackLoopCount", uint32(loopCount))

	return nil
}
