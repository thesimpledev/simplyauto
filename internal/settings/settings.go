// Package settings provides application settings storage using Windows Registry.
package settings

// Settings holds all application settings.
type Settings struct {
	// Hotkey bindings (virtual key codes)
	HotkeyAutoClicker uint16
	HotkeyRecord      uint16
	HotkeyPlayback    uint16
	HotkeyStop        uint16

	// Window settings
	AlwaysOnTop bool
}

// Default returns the default settings.
func Default() Settings {
	return Settings{
		HotkeyAutoClicker: 0x75, // F6
		HotkeyRecord:      0x78, // F9
		HotkeyPlayback:    0x79, // F10
		HotkeyStop:        0x7A, // F11
		AlwaysOnTop:       false,
	}
}
