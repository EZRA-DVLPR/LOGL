package ui

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func StartGUI() {
	a := app.New()
	w := a.NewWindow("New Window")

	w.SetContent(widget.NewLabel("Content in the window"))
	w.ShowAndRun()
}
