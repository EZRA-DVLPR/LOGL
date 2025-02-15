package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func StartGUI() {
	// setup for the entire gui portion of app
	// a is the whole app
	// w is the main window for interaction
	// w2 another window for other info that wont be interacted with
	//			think of it as an output terminal
	// im not sure if ill keep it in the final build
	a := app.New()
	w := a.NewWindow("Main window - GameList")
	w2 := a.NewWindow("Debug Window - GameList")

	// default window size accommodates changing of ASC-DESC without changing size of window
	w.Resize(fyne.NewSize(1140, 400))
	// the app will close when the main window (w) is closed
	w.SetMaster()

	// handle diagnostics...
	// w2.SetContent(widget.NewLabel("Debugging stuff..."))
	w2.SetContent(widget.NewButton("Open new window", func() {
		w3 := a.NewWindow("Third")
		w3.SetContent(widget.NewLabel("Third"))
		w3.Show()
	}))

	content := container.New(
		layout.NewVBoxLayout(),
		createMainWindowToolbar(w.Canvas()),
		createSearchBar(true),
		createDBRender(),
	)

	w.SetContent(content)
	// show all windows
	w.Show()
	w2.Show()

	// runloop for the app
	a.Run()
}

func updateTime(clock *widget.Label) {
	formatted := time.Now().Format("Time: 03:04:05")
	clock.SetText(formatted)
}

// creates the DBRender that will be display the DB
func createDBRender() (DBRender *widget.Label) {
	DBRender = widget.NewLabel("")
	updateTime(DBRender)

	go func() {
		for range time.Tick(time.Second) {
			updateTime(DBRender)
		}
	}()

	return
}
