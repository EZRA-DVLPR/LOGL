package ui

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"

	"github.com/EZRA-DVLPR/GameList/model"
)

// global vars holding the app and main window
var (
	a fyne.App
	w fyne.Window
)

func StartGUI() {
	version := "1.0.0"

	//	INFO:  set to default config location based on OS:
	// MacOS: $HOME/Library/Application Support
	// Windows: %AppData%
	// Linux: $XDG_CONFIG_HOME if non-empty, else $HOME/.config
	// Plan 9: $home/lib
	ucd, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Error fetching user config directory!", err)
	}
	// whatever path was decided above, append "LoGL" to it for usage for the app
	execPath := filepath.Join(ucd, "LOGL")

	// ensure LoGL dir exists
	if err := os.MkdirAll(execPath, os.ModePerm); err != nil {
		log.Fatal("Failed to create LoGL directory:", err)
	}
	os.Chdir(execPath)

	// set logging to open and write to a file
	logFile, err := setLogFile(version)
	if err != nil {
		fmt.Println("Error initializing logging process", err)
		os.Exit(1)
	}
	defer logFile.Close()

	a = app.NewWithID(".EZRA-DVLPR.GameList")
	w = a.NewWindow(fmt.Sprintf("Main window - GameList v%v", version))

	// set up prefs for app
	prefs := a.Preferences()

	// create all bindings here
	dbData := NewMyDataBindingEmpty()
	wWidth := binding.NewFloat()
	wHeight := binding.NewFloat()

	// load sort category from pref storage. default to "name" i.e. Game Name
	storedSortCategory := prefs.StringWithFallback("sort_category", "name")
	model.SetSortCategory(storedSortCategory)

	// load sort order from preferences storage. default to true (ASC)
	storedSortOrder := prefs.BoolWithFallback("sort_order", true)
	model.SetSortOrder(storedSortOrder)

	// load search sort from preferences storage. default to "All"
	storedSearchSort := prefs.StringWithFallback("search_source", "All")
	model.SetSearchSource(storedSearchSort)

	// default window size accommodates changing of "ASC"/"DESC" without changing size of window (1140, 400) (W,H)
	storedWWidth := prefs.FloatWithFallback("w_width", 1080)
	wWidth.Set(storedWWidth)
	storedWHeight := prefs.FloatWithFallback("w_height", 350)
	wHeight.Set(storedWHeight)

	wW, _ := wWidth.Get()
	wH, _ := wHeight.Get()

	w.Resize(fyne.NewSize(float32(wW), float32(wH)))

	// the app will close when the main window (w) is closed
	w.SetMaster()

	// load available themes from /themes dir
	availableThemes, err := loadAllThemes("themes")
	if err != nil {
		log.Fatal("Error loading themes from themes folder:", err)
	}

	// load saved text size from prefs
	ts := prefs.FloatWithFallback("text_size", 18)
	model.SetTextSize(ts)

	// load saved theme from prefs
	st := prefs.StringWithFallback("selected_theme", "Light")
	model.SetSelectedTheme(st)
	a.Settings().SetTheme(
		&CustomTheme{
			Theme:    theme.DefaultTheme(),
			textSize: float32(ts),
			colors:   availableThemes[st],
		},
	)

	// display the contents of the app
	content := container.NewBorder(
		// top is toolbar + searchbar
		container.NewVBox(
			createMainWindowToolbar(
				dbData,
				availableThemes,
			),
			createSearchBar(),
		),
		// dont render anything else in space besides DB
		nil, nil, nil,
		createDBRender(
			dbData,
			availableThemes,
		),
	)

	w.SetContent(content)
	w.Show()

	// when main window closes, save preferences for future sessions
	w.SetOnClosed(func() {
		// save sort type
		st, _ := model.GetSortCategory()
		prefs.SetString("sort_category", st)

		// save sort order
		so, _ := model.GetSortOrder()
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
		ss, _ := model.GetSearchSource()
		prefs.SetString("search_source", ss)

		// save text size
		ts, _ := model.GetTextSize()
		prefs.SetFloat("text_size", ts)

		// save theme selection
		sth, _ := model.GetSelectedTheme()
		prefs.SetString("selected_theme", sth)

		log.Println("Preferences saved:")
		log.Println("Sort Category:", st)
		log.Println("Sort Order:", so)
		log.Println("Sort Window Width:", wW)
		log.Println("Sort Window Height:", wH)
		log.Println("Search Source:", ss)
		log.Println("Text Size:", ts)
		log.Println("Selected Theme:", sth)
		log.Println("App closed!")
	})

	// runloop for the app
	a.Run()
}

// creates logfile based on: Version # and current time
func setLogFile(version string) (*os.File, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")

	logFileName := fmt.Sprintf("logs/GameList-V%s-%s.log", version, timestamp)

	// ensure logs dir exists
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// set logfile to be created in logs dir
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Set the log output to the file
	log.SetOutput(logFile)

	log.Println("App start!")
	return logFile, nil
}
