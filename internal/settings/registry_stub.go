//go:build !windows

package settings

// Load returns default settings on non-Windows platforms.
func Load() Settings {
	return Default()
}

// Save is a no-op on non-Windows platforms.
func Save(s Settings) error {
	return nil
}

// SaveHotkey is a no-op on non-Windows platforms.
func SaveHotkey(name string, vk uint16) error {
	return nil
}

// SaveAlwaysOnTop is a no-op on non-Windows platforms.
func SaveAlwaysOnTop(enabled bool) error {
	return nil
}
