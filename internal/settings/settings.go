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

	// Auto-clicker settings
	ClickIntervalHours  int
	ClickIntervalMins   int
	ClickIntervalSecs   int
	ClickIntervalMs     int
	ClickRandomEnabled  bool
	ClickRandomOffsetMs int
	ClickButton         string // "Left", "Right", "Middle"
	ClickType           string // "Single", "Double"
	ClickRepeatMode     string // "Until stopped", "Count"
	ClickRepeatCount    int

	// Playback settings
	PlaybackSpeed    string // "0.5x", "1x", "2x", "4x"
	PlaybackLoopMode string // "Once", "Count", "Continuous"
	PlaybackLoopCount int
}

// Default returns the default settings.
func Default() Settings {
	return Settings{
		HotkeyAutoClicker: 0x75, // F6
		HotkeyRecord:      0x78, // F9
		HotkeyPlayback:    0x79, // F10
		HotkeyStop:        0x7A, // F11
		AlwaysOnTop:       false,

		// Auto-clicker defaults
		ClickIntervalHours:  0,
		ClickIntervalMins:   0,
		ClickIntervalSecs:   1,
		ClickIntervalMs:     0,
		ClickRandomEnabled:  false,
		ClickRandomOffsetMs: 0,
		ClickButton:         "Left",
		ClickType:           "Single",
		ClickRepeatMode:     "Until stopped",
		ClickRepeatCount:    1,

		// Playback defaults
		PlaybackSpeed:     "1x",
		PlaybackLoopMode:  "Once",
		PlaybackLoopCount: 1,
	}
}
