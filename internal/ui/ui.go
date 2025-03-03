package ui

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func StartGUI() {
	a := app.NewWithID(".EZRA-DVLPR.GameList")
	w := a.NewWindow("Main window - GameList")

	// set up prefs for app
	prefs := a.Preferences()

	// create all bindings here
	sortCategory := binding.NewString()
	sortOrder := binding.NewBool()
	dbData := NewMyDataBindingEmpty()
	wWidth := binding.NewFloat()
	wHeight := binding.NewFloat()
	searchSource := binding.NewString()
	textSize := binding.NewFloat()
	selectedTheme := binding.NewString()

	// INFO: The following bindings do not persist through sessions
	searchText := binding.NewString()
	selectedRow := binding.NewInt()
	// dont highlight any row
	selectedRow.Set(-1)

	// load sort category from pref storage. default to "name"  i.e. Game Name
	storedSortCategory := prefs.StringWithFallback("sort_category", "name")
	sortCategory.Set(storedSortCategory)

	// load sort order from preferences storage. default to true (ASC)
	storedSortOrder := prefs.BoolWithFallback("sort_order", true)
	sortOrder.Set(storedSortOrder)

	// load search sort from preferences storage. default to "All"
	storedSearchSort := prefs.StringWithFallback("search source", "All")
	searchSource.Set(storedSearchSort)

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

	// load available themes
	themesDir := "themes"
	availableThemes, err := loadAllThemes(themesDir)
	if err != nil {
		log.Fatal("Error loading themes from themes folder:", err)
	}

	// load saved text size from prefs
	ts := prefs.FloatWithFallback("text_size", 18)
	textSize.Set(ts)

	// load saved theme from prefs
	// TODO: Base light/dark as default based on system settings
	st := prefs.StringWithFallback("selected_theme", "Light")
	selectedTheme.Set(st)

	customTheme := &CustomTheme{
		Theme:    theme.DefaultTheme(),
		textSize: float32(ts),
		colors:   availableThemes[st],
	}

	// Create theme selection buttons
	themeLabel := widget.NewLabel(availableThemes[st].Name)
	themeButtonContainer := container.NewHBox()
	for themeName, themeColors := range availableThemes {
		button := widget.NewButton(themeName, func(name string, colors ColorTheme) func() {
			return func() {
				customTheme.colors = colors
				themeLabel.SetText(name)
				selectedTheme.Set(name)

				a.Settings().SetTheme(customTheme)
			}
		}(themeName, themeColors))

		// Add a color indicator next to the button
		colorPreview := canvas.NewRectangle(hexToColor(themeColors.Primary))
		colorPreview.SetMinSize(fyne.NewSize(20, 20))

		themeOption := container.NewHBox(button, colorPreview)
		themeButtonContainer.Add(themeOption)
	}

	a.Settings().SetTheme(customTheme)

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
			createMainWindowToolbar(w.Canvas(), sortCategory, sortOrder, searchText, selectedRow, dbData, searchSource, textSize, selectedTheme, a, w),
			createSearchBar(searchText),
		),
		// dont render anything else in space besides DB
		themeButtonContainer, nil, themeLabel,
		// nil, nil, nil
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

		// save search source
		ss, _ := searchSource.Get()
		prefs.SetString("search_source", ss)

		// save text size
		ts, _ := textSize.Get()
		prefs.SetFloat("text_size", ts)

		// save theme selection
		st, _ = selectedTheme.Get()
		prefs.SetString("selected_theme", st)

		// debugging what values saved are:
		// ts, _ = textSize.Get()
		// st, _ = selectedTheme.Get()
		// log.Println(ts, st)
	})

	// runloop for the app
	a.Run()
}
