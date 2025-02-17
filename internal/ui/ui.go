package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

// setup for the entire gui portion of app
// a is the whole app
// w is the main window for interaction
// w2 another window for other info that wont be interacted with
//
//	think of it as an output terminal
//
// im not sure if ill keep it in the final build
func StartGUI() {
	a := app.NewWithID("GameListID") // FIX: IDK what the id does or how to use it properly
	w := a.NewWindow("Main window - GameList")
	// w2 := a.NewWindow("Debug Window - GameList")

	// set up prefs for app
	prefs := a.Preferences()

	// create all bindings here
	sortOrder := binding.NewBool()

	// load stored value from preferences storage. default to ASC if not found
	storedSortOrder := prefs.Bool("sort_order")
	if storedSortOrder {
		storedSortOrder = true
	}
	sortOrder.Set(storedSortOrder)

	// default window size accommodates changing of ASC-DESC without changing size of window
	w.Resize(fyne.NewSize(1140, 400))
	// the app will close when the main window (w) is closed
	w.SetMaster()

	// handle diagnostics...
	// w2.SetContent(widget.NewLabel("Debugging stuff..."))
	// w2.SetContent(widget.NewButton("Open new window", func() {
	// 	w3 := a.NewWindow("Third")
	// 	w3.SetContent(widget.NewLabel("Third"))
	// 	w3.Show()
	// }))

	//See diagram in documentation for clearer illustration
	//
	//--------toolbar--------
	//Search-TypingBoxSearch-
	//id|GameName|M|M+S|C
	//~~~~~~~~~~~~~~~~~~~~~~~
	//~~~~~~~~~~~~~~~~~~~~~~~

	// TODO: remove the need to pass the bindings
	// Instead make a function that will return all/some bindings to use inside and prevent passing unnecessary opts
	content := container.NewBorder(
		// top is toolbar + searchbar
		container.NewVBox(
			createMainWindowToolbar(w.Canvas(), sortOrder),
			createSearchBar(),
		),
		// dont render anything else in space besides DB
		nil, nil, nil,
		// default to display names ASC
		createDBRender("name", "ASC"),
	)

	// show all windows with their content
	w.SetContent(content)
	w.Show()
	// w2.Show()

	w.SetOnClosed(func() {
		// TODO: Make a function that will save all the values from the bindings to
		// the preferences list so the user can see them when they reopen app
		// i.e. persistent options saved b/n sessions
		val, _ := sortOrder.Get()
		prefs.SetBool("sort_order", val)
		fmt.Println("Saved last value for sort_order", val)
	})

	// runloop for the app
	a.Run()
}
