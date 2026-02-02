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

// SaveAutoClicker is a no-op on non-Windows platforms.
func SaveAutoClicker(hours, mins, secs, ms int, randomEnabled bool, randomOffset int,
	button, clickType, repeatMode string, repeatCount int) error {
	return nil
}

// SavePlayback is a no-op on non-Windows platforms.
func SavePlayback(speed, loopMode string, loopCount int) error {
	return nil
}
