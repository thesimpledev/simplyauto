// Package assets provides embedded application resources.
package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed logo.png
var logoPNG []byte

//go:embed LICENSE.txt
var licenseText string

//go:embed THIRD-PARTY-NOTICES.txt
var thirdPartyNotices string

// AppIcon returns the application icon as a Fyne resource.
func AppIcon() fyne.Resource {
	return fyne.NewStaticResource("logo.png", logoPNG)
}

// License returns the SimplyAuto MIT license text.
func License() string {
	return licenseText
}

// ThirdPartyNotices returns the concatenated licenses of all vendored
// dependencies, generated from vendor/ (see docs/store-submission.md).
func ThirdPartyNotices() string {
	return thirdPartyNotices
}
