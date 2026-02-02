// Package ui provides the graphical user interface.
package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"simplyauto/assets"
	simplyapp "simplyauto/internal/app"
	"simplyauto/internal/hooks"
)

const (
	AppTitle  = "SimplyAuto"
	AppWidth  = 450
	AppHeight = 500
)

type UI struct {
	app            fyne.App
	window         fyne.Window
	simplyApp      *simplyapp.App
	autoClickerTab *AutoClickerTab
	recorderTab    *RecorderTab
	settingsTab    *SettingsTab
	statusLabel    *widget.Label
	version        string
}

func New(simplyApp *simplyapp.App, version string) (*UI, error) {
	fyneApp := app.New()
	if fyneApp == nil {
		return nil, fmt.Errorf("failed to create Fyne application")
	}
	fyneApp.SetIcon(assets.AppIcon())

	window := fyneApp.NewWindow(AppTitle)
	if window == nil {
		return nil, fmt.Errorf("failed to create application window")
	}
	window.Resize(fyne.NewSize(AppWidth, AppHeight))
	window.SetIcon(assets.AppIcon())

	ui := &UI{
		app:       fyneApp,
		window:    window,
		simplyApp: simplyApp,
		version:   version,
	}

	window.SetOnClosed(func() {
		simplyApp.Cleanup()
	})

	ui.setupUI()
	ui.startEventLoop()

	return ui, nil
}

func (u *UI) setupUI() {
	u.autoClickerTab = NewAutoClickerTab(u.simplyApp, u.window)
	u.recorderTab = NewRecorderTab(u.simplyApp)
	u.recorderTab.SetWindow(u.window)
	u.settingsTab = NewSettingsTab(u.simplyApp, u.version)

	tabs := container.NewAppTabs(
		container.NewTabItem("Auto Clicker", u.autoClickerTab.Content()),
		container.NewTabItem("Macro Recorder", u.recorderTab.Content()),
		container.NewTabItem("Settings", u.settingsTab.Content()),
	)

	u.statusLabel = widget.NewLabel(u.getStatusText("Ready"))

	content := container.NewBorder(nil, u.statusLabel, nil, nil, tabs)
	u.window.SetContent(content)
}

func (u *UI) getHotkeyText(action simplyapp.HotkeyAction) string {
	binding := u.simplyApp.GetHotkeyBinding(action)
	if !binding.Bound {
		return "-"
	}
	return hooks.KeyName(binding.Key)
}

func (u *UI) getStatusText(status string) string {
	return fmt.Sprintf("%s | %s: Clicker | %s: Record | %s: Play | %s: Stop",
		status,
		u.getHotkeyText(simplyapp.HotkeyAutoClicker),
		u.getHotkeyText(simplyapp.HotkeyRecord),
		u.getHotkeyText(simplyapp.HotkeyPlayback),
		u.getHotkeyText(simplyapp.HotkeyStop),
	)
}

func (u *UI) startEventLoop() {
	go func() {
		for event := range u.simplyApp.EventChan {
			switch event.Type {
			case "autoclicker":
				u.autoClickerTab.UpdateState(event.Running, event.Count)
				if event.Running {
					u.updateStatus("Auto Clicker running...")
				} else {
					u.updateStatus("Ready")
				}
			case "recorder":
				u.recorderTab.UpdateRecordingState(event.Running, event.Count)
				if event.Running {
					u.updateStatus("Recording...")
				} else {
					u.updateStatus("Ready")
				}
			case "player":
				u.recorderTab.UpdatePlaybackState(event.Running, event.Progress, event.Total, event.Loop)
				if event.Running {
					u.updateStatus("Playing...")
				} else {
					u.updateStatus("Ready")
				}
			}
			u.recorderTab.UpdateFileButtons()
		}
	}()
}

func (u *UI) updateStatus(status string) {
	u.statusLabel.SetText(u.getStatusText(status))
}

func (u *UI) Run() {
	// Apply always-on-top after window is shown
	if u.simplyApp.Settings.AlwaysOnTop {
		go func() {
			time.Sleep(100 * time.Millisecond)
			SetWindowTopmost(AppTitle, true)
		}()
	}
	u.window.ShowAndRun()
}

func (u *UI) Window() fyne.Window {
	return u.window
}
