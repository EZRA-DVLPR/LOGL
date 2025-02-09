package ui

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
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
	w := a.NewWindow("Clock")
	w2 := a.NewWindow("Debug Window")

	// defaults for windows to be set
	w.Resize(fyne.NewSize(200, 200)) // set the default size of the window upon start
	w.SetMaster()                    // the app will close when the main window is closed

	// actually modify the content of the main window
	clock := widget.NewLabel("")
	updateTime(clock)

	w.SetContent(clock)
	go func() {
		for range time.Tick(time.Second) {
			updateTime(clock)
		}
	}()

	// handle diagnostics...
	// w2.SetContent(widget.NewLabel("Debugging stuff..."))
	w2.SetContent(widget.NewButton("Open new window", func() {
		w3 := a.NewWindow("Third")
		w3.SetContent(widget.NewLabel("Third"))
		w3.Show()
	}))

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
