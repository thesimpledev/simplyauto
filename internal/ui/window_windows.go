//go:build windows

package ui

import (
	"syscall"
	"unsafe"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procFindWindow   = user32.NewProc("FindWindowW")
	procSetWindowPos = user32.NewProc("SetWindowPos")
)

const (
	HWND_TOPMOST    = ^uintptr(0)      // -1
	HWND_NOTOPMOST  = ^uintptr(0) - 1  // -2
	SWP_NOMOVE      = 0x0002
	SWP_NOSIZE      = 0x0001
	SWP_NOACTIVATE  = 0x0010
)

func findWindowByTitle(title string) uintptr {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	return hwnd
}

func SetWindowTopmost(title string, topmost bool) {
	hwnd := findWindowByTitle(title)
	if hwnd == 0 {
		return
	}

	var insertAfter uintptr
	if topmost {
		insertAfter = HWND_TOPMOST
	} else {
		insertAfter = HWND_NOTOPMOST
	}

	procSetWindowPos.Call(
		hwnd,
		insertAfter,
		0, 0, 0, 0,
		SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE,
	)
}
