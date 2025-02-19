package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func createSearchBar() (searchBar *fyne.Container) {
	searchButton := createSearchButton()
	searchText := createSearchTextBox()

	// allows the searchText (where user types) fit into the entire rest of the Hspace
	searchBar = container.NewBorder(nil, nil, searchButton, nil, searchText)
	return searchBar
}

// INFO: search is only for game names
// TODO: Decide if it will search whenever there is l la change in the textbox or when the user hits enter
func createSearchButton() (searchButton *widget.Button) {
	searchButton = widget.NewButtonWithIcon("Search", theme.SearchIcon(), func() {
		log.Println("show search bar when typing into this after clicking or pressing hotkey")
	})

	return searchButton
}

func createSearchTextBox() (searchTextBox *widget.Entry) {
	searchTextBox = widget.NewEntry()
	searchTextBox.SetPlaceHolder("Search Game Names Here!")
	// get the changes to the text
	searchTextBox.OnChanged = func(newtext string) {
		log.Println(newtext)
	}
	return
}
