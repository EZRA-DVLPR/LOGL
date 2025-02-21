package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createSearchBar(searchText binding.String) (searchBar *fyne.Container) {
	// widget with icon and text for clarity to user
	searchSymbolText := container.NewHBox(
		widget.NewIcon(theme.SearchIcon()),
		widget.NewLabel("Search"),
	)

	// create the textbox for user input and attach the binding to it
	searchTextBox := widget.NewEntryWithData(searchText)
	searchTextBox.SetPlaceHolder("Search Game Names Here!")

	// fit searchTextBox to fill rest of space to the right of searchSymbolText
	searchBar = container.NewBorder(
		nil,
		nil,
		searchSymbolText,
		nil,
		searchTextBox,
	)
	return searchBar
}
