//go:build windows

package input

const dpiAwarenessContextPerMonitorAwareV2 = ^uintptr(4) + 1 // -4 as uintptr

func init() {
	proc := user32.NewProc("SetProcessDpiAwarenessContext")
	if proc.Find() != nil {
		return
	}
	proc.Call(dpiAwarenessContextPerMonitorAwareV2)
}
