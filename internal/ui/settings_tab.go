package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"simplyauto/internal/app"
	"simplyauto/internal/hooks"
	"simplyauto/internal/settings"
)

type SettingsTab struct {
	simplyApp               *app.App
	version                 string
	autoClickerHotkeySelect *widget.Select
	recordHotkeySelect      *widget.Select
	playHotkeySelect        *widget.Select
	stopHotkeySelect        *widget.Select
	statusLabels            map[app.HotkeyAction]*widget.Label
	alwaysOnTopCheck        *widget.Check
	content                 fyne.CanvasObject
}

func NewSettingsTab(simplyApp *app.App, version string) *SettingsTab {
	t := &SettingsTab{
		simplyApp:    simplyApp,
		version:      version,
		statusLabels: make(map[app.HotkeyAction]*widget.Label),
	}
	t.build()
	return t
}

func (t *SettingsTab) Content() fyne.CanvasObject {
	return t.content
}

func (t *SettingsTab) getKeyOptions() []string {
	keys := hooks.AvailableKeys()
	options := make([]string, len(keys))
	for i, k := range keys {
		options[i] = hooks.KeyName(k)
	}
	return options
}

func (t *SettingsTab) createHotkeySelect(action app.HotkeyAction) *widget.Select {
	binding := t.simplyApp.GetHotkeyBinding(action)
	currentKey := hooks.KeyName(binding.Key)

	sel := widget.NewSelect(t.getKeyOptions(), func(selected string) {
		key, ok := hooks.KeyFromName(selected)
		if !ok {
			return
		}
		label := t.statusLabels[action]
		if err := t.simplyApp.RebindHotkey(action, key); err != nil {
			if label != nil {
				label.SetText("Error: " + err.Error())
			}
		} else {
			if label != nil {
				label.SetText("Bound")
			}
		}
	})
	sel.SetSelected(currentKey)
	return sel
}

func (t *SettingsTab) createStatusLabel(action app.HotkeyAction) *widget.Label {
	binding := t.simplyApp.GetHotkeyBinding(action)
	text := "Bound"
	if !binding.Bound {
		if binding.ErrorMsg != "" {
			text = "Error: " + binding.ErrorMsg
		} else {
			text = "Not bound"
		}
	}
	label := widget.NewLabel(text)
	t.statusLabels[action] = label
	return label
}

func (t *SettingsTab) build() {
	t.autoClickerHotkeySelect = t.createHotkeySelect(app.HotkeyAutoClicker)
	t.recordHotkeySelect = t.createHotkeySelect(app.HotkeyRecord)
	t.playHotkeySelect = t.createHotkeySelect(app.HotkeyPlayback)
	t.stopHotkeySelect = t.createHotkeySelect(app.HotkeyStop)

	hotkeySection := container.NewVBox(
		widget.NewLabel("Hotkeys"),
		container.NewGridWithColumns(3,
			widget.NewLabel("Auto Clicker Toggle:"), t.autoClickerHotkeySelect, t.createStatusLabel(app.HotkeyAutoClicker),
			widget.NewLabel("Record Toggle:"), t.recordHotkeySelect, t.createStatusLabel(app.HotkeyRecord),
			widget.NewLabel("Start Playback:"), t.playHotkeySelect, t.createStatusLabel(app.HotkeyPlayback),
			widget.NewLabel("Stop All:"), t.stopHotkeySelect, t.createStatusLabel(app.HotkeyStop),
		),
		widget.NewSeparator(),
	)

	t.alwaysOnTopCheck = widget.NewCheck("Always on top", func(checked bool) {
		SetWindowTopmost(AppTitle, checked)
		settings.SaveAlwaysOnTop(checked)
	})
	// Load saved setting
	t.alwaysOnTopCheck.SetChecked(t.simplyApp.Settings.AlwaysOnTop)

	windowSection := container.NewVBox(
		widget.NewLabel("Window Options"),
		t.alwaysOnTopCheck,
		widget.NewSeparator(),
	)

	simplyAutoURL, _ := url.Parse("https://simplyauto.dev")
	theSimpleDevURL, _ := url.Parse("https://thesimpledev.com")

	aboutSection := container.NewVBox(
		widget.NewLabel("About"),
		widget.NewLabel("SimplyAuto v"+t.version),
		widget.NewLabel("Windows Auto Clicker & Macro Recorder"),
		widget.NewHyperlink("simplyauto.dev", simplyAutoURL),
		widget.NewHyperlink("thesimpledev.com", theSimpleDevURL),
	)

	t.content = container.NewVBox(
		hotkeySection,
		windowSection,
		aboutSection,
	)
}
