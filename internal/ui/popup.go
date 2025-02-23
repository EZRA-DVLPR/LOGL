package ui

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
	"github.com/EZRA-DVLPR/GameList/internal/scraper"
)

// window for popup that will be modified for the following functions
var w2 fyne.Window

func singleGameNameSearchPopup(a fyne.App, searchSource binding.String) {
	// if w2 already exists then focus it and complete task
	if w2 != nil {
		w2.RequestFocus()
		return
	}

	// define w2 properties
	w2 = a.NewWindow("Single Game Name Search")
	w2.Resize(fyne.NewSize(400, 80))

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Game Name to Search")
	w2.SetContent(
		container.New(
			layout.NewVBoxLayout(),
			entry,
			widget.NewButton("Begin Search", func() {
				// if entry is non-empty then perform search
				if entry.Text != "" {
					ss, _ := searchSource.Get()
					dbhandler.SearchAddToDB(entry.Text, ss)
				} else {
					log.Println("No game name given")
				}
				w2.Close()
			}),
		),
	)
	w2.SetOnClosed(func() {
		w2 = nil
	})
	w2.Show()
}

func manualEntryPopup(
	a fyne.App,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
	selectedRow binding.Int,
) {
	if w2 != nil {
		w2.RequestFocus()
		return
	}

	w2 = a.NewWindow("Manual Game Data Entry")
	w2.Resize(fyne.NewSize(400, 100))

	gamename := widget.NewEntry()
	main := widget.NewEntry()
	mainplus := widget.NewEntry()
	comp := widget.NewEntry()
	hltbURL := widget.NewEntry()
	completionatorURL := widget.NewEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{
				Text:   "Game Name",
				Widget: gamename,
			},
			{
				Text:   "Main (Hours)",
				Widget: main,
			},
			{
				Text:   "Main Plus Sides (Hours)",
				Widget: mainplus,
			},
			{
				Text:   "Completionist (Hours)",
				Widget: comp,
			},
			{
				Text:   "URL for HowLongToBeat",
				Widget: hltbURL,
			},
			{
				Text:   "URL for Completionator",
				Widget: completionatorURL,
			},
		},
		OnSubmit: func() {
			if gamename.Text == "" ||
				main.Text == "" ||
				mainplus.Text == "" ||
				comp.Text == "" {
				log.Println("Not enough game data given. Fill out top 4 fields")
			} else {
				if hltbURL.Text == "" {
					log.Println("No HLTB URL given")
				}
				if completionatorURL.Text == "" {
					log.Println("No Completionator URL given")
				}
				log.Println("Form submitted:", hltbURL.Text)
				log.Println("Form submitted:", completionatorURL.Text)

				// check if main, mainplus, comp are valid floats
				mainfl, err := strconv.ParseFloat(main.Text, 64)
				if err != nil {
					log.Println("Please provide proper decimal value for Main")
					w2.Close()
					return
				}
				mainplusfl, err := strconv.ParseFloat(mainplus.Text, 64)
				if err != nil {
					log.Println("Please provide proper decimal value for Main + Sides")
					w2.Close()
					return
				}
				compfl, err := strconv.ParseFloat(comp.Text, 64)
				if err != nil {
					log.Println("Please provide proper decimal value for Completionist")
					w2.Close()
					return
				}

				// insert the data into the db
				var newgame scraper.Game
				newgame.Name = gamename.Text
				newgame.Main = float32(mainfl)
				newgame.MainPlus = float32(mainplusfl)
				newgame.Comp = float32(compfl)
				newgame.HLTBUrl = hltbURL.Text
				newgame.CompletionatorUrl = completionatorURL.Text
				newgame.Favorite = 0

				dbhandler.AddToDB(newgame)
				updateDBData(sortCategory, sortOrder, searchText, dbData)

				// refresh the DB lol
				selectedRow.Set(0)
				selectedRow.Set(-1)
			}
			w2.Close()
		},
		OnCancel: func() {
			w2.Close()
		},
	}
	w2.SetContent(
		form,
	)
	w2.SetOnClosed(func() {
		w2 = nil
	})
	w2.Show()
}
