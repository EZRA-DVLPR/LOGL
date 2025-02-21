package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createSearchBar(userText binding.String) (searchBar *fyne.Container) {
	// widget with icon and text for clarity to user
	searchSymbolText := container.NewHBox(widget.NewIcon(theme.SearchIcon()), widget.NewLabel("Search"))

	// create the textbox for user input and attach the binding to it
	searchText := createSearchTextBox()
	searchText.Bind(userText)

	// allows searchText (where user types) to fit into the entire rest of the Hspace
	searchBar = container.NewBorder(
		nil,
		nil,
		searchSymbolText,
		nil,
		searchText,
	)
	return searchBar
}

// makes the search text box widget
// PERF: remove this function and put it inside the searchBar function
func createSearchTextBox() (searchTextBox *widget.Entry) {
	searchTextBox = widget.NewEntry()
	searchTextBox.SetPlaceHolder("Search Game Names Here!")
	return
}
