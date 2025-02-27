package ui

import (
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/EZRA-DVLPR/GameList/internal/dbhandler"
	"github.com/EZRA-DVLPR/GameList/internal/scraper"
)

// window for popup that will be modified for the following functions
var w2 fyne.Window

// confirmation window for updating/deleting all db data
var w3 fyne.Window

func singleGameNameSearchPopup(
	a fyne.App,
	searchSource binding.String,
	sortCategory binding.String,
	sortOrder binding.Bool,
	searchText binding.String,
	dbData *MyDataBinding,
	selectedRow binding.Int,
) {
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
				if strings.TrimSpace(entry.Text) != "" {
					log.Println("Search for game data beginning!")
					ss, _ := searchSource.Get()
					// search game data then add to db
					dbhandler.SearchAddToDB(entry.Text, ss)

					// update dbData
					updateDBData(sortCategory, sortOrder, searchText, dbData)
					forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
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
			if strings.TrimSpace(gamename.Text) == "" ||
				strings.TrimSpace(main.Text) == "" ||
				strings.TrimSpace(mainplus.Text) == "" ||
				strings.TrimSpace(comp.Text) == "" {
				log.Println("Not enough game data given. Fill out top 4 fields")
			} else {
				if hltbURL.Text == "" {
					log.Println("No HLTB URL given for manual entry for game", strings.TrimSpace(gamename.Text))
				}
				if completionatorURL.Text == "" {
					log.Println("No Completionator URL given for manual entry for game", strings.TrimSpace(gamename.Text))
				}

				// check if main, mainplus, comp are valid floats
				mainfl, err := strconv.ParseFloat(main.Text, 64)
				if err != nil {
					log.Println("Improper value for Main Story. Make sure its a valid decimal.")
					w2.Close()
					return
				}
				mainplusfl, err := strconv.ParseFloat(mainplus.Text, 64)
				if err != nil {
					log.Println("Improper value for Main + Sides. Make sure its a valid decimal.")
					w2.Close()
					return
				}
				compfl, err := strconv.ParseFloat(comp.Text, 64)
				if err != nil {
					log.Println("Improper value for Completionist. Make sure its a valid decimal.")
					w2.Close()
					return
				}

				// insert the data into the db
				var newgame scraper.Game
				newgame.Name = strings.TrimSpace(gamename.Text)
				newgame.Main = float32(mainfl)
				newgame.MainPlus = float32(mainplusfl)
				newgame.Comp = float32(compfl)
				newgame.HLTBUrl = strings.TrimSpace(hltbURL.Text)
				newgame.CompletionatorUrl = strings.TrimSpace(completionatorURL.Text)
				newgame.Favorite = 0

				dbhandler.AddToDB(newgame)
				forceRenderDB(sortCategory, sortOrder, searchText, dbData, selectedRow)
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

func settingsPopup(
	a fyne.App,
	w fyne.Window,
	searchSource binding.String,
) {
	// if w2 already exists then focus it and complete task
	if w2 != nil {
		w2.RequestFocus()
		return
	}

	// define w2 properties
	w2 = a.NewWindow("Settings Window")
	w2.Resize(fyne.NewSize(400, 600))

	w2.SetContent(
		container.New(
			layout.NewVBoxLayout(),
			searchSourceRadioWidget(searchSource),
			widget.NewSeparator(),
			themeSelector(),
			widget.NewSeparator(),
			textSlider(),
			widget.NewSeparator(),
			updateAllButton(a),
			widget.NewSeparator(),
			storageLocationSelector(w),
			widget.NewSeparator(),
			deleteAllButton(a),
		),
	)
	w2.SetOnClosed(func() {
		w2 = nil
	})
	w2.Show()
}

// radio for selection of sources
func searchSourceRadioWidget(searchSource binding.String) *fyne.Container {
	label := widget.NewLabel("Search Source Selection")

	radio := widget.NewRadioGroup(
		[]string{
			"All",
			"HLTB",
			"Completionator",
		},
		func(value string) { searchSource.Set(value) },
	)

	// set default to search source saved
	ss, _ := searchSource.Get()
	radio.SetSelected(ss)

	return container.New(
		layout.NewVBoxLayout(),
		label,
		radio,
	)
}

// selector for the theme of the application
// PERF: allow custom themes
// TODO: Add binding to have a default theme
func themeSelector() *fyne.Container {
	label := widget.NewLabel("Theme Selection")

	selector := widget.NewSelect([]string{"Light", "Dark"}, func(value string) {
		log.Println("Select set to", value)
	})
	selector.SetSelected("Light")

	return container.New(
		layout.NewVBoxLayout(),
		label,
		selector,
	)
}

func textSlider() *fyne.Container {
	label := widget.NewLabel("Text Size changer")

	slider := widget.NewSlider(0, 10)
	slider.OnChangeEnded = func(res float64) {
		log.Println("Slider was released at", res)
	}
	return container.New(
		layout.NewVBoxLayout(),
		label,
		slider,
	)
}

func updateAllButton(a fyne.App) *fyne.Container {
	label := widget.NewLabel("Press this button to update all games within the database")

	updateAll := widget.NewButton("Update All", func() {
		if w3 != nil {
			w3.RequestFocus()
			return
		}
		w3 = a.NewWindow("Confirm Deletion of entire Database")
		w3.Resize(fyne.NewSize(400, 200))
		w3.SetContent(container.New(
			layout.NewGridLayout(2),
			widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
				w3.Close()
			}),
			widget.NewButtonWithIcon("Update all data", theme.ConfirmIcon(), func() {
				log.Println("update all entries in db")
				w3.Close()
			}),
		))
		w3.Show()
		w3.SetOnClosed(func() {
			w3 = nil
		})
	})
	return container.New(
		layout.NewVBoxLayout(),
		label,
		updateAll,
	)
}

func storageLocationSelector(w fyne.Window) *fyne.Container {
	label := widget.NewLabel("Button to change location to store database and export files")

	storageDir := widget.NewButton("select dir", func() {
		w.RequestFocus()
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				log.Println("error:", err)
				w2.Close()
				return
			}
			if uri == nil {
				log.Println("no dir selected")
				w2.Close()
				return
			}
			log.Println("selected dir:", uri)
			w2.Close()
		}, w)
	})

	return container.New(
		layout.NewVBoxLayout(),
		label,
		storageDir,
	)
}

func deleteAllButton(a fyne.App) *fyne.Container {
	label := widget.NewLabel("Button to delete all games from data")

	deleteAll := widget.NewButton("Delete All", func() {
		if w3 != nil {
			w3.RequestFocus()
			return
		}
		w3 = a.NewWindow("Confirm Deletion of entire Database")
		w3.Resize(fyne.NewSize(400, 200))
		w3.SetContent(container.New(
			layout.NewGridLayout(2),
			widget.NewButtonWithIcon("Cancel", theme.CancelIcon(), func() {
				w3.Close()
			}),
			widget.NewButtonWithIcon("Delete all data", theme.ConfirmIcon(), func() {
				log.Println("drop table and shite")
				w3.Close()
			}),
		))
		w3.Show()
		w3.SetOnClosed(func() {
			w3 = nil
		})
	})
	return container.New(
		layout.NewVBoxLayout(),
		label,
		deleteAll,
	)
}
