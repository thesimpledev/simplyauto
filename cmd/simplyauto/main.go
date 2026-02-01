package main

import (
	"fmt"
	"os"
	"runtime"

	"golang.org/x/sys/windows"

	"simplyauto/internal/app"
	"simplyauto/internal/ui"
)

// Version is set via ldflags at build time
var Version = "dev"

func showMessageBox(title, message string, flags uint32) error {
	titlePtr, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return err
	}
	msgPtr, err := windows.UTF16PtrFromString(message)
	if err != nil {
		return err
	}
	_, err = windows.MessageBox(0, msgPtr, titlePtr, flags)
	return err
}

func main() {
	os.Setenv("FYNE_RENDER", "software")
	runtime.LockOSThread()

	simplyApp := app.New()
	defer simplyApp.Cleanup()

	if simplyApp.LogError != nil {
		showMessageBox("SimplyAuto - Warning",
			fmt.Sprintf("Could not create error log file: %v\n\nThe application does not have write privileges to its directory.", simplyApp.LogError),
			windows.MB_OK|windows.MB_ICONWARNING)
	}

	if err := simplyApp.RegisterDefaultHotkeys(); err != nil {
		simplyApp.Log.Printf("failed to register hotkeys: %v", err)
		if mbErr := showMessageBox("SimplyAuto - Hotkey Warning",
			fmt.Sprintf("Failed to register some hotkeys: %v\n\nAnother application may be using these keys. You can rebind them in Settings.", err),
			windows.MB_OK|windows.MB_ICONWARNING); mbErr != nil {
			simplyApp.Log.Printf("failed to show message box: %v", mbErr)
		}
	}

	appUI, err := ui.New(simplyApp)
	if err != nil {
		simplyApp.Log.Printf("failed to initialize UI: %v", err)
		if mbErr := showMessageBox("SimplyAuto - Error",
			fmt.Sprintf("Failed to initialize UI: %v", err),
			windows.MB_OK|windows.MB_ICONWARNING); mbErr != nil {
			simplyApp.Log.Printf("failed to show message box: %v", mbErr)
		}
		return
	}
	appUI.Run()
}
