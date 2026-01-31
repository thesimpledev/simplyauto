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
	"simplyauto/internal/autoclicker"
	"simplyauto/internal/hooks"
	"simplyauto/pkg/events"
)

type AutoClickerTab struct {
	app               *app.App
	window            fyne.Window
	hoursEntry        *widget.Entry
	minsEntry         *widget.Entry
	secsEntry         *widget.Entry
	msEntry           *widget.Entry
	randomCheck       *widget.Check
	randomEntry       *widget.Entry
	buttonSelect      *widget.RadioGroup
	clickSelect       *widget.RadioGroup
	repeatSelect      *widget.RadioGroup
	repeatEntry       *widget.Entry
	positionSelect    *widget.RadioGroup
	xEntry            *widget.Entry
	yEntry            *widget.Entry
	setPositionButton *widget.Button
	startButton       *widget.Button
	clicksLabel       *widget.Label
	content           fyne.CanvasObject
	positionCapture   *hooks.InputCapture
}

func NewAutoClickerTab(app *app.App, window fyne.Window) *AutoClickerTab {
	t := &AutoClickerTab{app: app, window: window}
	t.build()

	// Register callback so config is always applied before autoclicker starts
	app.OnAutoClickerStart = t.applyConfig

	return t
}

func (t *AutoClickerTab) Content() fyne.CanvasObject {
	return t.content
}

func (t *AutoClickerTab) build() {
	t.hoursEntry = widget.NewEntry()
	t.hoursEntry.SetText("0")
	t.hoursEntry.SetPlaceHolder("0")

	t.minsEntry = widget.NewEntry()
	t.minsEntry.SetText("0")
	t.minsEntry.SetPlaceHolder("0")

	t.secsEntry = widget.NewEntry()
	t.secsEntry.SetText("0")
	t.secsEntry.SetPlaceHolder("0")

	t.msEntry = widget.NewEntry()
	t.msEntry.SetText("100")
	t.msEntry.SetPlaceHolder("100")

	intervalRow := container.NewHBox(
		container.NewVBox(widget.NewLabel("Hours"), t.hoursEntry),
		container.NewVBox(widget.NewLabel("Mins"), t.minsEntry),
		container.NewVBox(widget.NewLabel("Secs"), t.secsEntry),
		container.NewVBox(widget.NewLabel("Ms"), t.msEntry),
	)

	t.randomEntry = widget.NewEntry()
	t.randomEntry.SetText("0")
	t.randomEntry.Disable()

	t.randomCheck = widget.NewCheck("Random offset +/-", func(checked bool) {
		if checked {
			t.randomEntry.Enable()
		} else {
			t.randomEntry.Disable()
		}
	})

	randomRow := container.NewHBox(t.randomCheck, t.randomEntry, widget.NewLabel("ms"))

	intervalSection := container.NewVBox(
		widget.NewLabel("Click Interval"),
		intervalRow,
		randomRow,
		widget.NewSeparator(),
	)

	t.buttonSelect = widget.NewRadioGroup([]string{"Left", "Right", "Middle"}, nil)
	t.buttonSelect.SetSelected("Left")
	t.buttonSelect.Horizontal = true

	t.clickSelect = widget.NewRadioGroup([]string{"Single", "Double"}, nil)
	t.clickSelect.SetSelected("Single")
	t.clickSelect.Horizontal = true

	clickSection := container.NewVBox(
		widget.NewLabel("Click Options"),
		container.NewHBox(widget.NewLabel("Mouse Button:"), t.buttonSelect),
		container.NewHBox(widget.NewLabel("Click Type:"), t.clickSelect),
		widget.NewSeparator(),
	)

	t.repeatEntry = widget.NewEntry()
	t.repeatEntry.SetText("1")
	t.repeatEntry.Disable()

	t.repeatSelect = widget.NewRadioGroup([]string{"Until stopped", "Count"}, func(s string) {
		if s == "Count" {
			t.repeatEntry.Enable()
		} else {
			t.repeatEntry.Disable()
		}
	})
	t.repeatSelect.SetSelected("Until stopped")

	repeatSection := container.NewVBox(
		widget.NewLabel("Click Repeat"),
		t.repeatSelect,
		container.NewHBox(widget.NewLabel("Repeat count:"), t.repeatEntry),
		widget.NewSeparator(),
	)

	t.xEntry = widget.NewEntry()
	t.xEntry.SetText("0")
	t.xEntry.Disable()

	t.yEntry = widget.NewEntry()
	t.yEntry.SetText("0")
	t.yEntry.Disable()

	t.setPositionButton = widget.NewButton("Set Position", func() {
		t.startPositionCapture()
	})
	t.setPositionButton.Disable()

	t.positionSelect = widget.NewRadioGroup([]string{"Current location", "Fixed position"}, func(s string) {
		if s == "Fixed position" {
			t.xEntry.Enable()
			t.yEntry.Enable()
			t.setPositionButton.Enable()
		} else {
			t.xEntry.Disable()
			t.yEntry.Disable()
			t.setPositionButton.Disable()
		}
	})
	t.positionSelect.SetSelected("Current location")

	positionSection := container.NewVBox(
		widget.NewLabel("Cursor Position"),
		t.positionSelect,
		container.NewHBox(widget.NewLabel("X:"), t.xEntry, widget.NewLabel("Y:"), t.yEntry, t.setPositionButton),
		widget.NewSeparator(),
	)

	t.startButton = widget.NewButton("START (F6)", func() {
		t.app.ToggleAutoClicker()
	})
	t.startButton.Importance = widget.HighImportance

	t.clicksLabel = widget.NewLabel("Clicks: 0")

	controlSection := container.NewVBox(
		layout.NewSpacer(),
		t.startButton,
		t.clicksLabel,
	)

	t.content = container.NewVBox(
		intervalSection,
		clickSection,
		repeatSection,
		positionSection,
		controlSection,
	)
}

func (t *AutoClickerTab) applyConfig() error {
	cfg := autoclicker.DefaultConfig()

	hours, _ := strconv.Atoi(t.hoursEntry.Text)
	mins, _ := strconv.Atoi(t.minsEntry.Text)
	secs, _ := strconv.Atoi(t.secsEntry.Text)
	ms, _ := strconv.Atoi(t.msEntry.Text)

	cfg.Interval = time.Duration(hours)*time.Hour +
		time.Duration(mins)*time.Minute +
		time.Duration(secs)*time.Second +
		time.Duration(ms)*time.Millisecond

	if t.randomCheck.Checked {
		cfg.RandomOffsetMs, _ = strconv.Atoi(t.randomEntry.Text)
	}

	switch t.buttonSelect.Selected {
	case "Left":
		cfg.Button = events.MouseLeft
	case "Right":
		cfg.Button = events.MouseRight
	case "Middle":
		cfg.Button = events.MouseMiddle
	}

	switch t.clickSelect.Selected {
	case "Single":
		cfg.ClickType = events.ClickSingle
	case "Double":
		cfg.ClickType = events.ClickDouble
	}

	switch t.repeatSelect.Selected {
	case "Until stopped":
		cfg.RepeatMode = autoclicker.RepeatUntilStopped
	case "Count":
		cfg.RepeatMode = autoclicker.RepeatCount
		cfg.RepeatCount, _ = strconv.Atoi(t.repeatEntry.Text)
	}

	switch t.positionSelect.Selected {
	case "Current location":
		cfg.PositionMode = autoclicker.PositionCurrent
	case "Fixed position":
		cfg.PositionMode = autoclicker.PositionFixed
		cfg.FixedX, _ = strconv.Atoi(t.xEntry.Text)
		cfg.FixedY, _ = strconv.Atoi(t.yEntry.Text)
	}

	if err := t.app.AutoClicker.SetConfig(cfg); err != nil {
		dialog.ShowError(err, t.window)
		return err
	}
	return nil
}

func (t *AutoClickerTab) setFieldsEnabled(enabled bool) {
	if enabled {
		t.hoursEntry.Enable()
		t.minsEntry.Enable()
		t.secsEntry.Enable()
		t.msEntry.Enable()
		t.randomCheck.Enable()
		if t.randomCheck.Checked {
			t.randomEntry.Enable()
		}
		t.buttonSelect.Enable()
		t.clickSelect.Enable()
		t.repeatSelect.Enable()
		if t.repeatSelect.Selected == "Count" {
			t.repeatEntry.Enable()
		}
		t.positionSelect.Enable()
		if t.positionSelect.Selected == "Fixed position" {
			t.xEntry.Enable()
			t.yEntry.Enable()
			t.setPositionButton.Enable()
		}
	} else {
		t.hoursEntry.Disable()
		t.minsEntry.Disable()
		t.secsEntry.Disable()
		t.msEntry.Disable()
		t.randomCheck.Disable()
		t.randomEntry.Disable()
		t.buttonSelect.Disable()
		t.clickSelect.Disable()
		t.repeatSelect.Disable()
		t.repeatEntry.Disable()
		t.positionSelect.Disable()
		t.xEntry.Disable()
		t.yEntry.Disable()
		t.setPositionButton.Disable()
	}
}

func (t *AutoClickerTab) UpdateState(running bool, clickCount int) {
	if running {
		t.startButton.SetText("STOP (F6)")
		t.startButton.Importance = widget.DangerImportance
	} else {
		t.startButton.SetText("START (F6)")
		t.startButton.Importance = widget.HighImportance
	}
	t.clicksLabel.SetText(fmt.Sprintf("Clicks: %d", clickCount))
	t.setFieldsEnabled(!running)
}

func (t *AutoClickerTab) startPositionCapture() {
	// Stop any existing capture
	if t.positionCapture != nil {
		t.positionCapture.Stop()
	}

	t.setPositionButton.SetText("Click anywhere...")
	t.setPositionButton.Disable()

	t.positionCapture = hooks.NewInputCapture(hooks.CaptureOptions{
		CaptureMouse:    true,
		CaptureKeyboard: false,
	})

	err := t.positionCapture.Start(func(event events.InputEvent) {
		// Only capture left mouse button down
		if event.Type == events.EventMouseLeftDown {
			// Update the position entries on the UI thread
			t.xEntry.SetText(strconv.Itoa(event.X))
			t.yEntry.SetText(strconv.Itoa(event.Y))

			// Stop capturing
			t.positionCapture.Stop()
			t.positionCapture = nil

			// Reset button text
			t.setPositionButton.SetText("Set Position")
			t.setPositionButton.Enable()
		}
	})

	if err != nil {
		t.setPositionButton.SetText("Set Position")
		t.setPositionButton.Enable()
		dialog.ShowError(err, t.window)
	}
}
