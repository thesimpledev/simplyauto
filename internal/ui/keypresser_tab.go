package ui

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"simplyauto/internal/app"
	"simplyauto/internal/keypresser"
	"simplyauto/internal/settings"
)

type KeyPresserTab struct {
	app          *app.App
	window       fyne.Window
	keysEntry    *widget.Entry
	hoursEntry   *widget.Entry
	minsEntry    *widget.Entry
	secsEntry    *widget.Entry
	msEntry      *widget.Entry
	randomCheck  *widget.Check
	randomEntry  *widget.Entry
	repeatSelect *widget.RadioGroup
	repeatEntry  *widget.Entry
	startButton  *widget.Button
	pressesLabel *widget.Label
	content      fyne.CanvasObject
}

func NewKeyPresserTab(app *app.App, window fyne.Window) *KeyPresserTab {
	t := &KeyPresserTab{app: app, window: window}
	t.build()

	// Register callback so config is always applied before the key presser starts
	app.OnKeyPresserStart = t.applyConfig

	return t
}

func (t *KeyPresserTab) Content() fyne.CanvasObject {
	return t.content
}

func (t *KeyPresserTab) build() {
	s := t.app.Settings

	t.keysEntry = widget.NewEntry()
	t.keysEntry.SetText(s.KeyPresserKeys)
	t.keysEntry.SetPlaceHolder("keys to press, e.g. a")
	t.keysEntry.OnChanged = func(_ string) { t.saveSettings() }

	keysSection := container.NewVBox(
		widget.NewLabel("Keys to Press"),
		t.keysEntry,
		widget.NewSeparator(),
	)

	t.hoursEntry = widget.NewEntry()
	t.hoursEntry.SetText(strconv.Itoa(s.KeyPresserIntervalHours))
	t.hoursEntry.SetPlaceHolder("0")
	t.hoursEntry.OnChanged = func(_ string) { t.saveSettings() }

	t.minsEntry = widget.NewEntry()
	t.minsEntry.SetText(strconv.Itoa(s.KeyPresserIntervalMins))
	t.minsEntry.SetPlaceHolder("0")
	t.minsEntry.OnChanged = func(_ string) { t.saveSettings() }

	t.secsEntry = widget.NewEntry()
	t.secsEntry.SetText(strconv.Itoa(s.KeyPresserIntervalSecs))
	t.secsEntry.SetPlaceHolder("0")
	t.secsEntry.OnChanged = func(_ string) { t.saveSettings() }

	t.msEntry = widget.NewEntry()
	t.msEntry.SetText(strconv.Itoa(s.KeyPresserIntervalMs))
	t.msEntry.SetPlaceHolder("100")
	t.msEntry.OnChanged = func(_ string) { t.saveSettings() }

	intervalRow := container.NewHBox(
		container.NewVBox(widget.NewLabel("Hours"), t.hoursEntry),
		container.NewVBox(widget.NewLabel("Mins"), t.minsEntry),
		container.NewVBox(widget.NewLabel("Secs"), t.secsEntry),
		container.NewVBox(widget.NewLabel("Ms"), t.msEntry),
	)

	t.randomEntry = widget.NewEntry()
	t.randomEntry.SetText(strconv.Itoa(s.KeyPresserRandomOffsetMs))
	t.randomEntry.OnChanged = func(_ string) { t.saveSettings() }
	if !s.KeyPresserRandomEnabled {
		t.randomEntry.Disable()
	}

	t.randomCheck = widget.NewCheck("Random offset +/-", func(checked bool) {
		if checked {
			t.randomEntry.Enable()
		} else {
			t.randomEntry.Disable()
		}
		t.saveSettings()
	})
	t.randomCheck.SetChecked(s.KeyPresserRandomEnabled)

	randomRow := container.NewHBox(t.randomCheck, t.randomEntry, widget.NewLabel("ms"))

	intervalSection := container.NewVBox(
		widget.NewLabel("Press Interval"),
		intervalRow,
		randomRow,
		widget.NewSeparator(),
	)

	t.repeatEntry = widget.NewEntry()
	t.repeatEntry.SetText(strconv.Itoa(s.KeyPresserRepeatCount))
	t.repeatEntry.OnChanged = func(_ string) { t.saveSettings() }
	if s.KeyPresserRepeatMode != "Count" {
		t.repeatEntry.Disable()
	}

	t.repeatSelect = widget.NewRadioGroup([]string{"Until stopped", "Count"}, func(sel string) {
		if sel == "Count" {
			t.repeatEntry.Enable()
		} else {
			t.repeatEntry.Disable()
		}
		t.saveSettings()
	})
	t.repeatSelect.SetSelected(s.KeyPresserRepeatMode)

	repeatSection := container.NewVBox(
		widget.NewLabel("Press Repeat"),
		t.repeatSelect,
		container.NewHBox(widget.NewLabel("Repeat count:"), t.repeatEntry),
		widget.NewSeparator(),
	)

	t.startButton = widget.NewButton("START (F6)", func() {
		t.app.ToggleKeyPresser()
	})
	t.startButton.Importance = widget.HighImportance

	t.pressesLabel = widget.NewLabel("Presses: 0")

	controlSection := container.NewVBox(
		layout.NewSpacer(),
		t.startButton,
		t.pressesLabel,
	)

	t.content = container.NewVBox(
		keysSection,
		intervalSection,
		repeatSection,
		controlSection,
	)
}

func (t *KeyPresserTab) saveSettings() {
	// Don't save during initialization when widgets aren't fully created
	if t.repeatEntry == nil || t.randomCheck == nil || t.repeatSelect == nil {
		return
	}

	hours, _ := strconv.Atoi(t.hoursEntry.Text)
	mins, _ := strconv.Atoi(t.minsEntry.Text)
	secs, _ := strconv.Atoi(t.secsEntry.Text)
	ms, _ := strconv.Atoi(t.msEntry.Text)
	randomOffset, _ := strconv.Atoi(t.randomEntry.Text)
	repeatCount, _ := strconv.Atoi(t.repeatEntry.Text)

	settings.SaveKeyPresser(
		t.keysEntry.Text,
		hours, mins, secs, ms,
		t.randomCheck.Checked, randomOffset,
		t.repeatSelect.Selected, repeatCount,
	)
}

func (t *KeyPresserTab) applyConfig() error {
	cfg := keypresser.DefaultConfig()

	hours, _ := strconv.Atoi(t.hoursEntry.Text)
	mins, _ := strconv.Atoi(t.minsEntry.Text)
	secs, _ := strconv.Atoi(t.secsEntry.Text)
	ms, _ := strconv.Atoi(t.msEntry.Text)

	cfg.Interval = time.Duration(hours)*time.Hour +
		time.Duration(mins)*time.Minute +
		time.Duration(secs)*time.Second +
		time.Duration(ms)*time.Millisecond

	cfg.Keys = t.keysEntry.Text

	if t.randomCheck.Checked {
		cfg.RandomOffsetMs, _ = strconv.Atoi(t.randomEntry.Text)
	}

	switch t.repeatSelect.Selected {
	case "Until stopped":
		cfg.RepeatMode = keypresser.RepeatUntilStopped
	case "Count":
		cfg.RepeatMode = keypresser.RepeatCount
		cfg.RepeatCount, _ = strconv.Atoi(t.repeatEntry.Text)
	}

	if err := t.app.KeyPresser.SetConfig(cfg); err != nil {
		dialog.ShowError(err, t.window)
		return err
	}
	return nil
}

func (t *KeyPresserTab) setFieldsEnabled(enabled bool) {
	if enabled {
		t.keysEntry.Enable()
		t.hoursEntry.Enable()
		t.minsEntry.Enable()
		t.secsEntry.Enable()
		t.msEntry.Enable()
		t.randomCheck.Enable()
		if t.randomCheck.Checked {
			t.randomEntry.Enable()
		}
		t.repeatSelect.Enable()
		if t.repeatSelect.Selected == "Count" {
			t.repeatEntry.Enable()
		}
	} else {
		t.keysEntry.Disable()
		t.hoursEntry.Disable()
		t.minsEntry.Disable()
		t.secsEntry.Disable()
		t.msEntry.Disable()
		t.randomCheck.Disable()
		t.randomEntry.Disable()
		t.repeatSelect.Disable()
		t.repeatEntry.Disable()
	}
}

func (t *KeyPresserTab) UpdateState(running bool, pressCount int) {
	if running {
		t.startButton.SetText("STOP (F6)")
		t.startButton.Importance = widget.DangerImportance
	} else {
		t.startButton.SetText("START (F6)")
		t.startButton.Importance = widget.HighImportance
	}
	t.pressesLabel.SetText(fmt.Sprintf("Presses: %d", pressCount))
	t.setFieldsEnabled(!running)
}
