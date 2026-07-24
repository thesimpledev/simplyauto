package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"simplyauto/internal/assets"
)

type HelpTab struct {
	content fyne.CanvasObject
}

func NewHelpTab(version string) *HelpTab {
	t := &HelpTab{}

	header := container.NewVBox(
		widget.NewLabel("SimplyAuto v"+version),
		widget.NewLabel("Free and open source under the MIT License."),
		widget.NewSeparator(),
	)

	license := widget.NewLabel(assets.License())
	license.Wrapping = fyne.TextWrapWord

	notices := widget.NewLabel(assets.ThirdPartyNotices())
	notices.Wrapping = fyne.TextWrapWord

	// Both sections start collapsed so the large notices text is only
	// laid out when opened.
	licenses := widget.NewAccordion(
		widget.NewAccordionItem("SimplyAuto License (MIT)", license),
		widget.NewAccordionItem("Third-Party Notices", notices),
	)

	t.content = container.NewBorder(header, nil, nil, nil, container.NewVScroll(licenses))
	return t
}

func (t *HelpTab) Content() fyne.CanvasObject {
	return t.content
}
