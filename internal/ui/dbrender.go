package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	// "fyne.io/fyne/v2/theme"
	// "fyne.io/fyne/v2/widget"
)

// creates the DBRender that will be display the DB
func createDBRender() (dbRender *fyne.Container) {
	// id | Name | Main | Main + Sides | Comp
	dbRender = container.New(
		layout.NewGridLayout(5),
		canvas.NewText("id", color.White),
		canvas.NewText("Name", color.White),
		canvas.NewText("Main", color.White),
		canvas.NewText("Main + Sides", color.White),
		canvas.NewText("Completionist", color.White),
		canvas.NewText("1.", color.White),
		canvas.NewText("Sample Game Name", color.White),
		canvas.NewText("-1", color.White),
		canvas.NewText("10 Hours", color.White),
		canvas.NewText("145.6 Hours", color.White),
	)

	return dbRender
}
