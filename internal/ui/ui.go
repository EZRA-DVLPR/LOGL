package ui

import (
	// "fmt"
	// "time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
)

func StartGUI() {
	a := app.NewWithID("GameListID") // FIX: IDK what the id does or how to use it properly
	w := a.NewWindow("Main window - GameList")

	// set up prefs for app
	prefs := a.Preferences()

	// create all bindings here
	sortCategory := binding.NewString()
	sortOrder := binding.NewBool()
	dbData := NewMyDataBindingEmpty()
	wWidth := binding.NewFloat()
	wHeight := binding.NewFloat()

	// INFO: The following bindings do not persist through sessions
	searchText := binding.NewString()
	selectedRow := binding.NewInt()
	// dont highlight any row
	selectedRow.Set(-1)

	// load sort category from pref storage. default to "name"  i.e. Game Name
	storedSortCategory := prefs.String("sort_category")
	if storedSortCategory == "" {
		storedSortCategory = "name"
	}
	sortCategory.Set(storedSortCategory)

	// load sort order from preferences storage. default to true (ASC)
	storedSortOrder := prefs.BoolWithFallback("sort_order", true)
	sortOrder.Set(storedSortOrder)

	// TODO: Handle default sizes of window when i finalize the length/size of the toolbar with icons
	// default window size accommodates changing of "ASC"/"DESC" without changing size of window (1140, 400) (W,H)
	// It seems that the first row doesnt render properly initially if the width is too great...
	storedWWidth := prefs.Float("w_width")
	if storedWWidth == 0 {
		storedWWidth = 1140
	}
	wWidth.Set(storedWWidth)

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

	//See diagram in documentation for clearer illustration
	//
	//--------toolbar--------
	//Search-TypingBoxSearch-
	//id|GameName|M|M+S|C----
	//~~~~~~~~~~~~~~~~~~~~~~~
	//~~~~~~~~~~~~~~~~~~~~~~~
	content := container.NewBorder(
		// top is toolbar + searchbar
		container.NewVBox(
			createMainWindowToolbar(w.Canvas(), sortCategory, sortOrder, searchText, selectedRow, dbData),
			createSearchBar(searchText),
		),
		// dont render anything else in space besides DB
		nil, nil, nil,
		// default to display names ASC
		createDBRender(selectedRow, sortCategory, sortOrder, searchText, dbData),
	)

	w.SetContent(content)
	w.Show()

	// when main window closes, save preferences for future session
	w.SetOnClosed(func() {
		// save sort type
		st, _ := sortCategory.Get()
		prefs.SetString("sort_category", st)

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
