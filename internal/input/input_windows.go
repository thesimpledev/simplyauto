//go:build windows

// Package input provides mouse and keyboard automation using Windows API.
package input

import (
	"syscall"
	"unsafe"
)

var (
	user32         = syscall.NewLazyDLL("user32.dll")
	procSetCursorPos = user32.NewProc("SetCursorPos")
	procSendInput    = user32.NewProc("SendInput")
)

// Input type constants
const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
)

// Mouse event flags
const (
	MOUSEEVENTF_MOVE       = 0x0001
	MOUSEEVENTF_LEFTDOWN   = 0x0002
	MOUSEEVENTF_LEFTUP     = 0x0004
	MOUSEEVENTF_RIGHTDOWN  = 0x0008
	MOUSEEVENTF_RIGHTUP    = 0x0010
	MOUSEEVENTF_MIDDLEDOWN = 0x0020
	MOUSEEVENTF_MIDDLEUP   = 0x0040
	MOUSEEVENTF_WHEEL      = 0x0800
	MOUSEEVENTF_ABSOLUTE   = 0x8000
)

// Keyboard event flags
const (
	KEYEVENTF_KEYDOWN   = 0x0000
	KEYEVENTF_KEYUP     = 0x0002
	KEYEVENTF_SCANCODE  = 0x0008
)

// MOUSEINPUT structure
type mouseInput struct {
	dx          int32
	dy          int32
	mouseData   uint32
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

// KEYBDINPUT structure
type keybdInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

// INPUT structure (union represented as mouse - keyboard uses same memory layout with padding)
// On 64-bit Windows, there must be 4 bytes of padding after dtype for proper alignment
// because the union contains ULONG_PTR (dwExtraInfo) which requires 8-byte alignment.
type inputUnion struct {
	dtype uint32
	_     uint32 // padding for 64-bit alignment
	mi    mouseInput
}

type keyboardInputUnion struct {
	dtype uint32
	_     uint32  // padding for 64-bit alignment
	ki    keybdInput
	__    [8]byte // padding to match mouseInput size (mouseInput has uintptr at end)
}

// Keycode maps key names to virtual key codes
var Keycode = map[string]uint16{
	"backspace":    0x08,
	"tab":          0x09,
	"enter":        0x0D,
	"shift":        0x10,
	"ctrl":         0x11,
	"alt":          0x12,
	"pause":        0x13,
	"capslock":     0x14,
	"escape":       0x1B,
	"space":        0x20,
	"pageup":       0x21,
	"pagedown":     0x22,
	"end":          0x23,
	"home":         0x24,
	"left":         0x25,
	"up":           0x26,
	"right":        0x27,
	"down":         0x28,
	"insert":       0x2D,
	"delete":       0x2E,
	"0":            0x30,
	"1":            0x31,
	"2":            0x32,
	"3":            0x33,
	"4":            0x34,
	"5":            0x35,
	"6":            0x36,
	"7":            0x37,
	"8":            0x38,
	"9":            0x39,
	"a":            0x41,
	"b":            0x42,
	"c":            0x43,
	"d":            0x44,
	"e":            0x45,
	"f":            0x46,
	"g":            0x47,
	"h":            0x48,
	"i":            0x49,
	"j":            0x4A,
	"k":            0x4B,
	"l":            0x4C,
	"m":            0x4D,
	"n":            0x4E,
	"o":            0x4F,
	"p":            0x50,
	"q":            0x51,
	"r":            0x52,
	"s":            0x53,
	"t":            0x54,
	"u":            0x55,
	"v":            0x56,
	"w":            0x57,
	"x":            0x58,
	"y":            0x59,
	"z":            0x5A,
	"lwin":         0x5B,
	"rwin":         0x5C,
	"numpad0":      0x60,
	"numpad1":      0x61,
	"numpad2":      0x62,
	"numpad3":      0x63,
	"numpad4":      0x64,
	"numpad5":      0x65,
	"numpad6":      0x66,
	"numpad7":      0x67,
	"numpad8":      0x68,
	"numpad9":      0x69,
	"multiply":     0x6A,
	"add":          0x6B,
	"subtract":     0x6D,
	"decimal":      0x6E,
	"divide":       0x6F,
	"f1":           0x70,
	"f2":           0x71,
	"f3":           0x72,
	"f4":           0x73,
	"f5":           0x74,
	"f6":           0x75,
	"f7":           0x76,
	"f8":           0x77,
	"f9":           0x78,
	"f10":          0x79,
	"f11":          0x7A,
	"f12":          0x7B,
	"numlock":      0x90,
	"scrolllock":   0x91,
	"lshift":       0xA0,
	"rshift":       0xA1,
	"lctrl":        0xA2,
	"rctrl":        0xA3,
	"lalt":         0xA4,
	"ralt":         0xA5,
	"semicolon":    0xBA,
	"equal":        0xBB,
	"comma":        0xBC,
	"minus":        0xBD,
	"period":       0xBE,
	"slash":        0xBF,
	"grave":        0xC0,
	"leftbracket":  0xDB,
	"backslash":    0xDC,
	"rightbracket": 0xDD,
	"quote":        0xDE,
}

// Move moves the mouse cursor to the specified position.
func Move(x, y int) {
	procSetCursorPos.Call(uintptr(x), uintptr(y))
}

// Click performs a mouse click at the current position.
func Click(button string, double bool) {
	Toggle(button, "down")
	Toggle(button, "up")
	if double {
		Toggle(button, "down")
		Toggle(button, "up")
	}
}

// Toggle presses or releases a mouse button.
func Toggle(button string, state string) {
	var flags uint32

	switch button {
	case "left":
		if state == "down" {
			flags = MOUSEEVENTF_LEFTDOWN
		} else {
			flags = MOUSEEVENTF_LEFTUP
		}
	case "right":
		if state == "down" {
			flags = MOUSEEVENTF_RIGHTDOWN
		} else {
			flags = MOUSEEVENTF_RIGHTUP
		}
	case "center", "middle":
		if state == "down" {
			flags = MOUSEEVENTF_MIDDLEDOWN
		} else {
			flags = MOUSEEVENTF_MIDDLEUP
		}
	default:
		return
	}

	input := inputUnion{
		dtype: INPUT_MOUSE,
		mi: mouseInput{
			dwFlags: flags,
		},
	}

	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

// ScrollDir scrolls the mouse wheel.
func ScrollDir(amount int, direction string) {
	var delta int32
	if direction == "up" {
		delta = int32(amount * 120)
	} else {
		delta = int32(-amount * 120)
	}

	input := inputUnion{
		dtype: INPUT_MOUSE,
		mi: mouseInput{
			dwFlags:   MOUSEEVENTF_WHEEL,
			mouseData: uint32(delta),
		},
	}

	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

// KeyDown presses a key.
func KeyDown(key string) {
	vk, ok := Keycode[key]
	if !ok {
		return
	}

	input := keyboardInputUnion{
		dtype: INPUT_KEYBOARD,
		ki: keybdInput{
			wVk:     vk,
			dwFlags: KEYEVENTF_KEYDOWN,
		},
	}

	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}

// KeyUp releases a key.
func KeyUp(key string) {
	vk, ok := Keycode[key]
	if !ok {
		return
	}

	input := keyboardInputUnion{
		dtype: INPUT_KEYBOARD,
		ki: keybdInput{
			wVk:     vk,
			dwFlags: KEYEVENTF_KEYUP,
		},
	}

	procSendInput.Call(1, uintptr(unsafe.Pointer(&input)), unsafe.Sizeof(input))
}
