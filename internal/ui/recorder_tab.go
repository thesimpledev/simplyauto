package ui

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"simplyauto/internal/app"
	"simplyauto/internal/recorder"
	intstorage "simplyauto/internal/storage"
)

// uriToPath converts a Fyne URI path to a native file path.
// On Windows, URI.Path() returns "/C:/path" which needs to be converted to "C:\path".
func uriToPath(uri fyne.URI) string {
	path := uri.Path()
	if runtime.GOOS == "windows" && len(path) > 2 && path[0] == '/' && path[2] == ':' {
		path = path[1:] // Remove leading slash: "/C:/path" -> "C:/path"
		path = strings.ReplaceAll(path, "/", "\\")
	}
	return path
}

type RecorderTab struct {
	app           *app.App
	window        fyne.Window
	recordButton  *widget.Button
	playButton    *widget.Button
	stopButton    *widget.Button
	speedSelect   *widget.Select
	loopSelect    *widget.RadioGroup
	loopEntry     *widget.Entry
	openButton    *widget.Button
	saveButton    *widget.Button
	fileLabel     *widget.Label
	eventsLabel   *widget.Label
	durationLabel *widget.Label
	progressBar   *widget.ProgressBar
	loopLabel     *widget.Label
	content       fyne.CanvasObject
}

func NewRecorderTab(app *app.App) *RecorderTab {
	t := &RecorderTab{app: app}
	t.build()

	// Register callback so config is always applied before playback starts
	app.OnPlaybackStart = t.applyPlaybackConfig

	return t
}

func (t *RecorderTab) SetWindow(w fyne.Window) {
	t.window = w
}

func (t *RecorderTab) Content() fyne.CanvasObject {
	return t.content
}

func (t *RecorderTab) build() {
	t.recordButton = widget.NewButton("RECORD (F9)", func() {
		t.app.ToggleRecording()
	})
	t.recordButton.Importance = widget.HighImportance

	t.playButton = widget.NewButton("PLAY (F10)", func() {
		t.app.TogglePlayback()
	})

	t.stopButton = widget.NewButton("STOP (F11)", func() {
		t.app.Stop()
	})
	t.stopButton.Importance = widget.DangerImportance

	controlRow := container.NewHBox(t.recordButton, t.playButton, t.stopButton)

	controlSection := container.NewVBox(
		widget.NewLabel("Recording Controls"),
		controlRow,
		widget.NewSeparator(),
	)

	t.speedSelect = widget.NewSelect([]string{"0.5x", "1x", "2x", "4x"}, nil)
	t.speedSelect.SetSelected("1x")

	t.loopEntry = widget.NewEntry()
	t.loopEntry.SetText("1")
	t.loopEntry.Disable()

	t.loopSelect = widget.NewRadioGroup([]string{"Once", "Count", "Continuous"}, func(s string) {
		if s == "Count" {
			t.loopEntry.Enable()
		} else {
			t.loopEntry.Disable()
		}
	})
	t.loopSelect.SetSelected("Once")

	playbackSection := container.NewVBox(
		widget.NewLabel("Playback Options"),
		container.NewHBox(widget.NewLabel("Speed:"), t.speedSelect),
		widget.NewLabel("Loop Mode:"),
		t.loopSelect,
		container.NewHBox(widget.NewLabel("Loop count:"), t.loopEntry),
		widget.NewSeparator(),
	)

	t.openButton = widget.NewButton("Open", func() {
		t.openFile()
	})

	t.saveButton = widget.NewButton("Save", func() {
		t.saveFile()
	})

	t.fileLabel = widget.NewLabel("No file loaded")

	fileSection := container.NewVBox(
		widget.NewLabel("File Operations"),
		container.NewHBox(t.openButton, t.saveButton),
		t.fileLabel,
		widget.NewSeparator(),
	)

	t.eventsLabel = widget.NewLabel("Events: 0")
	t.durationLabel = widget.NewLabel("Duration: 0.0s")
	t.progressBar = widget.NewProgressBar()
	t.progressBar.Hide()
	t.loopLabel = widget.NewLabel("")

	infoSection := container.NewVBox(
		widget.NewLabel("Recording Info"),
		container.NewHBox(t.eventsLabel, layout.NewSpacer(), t.durationLabel),
		t.progressBar,
		t.loopLabel,
	)

	t.content = container.NewVBox(
		controlSection,
		playbackSection,
		fileSection,
		infoSection,
	)
}

func (t *RecorderTab) applyPlaybackConfig() {
	var speed float64
	switch t.speedSelect.Selected {
	case "0.5x":
		speed = 0.5
	case "1x":
		speed = 1.0
	case "2x":
		speed = 2.0
	case "4x":
		speed = 4.0
	default:
		speed = 1.0
	}
	t.app.SetPlaybackSpeed(speed)

	var loopMode recorder.LoopMode
	var loopCount int
	switch t.loopSelect.Selected {
	case "Once":
		loopMode = recorder.LoopOnce
		loopCount = 1
	case "Count":
		loopMode = recorder.LoopCount
		loopCount, _ = strconv.Atoi(t.loopEntry.Text)
		if loopCount < 1 {
			loopCount = 1
		}
	case "Continuous":
		loopMode = recorder.LoopContinuous
		loopCount = 1
	}
	t.app.SetPlaybackLoop(loopMode, loopCount)
}

func (t *RecorderTab) UpdateRecordingState(recording bool, eventCount int) {
	if recording {
		t.recordButton.SetText("STOP RECORDING (F9)")
		t.recordButton.Importance = widget.DangerImportance
		t.playButton.Disable()
	} else {
		t.recordButton.SetText("RECORD (F9)")
		t.recordButton.Importance = widget.HighImportance
		t.playButton.Enable()
	}
	t.eventsLabel.SetText(fmt.Sprintf("Events: %d", eventCount))

	if !recording {
		if rec := t.app.GetCurrentRecording(); rec != nil {
			t.durationLabel.SetText(fmt.Sprintf("Duration: %.1fs", rec.Duration.Seconds()))
		}
	}
}

func (t *RecorderTab) UpdatePlaybackState(playing bool, progress, total, loop int) {
	if playing {
		t.playButton.SetText("PLAYING...")
		t.playButton.Disable()
		t.recordButton.Disable()
		t.progressBar.Show()

		if total > 0 {
			t.progressBar.SetValue(float64(progress) / float64(total))
		}
		t.loopLabel.SetText(fmt.Sprintf("Loop: %d", loop))
	} else {
		t.playButton.SetText("PLAY (F10)")
		t.playButton.Enable()
		t.recordButton.Enable()
		t.progressBar.Hide()
		t.loopLabel.SetText("")
	}
}

func (t *RecorderTab) UpdateFileButtons() {
	if t.app.IsIdle() {
		t.openButton.Enable()
	} else {
		t.openButton.Disable()
	}
}

func (t *RecorderTab) openFile() {
	if t.window == nil {
		return
	}

	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		path := uriToPath(reader.URI())
		if err := t.app.LoadRecording(path); err != nil {
			dialog.ShowError(err, t.window)
			return
		}

		t.fileLabel.SetText(reader.URI().Name())
		if rec := t.app.GetCurrentRecording(); rec != nil {
			t.eventsLabel.SetText(fmt.Sprintf("Events: %d", len(rec.Events)))
			t.durationLabel.SetText(fmt.Sprintf("Duration: %.1fs", rec.Duration.Seconds()))
		}
	}, t.window)

	fd.Show()
}

func (t *RecorderTab) saveFile() {
	if t.window == nil || !t.app.HasRecording() {
		return
	}

	fd := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		path := uriToPath(writer.URI())
		if err := t.app.SaveRecording(path); err != nil {
			dialog.ShowError(err, t.window)
			return
		}

		t.fileLabel.SetText(writer.URI().Name())
	}, t.window)

	fd.SetFileName("recording" + intstorage.FileExtension)
	fd.Show()
}
