//go:build !windows

package input

// Keycode maps key names to virtual key codes (stub for non-Windows)
var Keycode = map[string]uint16{}

func Move(x, y int)                      {}
func Click(button string, double bool)   {}
func Toggle(button string, state string) {}
func ScrollDir(amount int, dir string)   {}
func KeyDown(key string)                 {}
func KeyUp(key string)                   {}
func KeyDownVK(vk uint16)                {}
func KeyUpVK(vk uint16)                  {}
func KeyPressChar(ch rune)               {}
func CanType(ch rune) bool               { return false }
