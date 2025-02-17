package ui

import (
	// "fmt"
	// "time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	// "fyne.io/fyne/v2/widget"
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
	wWidth := binding.NewFloat()
	wHeight := binding.NewFloat()

	// load sort order from preferences storage. default to ASC
	storedSortOrder := prefs.Bool("sort_order")
	if storedSortOrder {
		storedSortOrder = true
	}
	sortOrder.Set(storedSortOrder)

	// default window size accommodates changing of ASC-DESC without changing size of window (1140, 400)
	// load screen width from pref storage. default to 1140
	storedWWidth := prefs.Float("w_width")
	if storedWWidth == 0 {
		storedWWidth = 1140
	}
	wWidth.Set(storedWWidth)

	// load screen height from pref storage. default to 400
	storedWHeight := prefs.Float("w_height")
	if storedWHeight == 0 {
		storedWHeight = 400
	}
	wHeight.Set(storedWHeight)

	wW, _ := wWidth.Get()
	wH, _ := wHeight.Get()

	w.Resize(fyne.NewSize(float32(wW), float32(wH)))

	// the app will close when the main window (w) is closed
	w.SetMaster()

	// handle diagnostics...
	// hello := widget.NewLabel("Debugging stuff...")
	// w2.SetContent(
	// 	container.NewVBox(
	// 		hello,
	// 		widget.NewButton("Hi!", func() {
	// 			hello.SetText("Welcome")
	// 		}),
	// 		// w3 := a.NewWindow("Third")
	// 		// w3.SetContent(widget.NewLabel("Third"))
	// 		// w3.Show()
	// 	))

	//See diagram in documentation for clearer illustration
	//
	//--------toolbar--------
	//Search-TypingBoxSearch-
	//id|GameName|M|M+S|C
	//~~~~~~~~~~~~~~~~~~~~~~~
	//~~~~~~~~~~~~~~~~~~~~~~~

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
		// save sort order
		so, _ := sortOrder.Get()
		prefs.SetBool("sort_order", so)

		// save screen size
		width := w.Content().Size().Width
		height := w.Content().Size().Height

		wWidth.Set(float64(width))
		wW, _ = wWidth.Get()
		prefs.SetFloat("w_width", wW)

		wHeight.Set(float64(height))
		wH, _ = wHeight.Get()
		prefs.SetFloat("w_height", wH)
	})

	// runloop for the app
	a.Run()
}
