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
	if v, _, err := key.GetIntegerValue("AlwaysOnTop"); err == nil {
		s.AlwaysOnTop = v != 0
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

	if err := key.SetDWordValue("HotkeyAutoClicker", uint32(s.HotkeyAutoClicker)); err != nil {
		return err
	}
	if err := key.SetDWordValue("HotkeyRecord", uint32(s.HotkeyRecord)); err != nil {
		return err
	}
	if err := key.SetDWordValue("HotkeyPlayback", uint32(s.HotkeyPlayback)); err != nil {
		return err
	}
	if err := key.SetDWordValue("HotkeyStop", uint32(s.HotkeyStop)); err != nil {
		return err
	}
	alwaysOnTopVal := uint32(0)
	if s.AlwaysOnTop {
		alwaysOnTopVal = 1
	}
	if err := key.SetDWordValue("AlwaysOnTop", alwaysOnTopVal); err != nil {
		return err
	}

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
