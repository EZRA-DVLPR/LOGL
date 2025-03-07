package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// struct that has popupmenu that changes when theme changes
type ThemeAwareMenu struct {
	CurrentPopup *widget.PopUpMenu
	MenuItems    []*fyne.MenuItem
	Canvas       fyne.Canvas
	Position     fyne.Position
}

// constructor
func NewThemeAwareMenu(items []*fyne.MenuItem, canvas fyne.Canvas) *ThemeAwareMenu {
	return &ThemeAwareMenu{
		MenuItems: items,
		Canvas:    canvas,
	}
}

// display menu at pos
func (tam *ThemeAwareMenu) Show(pos fyne.Position) {
	// Close any existing popup first
	if tam.CurrentPopup != nil {
		tam.CurrentPopup.Hide()
	}

	// Create a fresh popup with the current theme applied
	menu := fyne.NewMenu("", tam.MenuItems...)
	tam.CurrentPopup = widget.NewPopUpMenu(menu, tam.Canvas)
	tam.Position = pos
	tam.CurrentPopup.ShowAtPosition(pos)
}

// refresh, recreate, redisplay
func (tam *ThemeAwareMenu) Refresh() {
	if tam.CurrentPopup != nil && tam.CurrentPopup.Visible() {
		tam.CurrentPopup.Hide()
		tam.Show(tam.Position)
	}
}
