// Package assets provides embedded application resources.
package assets

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed logo.png
var logoPNG []byte

// AppIcon returns the application icon as a Fyne resource.
func AppIcon() fyne.Resource {
	return fyne.NewStaticResource("logo.png", logoPNG)
}
