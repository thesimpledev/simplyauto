//go:build !windows

package ui

func SetWindowTopmost(title string, topmost bool) {
	// No-op on non-Windows platforms
}
